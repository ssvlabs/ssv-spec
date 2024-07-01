package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// UponDecided returns decided msg if decided, nil otherwise
func (c *Controller) UponDecided(msg *ProcessingMessage) (*types.SignedSSVMessage, error) {
	if err := ValidateDecided(
		c.config,
		msg,
		c.CommitteeMember,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	// try to find instance
	inst := c.InstanceForHeight(msg.QBFTMessage.Height)
	prevDecided := inst != nil && inst.State.Decided
	isFutureDecided := msg.QBFTMessage.Height > c.Height

	if inst == nil {
		i := NewInstance(c.GetConfig(), c.CommitteeMember, c.Identifier, msg.QBFTMessage.Height)
		i.State.Round = msg.QBFTMessage.Round
		i.State.Decided = true
		i.State.DecidedValue = msg.SignedMessage.FullData
		err := i.State.CommitContainer.AddMsg(msg)
		if err != nil {
			return nil, err
		}
		c.StoredInstances.addNewInstance(i)
	} else if decided, _ := inst.IsDecided(); !decided {
		err := inst.State.CommitContainer.AddMsg(msg)
		if err != nil {
			return nil, err
		}
		inst.State.Decided = true
		inst.State.Round = msg.QBFTMessage.Round
		inst.State.DecidedValue = msg.SignedMessage.FullData
	} else { // decide previously, add if has more signers
		signers, _ := inst.State.CommitContainer.LongestUniqueSignersForRoundAndRoot(msg.QBFTMessage.Round, msg.QBFTMessage.Root)
		if len(msg.SignedMessage.OperatorIDs) > len(signers) {
			err := inst.State.CommitContainer.AddMsg(msg)
			if err != nil {
				return nil, err
			}
		}
	}

	if isFutureDecided {
		// bump height
		c.Height = msg.QBFTMessage.Height
	}

	if !prevDecided {
		return msg.SignedMessage, nil
	}
	return nil, nil
}

func ValidateDecided(
	config IConfig,
	msg *ProcessingMessage,
	share *types.CommitteeMember,
) error {

	isDecided, err := IsDecidedMsg(share, msg)
	if err != nil {
		return err
	}
	if !isDecided {
		return errors.New("not a decided msg")
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := baseCommitValidationVerifySignature(config, msg, msg.QBFTMessage.Height, share.Committee); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided")
	}

	r, err := HashDataRoot(msg.SignedMessage.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(r[:], msg.QBFTMessage.Root[:]) {
		return errors.New("H(data) != root")
	}

	return nil
}

// IsDecidedMsg returns true if signed commit has all quorum sigs
func IsDecidedMsg(share *types.CommitteeMember, msg *ProcessingMessage) (bool, error) {
	return share.HasQuorum(len(msg.SignedMessage.OperatorIDs)) && msg.QBFTMessage.MsgType == CommitMsgType, nil
}
