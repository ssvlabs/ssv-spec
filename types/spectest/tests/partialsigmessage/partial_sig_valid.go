package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigValid tests PostConsensusMessage sig == 96 bytes
func PartialSigValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, gloas.DataVersionGloas)

	return NewMsgSpecTest(
		"partial sig valid",
		testdoc.MsgSpecTestPartialSigValidDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		0,
	)
}
