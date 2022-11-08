package frost

import (
	"math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var thisCurve = curves.BLS12381G1()

func init() {
	types.InitBLS()
}

type FROST struct {
	network dkg.Network
	signer  types.DKGSigner
	storage dkg.Storage

	state *State
}

type State struct {
	identifier dkg.RequestID
	operatorID types.OperatorID
	sessionSK  *ecies.PrivateKey

	threshold    uint32
	currentRound ProtocolRound
	participant  *frost.DkgParticipant

	operators      []uint32
	operatorsOld   []uint32
	operatorShares map[uint32]*bls.SecretKey

	msgs            ProtocolMessageStore
	oldKeyGenOutput *dkg.KeyGenOutput
}

type ProtocolRound int

const (
	Uninitialized ProtocolRound = iota
	Preparation
	Round1
	Round2
	KeygenOutput
	Blame
)

var rounds = []ProtocolRound{
	Uninitialized,
	Preparation,
	Round1,
	Round2,
	KeygenOutput,
	Blame,
}

type ProtocolMessageStore map[ProtocolRound]map[uint32]*dkg.SignedMessage

func newProtocolMessageStore() ProtocolMessageStore {
	m := make(map[ProtocolRound]map[uint32]*dkg.SignedMessage)
	for _, round := range rounds {
		m[round] = make(map[uint32]*dkg.SignedMessage)
	}
	return m
}

func New(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
	init *dkg.Init,
) dkg.Protocol {

	fr := &FROST{
		network: network,
		signer:  signer,
		storage: storage,
		state: &State{
			identifier:     requestID,
			operatorID:     operatorID,
			threshold:      uint32(init.Threshold),
			currentRound:   Uninitialized,
			operators:      types.OperatorList(init.OperatorIDs).ToUint32List(),
			operatorShares: make(map[uint32]*bls.SecretKey),
			msgs:           newProtocolMessageStore(),
		},
	}

	return fr
}

// Temporary, TODO: Remove and use interface with Reshare
func NewResharing(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
	operatorsOld []types.OperatorID,
	init *dkg.Reshare,
	output *dkg.KeyGenOutput,
) dkg.Protocol {

	return &FROST{
		network: network,
		signer:  signer,
		storage: storage,

		state: &State{
			identifier:      requestID,
			operatorID:      operatorID,
			threshold:       uint32(init.Threshold),
			currentRound:    Uninitialized,
			operators:       types.OperatorList(init.OperatorIDs).ToUint32List(),
			operatorsOld:    types.OperatorList(operatorsOld).ToUint32List(),
			operatorShares:  make(map[uint32]*bls.SecretKey),
			msgs:            newProtocolMessageStore(),
			oldKeyGenOutput: output,
		},
	}
}

// TODO: If Reshare, confirm participating operators using qbft before kick-starting this process.
func (fr *FROST) Start() error {
	fr.state.currentRound = Preparation

	ctx := make([]byte, 16)
	if _, err := rand.Read(ctx); err != nil {
		return err
	}
	participant, err := frost.NewDkgParticipant(uint32(fr.state.operatorID), fr.state.threshold, string(ctx), thisCurve, fr.state.operators...)
	if err != nil {
		return errors.Wrap(err, "failed to initialize a dkg participant")
	}
	fr.state.participant = participant

	if !fr.needToRunCurrentRound() {
		return nil
	}

	k, err := ecies.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate session sk")
	}
	fr.state.sessionSK = k

	msg := &ProtocolMsg{
		Round: fr.state.currentRound,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}
	return fr.broadcastDKGMessage(msg)
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {

	if err := fr.validateSignedMessage(msg); err != nil {
		return false, nil, errors.Wrap(err, "failed to validate signed message")
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}
	if valid := protocolMessage.validate(); !valid {
		return false, nil, errors.New("failed to validate protocol message")
	}

	existingMessage, ok := fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)]

	if isBlameTypeInconsisstent := ok && !fr.haveSameRoot(existingMessage, msg); isBlameTypeInconsisstent {
		if err := fr.createAndBroadcastBlameOfInconsistentMessage(existingMessage, msg); err != nil {
			return false, nil, err
		}
		if blame, err := fr.processBlame(); err != nil {
			return false, nil, err
		} else {
			return true, &dkg.ProtocolOutcome{BlameOutput: blame}, nil
		}
	}

	fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)] = msg

	switch protocolMessage.Round {
	case Preparation:
		if fr.canProceedThisRound() {
			fr.state.currentRound = Round1
			if err := fr.processRound1(); err != nil {
				return false, nil, err
			}
		}
	case Round1:
		if fr.canProceedThisRound() {
			fr.state.currentRound = Round2
			if err := fr.processRound2(); err != nil {
				if err.Error() == "invalid share" {
					return true, &dkg.ProtocolOutcome{BlameOutput: err.(ErrInvalidShare).BlameOutput}, nil
				}
				return false, nil, err
			}
		}
	case Round2:
		if fr.canProceedThisRound() {
			fr.state.currentRound = KeygenOutput
			out, err := fr.processKeygenOutput()
			if err != nil {
				return false, nil, err
			}
			return true, &dkg.ProtocolOutcome{ProtocolOutput: out}, nil
		}
	case Blame:
		fr.state.currentRound = Blame
		blame, err := fr.processBlame()
		if err != nil {
			return false, nil, err
		}
		return true, &dkg.ProtocolOutcome{BlameOutput: blame}, nil
	default:
		return true, nil, dkg.ErrInvalidRound{}
	}

	return false, nil, nil
}

