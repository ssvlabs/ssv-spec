package runnertests

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name          string
	Pre           *dkg.Runner
	Messages      []*types.Message
	Outgoing      []*types.Message
	Output        *types.SignedDepositDataMsgBody
	KeySet        *testingutils.TestKeySet
	ExpectedError string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	pre := test.Pre

	var (
		lastErr, err error
		finished     bool
		output       *types.SignedDepositDataMsgBody
	)

	err = pre.Start()

	if err != nil {
		lastErr = err
	}
	for _, msg := range test.Messages {
		finished, _, err = pre.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
		if finished {
			break
		}
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
		//outgoing := test.Pre.Config.Network.(*testutils.MockNetwork).Broadcasted
		// TODO: Compare outgoing messages
		o, _ := test.Pre.SignSubProtocol.Output()
		output = &types.SignedDepositDataMsgBody{}
		output.Decode(o)
		require.True(t, proto.Equal(test.Output, output))
	}
}
