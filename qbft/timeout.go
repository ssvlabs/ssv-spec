package qbft

import (
	"github.com/pkg/errors"
)

type Timer interface {
	// TimeoutForRound will reset running timer if exists and will start a new timer for a specific round
	TimeoutForRound(round Round)
}

// RoundTimeout returns the number of seconds until next timeout for a give round
func (i *Instance) RoundTimeout(round Round) uint64 {
	return powerOf2(uint64(round))
}

func powerOf2(exp uint64) uint64 {
	if exp == 0 {
		return 1
	} else {
		return 2 * powerOf2(exp-1)
	}
}

func (i *Instance) UponRoundTimeout() error {
	newRound := i.State.Round + 1
	defer func() {
		i.State.Round = newRound
		i.State.ProposalAcceptedForCurrentRound = nil
		i.config.GetTimer().TimeoutForRound(i.State.Round)
	}()

	roundChange, err := CreateRoundChange(i.State, i.config, newRound, i.StartValue)
	if err != nil {
		return errors.Wrap(err, "could not generate round change msg")
	}

	if err := i.Broadcast(roundChange); err != nil {
		return errors.Wrap(err, "failed to broadcast round change message")
	}

	return nil
}