func (fr *FROST) canProceedThisRound() bool {
	// Note: for Resharing, Preparation (New Committee) -> Round1 (Old Committee) -> Round2 (New Committee)
	if fr.isResharing() && fr.state.currentRound == Round1 {
		return fr.allMessagesReceivedFor(Round1, fr.state.operatorsOld)
	}
	return fr.allMessagesReceivedFor(fr.state.currentRound, fr.state.operators)
}

func (fr *FROST) allMessagesReceivedFor(round ProtocolRound, operators []uint32) bool {
	for _, operatorID := range operators {
		if _, ok := fr.state.msgs[round][operatorID]; !ok {
			return false
		}
	}
	return true
}

func (fr *FROST) isResharing() bool {
	return len(fr.state.operatorsOld) > 0
}

func (fr *FROST) inOldCommittee() bool {
	for _, id := range fr.state.operatorsOld {
		if types.OperatorID(id) == fr.state.operatorID {
			return true
		}
	}
	return false
}

func (fr *FROST) inNewCommittee() bool {
	for _, id := range fr.state.operators {
		if types.OperatorID(id) == fr.state.operatorID {
			return true
		}
	}
	return false
}

func (fr *FROST) needToRunCurrentRound() bool {
	if !fr.isResharing() {
		return true // always run for new keygen
	}
	switch fr.state.currentRound {
	case Preparation, Round2, KeygenOutput:
		return fr.inNewCommittee()
	case Round1:
		return fr.inOldCommittee()
	default:
		return false
	}
}

func (fr *FROST) validateSignedMessage(msg *dkg.SignedMessage) error {
	if msg.Message.Identifier != fr.state.identifier {
		return errors.New("got mismatching identifier")
	}

	found, operator, err := fr.storage.GetDKGOperator(msg.Signer)
	if !found {
		return errors.New("unable to find signer")
	}
	if err != nil {
		return errors.Wrap(err, "unable to find signer")
	}

	root, err := msg.Message.GetRoot()
	if err != nil {
		return errors.Wrap(err, "failed to get root")
	}

	pk, err := crypto.Ecrecover(root, msg.Signature)
	if err != nil {
		return errors.Wrap(err, "unable to recover public key")
	}

	addr := common.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
	if addr != operator.ETHAddress {
		return errors.New("invalid signature")
	}
	return nil
}

func (fr *FROST) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
	msg, ok := fr.state.msgs[Preparation][operatorID]
	if !ok {
		return nil, errors.New("no session pk found for the operator")
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return nil, errors.Wrap(err, "failed to decode protocol msg")
	}

	sessionPK, err := ecies.NewPublicKeyFromBytes(protocolMessage.PreparationMessage.SessionPk)
	if err != nil {
		return nil, err
	}

	return ecies.Encrypt(sessionPK, data)
}

func (fr *FROST) toSignedMessage(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
	msgBytes, err := msg.Encode()
	if err != nil {
		return nil, err
	}

	bcastMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.state.identifier,
			Data:       msgBytes,
		},
		Signer: fr.state.operatorID,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.state.operatorID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.Errorf("operator with id %d not found", fr.state.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
	if err != nil {
		return nil, err
	}
	bcastMessage.Signature = sig
	return bcastMessage, nil
}

func (fr *FROST) broadcastDKGMessage(msg *ProtocolMsg) error {
	bcastMessage, err := fr.toSignedMessage(msg)
	if err != nil {
		return err
	}
	fr.state.msgs[fr.state.currentRound][uint32(fr.state.operatorID)] = bcastMessage
	return fr.network.BroadcastDKGMessage(bcastMessage)
}
