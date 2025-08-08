package testdoc

// DepositDataSpecTest
const DepositDataSpecTestType = "Beacon deposit data: validation of validator registration and signing root generation"

// Documentation for DepositDataSpecTest
const DepositDataTestDoc = "Generate a deposit data with the given validator public key, withdrawal credentials, and fork version, and verify the signing root"

// BeaconVoteEncodingSpecTest
const BeaconVoteEncodingSpecTestType = "Beacon vote encoding: validation of beacon vote encoding"

// Documentation for BeaconVoteEncodingSpecTest
const BeaconVoteEncodingTestDoc = "Test encoding and decoding of BeaconVote with hash tree root verification"

// CommitteeMemberTest
const CommitteeMemberTestType = "Committee member: validation of committee member quorum requirements"

// Documentation for CommitteeMemberTest
const HasQuorum3f1TestDoc = "Test that committee member has quorum with 3f+1 unique signers (full committee) should be accepted"
const HasQuorumTestDoc = "Test that committee member has quorum with 2f+1 unique signers should be accepted"
const NoQuorumDuplicateTestDoc = "Test that committee member has no quorum when signers include duplicates should be rejected"
const QuorumWithDuplicateTestDoc = "Test that committee member has quorum with 2f+1 unique signers even when duplicates are present should be accepted"

// DutySpecTest
const DutySpecTestType = "Duty: validation of beacon role to runner role conversion"

// Documentation for DutySpecTest
const MapAggregatorTestDoc = "Test mapping of BNRoleAggregator"
const MapAttesterTestDoc = "Test mapping of BNRoleAttester"
const MapProposerTestDoc = "Test mapping of BNRoleProposer"
const MapSyncCommitteeContributionTestDoc = "Test mapping of BNRoleSyncCommitteeContribution"
const MapSyncCommitteeTestDoc = "Test mapping of BNRoleSyncCommittee"
const MapUnknownRoleTestDoc = "Test mapping of an unknown role"
const MapValidatorRegistrationTestDoc = "Test mapping of BNRoleValidatorRegistration"
const MapVoluntaryExitTestDoc = "Test mapping of BNRoleVoluntaryExit"

// EncryptionSpecTest
const EncryptionSpecTestType = "Encryption: validation of public/private key operations and message encryption/decryption"

// Documentation for EncryptionSpecTest
const EncryptBLSSKTestDoc = "Test encryption of BLS secret key using RSA key pair"
const EncryptSimpleTestDoc = "Test simple RSA encryption of plain text message"

// MaxMsgSizeTest
const MaxMsgSizeTestType = "Max message size: validation of message size"

// Documentation for MaxMsgSizeTest
const MaxMsgSizeTestMaxConsensusDataDoc = "Test the maximum size of validator consensus data with maximum SSZ data"
const MaxMsgSizeTestMaxBeaconVoteDoc = "Test the maximum size of a beacon vote with source and target checkpoints"
const MaxMsgSizeTestMaxQBFTMessageWithNoJustificationDoc = "Test the maximum size of a QBFT message with no justifications"
const MaxMsgSizeTestMaxQBFTMessageWith1JustificationDoc = "Test the maximum size of a QBFT message with 1 justification"
const MaxMsgSizeTestMaxQBFTMessageWith2JustificationDoc = "Test the maximum size of a QBFT message with 2 justifications"

// StructureSizeTest
const StructureSizeTestType = "Structure size: validation of structure size"

