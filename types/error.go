package types

import (
	"github.com/pkg/errors"
)

const (
	UnmarshalSSZErrorCode int = iota + 1
	FutureMessageErrorCode
	ZeroCommitMessagesErrorCode
	NoRunningDutyErrorCode
	CommitMessageInvalidErrorCode
	CommitMessageTypeWrongErrorCode
	WrongMessageHeightErrorCode
	SignerIsNotInCommitteeErrorCode
	NonUniqueSignerErrorCode
	ValidatorRegistrationNoConsensusDataErrorCode
	ValidatorExitNoConsensusDataErrorCode
	UnknownDutyRoleDataErrorCode
	UnknownBlockVersionErrorCode
	IncorrectNumberOfSignaturesErrorCode
	EmptySignatureErrorCode
	NilSSVMessageErrorCode
	NoSignaturesErrorCode
	NoSignersErrorCode
	ZeroSignerNotAllowedErrorCode
	InconsistentSignersErrorCode
	NoPartialSigMessagesErrorCode
	NoRunnerForSlotErrorCode
	SkipConsensusMessageAsInstanceIsDecidedErrorCode
	SkipConsensusMessageAsConsensusHasFinishedErrorCode
	DecodeBeaconVoteErrorCode
	NoBeaconDutiesErrorCode
	NoValidatorSharesErrorCode
	AttestationSourceNotLessThanTargetErrorCode
	CheckpointMismatch
	MessageIDCommitteeIDMismatchErrorCode
	MessageTypeInvalidErrorCode
	MessageRoundInvalidErrorCode
	MessageIdentifierInvalidErrorCode
	ReconstructSignatureErrorCode
	SlashableAttestationErrorCode
	DecidedWrongInstanceErrorCode
	ValidatorRegistrationNoConsensusPhaseErrorCode
	ValidatorRegistrationNoPostConsensusPhaseErrorCode
	ValidatorExitNoConsensusPhaseErrorCode
	ValidatorExitNoPostConsensusPhaseErrorCode
	SSVMessageHasInvalidSignatureErrorCode
	DutyAlreadyPassedErrorCode
	WrongSigningRootErrorCode
	PartialSigInconsistentSignerErrorCode
	NoDecidedValueErrorCode
	NoRunningConsensusInstanceErrorCode
	ConsensusInstanceNotDecidedErrorCode
	PartialSigMessageInvalidSlotErrorCode
	PartialSigMessageFutureSlotErrorCode
	UnknownValidatorIndexErrorCode
	WrongRootsCountErrorCode
	DutyEpochTooFarFutureErrorCode
	WrongBeaconRoleTypeErrorCode
	WrongValidatorIndexErrorCode
	WrongValidatorPubkeyErrorCode
	InstanceStoppedProcessingMessagesErrorCode
	WrongMessageRoundErrorCode
	MessageAllowsOneSignerOnlyErrorCode
	NoProposalForCurrentRoundErrorCode
	PastRoundErrorCode
	ProposedDataMismatchErrorCode
	RoundChangeNoQuorumErrorCode
	ProposalInvalidErrorCode
	RootHashInvalidErrorCode
	MessageSignatureInvalidErrorCode
	QBFTValueInvalidErrorCode
	PrepareMessageInvalidErrorCode
	JustificationsNoQuorumInvalidErrorCode
	ProposalLeaderInvalidErrorCode
	InstanceAlreadyRunningErrorCode
	StartInstanceErrorCode
	TimeoutInstanceErrorCode
	AggCommAggAttCntMismatchErrorCode
	AggCommAggCommIdxCntMismatchErrorCode
	AggCommContributorsContributionsCntMismatchErrorCode
	AggCommSubnetNotInSCSubnetsErrorCode
)

type Error struct {
	Code int

	innerErr error
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:     code,
		innerErr: errors.New(message), // nolint:staticcheck
	}
}

func WrapError(code int, err error) *Error {
	return &Error{
		Code:     code,
		innerErr: err,
	}
}

func (e Error) Error() string {
	return e.innerErr.Error()
}

func (e Error) Is(target error) bool {
	var err *Error
	ok := errors.As(target, &err)
	return ok
}
