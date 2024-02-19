package validator

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type ValidatorTest struct {
	Name                   string
	Duties                 []*types.Duty
	Messages               []*types.SSVMessage
	OutputMessages         []*types.SSVMessage
	BeaconBroadcastedRoots []string
	ExpectedError          string
}

func (test *ValidatorTest) TestName() string {
	return "validator " + test.Name
}

func (test *ValidatorTest) Prepare() (*ssv.Validator, error) {
	ks := testingutils.Testing4SharesSet()

	v := testingutils.BaseValidator(ks)

	// Init all duty runners
	v.DutyRunners[types.BNRoleAttester] = testingutils.AttesterRunner(ks)
	v.DutyRunners[types.BNRoleAggregator] = testingutils.AggregatorRunner(ks)
	v.DutyRunners[types.BNRoleProposer] = testingutils.ProposerRunner(ks)
	v.DutyRunners[types.BNRoleSyncCommittee] = testingutils.SyncCommitteeRunner(ks)
	v.DutyRunners[types.BNRoleSyncCommitteeContribution] = testingutils.SyncCommitteeContributionRunner(ks)
	v.DutyRunners[types.BNRoleValidatorRegistration] = testingutils.ValidatorRegistrationRunner(ks)
	v.DutyRunners[types.BNRoleVoluntaryExit] = testingutils.VoluntaryExitRunner(ks)

	v.Network = v.DutyRunners[types.BNRoleAttester].GetNetwork()

	var lastErr error

	// Start each requested duty
	for _, duty := range test.Duties {
		lastErr = v.StartDuty(duty)
	}

	// process each message and store the last error
	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}
	return v, lastErr
}

// Run as an individual test
func (test *ValidatorTest) Run(t *testing.T) {

	v, lastErr := test.Prepare()

	// Check error
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	// test output message
	test.compareOutputMsgs(t, v)

	// test beacon broadcasted msgs
	test.compareBroadcastedBeaconMsgs(t, v)
}

func (test *ValidatorTest) compareBroadcastedBeaconMsgs(t *testing.T, v *ssv.Validator) {

	broadcastedRoots := make([]phase0.Root, 0)

	for _, role := range []types.BeaconRole{types.BNRoleAttester, types.BNRoleAggregator,
		types.BNRoleProposer,
		types.BNRoleSyncCommittee, types.BNRoleSyncCommitteeContribution,
		types.BNRoleValidatorRegistration, types.BNRoleVoluntaryExit} {

		broadcastedRoots = append(broadcastedRoots, v.DutyRunners[role].GetBeaconNode().(*testingutils.TestingBeaconNode).BroadcastedRoots...)
	}

	require.Len(t, broadcastedRoots, len(test.BeaconBroadcastedRoots))
	alreadyMatched := make([]bool, len(test.BeaconBroadcastedRoots))
	for _, r1 := range test.BeaconBroadcastedRoots {
		found := false
		for r2Idx, r2 := range broadcastedRoots {
			if alreadyMatched[r2Idx] {
				continue
			}
			if r1 == hex.EncodeToString(r2[:]) {
				found = true
				alreadyMatched[r2Idx] = true
				break
			}
		}
		require.Truef(t, found, "broadcasted beacon root not found")
	}
}

func (test *ValidatorTest) compareOutputMsgs(t *testing.T, v *ssv.Validator) {

	broadcastedMsgs := v.Network.(*testingutils.TestingNetwork).BroadcastedMsgs
	require.Len(t, broadcastedMsgs, len(test.OutputMessages))
	alreadyMatched := make([]bool, len(test.OutputMessages))
	for _, msg := range test.OutputMessages {
		found := false
		for index2, msg2 := range broadcastedMsgs {
			if alreadyMatched[index2] {
				continue
			}
			if msg.GetType() == msg2.GetType() && msg.GetID() == msg2.GetID() && bytes.Equal(msg.GetData(), msg2.GetData()) {
				found = true
				alreadyMatched[index2] = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Broadcasted message not found in output messages: %v", msg))
		}
	}
}

func (test *ValidatorTest) GetPostState() (interface{}, error) {
	validator, lastErr := test.Prepare()

	if len(test.ExpectedError) == 0 {
		if lastErr != nil {
			return nil, lastErr
		}
		return validator, nil
	} else {
		if lastErr == nil || (test.ExpectedError != lastErr.Error()) {
			return nil, lastErr
		}
		return validator, nil
	}
}