// Documentation for StructureSizeTest
const StructureSizeTestExpectedPartialSignatureMessageDoc = "Test the expected size of a single partial signature message"
const StructureSizeTestExpectedPartialSignatureMessagesDoc = "Test the expected size of partial signature messages collection with 1 message"
const StructureSizeTestExpectedPrepareQBFTMessageDoc = "Test the expected size of a prepare QBFT message with no justifications"
const StructureSizeTestExpectedCommitQBFTMessageDoc = "Test the expected size of a commit QBFT message with no justifications"
const StructureSizeTestExpectedRoundChangeQBFTMessageDoc = "Test the expected size of a round change QBFT message with 3 round change justifications"
const StructureSizeTestExpectedProposalQBFTMessageDoc = "Test the expected size of a proposal QBFT message with 3 round change and 3 prepare justifications"
const StructureSizeTestExpectedPrepareSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing a prepare QBFT message"
const StructureSizeTestExpectedCommitSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing a commit QBFT message"
const StructureSizeTestExpectedDecidedSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing a decided QBFT message with 3 signers"
const StructureSizeTestExpectedRoundChangeSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing a round change QBFT message with 3 justifications"
const StructureSizeTestExpectedProposalSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing a proposal QBFT message with 3 justifications and full data"
const StructureSizeTestExpectedPartialSignatureSignedSSVMessageDoc = "Test the expected size of a SignedSSVMessage containing partial signature messages with full data"
const StructureSizeTestMaxBeaconVoteDoc = "Test the maximum size of a beacon vote with source and target checkpoints"
const StructureSizeTestMaxConsensusDataDoc = "Test the maximum size of a validator consensus data with maximum SSZ data"
const StructureSizeTestMaxPartialSignatureMessageDoc = "Test the maximum size of a single partial signature message"
const StructureSizeTestMaxPartialSignatureMessagesDoc = "Test the maximum size of partial signature messages collection"
const StructureSizeTestMaxPartialSignatureMessagesForPreConsensusDoc = "Test the maximum size of partial signature messages for pre-consensus phase"
const StructureSizeTestMaxQBFTMessageWithNoJustificationDoc = "Test the maximum size of a QBFT message with no justifications"
const StructureSizeTestMaxQBFTMessageWith1JustificationDoc = "Test the maximum size of a QBFT message with 1 justification"
const StructureSizeTestMaxQBFTMessageWith2JustificationDoc = "Test the maximum size of a QBFT message with 2 justifications"
const StructureSizeTestMaxSignedSSVMessageFromQBFTWithNoJustificationDoc = "Test the maximum size of a SignedSSVMessage containing a QBFT message with no justifications"
const StructureSizeTestMaxSignedSSVMessageFromQBFTWith1JustificationDoc = "Test the maximum size of a SignedSSVMessage containing a QBFT message with 1 justification"
const StructureSizeTestMaxSignedSSVMessageFromQBFTWith2JustificationDoc = "Test the maximum size of a SignedSSVMessage containing a QBFT message with 2 justifications and full data"
const StructureSizeTestMaxSignedSSVMessageFromPartialSignatureMessagesDoc = "Test the maximum size of a SignedSSVMessage containing partial signature messages"
const StructureSizeTestMaxSSVMessageFromQBFTMessageDoc = "Test the maximum size of an SSVMessage containing a QBFT message with 2 justifications"
const StructureSizeTestMaxSSVMessageFromPartialSignatureMessagesDoc = "Test the maximum size of an SSVMessage containing partial signature messages"

// ShareTest
const ShareTestType = "Share: testing message signing and quorum verification"

// ShareTestType has no tests

// SSVMessageTest
const SSVMessageTestType = "SSV message validation: testing message ID ownership and routing logic"

// Documentation for SSVMessageTest
const SSVMessageTestBelongsDoc = "Test that message IDs with matching validator public key belong to the validator"
const SSVMessageTestDoesNotBelongDoc = "Test that message IDs with non-matching validator public key do not belong to the validator"

// SignedSSVMessageTest
const SignedSSVMessageTestType = "Signed SSV message: validation of signed SSV message"

// Documentation for SignedSSVMessageTest
const SignedSSVMessageTestValidDoc = "Test validation of a valid signed SSV message with proper RSA signature"
const SignedSSVMessageTestEmptySignatureDoc = "Test validation error for signed SSV message with empty signature"
const SignedSSVMessageTestSignersAndSignaturesWithDifferentLengthDoc = "Test validation error for signed SSV message with different number of signers and signatures"
const SignedSSVMessageTestNonUniqueSignerDoc = "Test validation error for signed SSV message with non-unique signers"
const SignedSSVMessageTestNoSignaturesDoc = "Test validation error for signed SSV message with no signatures"
const SignedSSVMessageTestZeroSignerDoc = "Test validation error for signed SSV message with zero signer ID"
const SignedSSVMessageTestNoSignersDoc = "Test validation error for signed SSV message with no signers"
const SignedSSVMessageTestNilSSVMessageDoc = "Test validation error for signed SSV message with nil SSVMessage"

