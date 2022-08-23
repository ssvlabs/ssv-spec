package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponDecided returns true if a decided messages was received.
func (i *Instance) UponDecided(signedDecided *SignedMessage, commitMsgContainer *MsgContainer) (bool, []byte, error) {
	if i.State.Decided {
		return true, i.State.DecidedValue, nil
	}

	if err := validateDecided(
		i.State,
		i.config,
		signedDecided,
		i.State.Height,
		i.State.Share.Committee,
		i.config.GetValueCheckF(),
	); err != nil {
		return false, nil, errors.Wrap(err, "invalid decided msg")
	}

	addMsg, err := commitMsgContainer.AddIfDoesntExist(signedDecided)
	if err != nil {
		return false, nil, errors.Wrap(err, "could not add commit msg to container")
	}
	if !addMsg {
		return false, nil, nil // UponCommit was already called
	}

	msgDecidedData, err := signedDecided.Message.GetCommitData()
	if err != nil {
		return false, nil, errors.Wrap(err, "could not get msg decided data")
	}

	return true, msgDecidedData.Data, nil
}

func validateDecided(
	state *State,
	config IConfig,
	signedDecided *SignedMessage,
	height Height,
	operators []*types.Operator,
	valCheck ProposedValueCheckF,
) error {
	if !isDecidedMsg(state, signedDecided) {
		return errors.New("not a decided msg")
	}

	if err := baseCommitValidation(config, signedDecided, height, operators); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	msgDecidedData, err := signedDecided.Message.GetCommitData()
	if err != nil {
		return errors.Wrap(err, "could not get msg decided data")
	}

	if err := valCheck(msgDecidedData.Data); err != nil {
		return errors.Wrap(err, "decided value invalid")
	}

	return baseCommitValidation(config, signedDecided, height, operators)
}

// returns true if signed commit has all quorum sigs
func isDecidedMsg(state *State, signedDecided *SignedMessage) bool {
	return state.Share.HasQuorum(len(signedDecided.Signers))
}
