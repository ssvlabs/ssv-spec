package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type MsgProcessingSpecTest struct {
	Name        string
	TestingNode *dkg.Node

	InputMessages  []*dkg.SignedMessage
	OutputMessages []*dkg.SignedMessage
	Output         map[types.OperatorID]*dkg.SignedOutput
	KeySet         *testingutils.TestKeySet
	ExpectedError  string

	LastMsgDelay *time.Duration
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	frost.DefaultTimeoutDuration = 2 * time.Second // to keep tests short

	testingutils.ResetRandSeed()
	node := test.TestingNode

	var lastErr error
	for idx, msg := range test.InputMessages {

		if idx+1 == len(test.InputMessages) && test.LastMsgDelay != nil {
			time.Sleep(*test.LastMsgDelay) // adds delay to last input msg to emulate timeout
		}

		byts, _ := msg.Encode()
		err := node.ProcessMessage(&types.SSVMessage{
			MsgType: types.DKGMsgType,
			Data:    byts,
		})

		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	// test output message
	broadcastedMsgs := node.GetConfig().GetNetwork().(*testingutils.TestingNetwork).GetBroadcastMessages()
	if len(test.OutputMessages) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			bMsg := broadcastedMsgs[i]
			require.Equal(t, types.DKGMsgType, bMsg.MsgType)
			sMsg := &dkg.SignedMessage{}
			require.NoError(t, sMsg.Decode(bMsg.Data))

			if sMsg.Message.MsgType == dkg.OutputMsgType {
				require.Equal(t, dkg.OutputMsgType, msg.Message.MsgType, "OutputMsgType expected")
				o1 := &dkg.SignedOutput{}
				require.NoError(t, o1.Decode(msg.Message.Data))

				o2 := &dkg.SignedOutput{}
				require.NoError(t, o2.Decode(sMsg.Message.Data))

				if o1.Data != nil && o2.Data != nil {
					es1 := o1.Data.EncryptedShare
					o1.Data.EncryptedShare = nil
					es2 := o2.Data.EncryptedShare
					o2.Data.EncryptedShare = nil

					s1, _ := types.Decrypt(test.KeySet.DKGOperators[msg.Signer].EncryptionKey, es1)
					s2, _ := types.Decrypt(test.KeySet.DKGOperators[msg.Signer].EncryptionKey, es2)
					require.Equal(t, s1, s2, "shares don't match")
					r1, _ := o1.Data.GetRoot()
					r2, _ := o2.Data.GetRoot()
					require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
				}

				if o1.KeySignData != nil && o2.KeySignData != nil {
					r1, _ := o1.KeySignData.GetRoot()
					r2, _ := o2.KeySignData.GetRoot()
					require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
				}
			} else {
				r1, _ := msg.GetRoot()
				r2, _ := sMsg.GetRoot()
				require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
			}

		}
	}

	streamed := node.GetConfig().Network.(*testingutils.TestingNetwork).DKGOutputs
	if len(test.Output) > 0 {
		require.Len(t, streamed, len(test.Output))
		for id, output := range test.Output {
			s := streamed[id]
			require.NotNilf(t, s, "output for operator %d not found", id)

			var r1, r2 []byte
			if output.KeySignData != nil && s.KeySignData != nil {
				r1, _ = output.KeySignData.GetRoot()
				r2, _ = s.KeySignData.GetRoot()
			} else {
				r1, _ = output.Data.GetRoot()
				r2, _ = s.Data.GetRoot()
			}
			require.EqualValues(t, r1, r2, fmt.Sprintf("output for operator %d roots not equal", id))
		}
	}
}
