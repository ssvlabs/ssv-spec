package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// processRound1 function can only be executed once all messages for the
// Preparation round have been received. It runs round 1 in the frost library and
// returns shares, commitments, and Schnorr's proof. These elements, including
// encrypted shares, are then serialized into a protocol message and broadcasted
// over the network. In the case of resharing, the secret from the old keygen
// output is used for splitting, while a random secret is generated and split
// using Shamir's secret sharing method for the new keygen.
func (fr *Instance) processRound1() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.SetCurrentRound(common.Round1)
	fr.state.roundTimer.StartRoundTimeoutTimer(fr.state.GetCurrentRound())

	if !fr.needToRunCurrentRound() {
		return false, nil, fr.state.participant.SkipRound1()
	}

	var skI []byte // secret to be shared, nil if new keygen, lagrange interpolation of own part of secret if resharing
	if fr.instanceParams.isResharing() {
		if skI, err = fr.partialInterpolate(); err != nil {
			return false, nil, err
		}
	}

	bCastMessage, p2pMessages, err := fr.state.participant.Round1(skI)
	if err != nil {
		return false, nil, err
	}

	// get bytes representation of commitment points
	commitments := make([][]byte, 0)
	for _, commitment := range bCastMessage.Verifiers.Commitments {
		commitments = append(commitments, commitment.ToAffineCompressed())
	}

	// get shares encrypted by operators
	shares := make(map[uint32][]byte)
	for _, operatorID := range fr.instanceParams.operators {
		if uint32(fr.instanceParams.operatorID) == operatorID {
			continue
		}

		shamirShare := p2pMessages[operatorID]
		share := &bls.SecretKey{}
		if err := share.Deserialize(shamirShare.Value); err != nil {
			return false, nil, err
		}
		fr.state.operatorShares[operatorID] = share

		encryptedShare, err := fr.state.encryptByOperatorID(operatorID, shamirShare.Value)
		if err != nil {
			return false, nil, err
		}
		shares[operatorID] = encryptedShare
	}

	msg := &ProtocolMsg{
		Round: common.Round1,
		Round1Message: &Round1Message{
			Commitment: commitments,
			ProofS:     bCastMessage.Wi.Bytes(),
			ProofR:     bCastMessage.Ci.Bytes(),
			Shares:     shares,
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	return false, nil, fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}

func (fr *Instance) partialInterpolate() ([]byte, error) {
	if !fr.instanceParams.isResharing() {
		return nil, nil
	}

	indices := make([]bls.Fr, fr.instanceParams.oldKeyGenOutput.Threshold)
	values := make([]bls.Fr, fr.instanceParams.oldKeyGenOutput.Threshold)
	for i, id := range fr.instanceParams.operatorsOld {
		(&indices[i]).SetInt64(int64(id))
		if types.OperatorID(id) == fr.instanceParams.operatorID {
			err := (&values[i]).Deserialize(fr.instanceParams.oldKeyGenOutput.Share.Serialize())
			if err != nil {
				return nil, err
			}
		} else {
			(&values[i]).SetInt64(0)
		}
	}

	skI := new(bls.Fr)
	if err := bls.FrLagrangeInterpolation(skI, indices, values); err != nil {
		return nil, err
	}
	return skI.Serialize(), nil
}
