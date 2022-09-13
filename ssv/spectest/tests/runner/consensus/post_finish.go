package consensus

// TODO<olegshmuelov>: CONSENSUS fix test
// PostFinish tests a valid commit msg after runner finished
/*func PostFinish() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GenerateConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances[0] = r.GetBaseRunner().State.RunningInstance
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "8d63f41e234c32466ac964be7f844016726b6658d27ff911fe327fcdaa33953f",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "sync committee",
				Runner: finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "6256487c77b59d2ade0b43a46ad89de4e94514d5a10d5c9d0b5d9721207e3332",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "ebf29b5d39134ddbd9ae8c1e4f21ed048d51caf73ff5bb5483c47f1686a6b17d",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "aec61000fa5edbfa6cfe6cfa5fd67572ded1fe933391cd59e80b6d486f60b86c",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
				DontStartDuty: true,
				ExpectedError: "failed processing consensus message: invalid consensus message: no running duty",
			},
			{
				Name:   "attester",
				Runner: finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.Message{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "1b16717b90bdfe9590fa5511717e6d1c3cda64bb5ed921d096ec1aa2987010ea",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: invalid consensus message: no running duty",
			},
		},
	}
}*/
