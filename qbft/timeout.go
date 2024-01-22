package qbft

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

var (
	quickTimeoutThreshold = Round(8) //nolint
	// quickTimeout is the timeout in seconds for the first 8 rounds
	quickTimeout int64 = 2 // 2 seconds
	// slowTimeout is the timeout in seconds for rounds after the first 8
	slowTimeout int64 = 120 // 2 minutes
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
	CurrentTime int64
}

func (t *RoundTimer) TimeoutForRound(round Round) int64 {
	switch t.Role {
	case types.BNRoleAttester, types.BNRoleSyncCommittee:
		return max(AttestationOrSyncCommitteeTimeout(round, t.Height, t.Network)-t.CurrentTime, 0)
	case types.BNRoleAggregator, types.BNRoleSyncCommitteeContribution:
		return max(AggregationOrContributionTimeout(round, t.Height, t.Network)-t.CurrentTime, 0)
	default:
		return DefaultTimeout(round)
	}
}

// RoundTimeout returns the unix epoch time (seconds) in which we should send a RC message
// Called for all beacon duties other than proposals
func RoundTimeout(round Round, height Height, baseDuration int64, network types.BeaconNetwork) int64 {
	// Calculate additional timeout in seconds based on round
	var additionalTimeout int64
	additionalTimeout = int64(min(round, quickTimeoutThreshold)) * quickTimeout
	if round > quickTimeoutThreshold {
		slowPortion := int64(round-quickTimeoutThreshold) * slowTimeout
		additionalTimeout += slowPortion
	}

	// Combine base duration and additional timeout
	timeoutDuration := baseDuration + additionalTimeout

	// Get the unix epoch start time of the duty seconds
	dutyStartTime := int64(network.EstimatedTimeAtSlot(phase0.Slot(height)))

	// Calculate the time until the duty should start plus the timeout duration
	return dutyStartTime + timeoutDuration
}

// AttestationOrSyncCommitteeTimeout returns the unix epoch time (seconds) in which we should send a RC message
func AttestationOrSyncCommitteeTimeout(round Round, height Height, network types.BeaconNetwork) int64 {
	return RoundTimeout(round, height, 4, network)
}

// AggregationOrContributionTimeout returns the unix epoch time (seconds) in which we should send a RC message
func AggregationOrContributionTimeout(r Round, height Height, network types.BeaconNetwork) int64 {
	return RoundTimeout(r, height, 8, network)
}

// DefaultTimeout returns the duration in seconds in which we should send a RC message
func DefaultTimeout(round Round) int64 {
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

// min returns the minimum of two values
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two values
func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
