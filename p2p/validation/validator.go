package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/libp2p/go-libp2p/core/peer"
)

type MessageValidation struct {
	selfPID           peer.ID
	signatureVerifier types.SignatureVerifier
}

func (mv *MessageValidation) ValidDomain(domain []byte) bool {
	return false
}
func (mv *MessageValidation) ValidRole(role types.RunnerRole) bool {
	return false
}
func (mv *MessageValidation) ActiveValidator(validatorPK []byte) bool {
	return false
}
func (mv *MessageValidation) ValidatorLiquidated(validatorPK []byte) bool {
	return false
}
func (mv *MessageValidation) ExistingClusterID(clusterID []byte) bool {
	return false
}
func (mv *MessageValidation) ValidConsensusMessageType(msgType qbft.MessageType) bool {
	return false
}
func (mv *MessageValidation) SignerBelongsToCommittee(signer types.OperatorID, msgID types.MessageID) bool {
	return false
}
func (mv *MessageValidation) ExpectedPartialSignatureTypeForRole(sigType types.PartialSigMsgType, msgID types.MessageID) bool {
	return false
}
func (mv *MessageValidation) ValidatorIndexBelongsToCommittee(validatorIndex phase0.ValidatorIndex, msgID types.MessageID) bool {
	return false
}
func (mv *MessageValidation) MatchedIdentifiers(msgID1 []byte, msgID2 []byte) bool {
	return false
}
func (mv *MessageValidation) UniqueSigners(signers []types.OperatorID) bool {
	return false
}
func (mv *MessageValidation) ValidSignersLengthForCommitMessage(signers []types.OperatorID) bool {
	return false
}
func (mv *MessageValidation) RoundBelongToAllowedSpread(peerID peer.ID, round qbft.Round, msgID types.MessageID) error {
	return nil
}
func (mv *MessageValidation) ValidFullDataRoot(fullData []byte, root [32]byte) bool {
	return false
}
func (mv *MessageValidation) ValidConsensusMessageCount(peerID peer.ID, msgID types.MessageID, height qbft.Height, msgType qbft.MessageType, round qbft.Round) error {
	return nil
}
func (mv *MessageValidation) IsLeader(signer types.OperatorID, height qbft.Height, round qbft.Round) bool {
	return false
}
func (mv *MessageValidation) ValidJustificationForProposal(signedSSVMessage *types.SignedSSVMessage) bool {
	return false
}
func (mv *MessageValidation) ValidJustificationForRoundChange(signedSSVMessage *types.SignedSSVMessage) bool {
	return false
}
func (mv *MessageValidation) HasSentDecidedWithSameNumberOfSigners(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) bool {
	return false
}
func (mv *MessageValidation) ValidRoundForRole(round qbft.Round, role types.RunnerRole) bool {
	return false
}
func (mv *MessageValidation) ValidPartialSigMessageCount(peerID peer.ID, msgID types.MessageID, partialSignatureMessages *types.PartialSignatureMessages) error {
	return nil
}
func (mv *MessageValidation) ValidDutySlot(peerID peer.ID, slot phase0.Slot, role types.RunnerRole) error {
	return nil
}
func (mv *MessageValidation) ValidRoleForConsensus(role types.RunnerRole) bool {
	return false
}
func (mv *MessageValidation) ValidNumberOfSignaturesForCommitteeDuty(senderID []byte, partialSignatureMessages *types.PartialSignatureMessages) bool {
	return false
}
func (mv *MessageValidation) ValidNumberOfDutiesPerEpoch(peerID peer.ID, msgID types.MessageID, slot phase0.Slot) bool {
	return false
}
func (mv *MessageValidation) NoTripleValidatorOccurrence(partialSignatureMessages *types.PartialSignatureMessages) bool {
	return false
}
func (mv *MessageValidation) ValidNumberOfAggregatorDutiesPerEpoch(peerID peer.ID, senderID []byte) bool {
	return false
}
func (mv *MessageValidation) HasProposerDuty(senderID []byte, slot phase0.Slot) bool {
	return false
}
func (mv *MessageValidation) HasSyncCommitteeDuty(senderID []byte, slot phase0.Slot) bool {
	return false
}
func (mv *MessageValidation) ValidNumberOfValidatorRegistrationDutiesPerEpoch(peerID peer.ID, slot phase0.Slot) bool {
	return false
}
func (mv *MessageValidation) ValidNumberOfVoluntaryExitDutiesPerEpoch(peerID peer.ID, slot phase0.Slot) bool {
	return false
}
func (mv *MessageValidation) CorrectTopic(msgID types.MessageID, topic string) bool {
	return false
}
func (mv *MessageValidation) ValidPartialSignatureType(signatureType types.PartialSigMsgType) bool {
	return false
}
func (mv *MessageValidation) PeerAlreadyAdvancedRound(peerID peer.ID, msgID types.MessageID, height qbft.Height, round qbft.Round) bool {
	return false
}
func (mv *MessageValidation) PeerHasSentProposalWithDifferentData(peerID peer.ID, msgID types.MessageID, height qbft.Height, round qbft.Round, fullData []byte) bool {
	return false
}
func (mv *MessageValidation) HasMoreSignersThanCommitteeSize(signers []types.OperatorID, msgID types.MessageID) bool {
	return false
}
