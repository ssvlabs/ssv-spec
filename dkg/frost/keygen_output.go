package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// processKeygenOutput function can be executed once all messages from round 2
// have been received. It reconstructs the Vk using the received VkShares and
// verifies that it is consistent. The keygen output is then created, including
// the Validator PK, operator public keys, and the operator's own share, and
// returned as a Protocol Outcome.
func (fr *Instance) processKeygenOutput() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.SetCurrentRound(common.KeygenOutput)
	fr.state.roundTimer.StartRoundTimeoutTimer(fr.state.GetCurrentRound())

	if !fr.needToRunCurrentRound() {
		return false, nil, nil
	}

	// verify shares by reconstructing consistent validatorPK over a sliding
	// window of t (threshold) operators
	reconstructed, err := fr.verifyShares()
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to verify shares")
	}
	reconstructedBytes := reconstructed.Serialize()

	// prepare keygen output
	out := &dkg.KeyGenOutput{
		Threshold: uint64(fr.instanceParams.threshold),
	}

	operatorPubKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.instanceParams.operators {
		msg, err := GetRound2Msg(fr.state.msgContainer, operatorID)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to retrieve round2 msg")
		}

		// set vk and secret share to output for this participant
		if operatorID == uint32(fr.instanceParams.operatorID) {
			out.ValidatorPK = msg.Vk
			sk := &bls.SecretKey{}
			if err := sk.Deserialize(fr.state.participant.SkShare.Bytes()); err != nil {
				return false, nil, err
			}
			out.Share = sk
		}

		// set operator public key for all operators
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(msg.VkShare); err != nil {
			return false, nil, err
		}
		operatorPubKeys[types.OperatorID(operatorID)] = pk
	}
	out.OperatorPubKeys = operatorPubKeys

	// assert validatorPK output from frost is equal to reconstructed Vk
	if !bytes.Equal(out.ValidatorPK, reconstructedBytes) {
		return false, nil, errors.New("can't reconstruct to the validator pk")
	}

	return true, &dkg.ProtocolOutcome{ProtocolOutput: out}, nil
}

// verifyShares uses a sliding window of size t (threshold) to iterate through
// operators and reconstruct Vk with each operator included at least once. Compare
// resulting Vk values to determine if all reconstructions are the same,
// indicating share validity.
func (fr *Instance) verifyShares() (*bls.G1, error) {

	var (
		quorumStart       = 0
		prevReconstructed = (*bls.G1)(nil)
	)

	for quorumEnd := int(fr.instanceParams.threshold); quorumEnd < len(fr.instanceParams.operators); quorumEnd++ {
		quorum := fr.instanceParams.operators[quorumStart:quorumEnd]
		currReconstructed, err := fr.reconstructValidatorPK(quorum)
		if err != nil {
			return nil, err
		}
		if prevReconstructed != nil && !currReconstructed.IsEqual(prevReconstructed) {
			return nil, errors.New("failed to create consistent public key from tshares")
		}
		prevReconstructed = currReconstructed
		quorumStart++
	}
	return prevReconstructed, nil
}

func (fr *Instance) reconstructValidatorPK(operators []uint32) (*bls.G1, error) {
	xVec, err := fr.getXVec(operators)
	if err != nil {
		return nil, err
	}

	yVec, err := fr.getYVec(operators)
	if err != nil {
		return nil, err
	}

	reconstructed := &bls.G1{}
	if err := bls.G1LagrangeInterpolation(reconstructed, xVec, yVec); err != nil {
		return nil, err
	}
	return reconstructed, nil
}

func (fr *Instance) getXVec(operators []uint32) ([]bls.Fr, error) {
	xVec := make([]bls.Fr, 0)
	for _, operator := range operators {
		x := bls.Fr{}
		x.SetInt64(int64(operator))
		xVec = append(xVec, x)
	}
	return xVec, nil
}

func (fr *Instance) getYVec(operators []uint32) ([]bls.G1, error) {
	yVec := make([]bls.G1, 0)
	for _, operatorID := range operators {
		msg, err := GetRound2Msg(fr.state.msgContainer, operatorID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve round2 msg")
		}
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(msg.VkShare); err != nil {
			return nil, errors.Wrap(err, "failed to deserialize public key")
		}
		y := bls.CastFromPublicKey(pk)
		yVec = append(yVec, *y)
	}
	return yVec, nil
}
