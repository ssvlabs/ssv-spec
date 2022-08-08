package spectest

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/consensus/aggregator"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/consensus/attester"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/consensus/proposer"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/consensus/synccommittee"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/consensus/synccommitteecontribution"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/messages"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	//postconsensus.ValidMessage(),
	//postconsensus.InvaliSignature(),
	//postconsensus.WrongSigningRoot(),
	//postconsensus.WrongBeaconChainSig(),
	//postconsensus.FutureConsensusState(),
	//postconsensus.PastConsensusState(),
	//postconsensus.MsgAfterReconstruction(),
	//postconsensus.DuplicateMsg(),

	messages.NoMsgs(),
	messages.InvalidMsg(),
	messages.InvalidContributionProofMetaData(),
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
	//messages.NoMessageSigners(),
	//messages.MultipleMessageSigners(),
	//messages.WrongMsgID(),
	//messages.UnknownMsgType(),
	//messages.NoDutyRunner(),

	//valcheck.WrongDutyPubKey(),

	attester.HappyFlow(),
	attester.SevenOperators(),
	//attestations.FarFutureDuty(),
	//attestations.DutySlotNotMatchingAttestationSlot(),
	//attestations.DutyCommitteeIndexNotMatchingAttestations(),
	//attestations.FarFutureAttestationTarget(),
	//attestations.AttestationSourceValid(),
	//attestations.DutyTypeWrong(),
	//attestations.AttestationDataNil(),
	//
	//processmsg.NoData(),
	//processmsg.InvalidConsensusMsg(),
	//processmsg.InvalidDecidedMsg(),
	//processmsg.InvalidPostConsensusMsg(),
	//processmsg.UnknownType(),
	//processmsg.WrongPubKey(),
	//processmsg.WrongBeaconType(),

	proposer.HappyFlow(),
	proposer.SevenOperators(),

	aggregator.HappyFlow(),
	aggregator.SevenOperators(),

	synccommittee.HappyFlow(),
	synccommittee.SevenOperators(),

	synccommitteecontribution.HappyFlow(),
	synccommitteecontribution.SevenOperators(),
}
