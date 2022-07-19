package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func SevenOperators() *MsgProcessingSpecTest {
	ks := testingutils.Testing7SharesSet()
	suite := testutils.TestSuiteSevenOperators()
	pre := testutils.SevenOperatorsInstance
	return &MsgProcessingSpecTest{
		Name:   "happy flow seven operators",
		Pre:    pre,
		KeySet: ks,
		Output: suite.MakeLocalKeyShare(1),
		Messages: []*keygen.ParsedMessage{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R1(5)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R1(6)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R1(7)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R2(5)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R2(6)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R2(7)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R3(5, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R3(6, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R3(7, 1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R4(5)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R4(6)).(*keygen.ParsedMessage),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R4(7)).(*keygen.ParsedMessage),
		},
	}
}
