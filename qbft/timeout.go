package qbft

import (
	"time"

	"github.com/pkg/errors"
)

// At the moment, these variables are useful for the node implementation
// TODO: It's usage should be updated in the PR https://github.com/ssvlabs/ssv-spec/pull/352
var (
	QuickTimeoutThreshold = Round(8)
	QuickTimeout          = 2 * time.Second
	SlowTimeout           = 2 * time.Minute
	// CutoffRound which round the instance should stop its timer and progress no further
	CutoffRound = 12 // stop processing attestations after 8*2+120*3 = 6.2 min (~ 1 epoch)
)

func (i *Instance) UponRoundTimeout() error {
	if !i.CanProcessMessages() {
		return errors.New("instance stopped processing timeouts")
	}

	newRound := i.State.Round + 1
	defer func() {
		i.State.Round = newRound
		i.State.ProposalAcceptedForCurrentRound = nil
		i.Config.GetTimer().TimeoutForRound(i.State.Round)
	}()

	roundChange, err := CreateRoundChange(i.State, i.signer, newRound, i.StartValue)
	if err != nil {
		return errors.Wrap(err, "could not generate round change msg")
	}

	if err := i.Broadcast(roundChange); err != nil {
		return errors.Wrap(err, "failed to broadcast round change message")
	}

	return nil
}
