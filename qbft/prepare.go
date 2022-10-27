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
	if i.State.ProposalAcceptedForCurrentRound == nil {
		return errors.New("no proposal accepted for prepare")
	}

	if err := validSignedPrepareForHeightRoundAndValue(
		i.config,
		signedPrepare,
		i.State.Height,
		i.State.Round,
		i.State.ProposalAcceptedForCurrentRound.Message.Input.Root[:],
		i.State.Share.Committee,
	); err != nil {
		return errors.Wrap(err, "invalid prepare msg")
	}

	addedMsg, err := prepareMsgContainer.AddFirstMsgForSignerAndRound(signedPrepare)
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

	commitMsg, err := CreateCommit(i.State, i.config, i.State.ProposalAcceptedForCurrentRound.Message.Input.Root)
	if err != nil {
		return errors.Wrap(err, "could not create commit msg")
	}

	commitEncoded, err := commitMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode commit message")
	}

	if err = i.Broadcast(commitEncoded, types.ConsensusCommitMsgType); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	return nil
}

func getRoundChangeJustification(state *State, config IConfig, prepareMsgContainer *MsgContainer) []*SignedMessage {
	if state.LastPreparedValue == nil {
		return nil
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	ret := make([]*SignedMessage, 0)
	for _, msg := range prepareMsgs {
		if err := validSignedPrepareForHeightRoundAndValue(
			config,
			msg,
			state.Height,
			state.LastPreparedRound,
			state.LastPreparedValue.Root[:],
			state.Share.Committee,
		); err == nil {
			ret = append(ret, msg)
		}
	}
	return ret
}

// validPreparesForHeightRoundAndValue returns an aggregated prepare msg for a specific Height and round
//func validPreparesForHeightRoundAndValue(
//	config IConfig,
//	prepareMessages []*SignedMessage,
//	height Height,
//	round Round,
//	value []byte,
//	operators []*types.Operator) *SignedMessage {
//	var aggregatedPrepareMsg *SignedMessage
//	for _, signedMsg := range prepareMessages {
//		if err := validSignedPrepareForHeightRoundAndValue(config, signedMsg, height, round, value, operators); err == nil {
//			if aggregatedPrepareMsg == nil {
//				aggregatedPrepareMsg = signedMsg
//			} else {
//				// TODO: check error
//				// nolint
//				aggregatedPrepareMsg.Aggregate(signedMsg)
//			}
//		}
//	}
//	return aggregatedPrepareMsg
//}

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
	if signedPrepare.Message.Height != height {
		return errors.New("msg Height wrong")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("msg round wrong")
	}

	if !bytes.Equal(signedPrepare.Message.Input.Root[:], value) {
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
func CreatePrepare(state *State, config IConfig, newRound Round, value [32]byte) (*SignedMessage, error) {
	msg := &Message{
		Height: state.Height,
		Round:  newRound,
		Input: &Data{
			Root:   value,
			Source: nil,
		},
	}

	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	return &SignedMessage{
		Message:   msg,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Signature: sig,
	}, nil
}
