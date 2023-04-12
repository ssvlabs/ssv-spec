package runner

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ssz "github.com/ferranbt/fastssz"
)

func getSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

func hexDecodeNoErr(h string) []byte {
	ret, err := hex.DecodeString(h)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func runnerHappyFlowStates(ks *testingutils.TestKeySet) []ssv.Runner {
	syncCommitteeContrib := testingutils.SyncCommitteeContributionRunner(ks)
	syncCommitteeContrib.GetBaseRunner().State = &ssv.State{}
	syncCommitteeContrib.GetBaseRunner().State.Finished = true
	syncCommitteeContrib.GetBaseRunner().State.DecidedValue = testingutils.TestSyncCommitteeContributionConsensusData
	syncCommitteeContrib.GetBaseRunner().State.StartingDuty = &testingutils.TestSyncCommitteeContributionConsensusData.Duty
	syncCommitteeContrib.GetBaseRunner().State.PreConsensusContainer = &ssv.PartialSigContainer{
		Quorum: 3,
		Signatures: map[string]map[types.OperatorID][]byte{
			"9989d3ab6c75aa22aef5d56898a930c3f67de00e241103071ecf84523c73fc1c": {
				1: hexDecodeNoErr("8ea233eb70fc9b3da15f2cc8a6859c1b0836b7d96453c606914e4b71ce6802fc1e2d2e7da5d32df8c2231652fc01158a194b95e44da3e26dbb737516f25aca7c3ee713721274fddbc829c3c526bfb3838682f6bb8a87f83b9be5e94e790e42ab"),
				2: hexDecodeNoErr("a4a926a9a3f4426c7b2c13d9035fdd6058a46fc1d8fc33cf91d96e97a27162330697c516ddae70ea64a14ddfaaa5a6920610fce020a72e1e32a6b883b983aeb6ce45da5eaad54987fb5b67487d7a3be12192663a4f27513bf53ce3c20bf82005"),
				3: hexDecodeNoErr("8973abad38fbb8560831d7770e461433a1a0157dbdf124f44f715d7d3317521ae38727fc601ffda793aa6995041b9b45041f48dd97b1d62401d2a462deb7182faba0a52585b60a053d57ee83b58808666b54424ad9bd5654b3f4257d681571cd"),
			},
			"a96747766942745b4a6447e3f414fcba5fede9727ee021a9576c9a670a986e53": {
				1: hexDecodeNoErr("831b8437cab11389d19df9adb77bf8afc33d9f15b2da91533b1f265c6ad3c4ae440d2b964344a45211acd6e3f5361bc40fe397581cede2c03f9422a17d41482cb41872762f1c0fb754c395cd3845a50427e8fa0ede4b8c40912af57cda8d2e67"),
				2: hexDecodeNoErr("8fa6354e03c6099934295e575bec9a1876be0cff3349b9daffd79d9e6e7ca5b9972646a44212336e6ffd083e1ae23c52148021e6235f2b892f625e210e2b55d64bb1a64091f1d947a47d0bd81d89949ac7dab8f52ea3ae5ad346f211b8a1005a"),
				3: hexDecodeNoErr("992ea601783f18e70428645e0acae09aec640ffb0646d1be4253ab7ddf92792b23506fa741a99300d89a6cb280f82de9145aaab377fa5827d3a95fb7a80fd8f87efd47ccd7004b33c30b2303e517c62e95266ef5400cf650eef762e62dcf6ebb"),
			},
			"f0f373fcb0097430b0915b34dc000a1ff0a2e993bc1c4bf90f0086749f3d269f": {
				1: hexDecodeNoErr("b42736ed260adcb7efec7b89763b8e25bb59239c26eddf8e82aab85b9ba54a0689951dbeb03c79a129dbfd9b36e0615e198861405c5e06f2c53dcfd1fe8ca35799452c06368feb73a4c5c3d2d1a37534bf734656da4d8dea319438a6c6120c5b"),
				2: hexDecodeNoErr("a37752b5cc25faa4f852eef362b918e7f31bbe3684b5c30299667829b53989c90791c66fe0f53e32a63c56f6f5d60b0a0b25bb9f48d440f8216e0a170f7ee44b5afdd22b5d8602186fd312e3f55dc05282976182b462a418001718ea1c083e79"),
				3: hexDecodeNoErr("97c2c77a73349420231c0681783020fd1086471ee39613ad96a92d8fb438b4f009cf56b0074f1e62fa32c354fec4b6c908bd8dcdcef5897664bff8516f12673ba2ce9eb28afe5128ede0f93fad242622ab5ccd9eaf2d173823c683116ef69d73"),
			},
		}}
	syncCommitteeContrib.GetBaseRunner().State.PostConsensusContainer = &ssv.PartialSigContainer{
		Quorum: 3,
		Signatures: map[string]map[types.OperatorID][]byte{
			"9f4156be6449811e4e6dfc7dc6c3f0adc0b3d8d4cb1cb1aa15925631ef77ea59": {
				1: hexDecodeNoErr("81ce4d4016e1a29a6fe0e7622a5ce880c70e39598921903b98de9026de74c8a7216dbb350a16a82bcfe0cd50d1075232168f32bd811f9f3747c7c806e1511d24465eb6d540644c25b872c70d39d6901b0e8f6de88b6590d10ae953021c56286b"),
				2: hexDecodeNoErr("a9ddfe47a9bb08fcd3c040dc3386b9a3dfd83b3c2b7d8b0cd36937285c92081738810afa6676dbfb4e593cdfa8335e0e0fad4341306b82be1245341e570dfe8a24727f761e1799a44fd79bf85cc2203b9bc15cab5111eb9c45a9419fe48ff31b"),
				3: hexDecodeNoErr("91a0d52c73177b4c6b04f302faa131f4332ff2913fc9de3127cbc7fb634b05f28a06c8a6241e92ab80da6bc9d030026108a9b81441137f89e2d5b0515cd3856221c87c9ccf43cf1fee917c2f706c53718e82049051705eb9cfb2871683f8dbc6"),
			},
			"a735f2e95e050090b43cf6c67c10645104eaf42285f9443ed5b9961d208e041f": {
				1: hexDecodeNoErr("81ea06365186cc9dbfda3dd85f598bc5ded84f31d4ff658810f9efbf02937f8b187bc806dcaa894fcafedd4a61fe629500d8faf4a71619eb6b726c02de06660e70bc24c71784d287b286a8d4b94d5b5af3994f5bdc53969156f4549417462b61"),
				2: hexDecodeNoErr("8668f29fe6bfa137fa32224678198394bec81e9e46469c2f09ab0c3b43cb62d0aa464f4a730f86e79d4be29454aee45d0df54577c27c0f4f394f19cda00c4c2f761c19306b044896d96bde48474300857bda984ac1b109f849786e20361e0dad"),
				3: hexDecodeNoErr("83fb9fe8964083189f9f41c0902de782967b8a7181bf14e21fc4c5c3770425d572a4cfaeee2280da62606c990caeddf3148cf65bb84ad1a08781c9e2535547509fedc3d629e348e98c0fbd2db14146bde41a1ad322ed0529f5f2d8e92a047b1b"),
			},
			"f48022de4342f2bea21fce105bae2e7563f7f1f11b709fb0a0ffbf252a96390c": {
				1: hexDecodeNoErr("85237d03d7cdd08ac5f364a79f14f582d7029d2d51dd27c7ec61b7c8ccd4f058c3d7540c410de4a155818fcafe329de00b84abc63d65c5b24c72064ea368314af7d8c3e6d516e8e09802b112c1c2fcd82cf63a9577c8e65374aa33f912d12a9a"),
				2: hexDecodeNoErr("89242fd0d79dfc15eb641db1b43468bc36364b973aa37eee056e015a5eac3dea8d055233cc43196b9e4f9d2aa190ec4900c43a30210701af84e0e2aff3221780e962fe75668d04d6922f7725495da00058e048555d7d364e493df403cb71eaff"),
				3: hexDecodeNoErr("90a6567a7d2eece80afe7f42c3cde6c519f5b5a61870f9c881aec5165314502b01b925d265f32bb3722e6999cc1248b91352eb647fa3e669f00b648a531e7b572a5eecac70b0d8ba620299fa11545a6794fbb2134885d34d7fb6361412df7449"),
			},
		}}
	syncCommitteeContrib.GetBaseRunner().State.RunningInstance = &qbft.Instance{
		StartValue: testingutils.TestSyncCommitteeContributionConsensusDataByts,
		State: &qbft.State{
			Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:     syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
			Round:  qbft.FirstRound,
			Height: qbft.FirstHeight,
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
				ks.Shares[1], types.OperatorID(1),
				syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
				testingutils.TestSyncCommitteeContributionConsensusDataByts,
			),
			LastPreparedRound: 1,
			LastPreparedValue: testingutils.TestSyncCommitteeContributionConsensusDataByts,
			Decided:           true,
			DecidedValue:      testingutils.TestSyncCommitteeContributionConsensusDataByts,
			ProposeContainer: &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
				qbft.FirstRound: {
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataByts,
					),
				},
			}},
			PrepareContainer: &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
				qbft.FirstRound: {
					testingutils.TestingPrepareMessageWithIdentifierAndRoot(
						ks.Shares[1], types.OperatorID(1),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
					testingutils.TestingPrepareMessageWithIdentifierAndRoot(
						ks.Shares[2], types.OperatorID(2),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
					testingutils.TestingPrepareMessageWithIdentifierAndRoot(
						ks.Shares[3], types.OperatorID(3),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
				},
			}},
			CommitContainer: &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
				qbft.FirstRound: {
					testingutils.TestingCommitMessageWithIdentifierAndRoot(
						ks.Shares[1], types.OperatorID(1),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
					testingutils.TestingCommitMessageWithIdentifierAndRoot(
						ks.Shares[2], types.OperatorID(2),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
					testingutils.TestingCommitMessageWithIdentifierAndRoot(
						ks.Shares[3], types.OperatorID(3),
						syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
						testingutils.TestSyncCommitteeContributionConsensusDataRoot),
				},
			}},
			RoundChangeContainer: &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{}},
		},
	}
	syncCommitteeContrib.GetBaseRunner().QBFTController.StoredInstances = append(syncCommitteeContrib.GetBaseRunner().QBFTController.StoredInstances, syncCommitteeContrib.GetBaseRunner().State.RunningInstance)

	syncCommittee := testingutils.SyncCommitteeRunner(ks)

	aggregator := testingutils.AggregatorRunner(ks)

	proposer := testingutils.ProposerRunner(ks)

	blindedProposer := testingutils.ProposerBlindedBlockRunner(ks)

	attester := testingutils.AttesterRunner(ks)

	validatorRegistration := testingutils.ValidatorRegistrationRunner(ks)

	return []ssv.Runner{
		syncCommitteeContrib,
		syncCommittee,
		aggregator,
		proposer,
		blindedProposer,
		attester,
		validatorRegistration,
	}
}

