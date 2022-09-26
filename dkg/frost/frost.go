package frost

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

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
	identifier   dkg.RequestID
	network      dkg.Network
	signer       types.DKGSigner
	storage      dkg.Storage
	threshold    uint32
	currentRound DKGRound

	operatorID types.OperatorID
	operators  []uint32
	party      *frost.DkgParticipant
	sessionSK  *ecies.PrivateKey

	msgs           map[DKGRound]map[uint32]*dkg.SignedMessage
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
		identifier: requestID,
		network:    network,
		signer:     signer,
		storage:    storage,
		operatorID: operatorID,

		msgs:           msgs,
		operatorShares: make(map[uint32]*bls.SecretKey),
	}
}

func (fr *FROST) Start(init *dkg.Init) error {
	otherOperators := make([]uint32, 0)
	for _, operatorID := range init.OperatorIDs {
		if fr.operatorID == operatorID {
			continue
		}
		otherOperators = append(otherOperators, uint32(operatorID))
	}

	operators := []uint32{uint32(fr.operatorID)}
	operators = append(operators, otherOperators...)
	fr.operators = operators

	pctx := make([]byte, 16)
	_, err := rand.Read(pctx)
	if err != nil {
		return err
	}

	party, err := frost.NewDkgParticipant(uint32(fr.operatorID), uint32(len(operators)), string(pctx), thisCurve, otherOperators...)
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

	fr.currentRound = Preparation

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

	bcastPrepMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(fr.operatorID),
		Signature: nil,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.operatorID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("operator with id %d not found", fr.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastPrepMessage, operator.ETHAddress)
	if err != nil {
		return err
	}

	bcastPrepMessage.Signature = sig
	fr.msgs[Preparation][uint32(fr.operatorID)] = bcastPrepMessage

	return fr.network.BroadcastDKGMessage(bcastPrepMessage)
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	if err := msg.Validate(); err != nil {
		return false, nil, errors.Wrap(err, "failed to validate message signature")
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "could not decode protocol msg")
	}

	if fr.msgs[protocolMessage.Round] == nil {
		fr.msgs[protocolMessage.Round] = make(map[uint32]*dkg.SignedMessage)
	}

	originalMessage, ok := fr.msgs[protocolMessage.Round][uint32(msg.Signer)]
	if ok {
		return false, nil, fr.createBlameTypeInconsistentMessageRequest(originalMessage, msg)
	}

	fr.msgs[protocolMessage.Round][uint32(msg.Signer)] = msg

	switch protocolMessage.Round {
	case Preparation:
		if fr.canProceedRound1() {
			fr.currentRound = Round1
			if err := fr.processRound1(); err != nil {
				return false, nil, err
			}
		}
	case Round1:
		if fr.canProceedRound2() {
			fr.currentRound = Round2
			if err := fr.processRound2(); err != nil {
				return false, nil, err
			}
		}
	case Round2:
		if fr.canProceedKeygenOutput() {
			if _, err := fr.verifyShares(); err != nil {
				return false, nil, errors.Wrapf(err, "failed to combine t+1 verification key share (vk)")
			}

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
	}

	return false, nil, nil
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
		if uint32(fr.operatorID) == operatorID {
			continue
		}

		share := &bls.SecretKey{}
		shamirShare := p2pMessages[operatorID]
		if err := share.Deserialize(shamirShare.Value); err != nil {
			return err
		}

		fr.operatorShares[operatorID] = share

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
		Signer:    types.OperatorID(uint32(fr.operatorID)),
		Signature: nil,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.operatorID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("operator with id %d not found", fr.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastRound1Message, operator.ETHAddress)
	if err != nil {
		return err
	}

	bcastRound1Message.Signature = sig
	fr.msgs[Round1][uint32(fr.operatorID)] = bcastRound1Message

	return fr.network.BroadcastDKGMessage(bcastRound1Message)
}

