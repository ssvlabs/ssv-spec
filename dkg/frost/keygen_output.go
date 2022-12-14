package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// processKeygenOutput verifies validatorPK from round2 and returns protocol outcome
func (fr *FROST) processKeygenOutput() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.currentRound = KeygenOutput

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
		Threshold: uint64(fr.config.threshold),
	}

	operatorPubKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.config.operators {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(fr.state.msgs[Round2][operatorID].Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		// set vk and secret share to output for this participant
		if operatorID == uint32(fr.config.operatorID) {
			out.ValidatorPK = protocolMessage.Round2Message.Vk
			sk := &bls.SecretKey{}
			if err := sk.Deserialize(fr.state.participant.SkShare.Bytes()); err != nil {
				return false, nil, err
			}
			out.Share = sk
		}

		// set operator public key for all operators
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
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

// verifyShares reconstructs Vk over a sliding window of t (threshold) operators
// and checks if reconstructed Vk is consistent
func (fr *FROST) verifyShares() (*bls.G1, error) {

	var (
		quorumStart       = 0
		prevReconstructed = (*bls.G1)(nil)
	)

	// Sliding window of quorum 0...threshold until n-threshold...n
	for quorumEnd := int(fr.config.threshold); quorumEnd < len(fr.config.operators); quorumEnd++ {
		quorum := fr.config.operators[quorumStart:quorumEnd]
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

func (fr *FROST) reconstructValidatorPK(operators []uint32) (*bls.G1, error) {
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

func (fr *FROST) getXVec(operators []uint32) ([]bls.Fr, error) {
	xVec := make([]bls.Fr, 0)
	for _, operator := range operators {
		x := bls.Fr{}
		x.SetInt64(int64(operator))
		xVec = append(xVec, x)
	}
	return xVec, nil
}

func (fr *FROST) getYVec(operators []uint32) ([]bls.G1, error) {
	yVec := make([]bls.G1, 0)
	for _, operator := range operators {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(fr.state.msgs[Round2][operator].Message.Data); err != nil {
			return nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		pk := &bls.PublicKey{}
		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
			return nil, errors.Wrap(err, "failed to deserialize public key")
		}
		y := bls.CastFromPublicKey(pk)
		yVec = append(yVec, *y)
	}
	return yVec, nil
}
