package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSSVMessageType tests an invalid SignedSSVMessageTest whose SSVMessage has an unknown MsgType
func InvalidSSVMessageType() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	ssvMsg := testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	ssvMsg.MsgType = types.MsgType(99)

	return NewSignedSSVMessageTest(
		"invalid ssvmessage type",
		testdoc.SignedSSVMessageTestInvalidSSVMessageTypeDoc,
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage:  ssvMsg,
			},
		},
		types.MessageTypeInvalidErrorCode,
		nil,
	)
}
