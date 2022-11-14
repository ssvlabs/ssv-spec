package frost

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func (fr *FROST) processRound1() error {

	if !fr.needToRunCurrentRound() {
		return fr.state.participant.SkipRound1()
	}

	var (
		skI []byte // secret to be shared, nil if new keygen, lagrange interpolation of own part of secret if resharing
		err error
	)

	if fr.isResharing() {
		skI, err = fr.partialInterpolate()
		if err != nil {
			return err
		}
	}

	bCastMessage, p2pMessages, err := fr.state.participant.Round1(skI)
	if err != nil {
		return err
	}

	// get bytes representation of commitment points
	commitments := make([][]byte, 0)
	for _, commitment := range bCastMessage.Verifiers.Commitments {
		commitments = append(commitments, commitment.ToAffineCompressed())
	}

	// encrypted shares by operators
	shares := make(map[uint32][]byte)
	for _, operatorID := range fr.state.operators {
		if uint32(fr.state.operatorID) == operatorID {
			continue
		}

		share := &bls.SecretKey{}
		shamirShare := p2pMessages[operatorID]

		if err := share.Deserialize(shamirShare.Value); err != nil {
			return err
		}

		fr.state.operatorShares[operatorID] = share

		encryptedShare, err := fr.encryptByOperatorID(operatorID, shamirShare.Value)
		if err != nil {
			return err
		}
		shares[operatorID] = encryptedShare
	}

	msg := &ProtocolMsg{
		Round: fr.state.currentRound,
		Round1Message: &Round1Message{
			Commitment: commitments,
			ProofS:     bCastMessage.Wi.Bytes(),
			ProofR:     bCastMessage.Ci.Bytes(),
			Shares:     shares,
		},
	}
	_, err = fr.broadcastDKGMessage(msg)
	return err
}

func (fr *FROST) partialInterpolate() ([]byte, error) {
	if !fr.isResharing() {
		return nil, nil
	}

	skI := new(bls.Fr)

	indices := make([]bls.Fr, fr.state.oldKeyGenOutput.Threshold)
	values := make([]bls.Fr, fr.state.oldKeyGenOutput.Threshold)
	for i, id := range fr.state.operatorsOld {
		(&indices[i]).SetInt64(int64(id))
		if types.OperatorID(id) == fr.state.operatorID {
			err := (&values[i]).Deserialize(fr.state.oldKeyGenOutput.Share.Serialize())
			if err != nil {
				return nil, err
			}
		} else {
			(&values[i]).SetInt64(0)
		}
	}

	if err := bls.FrLagrangeInterpolation(skI, indices, values); err != nil {
		return nil, err
	}
	return skI.Serialize(), nil
}
