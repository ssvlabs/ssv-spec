package dutyexe

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Wrong ID for SignedMessage
	getWrongID := func(role types.BeaconRole) []byte {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], role)
		return ret[:]
	}

	// Function to get decided message with wrong ID for role
	decidedMessage := func(role types.BeaconRole) *types.SignedSSVMessage {
		return testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
			[]types.OperatorID{1, 2, 3},
			testingutils.TestingDutySlot,
			getWrongID(role))
	}

	expectedError := "failed processing consensus message: invalid msg: message doesn't belong to Identifier"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "wrong duty pubkey",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee contribution",
				Runner:   testingutils.SyncCommitteeContributionRunner(ks),
				Duty:     &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{decidedMessage(types.BNRoleSyncCommitteeContribution)},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.SyncCommitteeRunner(ks),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				Messages:       []*types.SignedSSVMessage{decidedMessage(types.BNRoleSyncCommittee)},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:     "aggregator",
				Runner:   testingutils.AggregatorRunner(ks),
				Duty:     &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{decidedMessage(types.BNRoleAggregator)},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:           "attester",
				Runner:         testingutils.AttesterRunner(ks),
				Duty:           &testingutils.TestingAttesterDuty,
				Messages:       []*types.SignedSSVMessage{decidedMessage(types.BNRoleAttester)},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer (%s)", version.String()),
			Runner:   testingutils.ProposerRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{decidedMessage(types.BNRoleProposer)},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			ExpectedError: expectedError,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:   testingutils.ProposerBlindedBlockRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{decidedMessage(types.BNRoleProposer)},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			ExpectedError: expectedError,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
