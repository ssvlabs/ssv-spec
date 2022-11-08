package qbft

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponDecided returns decided msg if decided, nil otherwise
func (c *Controller) UponDecided(msg *SignedMessage) (*SignedMessage, error) {
	// decided msgs for past (already decided) instances will not decide again, just return
	if msg.Message.Height < c.Height {
		return nil, nil
	}

	if err := validateDecided(
		c.config,
		msg,
		c.Share,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	// did previously decide?
	inst := c.InstanceForHeight(msg.Message.Height)
	prevDecided := inst != nil && inst.State.Decided

	// Mark current instance decided
	if inst := c.InstanceForHeight(c.Height); inst != nil && !inst.State.Decided {
		inst.State.Decided = true
		if msg.Message.Round > inst.State.Round {
			inst.State.Round = msg.Message.Round
		}
		if c.Height == msg.Message.Height {
			inst.State.DecidedValue = msg.InputSource
		}
	}

	isFutureDecided := msg.Message.Height > c.Height
	if isFutureDecided {
		// add an instance for the decided msg
		i := NewInstance(c.GetConfig(), c.Share, c.Identifier, msg.Message.Height)
		i.State.Round = msg.Message.Round
		i.State.Decided = true
		i.State.DecidedValue = msg.InputSource
		c.StoredInstances.addNewInstance(i)

		// bump height
		c.Height = msg.Message.Height
	}

	if !prevDecided {
		if err := c.GetConfig().GetStorage().SaveHighestDecided(c.Identifier, msg); err != nil {
			// no need to fail processing the decided msg if failed to save
			fmt.Printf("%s\n", err.Error())
		}
		return msg, nil
	}
	return nil, nil
}

func validateDecided(
	config IConfig,
	signedDecided *SignedMessage,
	share *types.Share,
) error {
	if err := signedDecided.Validate(types.DecidedMsgType); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	input := &Data{
		Root:   signedDecided.Message.InputRoot,
		Source: signedDecided.InputSource,
	}
	err := input.Validate(nil)
	if err != nil {
		return errors.Wrap(err, "invalid input data")
	}

	if err := baseCommitValidation(config, signedDecided, signedDecided.Message.Height, share.Committee); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	return nil
}

// returns true if signed commit has all quorum sigs
func isDecidedMsg(share *types.Share, signedDecided *SignedMessage, msgType types.MsgType) bool {
	return msgType == types.DecidedMsgType && share.HasQuorum(len(signedDecided.Signers))
}
