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

// TestingKnownDomainTypes lists every domain DomainType.IsKnown accepts
var TestingKnownDomainTypes = []types.DomainType{
	types.GenesisMainnet, types.AlanMainnet, types.BooleMainnet,
	types.PrimusTestnet,
	types.ShifuTestnet, types.ShifuV2Testnet,
	types.JatoTestnet, types.JatoV2Testnet, types.JatoAlanTestnet,
}
var TestingUnknownDomainType = types.DomainType{0xaa, 0xbb, 0xcc, 0xdd}

// ==================================================
// Consensus Data - Invalid Types
// ==================================================

var EncodeConsensusDataTest = func(cd *types.ProposerConsensusData) []byte {
	encodedCD, _ := cd.Encode()
	return encodedCD
}

var EncodeAggregatorCommitteeConsensusDataTest = func(cd *types.AggregatorCommitteeConsensusData) []byte {
	encodedCD, _ := cd.Encode()
	return encodedCD
}

var TestConsensusUnkownDutyTypeData = &types.ProposerConsensusData{
	Duty:    TestingUnknownDutyType,
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ProposerConsensusData{
	Duty:    TestingWrongDutyPK,
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestConsensusWrongDutyPKDataByts, _ = TestConsensusWrongDutyPKData.Encode()

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgWrongID = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingWrongValidatorPubKey), types.RoleCommittee))
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
