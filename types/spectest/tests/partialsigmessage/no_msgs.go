package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)
	msg.Message.Messages = []*types.PartialSignatureMessage{}

	return &MsgSpecTest{
		Name:          "no messages",
		Messages:      []*types.SignedPartialSignatureMessage{msg},
		ExpectedError: "no PartialSignatureMessages messages",
	}
}
