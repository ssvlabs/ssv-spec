package dutyexe

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Correct ID for SSVMessage
	getID := func(role types.BeaconRole) types.MessageID {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret
	}
	// Wrong ID for SignedMessage
	getWrongID := func(role types.BeaconRole) []byte {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], role)
		return ret[:]
	}

	// Function to get decided message with wrong ID for role
	decidedMessage := func(role types.BeaconRole) *types.SSVMessage {
		signedMessage := testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
			testingutils.TestingDutySlot,
			getWrongID(role))

		btys, err := signedMessage.Encode()
		if err != nil {
			panic(err.Error())
		}

		return &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   getID(role),
			Data:    btys,
		}
	}

	expectedError := "failed processing consensus message: invalid msg: message doesn't belong to Identifier"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "wrong duty pubkey",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee contribution",
				Runner:   testingutils.SyncCommitteeContributionRunner(ks),
				Duty:     &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleSyncCommitteeContribution))},
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.SyncCommitteeRunner(ks),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleSyncCommittee))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  expectedError,
			},
			{
				Name:     "aggregator",
				Runner:   testingutils.AggregatorRunner(ks),
				Duty:     &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleAggregator))},
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:           "attester",
				Runner:         testingutils.AttesterRunner(ks),
				Duty:           &testingutils.TestingAttesterDuty,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleAttester))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
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
			Messages: []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleProposer))},
			OutputMessages: []*types.SignedPartialSignatureMessage{
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
			Messages: []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, decidedMessage(types.BNRoleProposer))},
			OutputMessages: []*types.SignedPartialSignatureMessage{
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
