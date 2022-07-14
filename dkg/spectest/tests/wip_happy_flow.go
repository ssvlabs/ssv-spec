package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WipHappyFlow tests a simple full happy flow until decided
func WipHappyFlow() *WipMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	suite := testutils.TestSuiteThreeOfFour()
	return &WipMsgProcessingSpecTest{
		Name:   "happy flow",
		KeySet: testingutils.Testing4SharesSet(),
		Output: &keygen.LocalKeyShare{
			Index:           1,
			Threshold:       2,
			ShareCount:      4,
			PublicKey:       suite.PublicKey,
			SecretShare:     suite.SecretShares[1],
			SharePublicKeys: suite.VkVec(),
		},
		Messages: []*keygen.ParsedMessage{
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)),
			testutils.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)),
			testutils.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)),
			testutils.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)),
			testutils.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)),
		},
	}
}
