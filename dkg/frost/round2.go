package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

// processRound2 verifies incoming shares from all operators and broadcasts round2
// message with validatorPK and its public key
func (fr *Instance) processRound2() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.currentRound = Round2

	if !fr.needToRunCurrentRound() {
		return false, nil, nil
	}

	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	for peerOpID, dkgMessage := range fr.state.msgContainer.AllMessagesForRound(Round1) {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		// prepare broadcast message
		verifiers := new(sharing.FeldmanVerifier)
		for _, commitmentBytes := range protocolMessage.Round1Message.Commitment {
			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
			if err != nil {
				return false, nil, errors.Wrap(err, "failed to decode commitment point")
			}
			verifiers.Commitments = append(verifiers.Commitments, commitment)
		}
		Wi, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofS)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to decode scalar")
		}
		Ci, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofR)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to decode scalar")
		}
		bcastMessage := &frost.Round1Bcast{
			Verifiers: verifiers,
			Wi:        Wi,
			Ci:        Ci,
		}
		bcast[peerOpID] = bcastMessage

		// prepare p2p message
		if uint32(fr.instanceParams.operatorID) == peerOpID {
			continue
		}

		encryptedShare := protocolMessage.Round1Message.Shares[uint32(fr.instanceParams.operatorID)]
		shareBytes, err := ecies.Decrypt(fr.state.sessionSK, encryptedShare)
		if err != nil {
			return fr.createAndBroadcastBlameOfInvalidShare(peerOpID)

		}
		share := &sharing.ShamirShare{
			Id:    uint32(fr.instanceParams.operatorID),
			Value: shareBytes,
		}
		if err = verifiers.Verify(share); err != nil {
			return fr.createAndBroadcastBlameOfInvalidShare(peerOpID)
		}
		p2psend[peerOpID] = share
	}

	bCastMessage, err := fr.state.participant.Round2(bcast, p2psend)
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: Round2,
		Round2Message: &Round2Message{
			Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
			VkShare: bCastMessage.VkShare.ToAffineCompressed(),
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	return false, nil, fr.config.network.BroadcastDKGMessage(bcastMsg)
}
