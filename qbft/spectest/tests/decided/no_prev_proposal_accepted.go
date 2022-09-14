package decided

// TODO<olegshmuelov>: DECIDED fix test
// NoPrevAcceptedProposal tests a commit msg received without a previous accepted proposal
/*func NoPrevAcceptedProposal() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = nil
	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Input: []byte{1, 2, 3, 4},
			}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "decided no previous accepted proposal",
		Pre:           pre,
		PostRoot:      "ed99ab91cac917c5bf9ff90eee30f21fe47d2e272d1f35d005dbdffef426ac02",
		InputMessages: msgs,
	}
}*/
