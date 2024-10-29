package testingutils

import (
	"crypto/rsa"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Mock values
var (
	MockPartialSignature  = [96]byte{1, 2, 3, 4}
	MockRSASignature      = [validation.RsaSignatureSize]byte{1, 2, 3, 4}
	DefaultKeySet         = Testing4SharesSet()
	DefaultValidatorIndex = phase0.ValidatorIndex(1)
	DefaultMsgID          = types.NewMsgID(TestingSSVDomainType, TestingValidatorPK[:], types.RoleProposer)

	// Map to get a KeySet for a certain CommitteeID, as defined in the testing TestingMessageValidator().NetworkDataFetcher
	KeySetForCommitteeID = map[types.CommitteeID]*TestKeySet{
		TestingCommitteeID:                      Testing4SharesSet(),
		TestingCommitteeIDWithSyncCommitteeDuty: Testing13SharesSet(),
	}
	// Map to get a valid ValidatorIndex that belongs to a CommitteeID, as defined in the testing TestingMessageValidator().NetworkDataFetcher
	ValidValidatorIndexForCommitteeID = map[types.CommitteeID]phase0.ValidatorIndex{
		TestingCommitteeID:                      1,
		TestingCommitteeIDWithSyncCommitteeDuty: ValidatorIndexWithSyncCommitteeDuty,
	}
)

// Encode SignedSSVMessage with no error
func EncodeMessage(msg *types.SignedSSVMessage) []byte {
	msgBytes, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	return msgBytes
}

// Encode qbft.Message with no error
func EncodeQbftMessage(msg *qbft.Message) []byte {
	msgBytes, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	return msgBytes
}

// Encode types.PartialSignatureMessages with no error
func EncodePartialSignatureMessage(msg *types.PartialSignatureMessages) []byte {
	msgBytes, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	return msgBytes
}

// Signs an types.SSVMessage with no error
func SignSSVMessage(sk *rsa.PrivateKey, ssvMsg *types.SSVMessage) []byte {
	signature, err := types.SignSSVMessage(sk, ssvMsg)
	if err != nil {
		panic(err)
	}
	return signature
}

// Get any valid psig type for a runner role
func ValidPartialSignatureTypeForRole(role types.RunnerRole) types.PartialSigMsgType {
	msgType := types.PostConsensusPartialSig
	if role == types.RoleValidatorRegistration {
		msgType = types.ValidatorRegistrationPartialSig
	}
	if role == types.RoleVoluntaryExit {
		msgType = types.VoluntaryExitPartialSig
	}
	return msgType
}

// Get default MessageID for role
func MessageIDForRole(role types.RunnerRole) types.MessageID {
	msgID := types.NewMsgID(TestingSSVDomainType, TestingValidatorPK[:], role)
	if role == types.RoleCommittee {
		msgID = types.NewMsgID(TestingSSVDomainType, TestingCommitteeID[:], role)
	}
	return msgID
}

// Get default MessageID for role
func MessageIDForRoleAndCommitteeID(role types.RunnerRole, committeeID types.CommitteeID) types.MessageID {
	msgID := types.NewMsgID(TestingSSVDomainType, TestingValidatorPK[:], role)
	if role == types.RoleCommittee {
		msgID = types.NewMsgID(TestingSSVDomainType, committeeID[:], role)
	}
	return msgID
}

// Returns receivedAt time for a given round. It uses the default slot qbft.FirstHeight(0)
func ReceivedAtForRound(round qbft.Round) time.Time {
	roundDelta := round - qbft.FirstRound
	if roundDelta <= 8 {
		receivedAt := time.Unix(NewTestingBeaconNode().GetBeaconNetwork().EstimatedTimeAtSlot(phase0.Slot(qbft.FirstHeight)), 0)
		receivedAt = receivedAt.Add(time.Duration(roundDelta) * qbft.QuickTimeout)
		return receivedAt
	} else {
		receivedAt := time.Unix(NewTestingBeaconNode().GetBeaconNetwork().EstimatedTimeAtSlot(phase0.Slot(qbft.FirstHeight)), 0)
		receivedAt = receivedAt.Add(8 * qbft.QuickTimeout).Add(time.Duration(roundDelta-8) * qbft.SlowTimeout)
		return receivedAt
	}
}

// Get the slot out of a message
func GetMessageSlot(data []byte) phase0.Slot {
	msg := &types.SignedSSVMessage{}
	err := msg.Decode(data)
	if err != nil {
		return 0
	}
	if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
		qbftMsg := &qbft.Message{}
		err := qbftMsg.Decode(msg.SSVMessage.Data)
		if err != nil {
			return 0
		}
		return phase0.Slot(qbftMsg.Height)
	}
	if msg.SSVMessage.MsgType == types.SSVPartialSignatureMsgType {
		pSigMsgs := &types.PartialSignatureMessages{}
		err := pSigMsgs.Decode(msg.SSVMessage.Data)
		if err != nil {
			return 0
		}
		return pSigMsgs.Slot
	}
	return 0
}

