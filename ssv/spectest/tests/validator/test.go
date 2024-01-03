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

// RunAsPartOfMultiTest runs the test as part of a MultiMsgProcessingSpecTest
func (test *ValidatorTest) RunAsPartOfMultiTest(t *testing.T) {

}

// Run as an individual test
func (test *ValidatorTest) Run(t *testing.T) {

	ks := testingutils.Testing4SharesSet()

	v := testingutils.BaseValidator(ks)

	v.DutyRunners[types.BNRoleAttester] = testingutils.AttesterRunner(ks)
	v.DutyRunners[types.BNRoleAggregator] = testingutils.AggregatorRunner(ks)
	v.DutyRunners[types.BNRoleProposer] = testingutils.ProposerRunner(ks)
	v.DutyRunners[types.BNRoleSyncCommittee] = testingutils.SyncCommitteeRunner(ks)
	v.DutyRunners[types.BNRoleSyncCommitteeContribution] = testingutils.SyncCommitteeContributionRunner(ks)
	v.DutyRunners[types.BNRoleValidatorRegistration] = testingutils.ValidatorRegistrationRunner(ks)
	v.DutyRunners[types.BNRoleVoluntaryExit] = testingutils.VoluntaryExitRunner(ks)

	v.Network = v.DutyRunners[types.BNRoleAttester].GetNetwork()

	var lastErr error

	for _, duty := range test.Duties {
		lastErr = v.StartDuty(duty)
	}

	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}

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
	for _, r1 := range test.BeaconBroadcastedRoots {
		found := false
		for _, r2 := range broadcastedRoots {
			if r1 == hex.EncodeToString(r2[:]) {
				found = true
				break
			}
		}
		require.Truef(t, found, "broadcasted beacon root not found")
	}
}

func (test *ValidatorTest) compareOutputMsgs(t *testing.T, v *ssv.Validator) {

	broadcastedMsgs := v.Network.(*testingutils.TestingNetwork).BroadcastedMsgs
	require.Len(t, broadcastedMsgs, len(test.OutputMessages))
	usedIndexes := make(map[int]bool)
	for _, msg := range test.OutputMessages {
		found := false
		for index2, msg2 := range broadcastedMsgs {
			if usedIndexes[index2] {
				continue
			}
			if bytes.Equal(msg.GetData(), msg2.GetData()) {
				if msg.GetType() == msg2.GetType() && msg.GetID() == msg2.GetID() {
					found = true
					usedIndexes[index2] = true
					break
				}
			}
		}
		if !found {
			panic(fmt.Sprintf("Broadcasted message not found in output messages: %v", msg))
		}
	}
}

func (test *ValidatorTest) GetPostState() (interface{}, error) {
	return nil, nil
}
