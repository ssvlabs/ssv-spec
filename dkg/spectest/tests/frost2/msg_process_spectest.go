package frost2

import (
	"fmt"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type MsgProcessingSpecTest struct {
	Name           string
	InputMessages  []*dkg.SignedMessage
	OutputMessages []*dkg.SignedMessage
	Output         map[types.OperatorID]*dkg.SignedOutput
	KeySet         *testingutils.TestKeySet
	Operator       *dkg.Operator
	NodeConfig     *dkg.Config
	ExpectedError  string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	testingutils.ResetRandSeed()

	node := dkg.NewNode(test.Operator, test.NodeConfig)

	var lastErr error
	network := test.NodeConfig.Network.(*testingutils.TestingNetwork)
	for _, msg := range test.InputMessages {
		byts, _ := msg.Encode()
		err := node.ProcessMessage(&types.SSVMessage{
			MsgType: types.DKGMsgType,
			Data:    byts,
		})

		if err != nil {
			lastErr = err
		}

		msgs := network.BroadcastedMsgs
		for _, msg := range msgs {
			sm := dkg.SignedMessage{}
			_ = sm.Decode(msg.Data)

			pm := frost.ProtocolMsg{}
			_ = pm.Decode(sm.Message.Data)

			fmt.Printf("%+v\n", pm)
		}

		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}