func (fr *FROST) processRound2() error {
	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	for operatorID, dkgMessage := range fr.msgs[Round1] {

		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
			return errors.Wrap(err, "could not decode protocol msg")
		}

		verifiers := new(sharing.FeldmanVerifier)
		for _, commitmentBytes := range protocolMessage.Round1Message.Commitment {
			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
			if err != nil {
				return err
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

		if uint32(fr.operatorID) == operatorID {
			continue
		}

		shareBytes, err := ecies.Decrypt(fr.sessionSK, protocolMessage.Round1Message.Shares[uint32(fr.operatorID)])
		if err != nil {
			return err
		}

		share := &sharing.ShamirShare{
			Id:    uint32(fr.operatorID),
			Value: shareBytes,
		}

		p2psend[operatorID] = share

		if err := verifiers.Verify(share); err != nil {
			return fr.createBlameTypeInvalidShareRequest(operatorID)
		}
	}

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
		Signer:    types.OperatorID(uint32(fr.operatorID)),
		Signature: nil,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.operatorID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("operator with id %d not found", fr.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastRound2Message, operator.ETHAddress)
	if err != nil {
		return err
	}

	bcastRound2Message.Signature = sig
	fr.msgs[Round2][uint32(fr.operatorID)] = bcastRound2Message

	return fr.network.BroadcastDKGMessage(bcastRound2Message)
}

func (fr *FROST) processKeygenOutput() (*dkg.KeyGenOutput, error) {
	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(fr.msgs[Round2][uint32(fr.operatorID)].Message.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode protocol msg")
	}

	vk := protocolMessage.Round2Message.Vk

	sk := &bls.SecretKey{}
	if err := sk.Deserialize(fr.party.SkShare.Bytes()); err != nil {
		return nil, err
	}

	operatorPubKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.operators {
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
			return nil, err
		}

		operatorPubKeys[types.OperatorID(operatorID)] = pk
	}

	out := &dkg.KeyGenOutput{
		Share:           sk,
		OperatorPubKeys: operatorPubKeys,
		ValidatorPK:     vk,
		Threshold:       uint64(fr.threshold),
	}
	return out, nil
}

func (fr *FROST) processBlame() (*dkg.BlameOutput, error) {
	for operatorID, msg := range fr.msgs[Blame] {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(msg.Message.Data); err != nil {
			return nil, fmt.Errorf("failed to decode blame data")
		}

		var (
			valid bool
			err   error
		)

		switch protocolMessage.BlameMessage.Type {
		case InvalidShare:
			valid, err = fr.processBlameTypeInvalidShare(operatorID, protocolMessage.BlameMessage)
			if err != nil {
				return nil, err
			}
		case InconsistentMessage:
			valid, err = fr.processBlameTypeInconsistentMessage(operatorID, protocolMessage.BlameMessage)
			if err != nil {
				return nil, err
			}
		}

		blameMessageBytes, err := protocolMessage.BlameMessage.Encode()
		if err != nil {
			return nil, err
		}

		blameOutput := &dkg.BlameOutput{
			Valid:        valid,
			BlameMessage: blameMessageBytes,
		}
		return blameOutput, nil
	}
	return nil, nil
}

