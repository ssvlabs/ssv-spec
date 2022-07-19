package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func ThirteenOperators() *MsgProcessingSpecTest {
	ks := testingutils.Testing13SharesSet()
	suite := testutils.TestSuiteThirteenOperators()
	pre := testutils.ThirteenOperatorsInstance
	return &MsgProcessingSpecTest{
		Name:   "happy flow thirteen operators",
		Pre:    pre,
		KeySet: ks,
		Output: suite.MakeLocalKeyShare(1),
		Messages: []*keygen.ParsedMessage{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R1(5)),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R1(6)),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R1(7)),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R1(8)),
			testutils.SignDKGMsg(ks.DKGOperators[9].SK, suite.R1(9)),
			testutils.SignDKGMsg(ks.DKGOperators[10].SK, suite.R1(10)),
			testutils.SignDKGMsg(ks.DKGOperators[11].SK, suite.R1(11)),
			testutils.SignDKGMsg(ks.DKGOperators[12].SK, suite.R1(12)),
			testutils.SignDKGMsg(ks.DKGOperators[13].SK, suite.R1(13)),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R2(5)),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R2(6)),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R2(7)),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R2(8)),
			testutils.SignDKGMsg(ks.DKGOperators[9].SK, suite.R2(9)),
			testutils.SignDKGMsg(ks.DKGOperators[10].SK, suite.R2(10)),
			testutils.SignDKGMsg(ks.DKGOperators[11].SK, suite.R2(11)),
			testutils.SignDKGMsg(ks.DKGOperators[12].SK, suite.R2(12)),
			testutils.SignDKGMsg(ks.DKGOperators[13].SK, suite.R2(13)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R3(5, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R3(6, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R3(7, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R3(8, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[9].SK, suite.R3(9, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[10].SK, suite.R3(10, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[11].SK, suite.R3(11, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[12].SK, suite.R3(12, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[13].SK, suite.R3(13, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R4(5)),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R4(6)),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R4(7)),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R4(8)),
			testutils.SignDKGMsg(ks.DKGOperators[9].SK, suite.R4(9)),
			testutils.SignDKGMsg(ks.DKGOperators[10].SK, suite.R4(10)),
			testutils.SignDKGMsg(ks.DKGOperators[11].SK, suite.R4(11)),
			testutils.SignDKGMsg(ks.DKGOperators[12].SK, suite.R4(12)),
			testutils.SignDKGMsg(ks.DKGOperators[13].SK, suite.R4(13)),
		},
	}
}
