package consensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureDecidedNoInstance tests processing a decided msg from a larger height with no running instance
// then returning an error and don't move to post consensus as it's not the same instance decided
func FutureDecidedNoInstance() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	getID := func(role types.RunnerRole) []byte {
		if role == types.RoleCommittee {
			opIDs := make([]types.OperatorID, len(ks.Committee()))
			for i, member := range ks.Committee() {
				opIDs[i] = member.Signer
			}
			committeeID := types.GetCommitteeID(opIDs)
			ret := types.NewMsgID(testingutils.TestingSSVDomainType, committeeID[:], role)
			return ret[:]
		}
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], role)
		return ret[:]
	}

	getDecidedMessage := func(role types.BeaconRole, height qbft.Height) *types.SSVMessage {
		signedMsg := testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
			height,
			getID(role),
		)

		byts, err := signedMsg.Encode()
		if err != nil {
			panic(err.Error())
		}

		ssvMessage := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   *(*types.MessageID)(getID(role)),
			Data:    byts,
		}

		return ssvMessage
	}

	expectedErr := "no runner found for message's slot"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future decided no running instance",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:           "attester",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingAttesterDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{getDecidedMessage(types.RoleCommittee, testingutils.TestingDutySlot+1)},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingSyncCommitteeDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{getDecidedMessage(types.RoleCommittee, testingutils.TestingDutySlot+1)},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:           "attester and sync committee",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingAttesterAndSyncCommitteeDuties,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{getDecidedMessage(types.RoleCommittee, testingutils.TestingDutySlot+1)},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:           "sync committee contribution",
				Runner:         testingutils.SyncCommitteeContributionRunner(ks),
				Duty:           &testingutils.TestingSyncCommitteeContributionDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleSyncCommitteeContribution, testingutils.TestingDutySlot+1))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:           "sync committee",
				Runner:         testingutils.SyncCommitteeRunner(ks),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleSyncCommittee, testingutils.TestingDutySlot+1))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:           "aggregator",
				Runner:         testingutils.AggregatorRunner(ks),
				Duty:           &testingutils.TestingAggregatorDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleAggregator, testingutils.TestingDutySlot+1))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:           "attester",
				Runner:         testingutils.AttesterRunner(ks),
				Duty:           &testingutils.TestingAttesterDuty,
				DontStartDuty:  true,
				Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleProposer, testingutils.TestingDutySlot+1))},
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:           fmt.Sprintf("proposer (%s)", version.String()),
			Runner:         testingutils.ProposerRunner(ks),
			Duty:           testingutils.TestingProposerDutyV(version),
			DontStartDuty:  true,
			Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleProposer, qbft.Height(testingutils.TestingDutySlotV(version))+1))},
			OutputMessages: []*types.SignedPartialSignatureMessage{},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:           fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:         testingutils.ProposerBlindedBlockRunner(ks),
			Duty:           testingutils.TestingProposerDutyV(version),
			DontStartDuty:  true,
			Messages:       []*types.SignedSSVMessage{testingutils.SignedSSVMessageF(ks, getDecidedMessage(types.BNRoleProposer, qbft.Height(testingutils.TestingDutySlotV(version))+1))},
			OutputMessages: []*types.SignedPartialSignatureMessage{},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
