package frost

import (
	"encoding/json"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/herumi/bls-eth-go-binary/bls"
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

	keygenMessage map[uint32]KeygenMessage
}

type KeygenMessage struct {
	Round1 Round1Message
	Round2 Round2Message
}

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

func (frost *FROST) Start(init *dkg.Init) error {
	// TODO: use init message to initialize number of operators and threshold
	bCastMessage, p2pMessages, err := frost.party.Round1(nil)
	if err != nil {
		return err
	}

	commitments := make([][]byte, 0)
	for _, commitment := range bCastMessage.Verifiers.Commitments {
		commitments = append(commitments, commitment.Scalar().Bytes())
	}

	shares := make([][]byte, 0)
	for _, operatorID := range frost.operators {
		shamirShare := p2pMessages[operatorID]

		share := bls.SecretKey{}
		err := share.Deserialize(shamirShare.Value)
		if err != nil {
			return err
		}

		frost.operatorShares[operatorID] = &share
		shares = append(shares, share.Serialize())
	}

	r1Message := &Round1Message{
		Commitment: commitments,
		Shares:     shares,
	}

	r1MessageBytes, err := json.Marshal(r1Message)
	if err != nil {
		return err
	}

	bcastRound1Message := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: frost.identifier,
			Data:       r1MessageBytes,
		},
		Signer:    types.OperatorID(frost.operatorID),
		Signature: nil,
	}

	return frost.network.BroadcastDKGMessage(bcastRound1Message)
}

func (frost *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	return false, nil, nil
}
