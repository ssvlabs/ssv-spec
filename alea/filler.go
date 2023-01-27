package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponFiller(signedFiller *SignedMessage, fillerMsgContainer *MsgContainer) error {

	// get data
	fillerData, err := signedFiller.Message.GetFillerData()
	if err != nil {
		errors.Wrap(err, "uponFiller: could not get filler data from signedFiller")
	}

	// Add message to container
	fillerMsgContainer.AddMsg(signedFiller)

	// get values from structure
	entries := fillerData.Entries
	priorities := fillerData.Priorities
	// proofs := fillerData.Proofs
	operatorID := fillerData.OperatorID

	// get queue of the node to which the filler message intends to add entries
	queue := i.State.VCBCState.queues[operatorID]

	// get local highest priority value
	_, localLastPriority := queue.PeekLast()

	// if message has entries with higher priority, store value
	for idx, priority := range priorities {
		if priority > localLastPriority {
			queue.Enqueue(entries[idx], priority)
		}
	}

	// signal that filler message was received (used for node to stop waiting in the recovery mechanism part)
	i.State.FillerMsgReceived += 1

	return nil
}

func CreateFiller(state *State, config IConfig, entries [][]*ProposalData, priorities []Priority, proofs []types.Signature, operatorID types.OperatorID) (*SignedMessage, error) {
	fillerData := &FillerData{
		Entries:    entries,
		Priorities: priorities,
		Proofs:     proofs,
		OperatorID: operatorID,
	}
	dataByts, err := fillerData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateFiller: could not encode filler data")
	}
	msg := &Message{
		MsgType:    FillerMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateFiller: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
