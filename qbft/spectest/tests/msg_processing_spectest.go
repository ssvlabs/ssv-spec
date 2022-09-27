package tests

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

// ChangeProposerFuncInstanceHeight tests with this height will return proposer operator ID 2
const ChangeProposerFuncInstanceHeight = 10

type MsgProcessingSpecTest struct {
	Name           string
	Pre            *qbft.Instance
	PostRoot       string
	InputMessages  []*qbft.SignedMessage
	OutputMessages []*qbft.SignedMessage
	ExpectedError  string
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	// a simple hack to change the proposer func
	if test.Pre.State.Height == ChangeProposerFuncInstanceHeight {
		test.Pre.GetConfig().(*qbft.Config).ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 2
		}
	}

	var lastErr error
	for _, msg := range test.InputMessages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	// test output message
	broadcastedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	if len(test.OutputMessages) > 0 || len(broadcastedMsgs) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			r1, _ := msg.GetRoot()

			msg2 := &qbft.SignedMessage{}
			require.NoError(t, msg2.Decode(broadcastedMsgs[i].Data))
			r2, _ := msg2.GetRoot()

			require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
		}
	}

	require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}

func (test *MsgProcessingSpecTest) TestName() string {
	return "qbft message processing " + test.Name
}
