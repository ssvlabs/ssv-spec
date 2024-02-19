package validator

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// BeaconRootBroadcasting tests the broadcast of a beacon root through a full attestation duty
func BeaconRootBroadcasting() tests.SpecTest {
	// KeySet
	ks := testingutils.Testing4SharesSet()

	// Message ID
	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	// Duty
	duty := testingutils.TestingAttesterDuty

	// QBFT Messages
	qbftMsgs := testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester)

	// Decided message
	decided := testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		qbft.Height(duty.Slot),
		msgID[:], testingutils.TestAttesterConsensusDataByts)
	decidedMsgByts, err := decided.Encode()
	if err != nil {
		panic(err.Error())
	}
	decidedMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    decidedMsgByts[:],
	}

	// Post-consensus messages
	postConsensusMsgs := make([]*types.SSVMessage, 0)
	opID := uint64(1)
	for opID <= ks.Threshold {
		postConsensusMessage := testingutils.PostConsensusAttestationMsg(ks.Shares[opID], opID, qbft.Height(duty.Slot))
		postConsensusMessageByts, err := postConsensusMessage.Encode()
		if err != nil {
			panic(err.Error())
		}
		postConsensusMsgs = append(postConsensusMsgs, &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   msgID,
			Data:    postConsensusMessageByts[:],
		})
		opID += 1
	}

	// Incoming Messages
	msgs := []*types.SSVMessage{}
	msgs = append(msgs, qbftMsgs...)
	msgs = append(msgs, decidedMsg)
	msgs = append(msgs, postConsensusMsgs...)

	// Output messages
	outMsgs := []*types.SSVMessage{
		qbftMsgs[0],
		qbftMsgs[1],
		qbftMsgs[4],
		decidedMsg,
		postConsensusMsgs[0],
	}

	// Beacon root
	beaconRoot := testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks))

	return &ValidatorTest{
		Name:                   "beacon root broadcasting",
		Duties:                 []*types.Duty{&duty},
		Messages:               msgs,
		OutputMessages:         outMsgs,
		BeaconBroadcastedRoots: []string{beaconRoot},
	}
}
