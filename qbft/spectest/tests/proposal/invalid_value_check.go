package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// InvalidValueCheck tests a proposal that doesn't pass value check
func InvalidValueCheck() *tests.MsgProcessingSpecTest {
	panic("implement")

	// TODO: implement passing invalid value

	//pre := testingutils.BaseInstance()
	//ks := testingutils.Testing4SharesSet()
	//msgs := []*qbft.SignedMessage{
	//	testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
	//		MsgType:    qbft.ProposalMsgType,
	//		Height:     qbft.FirstHeight,
	//		Round:      qbft.FirstRound,
	//		Identifier: []byte{1, 2, 3, 4},
	//		Root:       testingutils.ProposalDataBytes(testingutils.TestingInvalidValueCheck, nil, nil),
	//	}),
	//}
	//return &tests.MsgProcessingSpecTest{
	//	Name:           "invalid proposal value check",
	//	Pre:            pre,
	//	PostRoot:       "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
	//	InputMessages:  msgs,
	//	OutputMessages: []*qbft.SignedMessage{},
	//	ExpectedError:  "invalid signed message: proposal not justified: proposal value invalid: invalid value",
	//}
}
