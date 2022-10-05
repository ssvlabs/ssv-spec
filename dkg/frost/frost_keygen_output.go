package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (fr *FROST) processKeygenOutput() (*dkg.KeyGenOutput, error) {
	if _, err := fr.verifyShares(); err != nil {
		return nil, errors.Wrap(err, "failed to verify shares")
	}

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
	return out, nil
}

func (fr *FROST) verifyShares() ([]*bls.G1, error) {

	outputs := make([]*bls.G1, 0)

	for j := int(fr.state.threshold + 1); j < len(fr.state.operators); j++ {

		xVec := make([]bls.Fr, 0)
		yVec := make([]bls.G1, 0)

		for i := j - int(fr.state.threshold+1); i < j; i++ {
			operatorID := fr.state.operators[i]

			protocolMessage := &ProtocolMsg{}
			if err := protocolMessage.Decode(fr.state.msgs[Round2][operatorID].Message.Data); err != nil {
				return nil, errors.Wrap(err, "failed to decode protocol msg")
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
		if !outputs[i].IsEqual(outputs[i-1]) {
			return nil, errors.New("failed to create consistent public key from t+1 shares")
		}
	}

	return outputs, nil
}
