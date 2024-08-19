package valcheck

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type SpecTest struct {
	Name             string
	Network          types.BeaconNetwork
	RunnerRole       types.RunnerRole
	Duty             types.Duty
	Input            []byte
	SlashableSlots   map[string][]phase0.Slot               // map share pk to a list of slashable slots
	ValidatorsShares map[phase0.ValidatorIndex]*types.Share `json:"omitempty"` // Optional.
	// Specify validators shares for beacon vote value check
	ExpectedError string
	AnyError      bool
}

func (test *SpecTest) TestName() string {
	return test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	signer := testingutils.NewTestingKeyManager()
	if len(test.SlashableSlots) > 0 {
		signer = testingutils.NewTestingKeyManagerWithSlashableSlots(test.SlashableSlots)
	}

	check, err := test.valCheckF(signer)

	if err != nil && len(test.ExpectedError) > 0 {
		require.EqualError(t, err, test.ExpectedError)
		return
	}

	err = check(test.Input)

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

func (test *SpecTest) valCheckF(signer types.BeaconSigner) (qbft.ProposedValueCheckF, error) {
	pubKeyBytes := types.ValidatorPK(testingutils.TestingValidatorPubKey)
	switch test.RunnerRole {
	case types.RoleCommittee:
		return ssv.BeaconVoteValueCheckF(test.Duty.(*types.CommitteeDuty), signer, test.Network,
			test.ValidatorsShares)
	case types.RoleProposer:
		return ssv.ProposerValueCheckF(signer, test.Network, pubKeyBytes, testingutils.TestingValidatorIndex, nil),
			nil
	case types.RoleAggregator:
		return ssv.AggregatorValueCheckF(signer, test.Network, pubKeyBytes, testingutils.TestingValidatorIndex),
			nil
	case types.RoleSyncCommitteeContribution:
		return ssv.SyncCommitteeContributionValueCheckF(signer, test.Network, pubKeyBytes,
				testingutils.TestingValidatorIndex),
			nil
	default:
		panic("unknown role")
	}
}

func (tests *SpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

type MultiSpecTest struct {
	Name  string
	Tests []*SpecTest
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

// Create alias without duty
type SpecTestAlias struct {
	Name             string
	Network          types.BeaconNetwork
	RunnerRole       types.RunnerRole
	Input            []byte
	SlashableSlots   map[string][]phase0.Slot               // map share pk to a list of slashable slots
	ValidatorsShares map[phase0.ValidatorIndex]*types.Share `json:"omitempty"` // Optional.
	// Specify validators shares for beacon vote value check
	ExpectedError string
	AnyError      bool
	ValidatorDuty *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
	CommitteeDuty *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
}

func (t *SpecTest) MarshalJSON() ([]byte, error) {
	alias := SpecTestAlias{
		Name:             t.Name,
		Network:          t.Network,
		RunnerRole:       t.RunnerRole,
		Input:            t.Input,
		SlashableSlots:   t.SlashableSlots,
		ValidatorsShares: t.ValidatorsShares,
		ExpectedError:    t.ExpectedError,
		AnyError:         t.AnyError,
	}

	if t.Duty != nil {
		switch t.Duty.(type) {
		case *types.ValidatorDuty:
			alias.ValidatorDuty = t.Duty.(*types.ValidatorDuty)
		case *types.CommitteeDuty:
			alias.CommitteeDuty = t.Duty.(*types.CommitteeDuty)
		}
	}

	return json.Marshal(alias)
}

func (t *SpecTest) UnmarshalJSON(data []byte) error {
	var alias SpecTestAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	t.Name = alias.Name
	t.Network = alias.Network
	t.RunnerRole = alias.RunnerRole
	t.Input = alias.Input
	t.SlashableSlots = alias.SlashableSlots
	t.ValidatorsShares = alias.ValidatorsShares
	t.ExpectedError = alias.ExpectedError
	t.AnyError = alias.AnyError

	if alias.ValidatorDuty != nil {
		t.Duty = alias.ValidatorDuty
	} else if alias.CommitteeDuty != nil {
		t.Duty = alias.CommitteeDuty
	}

	return nil
}
