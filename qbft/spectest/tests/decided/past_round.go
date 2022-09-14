package decided

// TODO<olegshmuelov>: DECIDED fix test
// PastRound tests a decided msg for a past round
/*func PastRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 100

	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      5,
				Identifier: []byte{1, 2, 3, 4},
				Input: []byte{1, 2, 3, 4},
			}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "decided past round",
		Pre:            pre,
		PostRoot:       "d2978f3827f69c42778e9ac8d9676b992e3839cc5ed6d527b28b0d36889c4de2",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}*/
