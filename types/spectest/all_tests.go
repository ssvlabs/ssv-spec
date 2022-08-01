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
	encryption.SimpleEncrypt(),
	encryption.EncryptBLSSK(),
	
	consensusdata.Encoding(),
	beacon.DepositData(),
}
