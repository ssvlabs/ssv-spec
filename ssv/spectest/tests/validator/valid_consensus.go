package validator

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func ValidConsensus() tests.SpecTest {
	// KeySet
	ks := testingutils.Testing4SharesSet()

	// Message ID
	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	// Duty
	duty := testingutils.TestingAttesterDuty

	// data root
	dataRoot, err := qbft.HashDataRoot(testingutils.TestAttesterConsensusDataByts)
	if err != nil {
		panic(err.Error())
	}

	// Proposal message
	proposal := testingutils.TestingProposalMessageWithHeight(ks.Shares[1], 1, qbft.Height(duty.Slot))
	proposal.Message.Identifier = msgID[:]
	proposal.Message.Root = dataRoot
	proposal = testingutils.SignQBFTMsg(ks.Shares[1], 1, &proposal.Message)
	proposal.FullData = testingutils.TestAttesterConsensusDataByts
	proposalByts, err := proposal.Encode()
	if err != nil {
		panic(err.Error())
	}

	// Prepare message
	prepare := testingutils.TestingPrepareMessageWithHeight(ks.Shares[1], 1, qbft.Height(duty.Slot))
	prepare.Message.Identifier = msgID[:]
	prepare.Message.Root = dataRoot
	prepare = testingutils.SignQBFTMsg(ks.Shares[1], 1, &prepare.Message)
	prepareByts, err := prepare.Encode()
	if err != nil {
		panic(err.Error())
	}

	// Messages
	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    proposalByts[:],
		},
	}

	outMsgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    proposalByts[:],
		},
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    prepareByts[:],
		},
	}

	return &ValidatorTest{
		Name:                   "valid consensus",
		KeySet:                 ks,
		Duties:                 []*types.Duty{&duty},
		Messages:               msgs,
		OutputMessages:         outMsgs,
		BeaconBroadcastedRoots: []string{},
	}
}
