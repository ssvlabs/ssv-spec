package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponDecided returns error if could not process decided
func (c *Controller) UponDecided(msg *SignedMessage) (*SignedMessage, error) {
	// decided msgs for past (already decided) instances will not decide again, just return
	if msg.Message.Height < c.Height {
		return nil, nil
	}

	if err := validateDecided(
		msg.Message.Height,
		c.GenerateConfig(),
		msg,
		c.Share,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	// get decided value
	data, err := msg.Message.GetCommitData()
	if err != nil {
		return nil, errors.Wrap(err, "could not get decided data")
	}

	// if decided is for running instance (or higher), find and stop it
	if inst := c.InstanceForHeight(c.Height); inst != nil && !inst.State.Decided {
		inst.State.DecidedValue = data.Data
		inst.State.Decided = true
	}

	isFutureDecided := msg.Message.Height > c.Height
	if isFutureDecided {
		// add an instance for the decided msg
		i := NewInstance(c.GenerateConfig(), c.Share, c.Identifier, msg.Message.Height)
		i.State.Decided = true
		i.State.DecidedValue = data.Data
		c.StoredInstances.addNewInstance(i)

		// bump height
		c.Height = msg.Message.Height
	}

	if err := c.storage.SaveHighestDecided(msg); err != nil {
		return nil, errors.Wrap(err, "could not save decided")
	}
	return msg, nil
}

func validateDecided(
	height Height,
	config IConfig,
	signedDecided *SignedMessage,
	share *types.Share,
) error {
	if !isDecidedMsg(share, signedDecided) {
		return errors.New("not a decided msg")
	}

	if err := signedDecided.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := baseCommitValidation(config, signedDecided, height, share.Committee); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	msgDecidedData, err := signedDecided.Message.GetCommitData()
	if err != nil {
		return errors.Wrap(err, "could not get msg decided data")
	}
	if err := msgDecidedData.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided data")
	}

	valCheck := config.GetValueCheckF()
	if err := valCheck(msgDecidedData.Data); err != nil {
		return errors.Wrap(err, "decided value invalid")
	}

	return nil
}

// returns true if signed commit has all quorum sigs
func isDecidedMsg(share *types.Share, signedDecided *SignedMessage) bool {
	return share.HasQuorum(len(signedDecided.Signers))
}
