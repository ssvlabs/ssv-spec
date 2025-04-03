package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// DomainType and Fork Data
// ==================================================

var TestingSSVDomainType = types.JatoTestnet
var TestingForkData = types.ForkData{Epoch: TestingDutyEpoch, Domain: TestingSSVDomainType}

// ==================================================
// Consensus Data - Invalid Types
// ==================================================

var EncodeConsensusDataTest = func(cd *types.ValidatorConsensusData) []byte {
	encodedCD, _ := cd.Encode()
	return encodedCD
}

var TestConsensusUnkownDutyTypeData = &types.ValidatorConsensusData{
	Duty:    TestingUnknownDutyType,
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ValidatorConsensusData{
	Duty:    TestingWrongDutyPK,
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestConsensusWrongDutyPKDataByts, _ = TestConsensusWrongDutyPKData.Encode()

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgWrongID = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingWrongValidatorPubKey[:], types.RoleCommittee))
}

var ssvMsg = func(qbftMsg *types.SignedSSVMessage, postMsg *types.PartialSignatureMessages, msgID types.MessageID) *types.SSVMessage {

	if qbftMsg != nil {
		return &types.SSVMessage{
			MsgType: qbftMsg.SSVMessage.MsgType,
			MsgID:   msgID,
			Data:    qbftMsg.SSVMessage.Data,
		}
	}

	if postMsg != nil {
		msgType := types.SSVPartialSignatureMsgType
		data, err := postMsg.Encode()
		if err != nil {
			panic(err)
		}
		return &types.SSVMessage{
			MsgType: msgType,
			MsgID:   msgID,
			Data:    data,
		}
	}

	panic("msg type undefined")
}
