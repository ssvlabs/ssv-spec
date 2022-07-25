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
	LocalKeyShare *dkgtypes.LocalKeyShare
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
	protocol := func(init *dkgtypes.Init, operatorID types.OperatorID, identifier dkgtypes.RequestID) dkgtypes.Protocol {
		return testutils.MockProtocol{LocalKeyShare: test.LocalKeyShare}
	}
	network := testutils.NewMockNetwork()
	ks := testingutils.Testing13SharesSet()
	config := &dkgtypes.Config{
		Protocol:            protocol,
		BeaconNetwork:       types.PraterNetwork,
		Network:             network,
		Storage:             testutils.NewMockStorage(*ks),
		SignatureDomainType: types.PrimusTestnet,
		Signer: &testutils.MockSigner{
			SK:            ks.DKGOperators[1].SK,
			ETHAddress:    ks.DKGOperators[1].ETHAddress,
			EncryptionKey: ks.DKGOperators[1].EncryptionKey,
		},
	}
	node := dkg.NewNode(test.Operator, config)

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
	}
	outgoing := network.Broadcasted
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
