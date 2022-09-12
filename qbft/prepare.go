package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponPrepare(
	signedPrepare *SignedMessage,
	prepareMsgContainer,
	commitMsgContainer *MsgContainer,
) error {
	// TODO - if we receive a prepare before a proposal and return an error we will never process the prepare msg, we still need to add it to the container
	if i.State.ProposalAcceptedForCurrentRound == nil {
		return errors.New("no proposal accepted for prepare")
	}

	//acceptedProposalData, err := i.State.ProposalAcceptedForCurrentRound.Message.Input
	//if err != nil {
	//	return errors.Wrap(err, "could not get accepted proposal data")
	//}
	if err := validSignedPrepareForHeightRoundAndValue(
		i.config,
		signedPrepare,
		i.State.Height,
		i.State.Round,
		i.State.ProposalAcceptedForCurrentRound.Message.Input,
		//acceptedProposalData.Data,
		i.State.Share.Committee,
	); err != nil {
		return errors.Wrap(err, "invalid prepare msg")
	}

	addedMsg, err := prepareMsgContainer.AddIfDoesntExist(signedPrepare)
	if err != nil {
		return errors.Wrap(err, "could not add prepare msg to container")
	}
	if !addedMsg {
		return nil // uponPrepare was already called
	}

	if !HasQuorum(i.State.Share, prepareMsgContainer.MessagesForRound(i.State.Round)) {
		return nil // no quorum yet
	}

	if didSendCommitForHeightAndRound(i.State, commitMsgContainer) {
		return nil // already moved to commit stage
	}

	proposedValue := i.State.ProposalAcceptedForCurrentRound.Message.Input

	i.State.LastPreparedValue = proposedValue
	i.State.LastPreparedRound = i.State.Round

	commitMsg, err := CreateCommit(i.State, i.config, proposedValue)
	if err != nil {
		return errors.Wrap(err, "could not create commit msg")
	}

	commitEncoded, err := commitMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode prepare message")
	}

	msgID := types.PopulateMsgType(i.State.ID, types.ConsensusCommitMsgType)

	broadcastMsg := &types.Message{
		ID:   msgID,
		Data: commitEncoded,
	}

	if err = i.Broadcast(broadcastMsg); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	return nil
}

func getRoundChangeJustification(state *State, config IConfig, prepareMsgContainer *MsgContainer) []*SignedMessageHeader {
	if state.LastPreparedValue == nil {
		return nil
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	ret := make([]*SignedMessageHeader, 0)
	for _, msg := range prepareMsgs {
		if err := validSignedPrepareForHeightRoundAndValue(config, msg, state.Height, state.LastPreparedRound, state.LastPreparedValue, state.Share.Committee); err == nil {
			msgHeader, _ := msg.ToSignedMessageHeader()
			ret = append(ret, msgHeader)
		}
	}
	return ret
}

// validPreparesForHeightRoundAndValue returns an aggregated prepare msg for a specific Height and round
func validPreparesForHeightRoundAndValue(
	config IConfig,
	prepareMessages []*SignedMessage,
	height Height,
	round Round,
	value []byte,
	operators []*types.Operator) *SignedMessage {
	var aggregatedPrepareMsg *SignedMessage
	for _, signedMsg := range prepareMessages {
		if err := validSignedPrepareForHeightRoundAndValue(config, signedMsg, height, round, value, operators); err == nil {
			if aggregatedPrepareMsg == nil {
				aggregatedPrepareMsg = signedMsg
			} else {
				aggregatedPrepareMsg.Aggregate(signedMsg)
			}
		}
	}
	return aggregatedPrepareMsg
}

// validSignedPrepareForHeightRoundAndValue known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndValue(
	config IConfig,
	signedPrepare *SignedMessage,
	height Height,
	round Round,
	value []byte,
	operators []*types.Operator,
) error {
	//if signedPrepare.Message.MsgType != PrepareMsgType {
	//	return errors.New("prepare msg type is wrong")
	//}
	if signedPrepare.Message.Height != height {
		return errors.New("msg Height wrong")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("msg round wrong")
	}

	//prepareData, err := signedPrepare.Message.GetPrepareData()
	//if err != nil {
	//	return errors.Wrap(err, "could not get prepare data")
	//}
	//if err := prepareData.Validate(); err != nil {
	//	return errors.Wrap(err, "prepareData invalid")
	//}

	if bytes.Compare(signedPrepare.Message.Input, value) != 0 {
		return errors.New("prepare data != proposed data")
	}

	if len(signedPrepare.GetSigners()) != 1 {
		return errors.New("prepare msg allows 1 signer")
	}

	if err := signedPrepare.Signature.VerifyByOperators(signedPrepare, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "prepare msg signature invalid")
	}
	return nil
}

// validSignedPrepareHeaderForHeightRoundAndValue known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareHeaderForHeightRoundAndValue(
	config IConfig,
	signedPrepare *SignedMessageHeader,
	height Height,
	round Round,
	value []byte,
	operators []*types.Operator,
) error {
	if signedPrepare.Message.Height != height {
		return errors.New("msg Height wrong")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("msg round wrong")
	}

	//prepareData.Validate()
	// TODO<olegshmuelov> inputroot validate
	//if signedPrepare.Message.InputRoot == [32]byte{} {
	//	return errors.New("prepareData invalid")
	//}

	// TODO<olegshmuelov>: encode value to input root and compare?
	//if bytes.Compare(signedPrepare.Message.InputRoot, value) != 0 {
	//	return errors.New("prepare data != proposed data")
	//}

	if len(signedPrepare.GetSigners()) != 1 {
		return errors.New("prepare msg allows 1 signer")
	}

	// TODO<olegshmuelov>: do we need to do it for justifications?
	//if err := signedPrepare.Signature.VerifyByOperators(signedPrepare, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
	//	return errors.Wrap(err, "prepare msg signature invalid")
	//}
	return nil
}

// CreatePrepare
/**
Prepare(
                    signPrepare(
                        UnsignedPrepare(
                            |current.blockchain|,
                            newRound,
                            digest(m.proposedBlock)),
                        current.id
                        )
                );
*/
func CreatePrepare(state *State, config IConfig, newRound Round, value []byte) (*SignedMessage, error) {
	msg := &Message{
		Height: state.Height,
		Round:  newRound,
		Input:  value,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	prepareMsg := &SignedMessage{
		Message:   msg,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Signature: sig,
	}

	//prepareData := &PrepareData{
	//	Data: value,
	//}
	//dataByts, err := prepareData.Encode()
	//
	//msg := &Message{
	//	MsgType:    PrepareMsgType,
	//	Height:     state.Height,
	//	Round:      newRound,
	//	Identifier: state.ID,
	//	Data:       dataByts,
	//}
	//sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed signing prepare msg")
	//}
	//
	//signedMsg := &SignedMessage{
	//	Signature: sig,
	//	Signers:   []types.OperatorID{state.Share.OperatorID},
	//	Message:   msg,
	//}
	return prepareMsg, nil
}
