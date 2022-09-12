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

type MsgProcessingSpecTest struct {
	Name              string
	Pre               *qbft.Instance
	PostRoot          string
	InputMessages     []*qbft.SignedMessage
	OutputMessages    []*qbft.SignedMessage
	InputMessagesSIP  []*types.Message
	OutputMessagesSIP []*types.Message
	ExpectedError     string
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	var lastErr error
	for _, msg := range test.InputMessagesSIP {
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			lastErr = err
		}

		_, _, _, err := test.Pre.ProcessMsg(msg.GetID(), signedMsg)
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
	if len(test.OutputMessagesSIP) > 0 || len(broadcastedMsgs) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessagesSIP))

		for i, msg := range test.OutputMessagesSIP {
			msg1 := &qbft.SignedMessage{}
			require.NoError(t, msg1.Decode(msg.Data))
			r1, _ := msg1.GetRoot()

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