// FullHappyFlow  tests a full runner happy flow
func FullHappyFlow() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// register runners
	runnerStates := runnerHappyFlowStates(ks)
	roots := make([]string, 0)
	for _, runner := range runnerStates {
		r, err := runner.GetRoot()
		if err != nil {
			panic(err.Error())
		}
		roots = append(roots, hex.EncodeToString(r[:]))
		tests.RootRegister[hex.EncodeToString(r[:])] = runner
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					[]*types.SSVMessage{ // pre consensus
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
					append(
						// consensus
						testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
						[]*types.SSVMessage{ // post consensus
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
						}...,
					)...,
				),
				PostDutyRunnerStateRoot: roots[0], //"4987127ad389bb9d21500d447686f135a19f59ae10192e82bf052278853ad3d1",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "48c73f57659b69131467ef133ccb35d7de2fe96438d30bfa2b5ea63b19ead011",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					[]*types.SSVMessage{ // pre consensus
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
					append(
						testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator), // consensus
						[]*types.SSVMessage{ // post consensus
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
						}...,
					)...,
				),
				PostDutyRunnerStateRoot: "298bbb63d87a36eef30926c2c21baad6990db0f8fa03a83ca56b2463c7f0065c",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusData, ks, types.BNRoleProposer), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "76812c0f14ff09067547e9528730749b0c0090d1a4872689a0b8480d7b538884",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "proposer blinded block",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedBlockConsensusData, ks, types.BNRoleProposer), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "90755cc41b814519fd9fdd14bc82d239997ba51340c297f25f5f1552f27f66c7",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					}...,
				),
				PostDutyRunnerStateRoot: "9d55ff5721b21c5b99dd4b4bacb0acda0b674112fe3cec55cc6aeb04ad5dc2fc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
				},
				PostDutyRunnerStateRoot: "f36c8b537afaba0894dbc8c87cb94466d8ac2623e9283f1c584e3d544b5f2b88",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