func (fr *FROST) processBlameTypeInvalidShare(operatorID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {
	round1Message := &Round1Message{}
	if err := json.Unmarshal(blameMessage.BlameData[0], round1Message); err != nil {
		return false, err
	}

	verifiers := new(sharing.FeldmanVerifier)
	for _, commitmentBytes := range round1Message.Commitment {
		commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
		if err != nil {
			return false, err
		}
		verifiers.Commitments = append(verifiers.Commitments, commitment)
	}

	blamerSessionSK := ecies.NewPrivateKeyFromBytes(blameMessage.BlamerSessionSk)
	shareBytes, err := ecies.Decrypt(blamerSessionSK, round1Message.Shares[operatorID])
	if err != nil {
		return false, err
	}

	share := &sharing.ShamirShare{
		Id:    operatorID,
		Value: shareBytes,
	}

	if err := verifiers.Verify(share); err != nil {
		return false, err
	}
	return true, nil
}

func (fr *FROST) processBlameTypeInconsistentMessage(operatorID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {
	return true, nil
}

func (fr *FROST) canProceedRound1() bool {
	if fr.currentRound != Preparation {
		return false
	}

	for _, operatorID := range fr.operators {
		protocolMessage := &ProtocolMsg{}

		msg, ok := fr.msgs[Preparation][operatorID]
		if ok {
			if err := protocolMessage.Decode(msg.Message.Data); err != nil {
				return false
			}
			if protocolMessage.PreparationMessage == nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (fr *FROST) canProceedRound2() bool {
	if fr.currentRound != Round1 {
		return false
	}

	for _, operatorID := range fr.operators {
		protocolMessage := &ProtocolMsg{}

		msg, ok := fr.msgs[Round1][operatorID]
		if ok {
			if err := protocolMessage.Decode(msg.Message.Data); err != nil {
				return false
			}
			if protocolMessage.Round1Message == nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (fr *FROST) canProceedKeygenOutput() bool {
	if fr.currentRound != Round2 {
		return false
	}

	for _, operatorID := range fr.operators {
		protocolMessage := &ProtocolMsg{}

		msg, ok := fr.msgs[Round2][operatorID]
		if ok {
			if err := protocolMessage.Decode(msg.Message.Data); err != nil {
				return false
			}
			if protocolMessage.Round2Message == nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (fr *FROST) createBlameTypeInconsistentMessageRequest(originalMessage, newMessage *dkg.SignedMessage) error {

	blameData := make([][]byte, 0)
	originalMessageBytes, err := originalMessage.Encode()
	if err != nil {
		return err
	}
	newMessageBytes, err := newMessage.Encode()
	if err != nil {
		return err
	}
	blameData = append(blameData, originalMessageBytes, newMessageBytes)

	protocolMessage := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InconsistentMessage,
			TargetOperatorID: uint32(newMessage.Signer),
			BlameData:        blameData,
			BlamerSessionSk:  fr.sessionSK.Bytes(),
		},
	}

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	bcastBlameMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(fr.operatorID),
		Signature: nil,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.operatorID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("operator with id %d not found", fr.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastBlameMessage, operator.ETHAddress)
	if err != nil {
		return err
	}

	bcastBlameMessage.Signature = sig
	fr.msgs[Blame][uint32(fr.operatorID)] = bcastBlameMessage

	return fr.network.BroadcastDKGMessage(bcastBlameMessage)
}

func (fr *FROST) createBlameTypeInvalidShareRequest(operatorID uint32) error {
	round1Bytes, err := fr.msgs[Round1][operatorID].Encode()
	if err != nil {
		return err
	}

	blameData := make([][]byte, 0)
	blameData = append(blameData, round1Bytes)

	protocolMessage := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidShare,
			TargetOperatorID: operatorID,
			BlameData:        blameData,
			BlamerSessionSk:  fr.sessionSK.Bytes(),
		},
	}

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	bcastBlameMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.identifier,
			Data:       protocolMessageBytes,
		},
		Signer:    types.OperatorID(uint32(fr.operatorID)),
		Signature: nil,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.operatorID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("operator with id %d not found", fr.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastBlameMessage, operator.ETHAddress)
	if err != nil {
		return err
	}

	bcastBlameMessage.Signature = sig
	fr.msgs[Blame][uint32(fr.operatorID)] = bcastBlameMessage

	return fr.network.BroadcastDKGMessage(bcastBlameMessage)
}

func (fr *FROST) verifyShares() ([]*bls.G1, error) {

	outputs := make([]*bls.G1, 0)

	for j := int(fr.threshold + 1); j < len(fr.operators); j++ {

		xVec := make([]bls.Fr, 0)
		yVec := make([]bls.G1, 0)

		for i := j - int(fr.threshold+1); i < j; i++ {
			operatorID := fr.operators[i]

			protocolMessage := &ProtocolMsg{}
			if err := protocolMessage.Decode(fr.msgs[Round2][operatorID].Message.Data); err != nil {
				return nil, errors.Wrap(err, "could not decode protocol msg")
			}

			x := bls.Fr{}
			x.SetInt64(int64(operatorID))
			xVec = append(xVec, x)

			pk := &bls.PublicKey{}
			if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
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
		fmt.Printf("vk: %x\n", outputs[i].Serialize())
		if !outputs[i].IsEqual(outputs[i-1]) {
			return nil, fmt.Errorf("failed to create consistent public key from t+1 shares")
		}
	}

	return outputs, nil
}

func (fr *FROST) encryptForP2PSend(id uint32, data []byte) ([]byte, error) {
	msg, ok := fr.msgs[Preparation][id]
	if !ok {
		return nil, fmt.Errorf("no public key found for operator %d", id)
	}

	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode protocol msg")
	}

	pk, err := ecies.NewPublicKeyFromBytes(protocolMessage.PreparationMessage.SessionPk)
	if err != nil {
		return nil, err
	}

	return ecies.Encrypt(pk, data)
}
