package spectest

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/commit"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
	tests.SevenOperators(),
	tests.TenOperators(),
	tests.ThirteenOperators(),

	commit.HappyFlow(),
	commit.MultiSignerWithOverlap(),
	commit.MultiSignerNoOverlap(),
	commit.Decided(),
	commit.NoPrevAcceptedProposal(),
	commit.WrongHeight(),
	commit.WrongRound(),
	commit.ImparsableCommitData(),
	commit.WrongCommitData(),
	commit.WrongSignature(),
}
