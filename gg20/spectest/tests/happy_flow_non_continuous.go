package tests

import (
	dkgtu "github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/gg20/testutils"
	"github.com/bloxapp/ssv-spec/gg20/types"
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
		Messages: []*types.ParsedMessage{
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R1(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R1(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R1(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[8].SK, suite.R1(8)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R2(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R2(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R2(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[8].SK, suite.R2(8)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R3(2, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R3(4, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[8].SK, suite.R3(8, 1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[1].SK, suite.R4(1)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[2].SK, suite.R4(2)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[4].SK, suite.R4(4)).(*types.ParsedMessage),
			dkgtu.SignDKGMsg(ks.DKGOperators[8].SK, suite.R4(8)).(*types.ParsedMessage),
		},
	}
}
