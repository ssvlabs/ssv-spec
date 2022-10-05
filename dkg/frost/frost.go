package frost

import (
	"math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
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
	identifier  dkg.RequestID
	operatorID  types.OperatorID
	threshold   uint32
	sessionSK   *ecies.PrivateKey
	participant *frost.DkgParticipant

	currentRound   DKGRound
	msgs           map[DKGRound]map[uint32]*dkg.SignedMessage
	operators      []uint32
	operatorShares map[uint32]*bls.SecretKey
}

type DKGRound int

const (
	Preparation DKGRound = iota + 1
	Round1
	Round2
	Blame
)

func New(
	network dkg.Network,
	operatorID types.OperatorID,
	requestID dkg.RequestID,
	signer types.DKGSigner,
	storage dkg.Storage,
) dkg.KeyGenProtocol {

	msgs := make(map[DKGRound]map[uint32]*dkg.SignedMessage)
	msgs[Preparation] = make(map[uint32]*dkg.SignedMessage)
	msgs[Round1] = make(map[uint32]*dkg.SignedMessage)
	msgs[Round2] = make(map[uint32]*dkg.SignedMessage)
	msgs[Blame] = make(map[uint32]*dkg.SignedMessage)

	return &FROST{
		network: network,
		signer:  signer,
		storage: storage,

		state: &State{
			identifier:     requestID,
			operatorID:     operatorID,
			msgs:           msgs,
			operatorShares: make(map[uint32]*bls.SecretKey),
		},
	}
}

func (fr *FROST) Start(init *dkg.Init) error {

	otherOperators := make([]uint32, 0)
	for _, operatorID := range init.OperatorIDs {
		if fr.state.operatorID == operatorID {
			continue
		}
		otherOperators = append(otherOperators, uint32(operatorID))
	}

	operators := []uint32{uint32(fr.state.operatorID)}
	operators = append(operators, otherOperators...)
	fr.state.operators = operators

	ctx := make([]byte, 16)
	if _, err := rand.Read(ctx); err != nil {
		return err
	}

	participant, err := frost.NewDkgParticipant(uint32(fr.state.operatorID), uint32(len(operators)), string(ctx), thisCurve, otherOperators...)
	if err != nil {
		return errors.Wrap(err, "failed to initialize a dkg participant")
	}

	fr.state.participant = participant
	fr.state.threshold = uint32(init.Threshold)

	k, err := ecies.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate session sk")
	}
	fr.state.sessionSK = k

	fr.state.currentRound = Preparation
	msg := &ProtocolMsg{
		Round: Preparation,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}
	return fr.broadcastDKGMessage(msg)
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {

	if err := msg.Validate(); err != nil {
		return false, nil, errors.Wrap(err, "failed to validate message signature")
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}

	if valid := protocolMessage.validate(); !valid {
		return false, nil, errors.New("failed to validate protocol message")
	}

	if fr.state.msgs[protocolMessage.Round] == nil {
		fr.state.msgs[protocolMessage.Round] = make(map[uint32]*dkg.SignedMessage)
	}

	originalMessage, ok := fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)]
	if ok {
		return false, nil, fr.createBlameTypeInconsistentMessageRequest(originalMessage, msg)
	}

	fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)] = msg

	switch protocolMessage.Round {
	case Preparation:
		if fr.canProceedThisRound(Round1) {
			if err := fr.processRound1(); err != nil {
				return false, nil, err
			}
		}
	case Round1:
		if fr.canProceedThisRound(Round2) {
			if err := fr.processRound2(); err != nil {
				return false, nil, err
			}
		}
	case Round2:
		if fr.canProceedThisRound(-1) { // -1 checks if protocol has finished with round 2
			out, err := fr.processKeygenOutput()
			if err != nil {
				return false, nil, err
			}
			return true, out, nil
		}
	case Blame:
		out, err := fr.processBlame()
		if err != nil {
			return false, nil, err
		}
		return true, &dkg.KeyGenOutput{BlameOutout: out}, err
	default:
		return true, nil, dkg.ErrInvalidRound{}
	}

	return false, nil, nil
}

func (fr *FROST) canProceedThisRound(thisRound DKGRound) bool {

	if thisRound == Preparation {
		return true
	}

	prevRound := thisRound - 1
	if thisRound < Preparation {
		prevRound = Round2
	}

	if fr.state.currentRound != prevRound {
		return false
	}

	// received msgs from all operators for last round
	for _, operatorID := range fr.state.operators {
		if _, ok := fr.state.msgs[prevRound][operatorID]; !ok {
			return false
		}
	}

	return true
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
