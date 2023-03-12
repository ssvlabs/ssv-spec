package spectest

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/messages"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/consensus"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/proposer"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/synccommitteeaggregator"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/postconsensus"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/preconsensus"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck/valcheckattestations"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck/valcheckduty"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	runner.FullHappyFlow(),

	postconsensus.TooManyRoots(),
	postconsensus.TooFewRoots(),
	postconsensus.UnorderedExpectedRoots(),
	postconsensus.UnknownSigner(),
	postconsensus.InconsistentBeaconSigner(),
	postconsensus.PostFinish(),
	postconsensus.NoRunningDuty(),
	postconsensus.InvalidMessageSignature(),
	postconsensus.InvalidBeaconSignature(),
	postconsensus.DuplicateMsgDifferentRoots(),
	postconsensus.DuplicateMsg(),
	postconsensus.InvalidExpectedRoot(),
	postconsensus.PreDecided(),
	postconsensus.PostQuorum(),
	postconsensus.InvalidMessage(),
	postconsensus.InvalidMessageSlot(),
	postconsensus.ValidMessage(),
	postconsensus.ValidMessage7Operators(),
	postconsensus.ValidMessage10Operators(),
	postconsensus.ValidMessage13Operators(),
	postconsensus.Quorum(),
	postconsensus.Quorum7Operators(),
	postconsensus.Quorum10Operators(),
	postconsensus.Quorum13Operators(),
	postconsensus.InvalidDecidedValue(),

	newduty.ConsensusNotStarted(),
	newduty.NotDecided(),
	newduty.PostDecided(),
	newduty.Finished(),
	newduty.Valid(),
	newduty.PostWrongDecided(),
	newduty.PostInvalidDecided(),
	newduty.PostFutureDecided(),

	consensus.FutureDecided(),
	consensus.InvalidDecidedValue(),
	consensus.NoRunningDuty(),
	consensus.NoRunningConsensusInstance(),
	consensus.PostFinish(),
	consensus.PostDecided(),
	consensus.ValidDecided(),
	consensus.ValidDecided7Operators(),
	consensus.ValidDecided10Operators(),
	consensus.ValidDecided13Operators(),
	consensus.ValidMessage(),

	synccommitteeaggregator.SomeAggregatorQuorum(),
	synccommitteeaggregator.NoneAggregatorQuorum(),
	synccommitteeaggregator.AllAggregatorQuorum(),

	proposer.ProposeBlindedBlockDecidedRegular(),
	proposer.ProposeRegularBlockDecidedBlinded(),

	preconsensus.NoRunningDuty(),
	preconsensus.TooFewRoots(),
	preconsensus.TooManyRoots(),
	preconsensus.UnorderedExpectedRoots(),
	preconsensus.InvalidSignedMessage(),
	preconsensus.InvalidExpectedRoot(),
	preconsensus.DuplicateMsg(),
	preconsensus.DuplicateMsgDifferentRoots(),
	preconsensus.PostFinish(),
	preconsensus.PostDecided(),
	preconsensus.PostQuorum(),
	preconsensus.Quorum(),
	preconsensus.Quorum7Operators(),
	preconsensus.Quorum10Operators(),
	preconsensus.Quorum13Operators(),
	preconsensus.ValidMessage(),
	preconsensus.InvalidMessageSlot(),
	preconsensus.ValidMessage7Operators(),
	preconsensus.ValidMessage10Operators(),
	preconsensus.ValidMessage13Operators(),
	preconsensus.InconsistentBeaconSigner(),
	preconsensus.UnknownSigner(),
	preconsensus.InvalidBeaconSignature(),
	preconsensus.InvalidMessageSignature(),

	messages.EncodingAndRoot(),
	messages.NoMsgs(),
	messages.InvalidMsg(),
	messages.ValidContributionProofMetaData(),
	messages.SigValid(),
	messages.PartialSigValid(),
	messages.PartialRootValid(),
	messages.MessageSigner0(),
	messages.SignedMsgSigner0(),

	valcheckduty.WrongValidatorIndex(),
	valcheckduty.WrongValidatorPK(),
	valcheckduty.WrongDutyType(),
	valcheckduty.FarFutureDutySlot(),
	valcheckattestations.Slashable(),
	valcheckattestations.SourceHigherThanTarget(),
	valcheckattestations.FarFutureTarget(),
	valcheckattestations.CommitteeIndexMismatch(),
	valcheckattestations.SlotMismatch(),
	valcheckattestations.ConsensusDataNil(),
	valcheckattestations.Valid(),
}
