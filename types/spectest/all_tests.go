package spectest

import (
	"github.com/bloxapp/ssv-spec/types/spectest/tests/beacon"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/consensusdata"
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

	consensusdata.Encoding(),
	beacon.DepositData(),
}
