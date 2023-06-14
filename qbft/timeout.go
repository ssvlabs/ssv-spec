package qbft

import (
	"github.com/pkg/errors"
)

// CutoffRound which round the instance should stop its timer and progress no further
const CutoffRound = 20

func (i *Instance) UponRoundTimeout() error {
	if i.State.Round == CutoffRound {
		return errors.New("round > cutoff round")
	}

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
