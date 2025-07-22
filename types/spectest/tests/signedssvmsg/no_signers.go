package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoSigners tests an invalid SignedSSVMessageTest with no signers
func NoSigners() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return NewSignedSSVMessageTest(
		"no signers",
		testdoc.SignedSSVMessageTestNoSignersDoc,
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{},
				Signatures:  [][]byte{{1, 2, 3, 4}},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		"no signers",
		nil,
	)
}
