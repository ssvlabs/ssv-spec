package spectest

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/commit"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/messages"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/roundchange"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	messages.CommitDataEncoding(),
	messages.DecidedMsgEncoding(),
	messages.MsgNilIdentifier(),
	messages.MsgNonZeroIdentifier(),
	messages.MsgTypeUnknown(),
	messages.PrepareDataEncoding(),
	messages.ProposeDataEncoding(),
	messages.MsgDataNil(),
	messages.MsgDataNonZero(),
	//messages.RoundChangeDataEncoding(),
	messages.SignedMsgSigTooShort(),
	messages.SignedMsgSigTooLong(),
	messages.SignedMsgNoSigners(),
	messages.GetRoot(),
	messages.SignedMessageEncoding(),

	tests.HappyFlow(),
	tests.SevenOperators(),
	tests.TenOperators(),
	tests.ThirteenOperators(),

	commit.CurrentRound(),
	commit.FutureRound(),
	commit.PastRound(),
	commit.DuplicateMsg(),
	commit.HappyFlow(),
	commit.InvalidCommitData(),
	commit.PostDecided(),
	commit.WrongData(),
	commit.MultiSignerWithOverlap(),
	commit.MultiSignerNoOverlap(),
	commit.Decided(),
	commit.NoPrevAcceptedProposal(),
	commit.WrongHeight(),
	commit.ImparsableCommitData(),
	commit.WrongCommitData(),
	commit.WrongSignature(),

	roundchange.HappyFlow(),
	roundchange.PreviouslyPrepared(),
}
