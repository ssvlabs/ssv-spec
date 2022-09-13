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
	operatorID uint32
	operators  []uint32
	threshold  uint32
	party      *frost.DkgParticipant

	validatorPK    []byte
	operatorShares map[uint32]*bls.SecretKey

	round DKGRound
	msgs  map[DKGRound]map[uint32]*ProtocolMsg
}

type DKGRound int

const (
	Preparation DKGRound = iota + 1
	Round1
	Round2
	// Blame
)

func New(
	requestID dkg.RequestID,
	network dkg.Network,
	i, t uint32,
	operators []uint32,
) dkg.KeyGenProtocol {

	otherOperators := make([]uint32, 0)
	for _, operatorID := range operators {
		if operatorID == i {
			continue
		}
		otherOperators = append(otherOperators, operatorID)
	}

	ctx := "string to prevent replay attacks"
	party, _ := frost.NewDkgParticipant(i, uint32(len(operators)), ctx, thisCurve, otherOperators...)

	return &FROST{
		identifier: requestID,
		network:    network,
		operatorID: i,
		threshold:  t,
		operators:  operators,
		party:      party,
	}
}

func (fr *FROST) Start(init *dkg.Init) error {
	return nil
}

func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	protocolMessage := &ProtocolMsg{}
	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "could not decode protocol msg")
	}

	if fr.msgs[protocolMessage.Round] == nil {
		fr.msgs[protocolMessage.Round] = make(map[uint32]*ProtocolMsg)
	}

	fr.msgs[protocolMessage.Round][uint32(msg.Signer)] = protocolMessage

	switch protocolMessage.Round {
	case Preparation:
		if fr.canStartRound1() {
			// do round 1
		}
	case Round1:
		if fr.canStartRound2() {
			// do round 2
		}
	case Round2:
		if fr.hasFinishedRound2() {
			// return keygen output
		}
	}
	return false, nil, nil
}

func (fr *FROST) canStartRound1() bool {
	canStart := true
	// check if we have public keys of all participants
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

	r1MessageBytes, err := json.Marshal(r1Message)
	if err != nil {
		return err
	}

	protocolMessage := &ProtocolMsg{
		Round: Round1,
		Data:  r1MessageBytes,
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
		if protocolMessage.Round != Round1 {
			return errors.New("invalid protocol message")
		}

		r1Message, ok := protocolMessage.Data.(*Round1Message)
		if !ok {
			return errors.New("invalid data type")
		}

		verifiers := new(sharing.FeldmanVerifier)
		for _, commitmentBytes := range r1Message.Commitment {
			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
			if err != nil {
				return err
			}
			verifiers.Commitments = append(verifiers.Commitments, commitment)
		}

		Wi, _ := thisCurve.Scalar.SetBytes(r1Message.ProofS)
		Ci, _ := thisCurve.Scalar.SetBytes(r1Message.ProofR)

		bcastMessage := &frost.Round1Bcast{
			Verifiers: verifiers,
			Wi:        Wi,
			Ci:        Ci,
		}
		bcast[operatorID] = bcastMessage

		share := &sharing.ShamirShare{}
		if err := json.Unmarshal(r1Message.Shares[fr.operatorID], share); err != nil {
			return err
		}
		p2psend[operatorID] = share
	}

	bCastMessage, err := fr.party.Round2(bcast, p2psend)
	if err != nil {
		return err
	}

	r2Message := &Round2Message{
		Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
		VkShare: bCastMessage.VkShare.ToAffineCompressed(),
	}

	r2MessageBytes, err := json.Marshal(r2Message)
	if err != nil {
		return err
	}

	protocolMessage := &ProtocolMsg{
		Round: Round2,
		Data:  r2MessageBytes,
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
