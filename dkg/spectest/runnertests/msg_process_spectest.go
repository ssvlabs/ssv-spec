package runnertests

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name          string
	Operator      *dkgtypes.Operator
	Config        *dkgtypes.Config
	Messages      []*dkgtypes.Message
	Outgoing      []*dkgtypes.Message
	Output        map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage
	KeySet        *testingutils.TestKeySet
	ExpectedError string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	node := dkg.NewNode(test.Operator, test.Config)

	var (
		lastErr, err error
		output       map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage
	)

	if err != nil {
		lastErr = err
	}
	for _, msg := range test.Messages {
		node.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
		outgoing := test.Config.Network.(*testutils.MockNetwork).Broadcasted
		require.Equal(t, len(test.Outgoing), len(outgoing))
		for i, message := range outgoing {
			message.Signature = nil // Signature is not deterministic, so skip
			require.True(t, proto.Equal(outgoing[i], message))
		}
		output = test.Output
		require.Equal(t, len(test.Output), len(output))
		for id, message := range test.Output {
			require.True(t, proto.Equal(message, output[id]))
		}
	}
}
