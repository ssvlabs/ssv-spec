package spectest

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/consensus/attester"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/consensus/proposer"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/consensus/synccommittee"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/consensus/synccommitteecontribution"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/processmsg/preconsensus/randao"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	//randao.BaseTests(),
	randao.WrongSigner(),
	randao.UnknownRandaoSigner(),
	randao.UnknownSigner(),
	randao.WrongRandaoSigner(),
	randao.WrongRandaoRoot(),
	randao.MsgInvalid(),
	randao.DutyFinished(),
	randao.NoRunningDuty(),
	randao.PostQuorumMsg(),
	randao.ValidQuorum(),
	randao.Valid7Quorum(),
	randao.Valid10Quorum(),
	randao.Valid13Quorum(),
	randao.WrongSlot(),
	randao.MultiSigs(),
	randao.NoSigs(),
	randao.DuplicateMsgs(),
	//
	//////postconsensus.ValidMessage(),
	//////postconsensus.InvaliSignature(),
	//////postconsensus.WrongSigningRoot(),
	//////postconsensus.WrongBeaconChainSig(),
	//////postconsensus.FutureConsensusState(),
	//////postconsensus.PastConsensusState(),
	//////postconsensus.MsgAfterReconstruction(),
	//////postconsensus.DuplicateMsg(),
	////
	////messages.EncodingAndRoot(),
	////messages.NoMsgs(),
	////messages.InvalidMsg(),
	////messages.InvalidContributionProofMetaData(),
	////messages.ValidContributionProofMetaData(),
	////messages.SigValid(),
	////messages.SigTooShort(),
	////messages.SigTooLong(),
	////messages.PartialSigValid(),
	////messages.PartialSigTooShort(),
	////messages.PartialSigTooLong(),
	////messages.PartialRootValid(),
	////messages.PartialRootTooShort(),
	////messages.PartialRootTooLong(),
	////
	//////valcheck.WrongDutyPubKey(),
	////
	attester.HappyFlow(),
	//attester.SevenOperators(),
	//attester.TenOperators(),
	//attester.ThirteenOperators(),
	//attester.InvalidConsensusMsg(),
	//attester.ValidDecided(),
	//////attestations.FarFutureDuty(),
	//////attestations.DutySlotNotMatchingAttestationSlot(),
	//////attestations.DutyCommitteeIndexNotMatchingAttestations(),
	//////attestations.FarFutureAttestationTarget(),
	//////attestations.AttestationSourceValid(),
	//////attestations.DutyTypeWrong(),
	//////attestations.AttestationDataNil(),
	////
	////processmsg.UnknownRunner(),
	////processmsg.NoRunningDuty(),
	////processmsg.MsgNotBelonging(),
	////processmsg.NoData(),
	////processmsg.ValidDecidedConsensusMsg(),
	////processmsg.UnknownMsgType(),
	////
	proposer.HappyFlow(),
	//proposer.SevenOperators(),
	//
	//aggregator.HappyFlow(),
	//aggregator.SevenOperators(),
	//
	synccommittee.HappyFlow(),
	//synccommittee.SevenOperators(),
	//
	synccommitteecontribution.HappyFlow(),
	////synccommitteecontribution.SevenOperators(),
}
