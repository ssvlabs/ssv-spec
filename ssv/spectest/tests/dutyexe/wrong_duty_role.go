package dutyexe

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyRole tests decided value duty with wrong duty role (!= duty runner role)
func WrongDutyRole() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Correct ID for SSVMessage
	getID := func(role types.RunnerRole) types.MessageID {
		if role == types.RoleAggregatorCommittee {
			return testingutils.TestingAggregatorCommitteeMsgID
		}
		ret := types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), role)
		return ret
	}
	// Wrong ID for SignedMessage
	getWrongID := func(role types.RunnerRole) []byte {
		if role == types.RoleAggregatorCommittee {
			committee := make([]uint64, 0)
			for _, op := range ks.Committee() {
				committee = append(committee, op.Signer)
			}
			committeeID := types.GetCommitteeID(committee)
			ret := types.NewCommitteeMsgID(testingutils.TestingSSVDomainType, committeeID, types.RoleAggregatorCommittee+1)
			return ret[:]
		}
		ret := types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), role+1)
		return ret[:]
	}

	// Function to get decided message with wrong ID for role
	decidedMessage := func(role types.RunnerRole, slot phase0.Slot) *types.SignedSSVMessage {
		signedMessage := testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{1, 2, 3},
			qbft.Height(slot),
			//testingutils.TestingDutySlot,
			getWrongID(role))

		signedMessage.SSVMessage.MsgID = getID(role)

		sig1 := testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], signedMessage.SSVMessage).Signatures[0]
		sig2 := testingutils.SignedSSVMessageWithSigner(2, ks.OperatorKeys[2], signedMessage.SSVMessage).Signatures[0]
		sig3 := testingutils.SignedSSVMessageWithSigner(3, ks.OperatorKeys[3], signedMessage.SSVMessage).Signatures[0]

		signedMessage.Signatures = [][]byte{sig1, sig2, sig3}

		return signedMessage
	}

	expectedErrorCode := types.MessageIdentifierInvalidErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"wrong duty role",
		testdoc.DutyExeWrongDutyRoleDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee contribution",
				Runner:   testingutils.AggregatorCommitteeRunner(ks),
				Duty:     testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleAggregatorCommittee, testingutils.TestingSyncCommitteeContributionDuty.Slot)},
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedErrorCode: expectedErrorCode,
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("aggregator (%s)", version.String()),
			Runner:   testingutils.AggregatorCommitteeRunner(ks),
			Duty:     testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleAggregatorCommittee, testingutils.TestingDutySlotV(version))},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
			},
			ExpectedErrorCode: expectedErrorCode,
		},
		)
	}

	for _, version := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:     fmt.Sprintf("proposer (%s)", version.String()),
			Runner:   testingutils.ProposerRunner(ks),
			Duty:     testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{decidedMessage(types.RoleProposer, testingutils.TestingDutySlotV(version))},
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
			},
			ExpectedErrorCode: expectedErrorCode,
		})
	}

	return multiSpecTest
}
