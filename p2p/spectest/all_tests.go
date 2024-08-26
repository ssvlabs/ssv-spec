package spectest

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	validation "github.com/ssvlabs/ssv-spec/p2p/spectest/tests/validation"
)

var AllTests = []tests.TestF{
	// PubSub validation
	validation.PubSubNoData,
	validation.PubSubDataTooBig,
	validation.PubSubMalformed,

	// Syntax validation
	validation.NoSigners,
	validation.NoSignatures,
	validation.WrongRSASignatureSize,
	validation.SignersNotSorted,
	validation.ZeroSigner,
	validation.DuplicatedSigner,
	validation.SignersAndSignaturesWithDifferentLength,
	// validation.NilSSVMessage,
	validation.EmptyData,
	// validation.SSVMessageDataTooBig,
	validation.ConsensusUndecodableSSVMessageData,
	validation.PartialSignatureSSVMessageUndecodableData,
	validation.UndecodablePrepareJustificaiton,
	validation.UndecodableRoundChangeJustification,

	// General semantics validation
	validation.SignerNotInCommittee,
	validation.WrongDomain,
	validation.InvalidRole,
	validation.UnknownValidator,
	validation.ValidatorNotAttesting,
	validation.ValidatorLiquidated,
	validation.NonExistentCommitteeID,
	validation.IncorrectTopic,
	validation.DKGMessage,
	validation.UnknownMessageType,

	// Consensus semantics validation
	validation.NonDecidedWithMultipleSigners,
	validation.DecidedNotEnoughSigners,
	validation.PrepareWithFullData,
	validation.CommitWithFullData,
	validation.InvalidHash,
	validation.UnknownQBFTMessageType,
	validation.ZeroRound,
	validation.MismatchedIdentifier,

	// Consensus QBFT logic validation
	validation.SignerNotLeader,
	validation.DecidedWithSameSigners,
	validation.DuplicatedProposalWithDifferentData,
	validation.DuplicatedConsensusMessage,
	validation.RoundNotAllowedInTimeSpread,
	validation.RoundAlreadyAdvanced,
	validation.UnexpectedRoundChangeJustificationInPrepareMessage,
	validation.UnexpectedRoundChangeJustificationInCommitMessage,
	validation.UnexpectedPrepareJustificationInPrepareMessage,
	validation.UnexpectedPrepareJustificationInCommitMessage,
	validation.UnexpectedPrepareJustificationInRoundChangeMessage,

	// Consensus duty logic validation
	validation.UnexpectedConsensusMessageForValidatorRegistration,
	validation.UnexpectedConsensusMessageForVoluntaryExit,
	validation.RoundTooHigh,

	// Partial signature semantics validation
	validation.PartialSigMultipleSigners,
	validation.PartialSigWithFullData,
	validation.PartialSigInvalidType,
	validation.PartialSigTypeMismatch,
	validation.NoPartialSignatureMessages,
	// validation.WrongPartialSignatureSize,
	validation.PartialSigInconsistentSigners,
	validation.PartialSigValidatorIndexMismatch,
	validation.PartialSigValidatorIndexMismatchForCommitteeRole,

	// Partial signature duty logic validation
	validation.PartialSigTypeCount,
	validation.PartialSigTooManySignatures,
	validation.PartialSigTripleValidatorIndex,

	// General duty logic validation
	validation.SlotAlreadyAdvancedForConsensusMessage,
	validation.SlotAlreadyAdvancedForPartialSignatureMessage,
	validation.ConsensusMessageSentBeforeSlotStarts,
	validation.PartialSignatureMessageSentBeforeSlotStarts,
	validation.ConsensusMessageLateSlot,
	validation.PartialSignatureMessageLateSlot,
	validation.NoProposerDutyConsensusMessage,
	validation.NoProposerDutyPartialSignatureMessage,
	validation.NoSyncCommitteeContributionDutyConsensusMessage,
	validation.NoSyncCommitteeContributionDutyPartialSignatureMessage,
	validation.TooManyDutiesWithConsensusMessages,
	validation.TooManyDutiesWithPartialSignatureMessages,

	// Wrong signature
	validation.WrongSignatureConsensusMessage,
	validation.WrongSignaturePartialSignatureMessage,

	// Valid cases
	validation.HappyFlow,
}
