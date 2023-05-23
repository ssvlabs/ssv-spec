package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid partial pre consensus msg before duty starts
func NoRunningDuty() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "a1d866a9c085e355f93bd365fa67b7321ce780dc314bccc9b97f077a20988ae7",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing sync committee selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "f947dcc48dbae821f6a1bf1125d02a1bf7665704670dfca1bdb13300f130b308",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "194ee4fd2a60a8c1f4d4d51f4445e4176e199ac075d06614a121c0d13fc63e0b",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing validator registration message: invalid pre-consensus message: no running duty",
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
			},
			PostDutyRunnerStateRoot: noRunningDutyProposerSC(version).Root(),
			PostDutyRunnerState:     noRunningDutyProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			DontStartDuty:           true,
			ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SSVMessage{
				testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version)),
			},
			PostDutyRunnerStateRoot: noRunningDutyBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     noRunningDutyBlindedProposerSC(version).ExpectedState,
			OutputMessages:          []*types.SignedPartialSignatureMessage{},
			DontStartDuty:           true,
			ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
