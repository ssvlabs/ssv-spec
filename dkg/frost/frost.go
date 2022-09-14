package frost

import (
	"encoding/json"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/coinbase/kryptology/pkg/sharing"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var thisCurve = curves.BLS12381G1()

type FROST struct {
	identifier dkg.RequestID
	network    dkg.Network
	party      *frost.DkgParticipant

	operatorID uint32
	operators  []uint32
	threshold  uint32

	sessionPK []byte
	// validatorPK []byte

	operatorShares map[uint32]*bls.SecretKey
	ownShare       []byte

	round DKGRound
	msgs  map[DKGRound]map[uint32]*ProtocolMsg
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
	}
}

func (fr *FROST) Start(init *dkg.Init) error {
	// todo create participiant in round 1
	otherOperators := make([]uint32, 0)
	for _, operatorID := range init.OperatorIDs {
		if fr.operatorID == uint32(operatorID) {
			continue
		}
		otherOperators = append(otherOperators, uint32(operatorID))
	}

	operators := []uint32{fr.operatorID}
	operators = append(operators, otherOperators...)

	ctx := "string to prevent replay attacks"
	party, err := frost.NewDkgParticipant(fr.operatorID, uint32(len(operators)), ctx, thisCurve, otherOperators...)
	if err != nil {
		return err
	}
	fr.party = party
	fr.threshold = uint32(init.Threshold)
	fr.round = Preparation

	protocolMessage := &ProtocolMsg{
		PreparationMessage: &PreparationMessage{
			SessionPk: fr.sessionPK,
		},
	}

	protocolMessageBytes, err := protocolMessage.Encode()
	if err != nil {
		return err
	}

	if fr.msgs[Preparation] == nil {
		fr.msgs[Preparation] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Preparation][fr.operatorID] = protocolMessage

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

	if fr.msgs[fr.round] == nil {
		fr.msgs[fr.round] = make(map[uint32]*ProtocolMsg)
	}

	fr.msgs[fr.round][uint32(msg.Signer)] = protocolMessage

	switch fr.round {
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
			return true, out, err
		}
	}
	return false, nil, nil
}

func (fr *FROST) canStartRound1() bool {
	canStart := true
	for _, operatorID := range fr.operators {
		_, ok := fr.msgs[Preparation][operatorID]
		if !ok {
			canStart = false
		}
	}
	return canStart
}

func (fr *FROST) canStartRound2() bool {
	canStart := true
	for _, operatorID := range fr.operators {
		_, ok := fr.msgs[Round1][operatorID]
		if !ok {
			canStart = false
		}
	}
	return canStart
}

func (fr *FROST) hasFinishedRound2() bool {
	hasFinished := true
	for _, operatorID := range fr.operators {
		_, ok := fr.msgs[Round2][operatorID]
		if !ok {
			hasFinished = false
		}
	}
	return hasFinished
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
		share := bls.SecretKey{}

		shamirShare := p2pMessages[operatorID]
		if err := share.Deserialize(shamirShare.Value); err != nil {
			return err
		}
		fr.operatorShares[operatorID] = &share

		if fr.operatorID != operatorID {
			shares[operatorID] = share.Serialize()
		}
	}

	r1Message := &Round1Message{
		Commitment: commitments,
		ProofS:     bCastMessage.Wi.Bytes(),
		ProofR:     bCastMessage.Ci.Bytes(),
		Shares:     shares,
	}

	protocolMessage := &ProtocolMsg{
		Round1Message: r1Message,
	}

	if fr.msgs[Round1] == nil {
		fr.msgs[Round1] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Round1][fr.operatorID] = protocolMessage

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

func (fr *FROST) processRound2() error {
	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	for operatorID, protocolMessage := range fr.msgs[Round1] {
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

		share := &sharing.ShamirShare{}
		if err := json.Unmarshal(protocolMessage.Round1Message.Shares[fr.operatorID], share); err != nil {
			return err
		}
		p2psend[operatorID] = share
	}

	bCastMessage, err := fr.party.Round2(bcast, p2psend)
	if err != nil {
		return err
	}

	fr.ownShare = fr.party.SkShare.Bytes()

	r2Message := &Round2Message{
		Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
		VkShare: bCastMessage.VkShare.ToAffineCompressed(),
	}

	protocolMessage := &ProtocolMsg{
		Round2Message: r2Message,
	}

	if fr.msgs[Round2] == nil {
		fr.msgs[Round2] = make(map[uint32]*ProtocolMsg)
	}
	fr.msgs[Round2][fr.operatorID] = protocolMessage

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
	sk.Deserialize(fr.ownShare)

	operatorPublicKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.operators {
		pk := &bls.PublicKey{}
		pk.Deserialize(fr.msgs[Round2][operatorID].Round2Message.VkShare)

		operatorPublicKeys[types.OperatorID(operatorID)] = pk
	}

	// TODO: Use G1LagrangeInterpolation to check whether vkshares generate consistent vk
	out := &dkg.KeyGenOutput{
		Share:           sk,
		OperatorPubKeys: operatorPublicKeys,
		ValidatorPK:     vk,
		Threshold:       uint64(fr.threshold),
	}
	return out, nil
}
