package spectest

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	committeemultipleduty "github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee/multipleduty"
	committeesingleduty "github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee/singleduty"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/dutyexe"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/partialsigcontainer"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/consensus"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/proposer"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/synccommitteeaggregator"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/postconsensus"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/preconsensus"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck/valcheckattestations"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck/valcheckduty"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck/valcheckproposer"
)

var AllTests = []tests.TestF{
	runner.FullHappyFlow,

	postconsensus.TooManyRoots,
	postconsensus.TooFewRoots,
	postconsensus.UnorderedExpectedRoots,
	postconsensus.UnknownSigner,
	postconsensus.InconsistentBeaconSigner,
	postconsensus.PostFinish,
	postconsensus.NoRunningDuty,
	postconsensus.InvalidMessageSignature,
	postconsensus.InvalidBeaconSignatureInQuorum,
	postconsensus.DuplicateMsgDifferentRoots,
	postconsensus.DuplicateMsgDifferentRootsThenQuorum,
	postconsensus.DuplicateMsg,
	postconsensus.InvalidExpectedRoot,
	postconsensus.PreDecided,
	postconsensus.PostQuorum,
	postconsensus.InvalidMessage,
	postconsensus.InvalidOperatorSignature,
	postconsensus.InvalidMessageSlot,
	postconsensus.ValidMessage,
	postconsensus.ValidMessage7Operators,
	postconsensus.ValidMessage10Operators,
	postconsensus.ValidMessage13Operators,
	postconsensus.Quorum,
	postconsensus.Quorum7Operators,
	postconsensus.Quorum10Operators,
	postconsensus.Quorum13Operators,
	postconsensus.InvalidDecidedValue,
	postconsensus.InvalidThenQuorum,
	postconsensus.InvalidQuorumThenValidQuorum,
	postconsensus.InconsistentOperatorSigner,
	postconsensus.NilSSVMessage,
	postconsensus.InvalidValidatorIndex,
	postconsensus.PartialInvalidRootQuorumThenValidQuorum,
	postconsensus.PartialInvalidSigQuorumThenValidQuorum,
	postconsensus.MixedCommittees,

	newduty.ConsensusNotStarted,
	newduty.NotDecided,
	newduty.PostDecided,
	newduty.Finished,
	newduty.Valid,
	newduty.PostWrongDecided,
	newduty.PostInvalidDecided,
	newduty.PostFutureDecided,
	newduty.DuplicateDutyFinished,
	newduty.DuplicateDutyNotFinished,
	newduty.FirstHeight,

	committeesingleduty.StartDuty,
	committeesingleduty.StartNoDuty,
	committeesingleduty.ValidBeaconVote,
	committeesingleduty.WrongBeaconVote,
	committeesingleduty.Decided,
	committeesingleduty.HappyFlow,
	committeesingleduty.PastMessageDutyNotFinished,
	committeesingleduty.PastMessageDutyFinished,
	committeesingleduty.PastMessageDutyDoesNotExist,
	committeesingleduty.ProposalWithConsensusData,
	committeesingleduty.WrongMessageID,

	committeemultipleduty.SequencedDecidedDuties,
	committeemultipleduty.SequencedHappyFlowDuties,
	committeemultipleduty.ShuffledDecidedDuties,
	committeemultipleduty.ShuffledHappyFlowDutiesWithTheSameValidators,
	committeemultipleduty.ShuffledHappyFlowDutiesWithDifferentValidators,
	committeemultipleduty.FailedThanSuccessfulDuties,

	consensus.FutureDecidedNoInstance,
	consensus.FutureDecided,
	consensus.InvalidDecidedValue,
	consensus.FutureMessage,
	consensus.PastMessage,
	consensus.PostFinish,
	consensus.PostDecided,
	consensus.ValidDecided,
	consensus.ValidDecided7Operators,
	consensus.ValidDecided10Operators,
	consensus.ValidDecided13Operators,
	consensus.ValidMessage,
	consensus.InvalidSignature,
	consensus.DecidedSlashableAttestation,

	synccommitteeaggregator.SomeAggregatorQuorum,
	synccommitteeaggregator.NoneAggregatorQuorum,
	synccommitteeaggregator.AllAggregatorQuorum,

	proposer.ProposeBlindedBlockDecidedRegular,
	proposer.ProposeRegularBlockDecidedBlinded,
	proposer.BlindedRunnerAcceptsNormalBlock,
	proposer.NormalProposerAcceptsBlindedBlock,

	preconsensus.NoRunningDuty,
	preconsensus.TooFewRoots,
	preconsensus.TooManyRoots,
	preconsensus.UnorderedExpectedRoots,
	preconsensus.InvalidSignedMessage,
	preconsensus.InvalidOperatorSignature,
	preconsensus.InvalidExpectedRoot,
	preconsensus.DuplicateMsg,
	preconsensus.DuplicateMsgDifferentRoots,
	preconsensus.PostFinish,
	preconsensus.PostDecided,
	preconsensus.PostQuorum,
	preconsensus.Quorum,
	preconsensus.Quorum7Operators,
	preconsensus.Quorum10Operators,
	preconsensus.Quorum13Operators,
	preconsensus.ValidMessage,
	preconsensus.InvalidMessageSlot,
	preconsensus.ValidMessage7Operators,
	preconsensus.ValidMessage10Operators,
	preconsensus.ValidMessage13Operators,
	preconsensus.InconsistentBeaconSigner,
	preconsensus.UnknownSigner,
	preconsensus.InvalidBeaconSignatureInQuorum,
	preconsensus.InvalidMessageSignature,
	preconsensus.InvalidThenQuorum,
	preconsensus.InvalidQuorumThenValidQuorum,
	preconsensus.InconsistentOperatorSigner,
	preconsensus.NilSSVMessage,

	valcheckduty.WrongValidatorIndex,
	valcheckduty.WrongValidatorPK,
	valcheckduty.WrongDutyType,
	valcheckduty.FarFutureDutySlot,

	valcheckattestations.Slashable,
	valcheckattestations.SourceHigherThanTarget,
	valcheckattestations.FarFutureTarget,
	valcheckattestations.BeaconVoteDataNil,
	valcheckattestations.Valid,
	valcheckattestations.MinoritySlashable,
	valcheckattestations.MajoritySlashable,

	valcheckproposer.BlindedBlock,

	dutyexe.WrongDutyRole,
	dutyexe.WrongDutyPubKey,
	partialsigcontainer.OneSignature,
	partialsigcontainer.Quorum,
	partialsigcontainer.Duplicate,
	partialsigcontainer.DuplicateQuorum,
	partialsigcontainer.Invalid,
}
