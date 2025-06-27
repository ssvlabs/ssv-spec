package valcheck

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type SpecTest struct {
	Name              string
	Type              string
	Documentation     string
	Network           types.BeaconNetwork
	RunnerRole        types.RunnerRole
	DutySlot          phase0.Slot // DutySlot is used only for the RoleCommittee since the BeaconVoteValueCheckF requires the duty's slot
	Input             []byte
	SlashableSlots    map[string][]phase0.Slot // map share pk to a list of slashable slots
	ShareValidatorsPK []types.ShareValidatorPK `json:"omitempty"` // Optional. Specify validators shares for beacon vote value check
	ExpectedError     string
	AnyError          bool
}

func (test *SpecTest) TestName() string {
	return test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	signer := testingutils.NewTestingKeyManager()
	if len(test.SlashableSlots) > 0 {
		signer = testingutils.NewTestingKeyManagerWithSlashableSlots(test.SlashableSlots)
	}

	check := test.valCheckF(signer)

	err := check(test.Input)

	if test.AnyError {
		require.NotNil(t, err)
		return
	}
	if len(test.ExpectedError) > 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}
}

func (test *SpecTest) valCheckF(signer types.BeaconSigner) qbft.ProposedValueCheckF {
	pubKeyBytes := types.ValidatorPK(testingutils.TestingValidatorPubKey)

	shareValidatorsPK := test.ShareValidatorsPK
	if len(shareValidatorsPK) == 0 {
		keySet := testingutils.Testing4SharesSet()
		sharePK := keySet.Shares[1]
		sharePKBytes := sharePK.Serialize()
		shareValidatorsPK = []types.ShareValidatorPK{sharePKBytes}
	}
	switch test.RunnerRole {
	case types.RoleCommittee:
		return ssv.BeaconVoteValueCheckF(signer, test.DutySlot, shareValidatorsPK,
			testingutils.TestingDutyEpoch)
	case types.RoleProposer:
		return ssv.ProposerValueCheckF(signer, test.Network, pubKeyBytes, testingutils.TestingValidatorIndex, nil)
	case types.RoleAggregator:
		return ssv.AggregatorValueCheckF(signer, test.Network, pubKeyBytes, testingutils.TestingValidatorIndex)
	case types.RoleSyncCommitteeContribution:
		return ssv.SyncCommitteeContributionValueCheckF(signer, test.Network, pubKeyBytes, testingutils.TestingValidatorIndex)
	default:
		panic("unknown role")
	}
}

func (tests *SpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

type MultiSpecTest struct {
	Name          string
	Type          string
	Documentation string
	Tests         []*SpecTest
}

func (test *MultiSpecTest) TestName() string {
	return test.Name
}

func (test *MultiSpecTest) Run(t *testing.T) {
	for _, test := range test.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func (tests *MultiSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

func NewSpecTest(name, documentation string, network types.BeaconNetwork, role types.RunnerRole, dutySlot phase0.Slot, input []byte, slashableSlots map[string][]phase0.Slot, shareValidatorsPK []types.ShareValidatorPK, expectedError string, anyError bool) *SpecTest {
	return &SpecTest{
		Name:              name,
		Type:              "Value check: validations for input of different runner roles",
		Documentation:     documentation,
		Network:           network,
		RunnerRole:        role,
		DutySlot:          dutySlot,
		Input:             input,
		SlashableSlots:    slashableSlots,
		ShareValidatorsPK: shareValidatorsPK,
		ExpectedError:     expectedError,
		AnyError:          anyError,
	}
}

func NewMultiSpecTest(name, documentation string, tests []*SpecTest) *MultiSpecTest {
	return &MultiSpecTest{
		Name:          name,
		Type:          "Multi value check: multiple value check tests",
		Documentation: documentation,
		Tests:         tests,
	}
}
