package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
)

// ChangeProposerFuncInstanceHeight tests with this height will return proposer operator ID 2
const ChangeProposerFuncInstanceHeight = 1

type MsgProcessingSpecTest struct {
	Name               string
	Pre                *alea.Instance
	PostRoot           string
	InputMessages      []*alea.SignedMessage
	OutputMessages     []*alea.SignedMessage
	ExpectedError      string
	ExpectedTimerState *testingutils.TimerStateAlea
	DontRunAC          bool
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	// a simple hack to change the proposer func
	if test.Pre.State.Height == ChangeProposerFuncInstanceHeight {
		test.Pre.GetConfig().(*alea.Config).ProposerF = func(state *alea.State, round alea.Round) types.OperatorID {
			fmt.Println("Len:", len(state.Share.Committee))
			ans := int(round)%len(state.Share.Committee) + 1
			fmt.Println("ans:", ans)
			return types.OperatorID(ans)
		}
	}

	if !test.DontRunAC {
		go test.Pre.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)
	}
	var lastErr error
	for _, msg := range test.InputMessages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
	}

	test.Pre.State.StopAgreement = true

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	// if test.ExpectedTimerState != nil {
	// 	// checks round timer state
	// 	timer, ok := test.Pre.GetConfig().GetTimer().(*testingutils.TestAleaTimer)
	// 	if ok && timer != nil {
	// 		require.Equal(t, test.ExpectedTimerState.Timeouts, timer.State.Timeouts, "timer should have expected timeouts count")
	// 		require.Equal(t, test.ExpectedTimerState.Round, timer.State.Round, "timer should have expected round")
	// 	}
	// }

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	// test output message
	broadcastedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetworkAlea).BroadcastedMsgs

	// for _, broadcastedMsg := range broadcastedMsgs {
	// 	var message *alea.SignedMessage
	// 	fmt.Println(broadcastedMsg)
	// 	err := message.Decode(broadcastedMsg.Data)
	// 	if err != nil {
	// 		fmt.Println("failed to decode")
	// 	}
	// 	fmt.Println(message)
	// }

	if len(test.OutputMessages) > 0 || len(broadcastedMsgs) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			r1, _ := msg.GetRoot()

			msg2 := &alea.SignedMessage{}
			require.NoError(t, msg2.Decode(broadcastedMsgs[i].Data))
			r2, _ := msg2.GetRoot()

			require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
		}
	}

	require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}

func (test *MsgProcessingSpecTest) TestName() string {
	return "alea message processing " + test.Name
}
