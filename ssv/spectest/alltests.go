package spectest

import (
	//"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/synccommitteeaggregator"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/messages"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/consensus"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/synccommitteeaggregator"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/preconsensus"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	consensus.ValidDecided(),
	consensus.ValidDecided7Operators(),
	consensus.ValidDecided10Operators(),
	consensus.ValidDecided13Operators(),

	synccommitteeaggregator.SomeAggregatorQuorum(),
	synccommitteeaggregator.NoneAggregatorQuorum(),
	synccommitteeaggregator.AllAggregatorQuorum(),

	preconsensus.WrongExpectedRootsCount(),
	preconsensus.UnorderedExpectedRoots(),
	preconsensus.MultiBeaconSigsWrongSlot(),
	preconsensus.InvalidSignedMessage(),
	preconsensus.InvalidExpectedRoot(),
	preconsensus.DuplicateMsg(),
	preconsensus.PostFinish(),
	preconsensus.PostDecided(),
	preconsensus.PostQuorum(),
	preconsensus.Quorum(),
	preconsensus.Quorum7Operators(),
	preconsensus.Quorum10Operators(),
	preconsensus.Quorum130Operators(),
	preconsensus.ValidMessage(),
	preconsensus.ValidMessage7Operators(),
	preconsensus.ValidMessage10Operators(),
	preconsensus.ValidMessage13Operators(),
	preconsensus.UnknownBeaconSigner(),
	preconsensus.UnknownSigner(),
	preconsensus.InvalidBeaconSignature(),
	preconsensus.InvalidMessageSignature(),

	messages.EncodingAndRoot(),
	messages.NoMsgs(),
	messages.InvalidMsg(),
	messages.ValidContributionProofMetaData(),
	messages.SigValid(),
	messages.SigTooShort(),
	messages.SigTooLong(),
	messages.PartialSigValid(),
	messages.PartialSigTooShort(),
	messages.PartialSigTooLong(),
	messages.PartialRootValid(),
	messages.PartialRootTooShort(),
	messages.PartialRootTooLong(),
}
