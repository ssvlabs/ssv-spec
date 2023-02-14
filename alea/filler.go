package alea

import (
	"bytes"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponFiller(signedFiller *SignedMessage, fillerMsgContainer *MsgContainer) error {

	// get data
	fillerData, err := signedFiller.Message.GetFillerData()
	if err != nil {
		return errors.Wrap(err, "uponFiller: could not get filler data from signedFiller")
	}

	// Add message to container
	fillerMsgContainer.AddMsg(signedFiller)

	// get values from structure
	entries := fillerData.Entries
	priorities := fillerData.Priorities
	// proofs := fillerData.Proofs
	operatorID := fillerData.OperatorID

	// get queue of the node to which the filler message intends to add entries
	queue := i.State.VCBCState.Queues[operatorID]

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

func isValidFiller(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != FillerMsgType {
		return errors.New("msg type is not FillerMsgType")
	}
	if signedMsg.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	FillerData, err := signedMsg.Message.GetFillerData()
	if err != nil {
		return errors.Wrap(err, "could not get FillerData data")
	}
	if err := FillerData.Validate(); err != nil {
		return errors.Wrap(err, "FillerData invalid")
	}

	// author
	operatorID := FillerData.OperatorID
	InCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == operatorID {
			InCommittee = true
		}
	}
	if !InCommittee {
		return errors.New("author (OperatorID) doesn't exist in Committee")
	}

	// priority
	priorities := FillerData.Priorities
	for idx, priority := range priorities {
		if state.VCBCState.HasM(operatorID, priority) {
			if !state.VCBCState.EqualM(operatorID, priority, FillerData.Entries[idx]) {
				return errors.New("existing (priority,author) with different proposals")
			}
		}
	}

	// AggregatedMsg
	aggregatedMsgs := FillerData.AggregatedMsgs
	for idx, aggregatedMsg := range aggregatedMsgs {

		signedAggregatedMessage := &SignedMessage{}
		signedAggregatedMessage.Decode(aggregatedMsg)

		if err := signedAggregatedMessage.Signature.VerifyByOperators(signedAggregatedMessage, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
			return errors.Wrap(err, "aggregatedMsg signature invalid")
		}
		if len(signedAggregatedMessage.GetSigners()) < int(state.Share.Quorum) {
			return errors.New("aggregatedMsg signers don't reach quorum")
		}

		vcbcReadyData, err := signedAggregatedMessage.Message.GetVCBCReadyData()
		if err != nil {
			return errors.Wrap(err, "could not get VCBCReadyData from given aggregated message")
		}
		givenHash, err := GetProposalsHash(FillerData.Entries[idx])
		if err != nil {
			return errors.Wrap(err, "could not get hash from given proposals")
		}
		if !bytes.Equal(givenHash, vcbcReadyData.Hash) {
			return errors.New("hash of proposals given doesn't match hash in the VCBCReadyData of the aggregated message")
		}
		if vcbcReadyData.Author != FillerData.OperatorID {
			return errors.New("author given doesn't match author in the VCBCReadyData of the aggregated message")
		}
		if vcbcReadyData.Priority != FillerData.Priorities[idx] {
			return errors.New("priority given doesn't match priority in the VCBCReadyData of the aggregated message")
		}
	}

	return nil
}

func CreateFiller(state *State, config IConfig, entries [][]*ProposalData, priorities []Priority, aggregatedMsgs [][]byte, operatorID types.OperatorID) (*SignedMessage, error) {
	fillerData := &FillerData{
		Entries:        entries,
		Priorities:     priorities,
		AggregatedMsgs: aggregatedMsgs,
		OperatorID:     operatorID,
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
