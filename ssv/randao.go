package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (dr *Runner) SignRandaoPreConsensus(epoch spec.Epoch, slot spec.Slot, signer types.KeyManager) (*PartialSignatureMessages, error) {
	sig, r, err := signer.SignRandaoReveal(epoch, dr.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign partial randao reveal")
	}

	// generate partial sig for randao
	msg := &PartialSignatureMessage{
		Slot:             slot,
		PartialSignature: sig,
		SigningRoot:      ensureRoot(r),
		Signer:           dr.Share.OperatorID,
	}

	return &PartialSignatureMessages{
		Type:     RandaoPartialSig,
		Messages: []*PartialSignatureMessage{msg},
	}, nil
}

// ProcessRandaoMessage process randao msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *Runner) ProcessRandaoMessage(signedMsg *SignedPartialSignatureMessage) (bool, [][]byte, error) {
	if err := dr.validateRandaoMsg(signedMsg); err != nil {
		return false, nil, errors.Wrap(err, "invalid randao message")
	}

	roots := make([][]byte, 0)
	anyQuorum := false
	for _, msg := range signedMsg.Message.Messages {
		prevQuorum := dr.State.RandaoPartialSig.HasQuorum(msg.SigningRoot)

		if err := dr.State.RandaoPartialSig.AddSignature(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not add partial randao signature")
		}

		if prevQuorum {
			continue
		}

		quorum := dr.State.RandaoPartialSig.HasQuorum(msg.SigningRoot)
		if quorum {
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

// validateRandaoMsg returns nil if randao message is valid
func (dr *Runner) validateRandaoMsg(msg *SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.CurrentDuty.Slot); err != nil {
		return err
	}

	if len(msg.Message.Messages) != 1 {
		return errors.New("expecting 1 radano partial sig")
	}

	panic("verify beacon signing root is what we expect")

	return nil
}
