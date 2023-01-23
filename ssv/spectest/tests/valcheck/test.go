package valcheck

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type SpecTest struct {
	Name               string
	Network            types.BeaconNetwork
	BeaconRole         types.BeaconRole
	Input              []byte
	SlashableDataRoots [][]byte
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
	switch test.BeaconRole {
	case types.BNRoleAttester:
		return ssv.AttesterValueCheckF(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex, nil)
	case types.BNRoleProposer:
		return ssv.ProposerValueCheckF(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex, nil)
	case types.BNRoleAggregator:
		return ssv.AggregatorValueCheckF(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex)
	case types.BNRoleSyncCommittee:
		return ssv.SyncCommitteeValueCheckF(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex)
	case types.BNRoleSyncCommitteeContribution:
		return ssv.SyncCommitteeContributionValueCheckF(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex)
	default:
		panic("unknown role")
	}
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
