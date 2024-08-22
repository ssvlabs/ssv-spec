package validation

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Maximum amount of messages per type
const (
	maxPreConsensus  = 1
	maxProposals     = 1
	maxPrepares      = 1
	maxCommits       = 1
	maxRoundChanges  = 1
	maxPostConsensus = 1
)

// MessageCounts tracks the number of various message types received for validation
type MessageCounts struct {
	PreConsensus  int
	Proposal      int
	Prepare       int
	Commit        int
	RoundChange   int
	PostConsensus int
}

// Reset amounts for a new round. PreConsensus and PostConsensus are kept fixed
func (mc *MessageCounts) ResetForRound() {
	mc.Proposal = 0
	mc.Prepare = 0
	mc.Commit = 0
	mc.RoundChange = 0
}

// Checks if the provided consensus message exceeds the limits
func (c *MessageCounts) ValidateConsensusMessage(msgType qbft.MessageType, numSigners int) error {
	switch msgType {
	case qbft.ProposalMsgType:
		if c.Proposal >= maxProposals {
			err := ErrDuplicatedMessage
			return err
		}
	case qbft.PrepareMsgType:
		if c.Prepare >= maxPrepares {
			err := ErrDuplicatedMessage
			return err
		}
	case qbft.CommitMsgType:
		if numSigners == 1 {
			if c.Commit >= maxCommits {
				err := ErrDuplicatedMessage
				return err
			}
		}
	case qbft.RoundChangeMsgType:
		if c.RoundChange >= maxRoundChanges {
			err := ErrDuplicatedMessage
			return err
		}
	}

	return nil
}

// Checks if the provided partial signature message exceeds the limits
func (c *MessageCounts) ValidatePartialSignatureMessage(msgType types.PartialSigMsgType) error {
	switch msgType {
	case types.RandaoPartialSig, types.SelectionProofPartialSig, types.ContributionProofs, types.ValidatorRegistrationPartialSig, types.VoluntaryExitPartialSig:
		if c.PreConsensus >= maxPreConsensus {
			err := ErrInvalidPartialSignatureTypeCount
			return err
		}
	case types.PostConsensusPartialSig:
		if c.PostConsensus >= maxPostConsensus {
			err := ErrInvalidPartialSignatureTypeCount
			return err
		}
	}

	return nil
}

// Updates the registers based on the provided consensus message type
func (c *MessageCounts) RecordConsensusMessage(msgType qbft.MessageType, signers int) {
	switch msgType {
	case qbft.ProposalMsgType:
		c.Proposal++
	case qbft.PrepareMsgType:
		c.Prepare++
	case qbft.CommitMsgType:
		if signers == 1 {
			c.Commit++
		}
	case qbft.RoundChangeMsgType:
		c.RoundChange++
	}
}

// Updates the registers based on the provided partial signature message type
func (c *MessageCounts) RecordPartialSignatureMessage(msgType types.PartialSigMsgType) {
	switch msgType {
	case types.RandaoPartialSig, types.SelectionProofPartialSig, types.ContributionProofs, types.ValidatorRegistrationPartialSig, types.VoluntaryExitPartialSig:
		c.PreConsensus++
	case types.PostConsensusPartialSig:
		c.PostConsensus++
	}
}
