package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSSVMessageRole tests an invalid SignedSSVMessageTest whose SSVMessage MsgID encodes a negative role
func InvalidSSVMessageRole() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	ssvMsg := testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	ssvMsg.MsgID = types.NewCommitteeMsgID(testingutils.TestingSSVDomainType, testingutils.TestingCommitteeMember(ks).CommitteeID, types.RoleUnknown)

	return NewSignedSSVMessageTest(
		"invalid ssvmessage role",
		testdoc.SignedSSVMessageTestInvalidSSVMessageRoleDoc,
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage:  ssvMsg,
			},
		},
		types.SSVMessageInvalidRoleErrorCode,
		nil,
	)
}
