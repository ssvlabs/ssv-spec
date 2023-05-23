package beacon

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type DepositDataSpecTest struct {
	Name                  string
	ValidatorPK           string
	WithdrawalCredentials string
	ForkVersion           [4]byte
	ExpectedSigningRoot   string
}

func (test *DepositDataSpecTest) TestName() string {
	return test.Name
}

func (test *DepositDataSpecTest) Run(t *testing.T) {
	validatorPK, _ := hex.DecodeString(test.ValidatorPK)
	withdrawalCredentials, _ := hex.DecodeString(test.WithdrawalCredentials)

	r, _, err := types.GenerateETHDepositData(
		validatorPK,
		withdrawalCredentials,
		test.ForkVersion,
		types.DomainDeposit,
	)
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedSigningRoot, hex.EncodeToString(r))
}

// DepositData tests structuring of encoding data
func DepositData() *DepositDataSpecTest {
	return &DepositDataSpecTest{
		Name:                  "deposit data root and ssz",
		ValidatorPK:           "b3d50de8d77299da8d830de1edfb34d3ce03c1941846e73870bb33f6de7b8a01383f6b32f55a1d038a4ddcb21a765194",
		WithdrawalCredentials: "005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16",
		ForkVersion:           types.TestNetwork.ForkVersion(),
		ExpectedSigningRoot:   "500a3fd02057c19ec5e1e785d8cec46d0f3ee21185167a8c18996b91aed886b7",
	}
}
