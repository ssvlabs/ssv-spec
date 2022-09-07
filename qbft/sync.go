package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (c *Controller) processFutureMsg(msg *SignedMessage) (*SignedMessage, error) {
	if isDecidedMsg(c.Share, msg) {
		return c.UponFutureDecided(msg)
	}

	added, err := c.HigherReceivedMessages.AddFirstMsgForSignerAndRound(msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not add higher height msg")
	}
	if added && c.f1SyncTrigger() {
		// TODO should reset msg container? past msgs? all msgs?
		return nil, c.network.SyncHighestDecided(c.Identifier)
	}
	return nil, nil
}

// f1SyncTrigger returns true if received f+1 higher height messages from unique signers
func (c *Controller) f1SyncTrigger() bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range c.HigherReceivedMessages.AllMessaged() {
		for _, signer := range msg.GetSigners() {
			if _, found := uniqueSigners[signer]; !found {
				uniqueSigners[signer] = true
			}
		}
	}
	return c.Share.HasPartialQuorum(len(uniqueSigners))
}
