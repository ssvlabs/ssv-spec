package frost

import (
	"fmt"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var thisCurve = curves.BLS12381G1()

func init() {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)
}

type FROST struct {
	identifier dkg.RequestID
	network    dkg.Network
	threshold  uint32
	round      DKGRound

	operatorID uint32
	operators  []uint32
	party      *frost.DkgParticipant
	sessionSK  *ecies.PrivateKey

	opShareLock    *sync.Mutex
	operatorShares map[uint32]*bls.SecretKey

	msgLock *sync.Mutex
	msgs    map[DKGRound]map[uint32]*ProtocolMsg
}

type DKGRound int

const (
	Preparation DKGRound = iota + 1
	Round1
	Round2
	Blame
)

func New(
	requestID dkg.RequestID,
	network dkg.Network,
	operatorID uint32,
) dkg.KeyGenProtocol {

	return &FROST{
		identifier: requestID,
		network:    network,
		operatorID: operatorID,

		msgLock:        &sync.Mutex{},
		msgs:           make(map[DKGRound]map[uint32]*ProtocolMsg),
		opShareLock:    &sync.Mutex{},
		operatorShares: make(map[uint32]*bls.SecretKey),
	}
}

func (fr *FROST) Start(init *dkg.Init) error {
	otherOperators := make([]uint32, 0)
	for _, operatorID := range init.OperatorIDs {
		if fr.operatorID == uint32(operatorID) {
			continue
		}
		otherOperators = append(otherOperators, uint32(operatorID))
	}

	operators := []uint32{fr.operatorID}
	operators = append(operators, otherOperators...)
	fr.operators = operators

	ctx := "string to prevent replay attacks"
	party, err := frost.NewDkgParticipant(fr.operatorID, uint32(len(operators)), ctx, thisCurve, otherOperators...)
	if err != nil {
		return err
	}

	fr.party = party
	fr.threshold = uint32(init.Threshold)

	k, err := ecies.GenerateKey()
	if err != nil {
		return err
	}
	fr.sessionSK = k

	fr.round = Preparation
	protocolMessage := &ProtocolMsg{
		Round: Preparation,
		PreparationMessage: &PreparationMessage{
			SessionPk: k.PublicKey.Bytes(true),
		},
	}

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	fr.msgLock.Lock()
	if fr.msgs[Preparation] == nil {
		fr.msgs[Preparation] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Preparation][fr.operatorID] = protocolMessage
	fr.msgLock.Unlock()

	bcastPrepMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(fr.operatorID),
		Signature: nil,
	}

	return fr.network.BroadcastDKGMessage(bcastPrepMessage)
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "could not decode protocol msg")
	}

	fr.msgLock.Lock()
	if fr.msgs[protocolMessage.Round] == nil {
		fr.msgs[protocolMessage.Round] = make(map[uint32]*ProtocolMsg)
	}

	fr.msgs[protocolMessage.Round][uint32(msg.Signer)] = protocolMessage
	fr.msgLock.Unlock()

	switch protocolMessage.Round {
	case Preparation:
		if fr.canStartRound1() {
			fr.round = Round1
			if err := fr.processRound1(); err != nil {
				return false, nil, err
			}
		}
	case Round1:
		if fr.canStartRound2() {
			fr.round = Round2
			if err := fr.processRound2(); err != nil {
				return false, nil, err
			}
		}
	case Round2:
		if fr.hasFinishedRound2() {
			out, err := fr.getKeygenOutput()
			if err != nil {
				return false, out, err
			}
			return true, out, nil
		}
	}
	return false, nil, nil
}

func (fr *FROST) canStartRound1() bool {
	if fr.round != Preparation {
		return false
	}

	fr.msgLock.Lock()
	defer fr.msgLock.Unlock()
	for _, operatorID := range fr.operators {
		msg, ok := fr.msgs[Preparation][operatorID]
		if !ok || msg.PreparationMessage == nil {
			return false
		}
	}
	return true
}

func (fr *FROST) encryptForP2PSend(id uint32, share []byte) ([]byte, error) {
	fr.msgLock.Lock()
	defer fr.msgLock.Unlock()
	msg, ok := fr.msgs[Preparation][id]
	if !ok {
		return nil, fmt.Errorf("no public key found for operator %d", id)
	}

	pk, err := ecies.NewPublicKeyFromBytes(msg.PreparationMessage.SessionPk)
	if err != nil {
		return nil, err
	}

	return ecies.Encrypt(pk, share)
}

func (fr *FROST) processRound1() error {
	bCastMessage, p2pMessages, err := fr.party.Round1(nil)
	if err != nil {
		return err
	}

	commitments := make([][]byte, 0)
	for _, commitment := range bCastMessage.Verifiers.Commitments {
		commitments = append(commitments, commitment.ToAffineCompressed())
	}

	shares := make(map[uint32][]byte)
	for _, operatorID := range fr.operators {
		if fr.operatorID == operatorID {
			continue
		}

		share := &bls.SecretKey{}
		shamirShare := p2pMessages[operatorID]
		if err := share.Deserialize(shamirShare.Value); err != nil {
			return err
		}

		fr.opShareLock.Lock()
		fr.operatorShares[operatorID] = share
		fr.opShareLock.Unlock()

		encryptedShare, err := fr.encryptForP2PSend(operatorID, shamirShare.Value)
		if err != nil {
			return err
		}
		shares[operatorID] = encryptedShare
	}

	protocolMessage := &ProtocolMsg{
		Round: Round1,
		Round1Message: &Round1Message{
			Commitment: commitments,
			ProofS:     bCastMessage.Wi.Bytes(),
			ProofR:     bCastMessage.Ci.Bytes(),
			Shares:     shares,
		},
	}

	fr.msgLock.Lock()
	if fr.msgs[Round1] == nil {
		fr.msgs[Round1] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Round1][fr.operatorID] = protocolMessage
	fr.msgLock.Unlock()

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	bcastRound1Message := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(fr.operatorID),
		Signature: nil,
	}

	return fr.network.BroadcastDKGMessage(bcastRound1Message)
}

