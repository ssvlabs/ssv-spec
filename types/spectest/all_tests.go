package spectest

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/aggregatorcommitteeconsensusdata"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/committeemember"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/maxmsgsize"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata/proposer"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
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
	partialsigmessage.InconsistentSignedMessage(),

	share.Encoding(),

	committeemember.HasQuorum(),
	committeemember.HasQuorum3f1(),
	committeemember.NoQuorumDuplicate(),
	committeemember.QuorumWithDuplicate(),

	encryption.SimpleEncrypt(),
	encryption.EncryptBLSSK(),

	proposerconsensusdata.InvalidDuty(),

	proposerconsensusdata.ProposerConsensusDataEncoding(),
	proposerconsensusdata.BlindedProposerConsensusDataEncoding(),
	proposerconsensusdata.CapellaBlockValidation(),
	proposerconsensusdata.CapellaBlindedBlockValidation(),
	proposerconsensusdata.ProposerNoJustifications(),
	proposerconsensusdata.InvalidCapellaBlindedBlockValidation(),
	proposerconsensusdata.InvalidCapellaBlockValidation(),
	proposerconsensusdata.DenebBlockValidation(),
	proposerconsensusdata.DenebBlindedBlockValidation(),
	proposerconsensusdata.InvalidDenebBlockValidation(),
	proposerconsensusdata.InvalidDenebBlindedBlockValidation(),
	proposerconsensusdata.ElectraBlockValidation(),
	proposerconsensusdata.ElectraBlindedBlockValidation(),
	proposerconsensusdata.InvalidElectraBlockValidation(),
	proposerconsensusdata.InvalidElectraBlindedBlockValidation(),
	proposerconsensusdata.FuluBlockValidation(),
	proposerconsensusdata.FuluBlindedBlockValidation(),
	proposerconsensusdata.InvalidFuluBlockValidation(),
	proposerconsensusdata.InvalidFuluBlindedBlockValidation(),

	proposerconsensusdata.WrongDutyTypeValidatorRegistration(),
	proposerconsensusdata.WrongDutyTypeVoluntaryExit(),

	aggregatorcommitteeconsensusdata.Phase0AggregatorConsensusDataEncoding(),
	aggregatorcommitteeconsensusdata.ElectraAggregatorConsensusDataEncoding(),
	aggregatorcommitteeconsensusdata.SyncCommitteeContributionConsensusDataEncoding(),

	aggregatorcommitteeconsensusdata.Phase0AggregatorValidation(),
	aggregatorcommitteeconsensusdata.Phase0AggregatorNoJustifications(),
	aggregatorcommitteeconsensusdata.ElectraAggregatorValidation(),
	aggregatorcommitteeconsensusdata.ElectraAggregatorNoJustifications(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationCommitteeIndexesLength(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationNoValidators(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationDuplicateCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationMissingCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationUnusedCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationPhase0AttestationDecoding(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationElectraAttestationDecoding(),

	aggregatorcommitteeconsensusdata.SyncCommitteeContributionValidation(),
	aggregatorcommitteeconsensusdata.SyncCommitteeContributionNoJustifications(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionDuplicatedSubnet(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionMissingSubnet(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionUnusedSubnet(),

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
	signedssvmsg.NoSigners(),
	signedssvmsg.NoSignatures(),
	signedssvmsg.SignersAndSignaturesWithDifferentLength(),
	signedssvmsg.NonUniqueSigner(),

	duty.MapAttester(),
	duty.MapProposer(),
	duty.MapAggregator(),
	duty.MapSyncCommittee(),
	duty.MapSyncCommitteeContribution(),
	duty.MapValidatorRegistration(),
	duty.MapVoluntaryExit(),
	duty.MapUnknownRole(),

	beaconvote.BeaconVoteEncoding(),

	maxmsgsize.MaxConsensusData(),
	maxmsgsize.MaxBeaconVote(),
	maxmsgsize.MaxAggregatorCommitteeConsensusData(),
	maxmsgsize.MaxPhase0Attestation(),
	maxmsgsize.MaxElectraAttestation(),
	maxmsgsize.MaxQBFTMessageWithNoJustification(),
	maxmsgsize.MaxQBFTMessageWith1Justification(),
	maxmsgsize.MaxQBFTMessageWith2Justification(),
	maxmsgsize.MaxPartialSignatureMessage(),
	maxmsgsize.MaxPartialSignatureMessages(),
	maxmsgsize.MaxPartialSignatureMessagesForPreConsensus(),
	maxmsgsize.MaxSSVMessageFromQBFTMessage(),
	maxmsgsize.MaxSSVMessageFromPartialSignatureMessage(),
	maxmsgsize.MaxSignedSSVMessageFromQBFTMessageWithNoJustification(),
	maxmsgsize.MaxSignedSSVMessageFromQBFTMessageWith1Justification(),
	maxmsgsize.MaxSignedSSVMessageFromQBFTMessageWith2Justification(),
	maxmsgsize.MaxSignedSSVMessageFromPartialSignatureMessages(),

	maxmsgsize.ExpectedPrepareQBFTMessage(),
	maxmsgsize.ExpectedCommitQBFTMessage(),
	maxmsgsize.ExpectedRoundChangeQBFTMessage(),
	maxmsgsize.ExpectedProposalQBFTMessage(),

	maxmsgsize.ExpectedPartialSignatureMessage(),
	maxmsgsize.ExpectedPartialSignatureMessages(),

	maxmsgsize.ExpectedPrepareSignedSSVMessage(),
	maxmsgsize.ExpectedCommitSignedSSVMessage(),
	maxmsgsize.ExpectedDecidedSignedSSVMessage(),
	maxmsgsize.ExpectedRoundChangeSignedSSVMessage(),
	maxmsgsize.ExpectedProposalSignedSSVMessage(),
	maxmsgsize.ExpectedPartialSignatureSignedSSVMessage(),
}
