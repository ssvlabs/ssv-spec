package prepare

// TODO<olegshmuelov>: irrelevant test
// ImparsableProposalData tests a prepare msg received with imparsable data
/*func ImparsableProposalData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})

	prepareMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "imparsable prepare data",
		Pre:      pre,
		PostRoot: "be41977d818071451988105377df7c5ccf89ecc05ddf033b7b3b83d89f52d530",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareMsgEncoded,
			},
		},
		ExpectedError: "invalid prepare msg: could not get prepare data: could not decode prepare data from message: invalid character '\\x01' looking for beginning of value",
	}
}*/
