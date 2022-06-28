package spectest

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/commit"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/messages"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/prepare"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/proposal"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/roundchange"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	proposal.HappyFlow(),
	proposal.NotPreparedPreviouslyJustification(),
	proposal.PreparedPreviouslyJustification(),
	proposal.DifferentJustifications(),
	proposal.JustificationsNotHeighest(),
	proposal.JustificationsValueNotJustified(),
	proposal.DuplicateMsg(),
	proposal.FirstRoundJustification(),
	proposal.FutureRoundNoAcceptedProposal(),
	proposal.FutureRoundAcceptedProposal(),
	proposal.PastRound(),
	proposal.ImparsableProposalData(),
	proposal.InvalidRoundChangeJustificationPrepared(),
	proposal.InvalidRoundChangeJustification(),
	proposal.PreparedPreviouslyNoRCJustificationQuorum(),
	proposal.NoRCJustification(),
	proposal.PreparedPreviouslyNoPrepareJustificationQuorum(),
	proposal.PreparedPreviouslyDuplicatePrepareMsg(),
	proposal.PreparedPreviouslyDuplicateRCMsg(),
	proposal.DuplicateRCMsg(),
	proposal.InvalidPrepareJustificationValue(),
	proposal.InvalidPrepareJustificationRound(),
	proposal.InvalidProposalData(),
	proposal.InvalidValueCheck(),
	proposal.MultiSigner(),
	proposal.PostDecided(),
	proposal.PostPrepared(),
	proposal.SecondProposalForRound(),
	proposal.WrongHeight(),
	proposal.WrongProposer(),
	proposal.WrongSignature(),

	messages.RoundChangeDataInvalidJustifications(),
	messages.RoundChangeDataInvalidPreparedRound(),
	messages.RoundChangeDataInvalidPreparedValue(),
	messages.RoundChangePrePreparedJustifications(),
	messages.RoundChangeNotPreparedJustifications(),
	messages.CommitDataEncoding(),
	messages.DecidedMsgEncoding(),
	messages.MsgNilIdentifier(),
	messages.MsgNonZeroIdentifier(),
	messages.MsgTypeUnknown(),
	messages.PrepareDataEncoding(),
	messages.ProposeDataEncoding(),
	messages.MsgDataNil(),
	messages.MsgDataNonZero(),
	messages.SignedMsgSigTooShort(),
	messages.SignedMsgSigTooLong(),
	messages.SignedMsgNoSigners(),
	messages.GetRoot(),
	messages.SignedMessageEncoding(),
	messages.CreateProposal(),
	messages.CreateProposalPreviouslyPrepared(),
	messages.CreateProposalNotPreviouslyPrepared(),
	messages.CreatePrepare(),
	messages.CreateCommit(),
	messages.CreateRoundChange(),
	messages.CreateRoundChangePreviouslyPrepared(),
	messages.RoundChangeDataEncoding(),

	tests.HappyFlow(),
	tests.SevenOperators(),
	tests.TenOperators(),
	tests.ThirteenOperators(),

	prepare.DuplicateMsg(),
	prepare.HappyFlow(),
	prepare.ImparsableProposalData(),
	prepare.InvalidPrepareData(),
	prepare.MultiSigner(),
	prepare.NoPreviousProposal(),
	prepare.OldRound(),
	prepare.FutureRound(),
	prepare.PostDecided(),
	prepare.WrongData(),
	prepare.WrongHeight(),
	prepare.WrongSignature(),

	commit.CurrentRound(),
	commit.FutureRound(),
	commit.PastRound(),
	commit.DuplicateMsg(),
	commit.HappyFlow(),
	commit.InvalidCommitData(),
	commit.PostDecided(),
	commit.WrongData1(),
	commit.WrongData2(),
	commit.MultiSignerWithOverlap(),
	commit.MultiSignerNoOverlap(),
	commit.Decided(),
	commit.NoPrevAcceptedProposal(),
	commit.WrongHeight(),
	commit.ImparsableCommitData(),
	commit.WrongSignature(),

	roundchange.HappyFlow(),
	roundchange.PreviouslyPrepared(),
	roundchange.F1Speedup(),
	roundchange.F1SpeedupPrepared(),
}
