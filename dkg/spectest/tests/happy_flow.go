package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	suite := testutils.TestSuiteFourOperators()
	pre := testutils.BaseInstance
	return &MsgProcessingSpecTest{
		Name:   "happy flow",
		Pre:    pre,
		KeySet: ks,
		Output: suite.MakeLocalKeyShare(1),
		Messages: []*keygen.ParsedMessage{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*keygen.ParsedMessage),
		},
	}
}
