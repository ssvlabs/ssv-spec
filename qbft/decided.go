package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponFutureDecided returns error if could not process decided
func (c *Controller) UponFutureDecided(msg *SignedMessage) (*SignedMessage, error) {
	if msg.Message.Height <= c.Height {
		return nil, errors.New("invalid height")
	}

	if err := validateDecided(
		msg.Message.Height,
		c.GenerateConfig(),
		msg,
		c.Share.Committee,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	// stop any running instance
	inst := c.InstanceForHeight(c.Height)
	if inst != nil {
		inst.State.Decided = true
	}

	// get decided value
	data, err := msg.Message.GetCommitData()
	if err != nil {
		return nil, errors.Wrap(err, "could not get decided data")
	}

	// add an instance for the decided msg
	i := NewInstance(c.GenerateConfig(), c.Share, c.Identifier, msg.Message.Height)
	i.State.Decided = true
	i.State.DecidedValue = data.Data
	c.StoredInstances.addNewInstance(i)

	// bump height
	c.Height = msg.Message.Height

	if err := c.storage.SaveHighestDecided(msg); err != nil {
		return nil, errors.Wrap(err, "could not save decided")
	}

	return msg, nil
}

func validateDecided(
	height Height,
	config IConfig,
	signedDecided *SignedMessage,
	operators []*types.Operator,
) error {
	if err := baseCommitValidation(config, signedDecided, height, operators); err != nil {
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
