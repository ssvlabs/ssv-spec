package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

// processRound2 function can be executed once all messages for round 1 have been
// received. It uses the Fieldman verifier (VSS) to verify the shares from the
// commitments. If an invalid share is detected, a blame request is generated,
// otherwise the function proceeds to run round 2 in the frost library, which
// returns the Vk (public key) and Vkshare (public key of its own share). These
// values are serialized into a Protocol Message and broadcasted over the network.
func (fr *Instance) processRound2() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {

	if !fr.canProceedThisRound() {
		return false, nil, nil
	}
	fr.state.SetCurrentRound(common.Round2)
	fr.state.roundTimer.StartRoundTimeoutTimer(fr.state.GetCurrentRound())

	if !fr.needToRunCurrentRound() {
		return false, nil, nil
	}

	bcast := make(map[uint32]*frost.Round1Bcast)
	p2psend := make(map[uint32]*sharing.ShamirShare)

	for peerOpID, dkgMessage := range fr.state.msgContainer.AllMessagesForRound(common.Round1) {
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
		Round: common.Round2,
		Round2Message: &Round2Message{
			Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
			VkShare: bCastMessage.VkShare.ToAffineCompressed(),
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	return false, nil, fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}
