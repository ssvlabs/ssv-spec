package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SigValid tests SignedPostConsensusMessage sig == 96 bytes
func SigValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, gloas.DataVersionGloas)

	return NewMsgSpecTest(
		"sig valid",
		testdoc.MsgSpecTestSigValidDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		0,
	)
}
