package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Messages = []*types.PartialSignatureMessage{}

	return &MsgSpecTest{
		Name:          "no messages",
		Messages:      []*types.PartialSignatureMessages{msg},
		ExpectedError: "no PartialSignatureMessages messages",
	}
}
