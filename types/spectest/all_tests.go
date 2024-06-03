package spectest

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/committeemember"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/consensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/consensusdata/proposer"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssz"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	ssvmsg.Encoding(),
	ssvmsg.MsgIDBelongs(),
	ssvmsg.MsgIDDoesntBelongs(),

	partialsigmessage.Encoding(),
	partialsigmessage.InvalidMsg(),
	partialsigmessage.MessageSigner0(),
	partialsigmessage.NoMsgs(),
	partialsigmessage.SigValid(),
	partialsigmessage.PartialSigValid(),
	partialsigmessage.PartialRootValid(),
	partialsigmessage.ValidContributionProofMetaData(),
	partialsigmessage.InconsistentSignedMessage(),

	share.Encoding(),

	committeemember.HasQuorum(),
	committeemember.HasQuorum3f1(),
	committeemember.NoQuorumDuplicate(),
	committeemember.QuorumWithDuplicate(),

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
	consensusdata.DenebBlockValidation(),
	consensusdata.DenebBlindedBlockValidation(),
	consensusdata.InvalidDenebBlockValidation(),
	consensusdata.InvalidDenebBlindedBlockValidation(),

	consensusdata.AggregatorConsensusDataEncoding(),
	consensusdata.AggregatorValidation(),
	consensusdata.AggregatorNoJustifications(),
	consensusdata.InvalidAggregatorValidation(),

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

	signedssvmsg.Encoding(),
	signedssvmsg.Valid(),
	signedssvmsg.NilSSVMessage(),
	signedssvmsg.EmptySignature(),
	signedssvmsg.ZeroSigner(),

	duty.MapAttester(),
	duty.MapProposer(),
	duty.MapAggregator(),
	duty.MapSyncCommittee(),
	duty.MapSyncCommitteeContribution(),
	duty.MapValidatorRegistration(),
	duty.MapVoluntaryExit(),
	duty.MapUnknownRole(),

	beaconvote.BeaconVoteEncoding(),
}
