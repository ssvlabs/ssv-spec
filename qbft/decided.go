package qbft

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// UponDecided returns decided msg if decided, nil otherwise
func (c *Controller) UponDecided(signedMessage *types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	if err := ValidateDecided(
		c.config,
		signedMessage,
		c.Share,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	// Decode
	message := &Message{}
	if err := message.Decode(signedMessage.SSVMessage.Data); err != nil {
		return nil, errors.Wrap(err, "Could not decode decided Message")
	}

	// try to find instance
	inst := c.InstanceForHeight(message.Height)
	prevDecided := inst != nil && inst.State.Decided
	isFutureDecided := message.Height > c.Height

	if inst == nil {
		i := NewInstance(c.GetConfig(), c.Share, c.Identifier, message.Height)
		i.State.Round = message.Round
		i.State.Decided = true
		i.State.DecidedValue = message.FullData
		i.State.CommitContainer.AddMsg(signedMessage)
		c.StoredInstances.addNewInstance(i)
	} else if decided, _ := inst.IsDecided(); !decided {
		inst.State.Decided = true
		inst.State.Round = message.Round
		inst.State.DecidedValue = message.FullData
		inst.State.CommitContainer.AddMsg(signedMessage)
	} else { // decide previously, add if has more signers
		signers, _ := inst.State.CommitContainer.LongestUniqueSignersForRoundAndRoot(message.Round, message.Root)
		if len(signedMessage.GetOperatorIDs()) > len(signers) {
			inst.State.CommitContainer.AddMsg(signedMessage)
		}
	}

	if isFutureDecided {
		// bump height
		c.Height = message.Height
	}

	if !prevDecided {
		return signedMessage, nil
	}
	return nil, nil
}

func ValidateDecided(
	config IConfig,
	signedDecided *types.SignedSSVMessage,
	share *types.Share,
) error {

	// Decode
	message := &Message{}
	if err := message.Decode(signedDecided.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "Could not decode decided Message to validate")
	}

	if !IsDecidedMsg(share, signedDecided) {
		return errors.New("not a decided msg")
	}

	if err := signedDecided.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := baseCommitValidation(config, signedDecided, message.Height, share.Committee); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := signedDecided.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided")
	}

	r, err := HashDataRoot(message.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(r[:], message.Root[:]) {
		return errors.New("H(data) != root")
	}

	return nil
}

// IsDecidedMsg returns true if signed commit has all quorum sigs
func IsDecidedMsg(share *types.Share, signedDecided *types.SignedSSVMessage) bool {
	// Decode
	message := &Message{}
	if err := message.Decode(signedDecided.SSVMessage.Data); err != nil {
		return false
	}

	return share.HasQuorum(len(signedDecided.GetOperatorIDs())) && message.MsgType == CommitMsgType
}
