package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages = []*types.PartialSignatureMessage{}

	return &MsgSpecTest{
		Name:          "no messages",
		Messages:      []*types.SignedPartialSignatureMessage{msg},
		ExpectedError: "no PartialSignatureMessages messages",
	}
}
