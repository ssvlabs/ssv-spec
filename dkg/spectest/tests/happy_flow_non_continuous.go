package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func HappyFlowNonContinuous() *MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSetNonContinuous()
	suite := testutils.TestSuiteFourOperatorsNonContinuous()
	pre := testutils.BaseInstanceNonContinuous
	return &MsgProcessingSpecTest{
		Name:   "happy flow non-continuous",
		Pre:    pre,
		KeySet: ks,
		Output: suite.MakeLocalKeyShare(1),
		Messages: []*keygen.ParsedMessage{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R1(8)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R2(8)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R3(8, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R4(8)).(*keygen.ParsedMessage),
		},
	}
}
