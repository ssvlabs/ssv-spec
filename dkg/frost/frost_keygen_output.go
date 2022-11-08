package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (fr *FROST) processKeygenOutput() (*dkg.KeyGenOutput, error) {

	if !fr.needToRunCurrentRound() {
		return nil, nil
	}

	reconstructed, err := fr.verifyShares()
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify shares")
	}

	reconstructedBytes := reconstructed.Serialize()

	out := &dkg.KeyGenOutput{
		Threshold: uint64(fr.state.threshold),
	}

	operatorPubKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range fr.state.operators {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(fr.state.msgs[Round2][operatorID].Message.Data); err != nil {
			return nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		if operatorID == uint32(fr.state.operatorID) {
			sk := &bls.SecretKey{}
			if err := sk.Deserialize(fr.state.participant.SkShare.Bytes()); err != nil {
				return nil, err
			}

			out.Share = sk
			out.ValidatorPK = protocolMessage.Round2Message.Vk
		}

		pk := &bls.PublicKey{}
		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
			return nil, err
		}

		operatorPubKeys[types.OperatorID(operatorID)] = pk
	}

	out.OperatorPubKeys = operatorPubKeys

	if !bytes.Equal(out.ValidatorPK, reconstructedBytes) {
		return nil, errors.New("can't reconstruct to the validator pk")
	}

	return out, nil
}

func (fr *FROST) verifyShares() (*bls.G1, error) {

	var (
		quorumStart       = 0
		quorumEnd         = int(fr.state.threshold)
		prevReconstructed = (*bls.G1)(nil)
	)

	// Sliding window of quorum 0...threshold until n-threshold...n
	for quorumEnd < len(fr.state.operators) {
		quorum := fr.state.operators[quorumStart:quorumEnd]
		currReconstructed, err := fr.verifyShare(quorum)
		if err != nil {
			return nil, err
		}
		if prevReconstructed != nil && !currReconstructed.IsEqual(prevReconstructed) {
			return nil, errors.New("failed to create consistent public key from tshares")
		}
		prevReconstructed = currReconstructed
		quorumStart++
		quorumEnd++
	}
	return prevReconstructed, nil
}

func (fr *FROST) verifyShare(operators []uint32) (*bls.G1, error) {
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
