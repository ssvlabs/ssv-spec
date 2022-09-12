package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (c *Controller) UponFutureMsg(msg *SignedMessage) (*SignedMessage, error) {
	if err := validateFutureMsg(c.GenerateConfig(), msg, c.Share.Committee); err != nil {
		return nil, errors.Wrap(err, "invalid future msg")
	}
	if err := c.verifyAndAddHigherHeightMsg(msg); err != nil {
		return nil, errors.Wrap(err, "failed adding higher height msg")
	}
	if c.f1SyncTrigger() {
		return nil, c.network.SyncHighestDecided(c.Identifier)
	}
	return nil, nil
}

func validateFutureMsg(
	config IConfig,
	msg *SignedMessage,
	operators []*types.Operator,
) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	// verify signature
	if err := msg.Signature.VerifyByOperators(msg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "commit msg signature invalid")
	}

	return nil
}

// verifyAndAddHigherHeightMsg verifies msg, cleanup queue and adds the message if unique signer
func (c *Controller) verifyAndAddHigherHeightMsg(msg *SignedMessage) error {
	if err := msg.Signature.VerifyByOperators(msg, c.Domain, types.QBFTSignatureType, c.Share.Committee); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}
	if len(msg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	// cleanup lower height msgs
	cleanedQueue := make(map[types.OperatorID]Height)
	signerExists := false
	for signer, height := range c.FutureMsgsContainer {
		if height <= c.Height {
			continue
		}

		if signer == msg.GetSigners()[0] {
			signerExists = true
		}
		cleanedQueue[signer] = height
	}

	if !signerExists {
		cleanedQueue[msg.GetSigners()[0]] = msg.Message.Height
	}
	c.FutureMsgsContainer = cleanedQueue
	return nil
}

// f1SyncTrigger returns true if received f+1 higher height messages from unique signers
func (c *Controller) f1SyncTrigger() bool {
	return c.Share.HasPartialQuorum(len(c.FutureMsgsContainer))
}
