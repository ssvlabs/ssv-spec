package spectest

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/blindedexecutionpayloadenvelope"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedbeaconblock"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedexecutionpayloadbid"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/aggregatorcommitteeconsensusdata"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/builderrequestauth"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/committeemember"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/maxmsgsize"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/payloadattestationdata"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/payloadattestationmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata/proposer"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerpreferences"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedbuilderrequestauth"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedproposerpreferences"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
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
	proposerconsensusdata.GloasBlockValidation(),
	proposerconsensusdata.InvalidGloasBlockValidation(),
	proposerconsensusdata.NonGloasVersionValidation(),
	proposerconsensusdata.ProposerNoJustifications(),

	proposerconsensusdata.WrongDutyTypeVoluntaryExit(),

	aggregatorcommitteeconsensusdata.AggregatorConsensusDataEncoding(),
	aggregatorcommitteeconsensusdata.SyncCommitteeContributionConsensusDataEncoding(),

	aggregatorcommitteeconsensusdata.AggregatorValidation(),
	aggregatorcommitteeconsensusdata.AggregatorNoJustifications(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationCommitteeIndexesLength(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationNoValidators(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationDuplicateCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationMissingCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationUnusedCommitteeIndex(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationAttestationDecoding(),
	aggregatorcommitteeconsensusdata.InvalidAggregatorValidationNonGloasVersion(),

	aggregatorcommitteeconsensusdata.SyncCommitteeContributionValidation(),
	aggregatorcommitteeconsensusdata.SyncCommitteeContributionNoJustifications(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionDuplicatedSubnet(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionMissingSubnet(),
	aggregatorcommitteeconsensusdata.InvalidSyncCommitteeContributionUnusedSubnet(),

	consensusdataproposer.VersionedBlockValidation(),
	consensusdataproposer.VersionedBlockUnknownVersion(),
	consensusdataproposer.VersionedBlockConsensusDataNil(),

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
	duty.MapVoluntaryExit(),
	duty.MapPTCAttester(),
	duty.MapProposerPreferences(),
	duty.MapUnknownRole(),

	beaconvote.BeaconVoteEncoding(),
	payloadattestationdata.PayloadAttestationDataEncoding(),
	payloadattestationmessage.PayloadAttestationMessageEncoding(),
	signedbeaconblock.SignedBeaconBlockEncoding(),
	signedbeaconblock.SignedBeaconBlockDevnet6Encoding(),
	signedexecutionpayloadbid.SignedExecutionPayloadBidEncoding(),
	blindedexecutionpayloadenvelope.BlindedExecutionPayloadEnvelopeEncoding(),
	proposerpreferences.ProposerPreferencesEncoding(),
	signedproposerpreferences.SignedProposerPreferencesEncoding(),
	builderrequestauth.BuilderRequestAuthEncoding(),
	signedbuilderrequestauth.SignedBuilderRequestAuthEncoding(),

	maxmsgsize.MaxConsensusData(),
	maxmsgsize.MaxBeaconVote(),
	maxmsgsize.MaxAggregatorCommitteeConsensusData(),
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
