package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// processRound1 splits secret into shares between ooperators and broadcasts round1
// message with encrypted shares, commitments and Schnorr proof values
func (fr *FROST) processRound1() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.currentRound = Round1

	if !fr.needToRunCurrentRound() {
		return false, nil, fr.state.participant.SkipRound1()
	}

	var skI []byte // secret to be shared, nil if new keygen, lagrange interpolation of own part of secret if resharing
	if fr.config.isResharing() {
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
	for _, operatorID := range fr.config.operators {
		if uint32(fr.config.operatorID) == operatorID {
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
		Round: Round1,
		Round1Message: &Round1Message{
			Commitment: commitments,
			ProofS:     bCastMessage.Wi.Bytes(),
			ProofR:     bCastMessage.Ci.Bytes(),
			Shares:     shares,
		},
	}

	_, err = fr.broadcastDKGMessage(msg)
	return false, nil, err
}

func (fr *FROST) partialInterpolate() ([]byte, error) {
	if !fr.config.isResharing() {
		return nil, nil
	}

	indices := make([]bls.Fr, fr.config.oldKeyGenOutput.Threshold)
	values := make([]bls.Fr, fr.config.oldKeyGenOutput.Threshold)
	for i, id := range fr.config.operatorsOld {
		(&indices[i]).SetInt64(int64(id))
		if types.OperatorID(id) == fr.config.operatorID {
			err := (&values[i]).Deserialize(fr.config.oldKeyGenOutput.Share.Serialize())
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
