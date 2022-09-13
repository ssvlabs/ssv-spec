package consensus

// TODO<olegshmuelov>: CONSENSUS fix test
// ValidDecided10Operators tests a valid decided value (10 operators)
//func ValidDecided10Operators() *tests.MultiMsgProcessingSpecTest {
//	ks := testingutils.Testing10SharesSet()
//	return &tests.MultiMsgProcessingSpecTest{
//		Name: "consensus valid decided 10 operators",
//		Tests: []*tests.MsgProcessingSpecTest{
//			{
//				Name:                    "sync committee contribution",
//				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
//				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
//				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusDataByts, ks, types.BNRoleSyncCommitteeContribution),
//				PostDutyRunnerStateRoot: "dbb07e30b4a7d7982d0234389106bbcc354f5ba9b152789cce30648151509c80",
//				OutputMessages: []*ssv.SignedPartialSignatureMessage{
//					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
//					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
//				},
//			},
//			{
//				Name:                    "sync committee",
//				Runner:                  testingutils.SyncCommitteeRunner(ks),
//				Duty:                    testingutils.TestingSyncCommitteeDuty,
//				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusDataByts, ks, types.BNRoleSyncCommittee),
//				PostDutyRunnerStateRoot: "dc07f90ba8026399107bfac6e9199ca5d5a321c20bc3331d015b44de860e0df0",
//				OutputMessages: []*ssv.SignedPartialSignatureMessage{
//					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
//				},
//			},
//			{
//				Name:                    "aggregator",
//				Runner:                  testingutils.AggregatorRunner(ks),
//				Duty:                    testingutils.TestingAggregatorDuty,
//				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusDataByts, ks, types.BNRoleAggregator),
//				PostDutyRunnerStateRoot: "5d1cc9a17435ec5c730d011d3400fcf428e03f2a76aecfc0ee776b3d35037bea",
//				OutputMessages: []*ssv.SignedPartialSignatureMessage{
//					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
//					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
//				},
//			},
//			{
//				Name:                    "proposer",
//				Runner:                  testingutils.ProposerRunner(ks),
//				Duty:                    testingutils.TestingProposerDuty,
//				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusDataByts, ks, types.BNRoleProposer),
//				PostDutyRunnerStateRoot: "509ba224a6ac57f606870bf36caca6164ffb2d9c061f071e3164b71e4fd6f7ec",
//				OutputMessages: []*ssv.SignedPartialSignatureMessage{
//					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
//					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
//				},
//			},
//			{
//				Name:                    "attester",
//				Runner:                  testingutils.AttesterRunner(ks),
//				Duty:                    testingutils.TestingAttesterDuty,
//				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusDataByts, ks, types.BNRoleAttester),
//				PostDutyRunnerStateRoot: "896ec3deedcc8ac7fb74345d9056b52783e30b120dc34b016dd7bccb1e0ac5ff",
//				OutputMessages: []*ssv.SignedPartialSignatureMessage{
//					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
//				},
//			},
//		},
//	}
//}
