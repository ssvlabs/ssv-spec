package valcheck

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type SpecTest struct {
	Name               string
	Network            types.BeaconNetwork
	RunnerRole         types.RunnerRole
	Input              []byte
	SlashableDataRoots map[string][][]byte      // map share pk to a list of slashable data roots
	ShareValidatorsPK  []types.ShareValidatorPK `json:"omitempty"` // Optional. Specify validators shares for beacon vote value check
	ExpectedError      string
	AnyError           bool
}

func (test *SpecTest) TestName() string {
	return test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	signer := testingutils.NewTestingKeyManager()
	if len(test.SlashableDataRoots) > 0 {
		signer = testingutils.NewTestingKeyManagerWithSlashableRoots(test.SlashableDataRoots)
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
		return ssv.BeaconVoteValueCheckF(signer, testingutils.TestingDutySlot, shareValidatorsPK,
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
