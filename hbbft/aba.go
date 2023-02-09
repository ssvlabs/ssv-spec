package hbbft

import (
	"github.com/pkg/errors"
)

func (i *Instance) StartABA(vote byte) (byte, error) {
	// current ABA state (associated with current ACS)
	abaState := i.State.HBBFTState.GetCurrentABAState()

	// set ABA's input value
	abaState.setVInput(abaState.Round, vote)

	// broadcast INIT message with input vote
	initMsg, err := CreateABAInit(i.State, i.config, vote, abaState.Round, i.State.HBBFTState.GetRound())
	if err != nil {
		return byte(2), errors.Wrap(err, "StartABA: failed to create ABA Init message")
	}
	i.Broadcast(initMsg)

	// update sent flag
	abaState.setSentInit(abaState.Round, vote, true)

	// process own init msg
	i.uponABAInit(initMsg)

	// wait until channel Terminate receives a signal
	for {
		if abaState.Terminate {
			break
		}
	}

	// abaState.Terminate = false

	// returns the decided value
	return abaState.Vdecided, nil
}
