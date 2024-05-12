package dutyexe

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Correct ID for SSVMessage
	getID := func(role types.RunnerRole) types.MessageID {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret
	}
	// Wrong ID for SignedMessage
	getWrongID := func(role types.RunnerRole) []byte {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], role)
		return ret[:]
	}

	// Function to get decided message with wrong ID for role
	decidedMessage := func(role types.RunnerRole) *types.SignedSSVMessage {
		signedMessage := testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{1, 2, 3},
			testingutils.TestingDutySlot,
			getWrongID(role))

		signedMessage.SSVMessage.MsgID = getID(role)

		sig1 := testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], signedMessage.SSVMessage).Signatures[0]
		sig2 := testingutils.SignedSSVMessageWithSigner(2, ks.OperatorKeys[2], signedMessage.SSVMessage).Signatures[0]
		sig3 := testingutils.SignedSSVMessageWithSigner(3, ks.OperatorKeys[3], signedMessage.SSVMessage).Signatures[0]

		signedMessage.Signatures = [][]byte{sig1, sig2, sig3}

		return signedMessage
	}

	expectedError := "failed processing consensus message: invalid msg: message doesn't belong to Identifier"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "wrong duty pubkey",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee contribution",
				Runner:   testingutils.SyncCommitteeContributionRunner(ks),
				Duty:     &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleSyncCommitteeContribution)},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
			{
				Name:     "aggregator",
				Runner:   testingutils.AggregatorRunner(ks),
				Duty:     &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleAggregator)},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedError,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer (%s)", version.String()),
			Runner:   testingutils.ProposerRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleProposer)},
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
			Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleProposer)},
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
