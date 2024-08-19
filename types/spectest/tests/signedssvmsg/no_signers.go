package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoSigners tests an invalid SignedSSVMessageTest with no signers
func NoSigners() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "no signers",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{},
				Signatures:  [][]byte{{1, 2, 3, 4}},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "no signers",
	}
}
