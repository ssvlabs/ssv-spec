package beacon

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
	"github.com/stretchr/testify/require"
)

type DepositDataSpecTest struct {
	Name                  string
	Type                  string
	Documentation         string
	ValidatorPK           string
	WithdrawalCredentials string
	ForkVersion           types.ForkVersion
	ExpectedSigningRoot   string
}

func (test *DepositDataSpecTest) TestName() string {
	return test.Name
}

func (test *DepositDataSpecTest) Run(t *testing.T) {
	validatorPK, err := hex.DecodeString(test.ValidatorPK)
	require.NoError(t, err)
	withdrawalCredentials, err := hex.DecodeString(test.WithdrawalCredentials)
	require.NoError(t, err)

	r, _, err := testingutils.GenerateETHDepositData(
		validatorPK,
		withdrawalCredentials,
		phase0.Version(test.ForkVersion[:]),
		types.DomainDeposit,
	)
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedSigningRoot, hex.EncodeToString(r))

	comparable2.CompareWithJson(t, test, test.TestName(), reflect.TypeOf(test).String())
}

// creates a new DepositDataSpecTest with the Type field automatically set
func NewDepositDataSpecTest(name, documentation, validatorPK, withdrawalCredentials string, forkVersion [4]byte, expectedSigningRoot string) *DepositDataSpecTest {
	return &DepositDataSpecTest{
		Name:                  name,
		Type:                  testdoc.DepositDataSpecTestType,
		Documentation:         documentation,
		ValidatorPK:           validatorPK,
		WithdrawalCredentials: withdrawalCredentials,
		ForkVersion:           forkVersion,
		ExpectedSigningRoot:   expectedSigningRoot,
	}
}

// DepositData tests structuring of encoding data
func DepositData() *DepositDataSpecTest {
	return NewDepositDataSpecTest(
		"deposit data root and ssz",
		testdoc.DepositDataTestDoc,
		"b3d50de8d77299da8d830de1edfb34d3ce03c1941846e73870bb33f6de7b8a01383f6b32f55a1d038a4ddcb21a765194",
		"005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16",
		types.MainNetwork.ForkVersion(),
		"69d2af2fd5870077e45f574087a38f476ac3b0f680a511767fb1b0f17f8c4cbd",
	)
}
