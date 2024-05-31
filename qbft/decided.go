package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// UponDecided returns decided msg if decided, nil otherwise
func (c *Controller) UponDecided(signedMsg *types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	if err := ValidateDecided(
		c.config,
		signedMsg,
		c.SharedValidator,
	); err != nil {
		return nil, errors.Wrap(err, "invalid decided msg")
	}

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return nil, err
	}

	// try to find instance
	inst := c.InstanceForHeight(msg.Height)
	prevDecided := inst != nil && inst.State.Decided
	isFutureDecided := msg.Height > c.Height

	if inst == nil {
		i := NewInstance(c.GetConfig(), c.SharedValidator, c.Identifier, msg.Height)
		i.State.Round = msg.Round
		i.State.Decided = true
		i.State.DecidedValue = signedMsg.FullData
		err := i.State.CommitContainer.AddMsg(signedMsg)
		if err != nil {
			return nil, err
		}
		c.StoredInstances.addNewInstance(i)
	} else if decided, _ := inst.IsDecided(); !decided {
		err := inst.State.CommitContainer.AddMsg(signedMsg)
		if err != nil {
			return nil, err
		}
		inst.State.Decided = true
		inst.State.Round = msg.Round
		inst.State.DecidedValue = signedMsg.FullData
	} else { // decide previously, add if has more signers
		signers, _ := inst.State.CommitContainer.LongestUniqueSignersForRoundAndRoot(msg.Round, msg.Root)
		if len(signedMsg.GetOperatorIDs()) > len(signers) {
			err := inst.State.CommitContainer.AddMsg(signedMsg)
			if err != nil {
				return nil, err
			}
		}
	}

	if isFutureDecided {
		// bump height
		c.Height = msg.Height
	}

	if !prevDecided {
		return signedMsg, nil
	}
	return nil, nil
}

func ValidateDecided(
	config IConfig,
	signedDecided *types.SignedSSVMessage,
	share *types.SharedValidator,
) error {

	isDecided, err := IsDecidedMsg(share, signedDecided)
	if err != nil {
		return err
	}
	if !isDecided {
		return errors.New("not a decided msg")
	}

	if err := signedDecided.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	msg, err := DecodeMessage(signedDecided.SSVMessage.Data)
	if err != nil {
		return err
	}

	if err := baseCommitValidationVerifySignature(config, signedDecided, msg.Height, share.Committee); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if err := signedDecided.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided")
	}

	r, err := HashDataRoot(signedDecided.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(r[:], msg.Root[:]) {
		return errors.New("H(data) != root")
	}

	return nil
}

// IsDecidedMsg returns true if signed commit has all quorum sigs
func IsDecidedMsg(share *types.SharedValidator, signedDecided *types.SignedSSVMessage) (bool, error) {

	msg, err := DecodeMessage(signedDecided.SSVMessage.Data)
	if err != nil {
		return false, err
	}

	return share.HasQuorum(len(signedDecided.GetOperatorIDs())) && msg.MsgType == CommitMsgType, nil
}
