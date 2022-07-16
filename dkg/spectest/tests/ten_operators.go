package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func TenOperators() *MsgProcessingSpecTest {
	ks := testingutils.Testing10SharesSet()
	suite := testutils.TestSuiteTenOperators()
	pre := testutils.TenOperatorsInstance
	return &MsgProcessingSpecTest{
		Name:   "happy flow",
		Pre:    pre,
		KeySet: ks,
		Output: &keygen.LocalKeyShare{
			Index:           1,
			Threshold:       6,
			ShareCount:      10,
			PublicKey:       suite.PublicKey,
			SecretShare:     suite.SecretShares[1],
			SharePublicKeys: suite.VkVec(),
		},
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
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[5].SK, suite.R3(5, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[6].SK, suite.R3(6, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[7].SK, suite.R3(7, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[8].SK, suite.R3(8, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[9].SK, suite.R3(9, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[10].SK, suite.R3(10, 1)),
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
		},
	}
}