// MsgSpecTest
const MsgSpecTestType = "Partial signature messages: validation of partial signature messages"

// Documentation for MsgSpecTest
const MsgSpecTestInvalidMsgDoc = "Test validation error when partial signature messages contain invalid message with inconsistent signers"
const MsgSpecTestInconsistentSignedMessageDoc = "Test validation error when signed partial signature message contains messages from different signers"
const MsgSpecTestMessageSigner0Doc = "Test validation error when partial signature message has signer ID 0 which is not allowed"
const MsgSpecTestNoMsgsDoc = "Test validation error when partial signature messages contain no messages"
const MsgSpecTestPartialRootValidDoc = "Test validation of partial signature message with 32-byte signing root"
const MsgSpecTestPartialSigValidDoc = "Test validation of partial signature message with 96-byte signature length"
const MsgSpecTestSigValidDoc = "Test validation of signed post consensus message with 96-byte signature length"
const MsgSpecTestValidContributionProofMetaDataDoc = "Test validation of partial signature message with contribution proof metadata type"

// PartialSignatureMessageEncodingTest
const PartialSignatureMessageEncodingTestType = "Partial signature messages encoding: validation of partial signature messages encoding"

// Documentation for PartialSignatureMessageEncodingTest
const PartialSignatureMessageEncodingTestDoc = "Test encoding and decoding of partial signature messages with hash tree root verification"

// SSZSpecTest
const SSZSpecTestType = "SSZ: validation of SSZ encoding and decoding"

// Documentation for SSZSpecTest
const SSZSpecTestWithdrawalsMarshalingDoc = "Test SSZ marshaling and hash tree root calculation of Capella withdrawals"

// ValidatorConsensusDataTest
const ValidatorConsensusDataTestType = "Validator consensus data: validation of validator consensus data"

// AggregatorConsensusDataTest
const AggregatorConsensusDataTestType = "Aggregator consensus data: validation of aggregator committee consensus data"

// Documentation for ValidatorConsensusDataTest
const ValidatorConsensusDataTestInvalidDenebBlockDoc = "Test validation error for invalid consensus data with empty Deneb block data"
const ValidatorConsensusDataTestInvalidDenebBlindedBlockDoc = "Test validation error for invalid consensus data with empty Deneb blinded block data"
const ValidatorConsensusDataTestInvalidCapellaBlockDoc = "Test validation error for invalid consensus data with empty Capella block data"
const ValidatorConsensusDataTestInvalidCapellaBlindedBlockDoc = "Test validation error for invalid consensus data with empty Capella blinded block data"
const ValidatorConsensusDataTestInvalidElectraBlockDoc = "Test validation error for invalid consensus data with empty Electra block data"
const ValidatorConsensusDataTestInvalidElectraBlindedBlockDoc = "Test validation error for invalid consensus data with empty Electra blinded block data"
const ValidatorConsensusDataTestDenebBlindedBlockDoc = "Test validation of valid consensus data with Deneb blinded block"
const ValidatorConsensusDataTestCapellaBlindedBlockDoc = "Test validation of valid consensus data with Capella blinded block"
const ValidatorConsensusDataTestCapellaBlockDoc = "Test validation of valid consensus data with Capella block"
const ValidatorConsensusDataTestElectraBlindedBlockDoc = "Test validation of valid consensus data with Electra blinded block"
const ValidatorConsensusDataTestPhase0AggregatorNoJustificationsDoc = "Test phase0 aggregator consensus data with no pre-consensus justifications"
const ValidatorConsensusDataTestElectraAggregatorNoJustificationsDoc = "Test Electra aggregator consensus data with no pre-consensus justifications"
const ValidatorConsensusDataTestPhase0AggregatorValidationDoc = "Test validation of valid consensus data with Phase0 AggregateAndProof"
const ValidatorConsensusDataTestElectraAggregatorValidationDoc = "Test validation of valid consensus data with Electra AggregateAndProof"
const ValidatorConsensusDataTestProposerNoJustificationsDoc = "Test proposer consensus data with no pre-consensus justifications"
const ValidatorConsensusDataTestSyncCommitteeContributionValidationDoc = "Test validation of valid consensus data with sync committee contribution"
const ValidatorConsensusDataTestVoluntaryExitDoc = "Test validation error for voluntary exit consensus data which has no consensus data"
const ValidatorConsensusDataTestValidatorRegistrationDoc = "Test validation error for validator registration consensus data which has no consensus data"
const ValidatorConsensusDataTestDenebBlockDoc = "Test validation of valid consensus data with Deneb block"
const ValidatorConsensusDataTestElectraBlockDoc = "Test validation of valid consensus data with Electra block"
const ValidatorConsensusDataTestInvalidPhase0AggregatorDoc = "Test validation error for invalid consensus data with Phase0 AggregateAndProof using incorrect data"
const ValidatorConsensusDataTestInvalidElectraAggregatorDoc = "Test validation error for invalid consensus data with Electra AggregateAndProof using incorrect data"
const ValidatorConsensusDataTestInvalidDutyDoc = "Test validation error for consensus data with unknown duty role"
const ValidatorConsensusDataTestInvalidSyncCommitteeContributionDoc = "Test validation error for invalid consensus data with sync committee contribution using incorrect data"

