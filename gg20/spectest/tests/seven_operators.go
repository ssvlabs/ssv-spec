package tests

import (
	dkgtu "github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/gg20/testutils"
	"github.com/bloxapp/ssv-spec/gg20/types"
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
		Messages: []*types.ParsedMessage{
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R1(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[5].SK, suite.R1(5)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[6].SK, suite.R1(6)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[7].SK, suite.R1(7)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R2(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[5].SK, suite.R2(5)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[6].SK, suite.R2(6)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[7].SK, suite.R2(7)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R3(3, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[5].SK, suite.R3(5, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[6].SK, suite.R3(6, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[7].SK, suite.R3(7, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[3].SK, suite.R4(3)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[5].SK, suite.R4(5)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[6].SK, suite.R4(6)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[7].SK, suite.R4(7)).(*types.ParsedMessage),
		},
	}
}
