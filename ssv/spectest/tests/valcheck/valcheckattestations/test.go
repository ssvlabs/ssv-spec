package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type AttestationValCheckSpecTest struct {
	Name                          string
	Network                       types.BeaconNetwork
	Input                         []byte
	SlashableAttestationDataRoots [][]byte
	ExpectedError                 string
	AnyError                      bool
}

func (test *AttestationValCheckSpecTest) TestName() string {
	return test.Name
}

func (test *AttestationValCheckSpecTest) Run(t *testing.T) {
	signer := testingutils.NewTestingKeyManager()
	if len(test.SlashableAttestationDataRoots) > 0 {
		signer = testingutils.NewTestingKeyManagerWithSlashableRoots(test.SlashableAttestationDataRoots)
	}

	check := ssv.BeaconAttestationValueCheck(signer, test.Network, testingutils.TestingValidatorPubKey[:], testingutils.TestingValidatorIndex)

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
