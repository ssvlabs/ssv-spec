package futuremsg

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	typescomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

type ControllerSyncSpecTest struct {
	Name                 string
	InputMessages        []*qbft.SignedMessage
	SyncDecidedCalledCnt int
	ControllerPostRoot   string
	ControllerPostState  types.Root `json:"-"` // Field is ignored by encoding/json
	ExpectedError        string
	SkipInstanceStart    bool
}

func (test *ControllerSyncSpecTest) TestName() string {
	return "qbft controller sync " + test.Name
}

func (test *ControllerSyncSpecTest) Run(t *testing.T) {
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)

	if !test.SkipInstanceStart {
		err := contr.StartNewInstance(qbft.FirstHeight, []byte{1, 2, 3, 4})
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	var lastErr error
	for _, msg := range test.InputMessages {
		_, err := contr.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
	}

	syncedDecidedCnt := config.GetNetwork().(*testingutils.TestingNetwork).SyncHighestDecidedCnt
	require.EqualValues(t, test.SyncDecidedCalledCnt, syncedDecidedCnt)

	r, err := contr.GetRoot()
	require.NoError(t, err)
	if test.ControllerPostRoot != hex.EncodeToString(r[:]) {
		diff := typescomparable.PrintDiff(contr, test.ControllerPostState)
		require.Fail(t, "post state not equal", diff)
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSyncSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}
