package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

func (fr *FROST) processRound2() (*dkg.ProtocolOutcome, error) {

	if !fr.needToRunCurrentRound() {
		return nil, nil
	}

	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	for peerOID, dkgMessage := range fr.state.msgs[Round1] {

		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
			return nil, errors.Wrap(err, "failed to decode protocol msg")
		}
		verifiers := new(sharing.FeldmanVerifier)
		for _, commitmentBytes := range protocolMessage.Round1Message.Commitment {
			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode commitment point")
			}
			verifiers.Commitments = append(verifiers.Commitments, commitment)
		}

		Wi, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofS)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode scalar")
		}
		Ci, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofR)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode scalar")
		}

		bcastMessage := &frost.Round1Bcast{
			Verifiers: verifiers,
			Wi:        Wi,
			Ci:        Ci,
		}
		bcast[peerOID] = bcastMessage

		if uint32(fr.state.operatorID) == peerOID {
			continue
		}

		encryptedShare := protocolMessage.Round1Message.Shares[uint32(fr.state.operatorID)]
		shareBytes, err := ecies.Decrypt(fr.state.sessionSK, encryptedShare)
		if err != nil {
			fr.state.currentRound = Blame
			return fr.createAndBroadcastBlameOfInvalidShare(peerOID)
		}

		share := &sharing.ShamirShare{
			Id:    uint32(fr.state.operatorID),
			Value: shareBytes,
		}

		p2psend[peerOID] = share

		err = verifiers.Verify(share)
		if err != nil {
			fr.state.currentRound = Blame
			return fr.createAndBroadcastBlameOfInvalidShare(peerOID)
		}
	}

	bCastMessage, err := fr.state.participant.Round2(bcast, p2psend)
	if err != nil {
		return nil, err
	}

	msg := &ProtocolMsg{
		Round: fr.state.currentRound,
		Round2Message: &Round2Message{
			Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
			VkShare: bCastMessage.VkShare.ToAffineCompressed(),
		},
	}
	_, err = fr.broadcastDKGMessage(msg)
	return nil, err
}
