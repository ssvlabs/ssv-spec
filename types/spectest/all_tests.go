package spectest

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/committeemember"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/maxmsgsize"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssz"
	validatorconsensusdata "github.com/ssvlabs/ssv-spec/types/spectest/tests/validatorconsensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/validatorconsensusdata/proposer"
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

	validatorconsensusdata.InvalidDuty(),

	validatorconsensusdata.ProposerConsensusDataEncoding(),
	validatorconsensusdata.BlindedProposerConsensusDataEncoding(),
	validatorconsensusdata.CapellaBlockValidation(),
	validatorconsensusdata.CapellaBlindedBlockValidation(),
	validatorconsensusdata.ProposerNoJustifications(),
	validatorconsensusdata.InvalidCapellaBlindedBlockValidation(),
	validatorconsensusdata.InvalidCapellaBlockValidation(),
	validatorconsensusdata.DenebBlockValidation(),
	validatorconsensusdata.DenebBlindedBlockValidation(),
	validatorconsensusdata.InvalidDenebBlockValidation(),
	validatorconsensusdata.InvalidDenebBlindedBlockValidation(),

	validatorconsensusdata.AggregatorConsensusDataEncoding(),
	validatorconsensusdata.AggregatorValidation(),
	validatorconsensusdata.AggregatorNoJustifications(),
	validatorconsensusdata.InvalidAggregatorValidation(),

	validatorconsensusdata.SyncCommitteeContributionConsensusDataEncoding(),
	validatorconsensusdata.SyncCommitteeContributionValidation(),
	validatorconsensusdata.SyncCommitteeContributionNoJustifications(),
	validatorconsensusdata.InvalidSyncCommitteeContributionValidation(),

	validatorconsensusdata.ValidatorRegistration(),
	validatorconsensusdata.VoluntaryExit(),

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

	maxmsgsize.MaxConsensusData(),
	maxmsgsize.MaxBeaconVote(),
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
