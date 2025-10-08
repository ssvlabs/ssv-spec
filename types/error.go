package types

const (
	UnmarshalSSZErrorCode int = iota + 1
	FutureMessageErrorCode
	ZeroCommitMessagesErrorCode
	NoRunningDutyErrorCode
	SignedCommitIsInvalidErrorCode
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
	NoPartialSigMsgsErrorCode
	NoRunnerForSlotErrorCode
	SkipConsensusMsgAsInstanceIsDecidedErrorCode
	SkipConsensusMsgAsConsensusHasFinishedErrorCode
	DecodeBeaconVoteErrorCode
	NoBeaconDutiesErrorCode
	NoValidatorSharesErrorCode
	AttestationSourceNotLessThanTargetErrorCode
	AttestationTargetEpochTooFarFutureErrorCode
	MsgIDCommitteeIDMismatchErrorCode
	MsgDoesNotBelongToIndetifierErrorCode
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
	InvalidPartialSigMsgSlot
	FuturePartialSigMsgSlot
	UnknownValidatorIndexErrorCode
	WrongRootsCountErrorCode
	DutyEpochTooFarFutureErrorCode
	WrongBeaconRoleTypeErrorCode
	WrongValidatorIndexErrorCode
	WrongValidatorPubkeyErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
	//TodoErrorCode
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