// Documentation for AggregatorConsensusDataTest
const AggregatorConsensusDataTestSyncCommitteeContributionNoJustificationsDoc = "Test sync committee contribution consensus data with no sync committee contribution pre-consensus justifications"

// ValidatorConsensusDataEncodingTest
const ValidatorConsensusDataEncodingTestType = "Validator consensus data encoding"

// Documentation for ValidatorConsensusDataEncodingTest
const ValidatorConsensusDataEncodingTestProposerDoc = "Test encoding and decoding of proposer consensus data with Capella blinded block"
const ValidatorConsensusDataEncodingTestBlindedProposerDoc = "Test encoding and decoding of blinded proposer consensus data with Capella blinded block"
const ValidatorConsensusDataEncodingTestPhase0AggregatorDoc = "Test encoding and decoding of phase0 aggregator consensus data"
const ValidatorConsensusDataEncodingTestElectraAggregatorDoc = "Test encoding and decoding of Electra aggregator consensus data"
const ValidatorConsensusDataEncodingTestSyncCommitteeContributionDoc = "Test encoding and decoding of sync committee contribution consensus data"

// ShareEncodingTest
const ShareEncodingTestType = "Share encoding"

// Documentation for ShareEncodingTest
const ShareEncodingTestDoc = "Test encoding and decoding of share with hash tree root verification"

// SSVMessageEncodingTest
const SSVMessageEncodingTestType = "SSV message encoding"

// Documentation for SSVMessageEncodingTest
const SSVMessageEncodingTestDoc = "Test encoding and decoding of SSVMessage with hash tree root verification"

// SignedSSVMessageEncodingTest
const SignedSSVMessageEncodingTestType = "Signed SSV message encoding: validation of signed SSV message encoding"

// Documentation for SignedSSVMessageEncodingTest
const SignedSSVMessageEncodingTestDoc = "Test encoding of a signed SSV message to bytes"

// ProposerSpecTest
const ProposerSpecTestType = "Proposer: validation of proposer consensus data"

// Documentation for ProposerSpecTest
const ProposerSpecTestVersionedBlockValidationDoc = "Test validation of valid consensus data with versioned Deneb block"
const ProposerSpecTestVersionedBlindedBlockValidationDoc = "Test validation of valid consensus data with versioned Deneb blinded block"
const ProposerSpecTestVersionedBlockUnknownVersionDoc = "Test validation error for consensus data with unknown block version"
const ProposerSpecTestVersionedBlindedBlockUnknownVersionDoc = "Test validation error for consensus data with unknown blinded block version"
const ProposerSpecTestVersionedBlockConsensusDataNilDoc = "Test validation error for consensus data with nil block data"
const ProposerSpecTestVersionedBlindedBlockConsensusDataNilDoc = "Test validation error for consensus data with nil blinded block data"