// Returns a default signed consensus message for a specific slot with a given MsgID
var ConsensusMsgForSlot = func(slot phase0.Slot, msgID types.MessageID, ks *TestKeySet) *types.SignedSSVMessage {
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data: EncodeQbftMessage(&qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Round:      qbft.FirstRound,
			Identifier: msgID[:],
			Root:       [32]byte{1},
			Height:     qbft.Height(slot),
		}),
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{SignSSVMessage(ks.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}
}

// Returns a prepare consensus message for a certain round, with a given MessageID
var ConsensusMessageForRound = func(round qbft.Round, msgID types.MessageID) *types.SignedSSVMessage {
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data: EncodeQbftMessage(&qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Round:      round,
			Height:     qbft.FirstHeight,
			Root:       TestingQBFTRootData,
			Identifier: msgID[:],
		}),
	}
	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{SignSSVMessage(DefaultKeySet.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}
}

// Returns a default signed partial signature message for a specific slot with a given MsgID and validator index
var PartialSignatureMsgForSlot = func(slot phase0.Slot, msgID types.MessageID, validatorIndex phase0.ValidatorIndex, ks *TestKeySet) *types.SignedSSVMessage {

	role := msgID.GetRoleType()

	msgType := ValidPartialSignatureTypeForRole(role)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data: EncodePartialSignatureMessage(&types.PartialSignatureMessages{
			Type: msgType,
			Slot: slot,
			Messages: []*types.PartialSignatureMessage{
				{
					PartialSignature: MockPartialSignature[:],
					SigningRoot:      [32]byte{1},
					Signer:           1,
					ValidatorIndex:   validatorIndex,
				},
			},
		}),
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{SignSSVMessage(ks.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}
}

// Returns a signed partial signature message for a specific signature type and role
var PartialSignatureMsgForSignatureTypeRoleAndRoot = func(pSigType types.PartialSigMsgType, role types.RunnerRole, root [32]byte) *types.SignedSSVMessage {

	msgID := MessageIDForRole(role)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data: EncodePartialSignatureMessage(&types.PartialSignatureMessages{
			Type: pSigType,
			Slot: 0,
			Messages: []*types.PartialSignatureMessage{
				{
					PartialSignature: MockPartialSignature[:],
					SigningRoot:      root,
					Signer:           1,
					ValidatorIndex:   1,
				},
			},
		}),
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{SignSSVMessage(DefaultKeySet.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}
}

// Returns a signed partial signature message for a specific signature type and role
var PartialSignatureMsgForSignatureTypeAndRole = func(pSigType types.PartialSigMsgType, role types.RunnerRole) *types.SignedSSVMessage {
	return PartialSignatureMsgForSignatureTypeRoleAndRoot(pSigType, role, [32]byte{1})
}

// Returns a partial signature message given a the number of signatures
var PartialSignatureMsgForNumSignatures = func(numSignatures int, role types.RunnerRole, numCommitteeValidators int) *types.SignedSSVMessage {

	msgID := MessageIDForRole(role)
	msgType := ValidPartialSignatureTypeForRole(role)

	pSigMsgs := make([]*types.PartialSignatureMessage, 0)
	for i := 0; i < numSignatures; i++ {

		validatorIndex := DefaultValidatorIndex
		if role == types.RoleCommittee {
			// For the committee role, keep switching the validator index
			validatorIndex = phase0.ValidatorIndex(1 + (i % numCommitteeValidators))
		}
		pSigMsgs = append(pSigMsgs, &types.PartialSignatureMessage{
			PartialSignature: MockPartialSignature[:],
			SigningRoot:      [32]byte{1},
			Signer:           1,
			ValidatorIndex:   validatorIndex,
		})
	}
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data: EncodePartialSignatureMessage(&types.PartialSignatureMessages{
			Type:     msgType,
			Slot:     0,
			Messages: pSigMsgs,
		}),
	}
	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{SignSSVMessage(DefaultKeySet.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}
}
