package tests

import (
	dkgtu "github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/gg20/testutils"
	"github.com/bloxapp/ssv-spec/gg20/types"
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
		Messages: []*types.ParsedMessage{
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*types.ParsedMessage),
		},
	}
}
