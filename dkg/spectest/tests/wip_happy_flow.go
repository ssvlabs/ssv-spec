package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WipHappyFlow tests a simple full happy flow until decided
func WipHappyFlow() *WipMsgProcessingSpecTest {
	suite := testutils.TestSuiteThreeOfFour()
	return &WipMsgProcessingSpecTest{
		Name:   "happy flow",
		KeySet: testingutils.Testing4SharesSet(),
		Output: &suite[0].LocalKeyShare,
		Messages: []*keygen.ParsedMessage{
			&suite[0].R1Message,
			&suite[1].R1Message,
			&suite[2].R1Message,
			&suite[3].R1Message,
			&suite[0].R2Message,
			&suite[1].R2Message,
			&suite[2].R2Message,
			&suite[3].R2Message,
			&suite[1].R3Messages[0],
			&suite[2].R3Messages[0],
			&suite[3].R3Messages[0],
			&suite[0].R4Message,
			&suite[1].R4Message,
			&suite[2].R4Message,
			&suite[3].R4Message,
		},
	}
}