func (fr *FROST) canStartRound2() bool {
	if fr.round != Round1 {
		return false
	}

	fr.msgLock.Lock()
	defer fr.msgLock.Unlock()
	for _, operatorID := range fr.operators {
		if fr.operatorID == operatorID {
			continue
		}

		msg, ok := fr.msgs[Round1][operatorID]
		if !ok || (ok && msg.Round1Message == nil) {
			return false
		}
	}
	return true
}

func (fr *FROST) hasFinishedRound2() bool {
	if fr.round != Round2 {
		return false
	}

	fr.msgLock.Lock()
	defer fr.msgLock.Unlock()
	for _, operatorID := range fr.operators {
		msg, ok := fr.msgs[Round2][operatorID]
		if !ok || msg.Round2Message == nil {
			return false
		}
	}
	return true
}

func (fr *FROST) processRound2() error {
	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	fr.msgLock.Lock()
	for operatorID, protocolMessage := range fr.msgs[Round1] {

		verifiers := new(sharing.FeldmanVerifier)
		for _, commitmentBytes := range protocolMessage.Round1Message.Commitment {
			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
			if err != nil {
				return errors.Wrapf(err, "failed to parse commitment for operator %d", fr.operatorID)
			}
			verifiers.Commitments = append(verifiers.Commitments, commitment)
		}

		Wi, _ := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofS)
		Ci, _ := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofR)

		bcastMessage := &frost.Round1Bcast{
			Verifiers: verifiers,
			Wi:        Wi,
			Ci:        Ci,
		}
		bcast[operatorID] = bcastMessage

		if fr.operatorID == operatorID {
			continue
		}

		shareBytes, err := ecies.Decrypt(fr.sessionSK, protocolMessage.Round1Message.Shares[fr.operatorID])
		if err != nil {
			return err
		}

		share := &sharing.ShamirShare{
			Id:    fr.operatorID,
			Value: shareBytes,
		}

		p2psend[operatorID] = share
	}
	fr.msgLock.Unlock()

	bCastMessage, err := fr.party.Round2(bcast, p2psend)
	if err != nil {
		return err
	}

	protocolMessage := &ProtocolMsg{
		Round: Round2,
		Round2Message: &Round2Message{
			Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
			VkShare: bCastMessage.VkShare.ToAffineCompressed(),
		},
	}

	fr.msgLock.Lock()
	if fr.msgs[Round2] == nil {
		fr.msgs[Round2] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Round2][fr.operatorID] = protocolMessage
	fr.msgLock.Unlock()

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	bcastRound2Message := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(fr.operatorID),
		Signature: nil,
	}

	return fr.network.BroadcastDKGMessage(bcastRound2Message)
}

func (fr *FROST) getKeygenOutput() (*dkg.KeyGenOutput, error) {
	if fr.round != Round2 {
		return nil, dkg.ErrInvalidRound{}
	}

	vk := fr.msgs[Round2][fr.operatorID].Round2Message.Vk

	sk := &bls.SecretKey{}
	if err := sk.Deserialize(fr.party.SkShare.Bytes()); err != nil {
		return nil, err
	}

	operatorPublicKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.operators {
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(fr.msgs[Round2][operatorID].Round2Message.VkShare); err != nil {
			return nil, err
		}

		operatorPublicKeys[types.OperatorID(operatorID)] = pk
	}

	outputs := make([]*bls.G1, 0)
	for j := int(fr.threshold + 1); j < len(fr.operators); j++ {
		xVec := make([]bls.Fr, 0)
		yVec := make([]bls.G1, 0)
		for i := j - int(fr.threshold); i < j; i++ {
			operatorID := fr.operators[i]

			x := bls.Fr{}
			x.SetInt64(int64(operatorID))
			xVec = append(xVec, x)

			pk := &bls.PublicKey{}
			if err := pk.Deserialize(fr.msgs[Round2][operatorID].Round2Message.VkShare); err != nil {
				return nil, err
			}

			y := bls.CastFromPublicKey(pk)
			yVec = append(yVec, *y)
		}

		out := &bls.G1{}
		if err := bls.G1LagrangeInterpolation(out, xVec, yVec); err != nil {
			return nil, err
		}

		outputs = append(outputs, out)
	}

	for i := 1; i < len(outputs); i++ {

		if !outputs[i].IsEqual(outputs[i-1]) {
			return nil, fmt.Errorf("failed to create consistent public keys from t+1 shares")
		}
	}

	out := &dkg.KeyGenOutput{
		Share:           sk,
		OperatorPubKeys: operatorPublicKeys,
		ValidatorPK:     vk,
		Threshold:       uint64(fr.threshold),
	}
	return out, nil
}
