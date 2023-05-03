package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/types/spectest/tests/beacon"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/consensusdata"
	consensusdataproposer "github.com/bloxapp/ssv-spec/types/spectest/tests/consensusdata/proposer"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/ssz"

	"github.com/bloxapp/ssv-spec/types/spectest/tests/encryption"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/share"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/ssvmsg"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	ssvmsg.Encoding(),

	partialsigmessage.Encoding(),

	share.Encoding(),
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
	consensusdata.Encoding(),
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
	// TODO: add new consensusdata tests

	consensusdataproposer.VersionedBlockValidation(),
	consensusdataproposer.VersionedBlindedBlockValidation(),
	consensusdataproposer.VersionedBlockUnknownVersion(),
	consensusdataproposer.VersionedBlindedBlockUnknownVersion(),
	consensusdataproposer.VersionedBlockConsensusDataNil(),
	consensusdataproposer.VersionedBlindedBlockConsensusDataNil(),

	ssz.SSZWithdrawalsMarshaling(),

	beacon.DepositData(),
}
