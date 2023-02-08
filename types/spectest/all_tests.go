package spectest

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/spectest/tests/beacon"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/spectest/tests/consensusdata"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/spectest/tests/encryption"
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

	consensusdata.Encoding(),
	consensusdata.DutyNil(),
	consensusdata.ValidAttester(),
	consensusdata.NoAttestationData(),
	consensusdata.ValidAggregator(),
	consensusdata.NoAggregatorData(),
	consensusdata.ValidProposer(),
	consensusdata.ValidBlindedProposer(),
	consensusdata.NilBlock(),
	consensusdata.BlockAndBlindedBlock(),
	consensusdata.ValidSyncCommitteeAggregator(),
	consensusdata.EmptySyncCommitteeAggregator(),

	beacon.DepositData(),
}
