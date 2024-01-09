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

	consensusdata.InvalidDuty(),

	consensusdata.ProposerConsensusDataEncoding(),
	consensusdata.BlindedProposerConsensusDataEncoding(),
	consensusdata.CapellaBlockValidation(),
	consensusdata.CapellaBlindedBlockValidation(),
	consensusdata.ProposerNoJustifications(),
	consensusdata.InvalidCapellaBlindedBlockValidation(),
	consensusdata.InvalidCapellaBlockValidation(),

	consensusdata.AttestationConsensusDataEncoding(),
	consensusdata.AttestationValidation(),
	consensusdata.AttesterJustifications(),
	consensusdata.InvalidAttestationValidation(),

	consensusdata.AggregatorConsensusDataEncoding(),
	consensusdata.AggregatorValidation(),
	consensusdata.AggregatorNoJustifications(),
	consensusdata.InvalidAggregatorValidation(),

	consensusdata.SyncCommitteeConsensusDataEncoding(),
	consensusdata.SyncCommitteeJustifications(),
	consensusdata.SyncCommitteeNoJustifications(),
	consensusdata.InvalidSyncCommitteeBlockValidation(),

	consensusdata.SyncCommitteeContributionConsensusDataEncoding(),
	consensusdata.SyncCommitteeContributionValidation(),
	consensusdata.SyncCommitteeContributionNoJustifications(),
	consensusdata.InvalidSyncCommitteeContributionValidation(),

	consensusdata.ValidatorRegistration(),
	consensusdata.VoluntaryExit(),

	consensusdataproposer.VersionedBlockValidation(),
	consensusdataproposer.VersionedBlindedBlockValidation(),
	consensusdataproposer.VersionedBlockUnknownVersion(),
	consensusdataproposer.VersionedBlindedBlockUnknownVersion(),
	consensusdataproposer.VersionedBlockConsensusDataNil(),
	consensusdataproposer.VersionedBlindedBlockConsensusDataNil(),

	ssz.SSZWithdrawalsMarshaling(),

	beacon.DepositData(),
}
