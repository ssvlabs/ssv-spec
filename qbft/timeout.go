package qbft

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

var (
	quickTimeoutThreshold = Round(8) //nolint
	// quickTimeout is the timeout in seconds for the first 8 rounds
	quickTimeout uint64 = 2 // 2 seconds
	// slowTimeout is the timeout in seconds for rounds after the first 8
	slowTimeout uint64 = 120 // 2 minutes
	// CutoffRound which round the instance should stop its timer and progress no further
	CutoffRound = 15 // stop processing instances after 8*2+120*6 = 14.2 min (~ 2 epochs)
)

type RoundTimer struct {
	// role is the role of the QBFT instance
	Role types.BeaconRole
	// height is the current height of the instance
	Height Height
	// network is the beacon network
	Network types.BeaconNetwork
	//current unix epoch time in seconds
	CurrentTime uint64
}

func (t *RoundTimer) TimeoutForRound(round Round) uint64 {
	switch t.Role {
	case types.BNRoleAttester | types.BNRoleSyncCommittee:
		return AttestationOrSyncCommitteeTimeout(round, t.Height, t.Network) - t.CurrentTime
	case types.BNRoleAggregator | types.BNRoleSyncCommitteeContribution:
		return AggregationOrContributionTimeout(round, t.Height, t.Network) - t.CurrentTime
	default:
		return DefaultTimeout(round)
	}
}

// RoundTimeout returns the unix epoch time (seconds) in which we should send a RC message
// Called for all beacon duties other than proposals
func RoundTimeout(round Round, height Height, baseDuration uint64, network types.BeaconNetwork) uint64 {
	// Calculate additional timeout in seconds based on round
	var additionalTimeout uint64
	additionalTimeout = uint64(round) * quickTimeout
	if round > quickTimeoutThreshold {
		slowPortion := uint64(round-quickTimeoutThreshold) * slowTimeout
		additionalTimeout += slowPortion
	}

	// Combine base duration and additional timeout
	timeoutDuration := baseDuration + additionalTimeout

	// Get the unix epoch start time of the duty seconds
	dutyStartTime := uint64(network.EstimatedTimeAtSlot(phase0.Slot(height)))

	// Calculate the time until the duty should start plus the timeout duration
	return dutyStartTime + timeoutDuration
}

// AttestationOrSyncCommitteeTimeout returns the unix epoch time (seconds) in which we should send a RC message
func AttestationOrSyncCommitteeTimeout(round Round, height Height, network types.BeaconNetwork) uint64 {
	return RoundTimeout(round, height, 4, network)
}

// AggregationOrContributionTimeout returns the unix epoch time (seconds) in which we should send a RC message
func AggregationOrContributionTimeout(r Round, height Height, network types.BeaconNetwork) uint64 {
	return RoundTimeout(r, height, 8, network)
}

// DefaultTimeout returns the duration in seconds in which we should send a RC message
func DefaultTimeout(round Round) uint64 {
	if round <= quickTimeoutThreshold {
		return quickTimeout
	} else {
		return slowTimeout
	}
}

func (i *Instance) UponRoundTimeout() error {
	if !i.CanProcessMessages() {
		return errors.New("instance stopped processing timeouts")
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
