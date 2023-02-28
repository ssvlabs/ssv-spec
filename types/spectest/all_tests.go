package spectest

import (
	"github.com/bloxapp/ssv-spec/types/spectest/tests/beacon"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/encryption"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	//share.Encoding(),
	//share.HasQuorum(),
	//share.QuorumWithDuplicate(),
	//share.HasQuorum3f1(),
	//share.NoQuorum(),
	//share.NoQuorumDuplicate(),
	//share.HasPartialQuorum(),
	//share.PartialQuorumWithDuplicate(),
	//share.HasPartialQuorum2f1(),
	//share.NoPartialQuorum(),
	//share.NoPartialQuorumDuplicate(),

	encryption.SimpleEncrypt(),
	encryption.EncryptBLSSK(),

	// TODO: consider removing tests below
	//consensusdata.Encoding(),
	//consensusdata.DutyNil(),
	//consensusdata.ValidAttester(),
	//consensusdata.NoAttestationData(),
	//consensusdata.ValidAggregator(),
	//consensusdata.NoAggregatorData(),
	//consensusdata.ValidProposer(),
	//consensusdata.ValidBlindedProposer(),
	//consensusdata.NilBlock(),
	//consensusdata.BlockAndBlindedBlock(),
	//consensusdata.ValidSyncCommitteeAggregator(),
	//consensusdata.EmptySyncCommitteeAggregator(),
	//consensusdata.AttesterJustifications(),
	//consensusdata.AttesterNoValidation(),
	//consensusdata.ProposerJustifications(),
	//consensusdata.ProposerNoJustifications(),
	//consensusdata.AggregatorJustifications(),
	//consensusdata.AggregatorNoJustifications(),
	//consensusdata.SyncCommitteeJustifications(),
	//consensusdata.SyncCommitteeNoJustifications(),
	//consensusdata.SyncCommitteeContributionJustifications(),
	//consensusdata.SyncCommitteeNoContributionJustifications(),
	//consensusdata.ValidatorRegistration(),
	// TODO: add new consensusdata tests

	beacon.DepositData(),
}
