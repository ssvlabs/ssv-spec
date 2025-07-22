package testdoc

// Test type strings for ssv/spectest/tests
const MsgProcessingSpecTestType = "SSV message processing: validation of message handling and state transitions"
const MultiMsgProcessingSpecTestType = "SSV multi message processing: multiple message processing tests"

// Test type strings for ssv/spectest/tests/committee
const CommitteeSpecTestType = "Committee: validation of committee member selection and quorum requirements"
const MultiCommitteeSpecTestType = "Multi committee: multiple committee tests"

// Test type strings for ssv/spectest/tests/partialsigcontainer
const PartialSigContainerTestType = "Partial signature container: validation of signature aggregation and quorum verification"

// Test type strings for ssv/spectest/tests/valcheck
const ValCheckSpecTestType = "Value check: validations for input of different runner roles"
const MultiValCheckSpecTestType = "Multi value check: multiple value check tests"

// Test type strings for ssv/spectest/tests/runner/construction
const RunnerConstructionSpecTestType = "Runner construction: validation of runner construction"

// Test type strings for ssv/spectest/tests/runner/duties
const RunnerDutiesSpecTestType = "Start new runner duty: validation of new duties"
const MultiRunnerDutiesSpecTestType = "Multi start new runner duty: multiple start new runner duty tests"

// Test type strings for ssv/spectest/tests/runner/duties/synccommitteeaggregator
const SyncCommitteeAggregatorProofSpecTestType = "Sync committee aggregator proof: validation of sync committee aggregator proof"

// Documentation for happy flow tests
const HappyFlowDoc = "Tests a full runner happy flow"

// Documentation for committee/multipleduty tests
const CommitteeSequencedHappyFlowDutiesDoc = "Tests complete happy flow execution for sequences of duties including consensus and post-consensus phases"
const CommitteeShuffledDecidedDutiesDoc = "Tests decision phase for sequences of duties with shuffled input messages while preserving order between duty messages"
const CommitteeShuffledHappyFlowDutiesWithDifferentValidatorsDoc = "Tests complete happy flow execution for sequences of duties with shuffled input messages while preserving order between duty messages for different validators"
const CommitteeShuffledHappyFlowDutiesWithSameValidatorsDoc = "Tests complete happy flow execution for sequences of duties with shuffled input messages while preserving order between duty messages for the same validators"
const CommitteeFailedThanSuccessfulDutiesDoc = "Tests sequences where some duties fail (no post-consensus) followed by successful duties with complete happy flow"
const CommitteeSequencedDecidedDutiesDoc = "Tests decision phase for sequences of duties without completing post-consensus"

// Documentation for committee/singleduty tests
const CommitteeValidBeaconVoteDoc = "Tests committee behavior when processing proposal messages with valid beacon vote data"
const CommitteeWrongBeaconVoteDoc = "Tests committee behavior when processing proposal messages with invalid beacon vote data (source >= target)"
const CommitteeWrongMessageIDDoc = "Tests committee behavior when processing messages with incorrect message IDs that don't match the committee ID"
const CommitteePastMsgDutyNotFinishedDoc = "Tests committee behavior when processing past proposal messages for a duty that has not finished (consensus already finished)"
const CommitteeProposalWithConsensusDataDoc = "Tests committee behavior when processing proposal messages with ValidatorConsensusData instead of BeaconVote objects"
const CommitteeStartDutyDoc = "Tests basic duty starting for attestations and sync committees without consensus messages"
const CommitteeStartNoDutyDoc = "Tests committee behavior when starting with an empty duty (no attestation or sync committee duties)"
const CommitteeStartWithNoSharesForDutyDoc = "Tests committee behavior when starting a duty for validators that the committee doesn't have shares for"
const CommitteeHappyFlowDoc = "Tests complete duty execution flow including consensus and post-consensus phases for attestations and sync committees"
const CommitteeMissingSomeSharesDoc = "Tests complete duty execution for a committee that only has shares for a fraction of the duty's validators"
const CommitteePastMsgDutyDoesNotExistDoc = "Tests committee behavior when processing past proposal messages for a duty that doesn't exist"
const CommitteePastMsgDutyFinishedDoc = "Tests committee behavior when processing past messages for a duty that has finished"
const CommitteeDecidedDoc = "Tests committee runner decision phase for attestations and sync committees without completing post-consensus"

// Documentation for dutyexe tests
const DutyExeWrongDutyPubKeyDoc = "Tests decided value with duty validator pubkey != the duty runner's pubkey"
const DutyExeWrongDutyRoleDoc = "Tests behavior when processing decided value duty with wrong duty role (!= duty runner role)"

// Documentation for partialsigcontainer tests
const PartialSigContainerDuplicateDoc = "Tests partial signature container with duplicate signatures (below quorum)"
const PartialSigContainerDuplicateQuorumDoc = "Tests partial signature container with duplicate signatures but still achieving quorum"
const PartialSigContainerInvalidDoc = "Tests partial signature container with invalid signatures"
const PartialSigContainerOneSignatureDoc = "Tests partial signature container with only one signature (below quorum)"
const PartialSigContainerQuorumDoc = "Tests partial signature container quorum and signature reconstruction"

// Documentation for runner/consensus tests
const ConsensusValidDecided7OperatorsDoc = "Tests valid consensus decided message processing with 7 operators"
const ConsensusValidMessageDoc = "Tests valid consensus message processing across different runner types"
const ConsensusZeroSignerDoc = "Tests consensus message processing with signer ID 0, expecting error"
const ConsensusValidDecidedDoc = "Tests valid consensus decided message processing across different runner types"
const ConsensusValidDecided10OperatorsDoc = "Tests valid consensus decided message processing with 10 operators"
const ConsensusValidDecided13OperatorsDoc = "Tests valid consensus decided message processing with 13 operators"
const ConsensusNoSignaturesDoc = "Tests consensus message processing with no signatures, expecting error"
const ConsensusNoSignersDoc = "Tests consensus message processing with no signers, expecting error"
const ConsensusNonUniqueSignersDoc = "Tests consensus message processing with non-unique signers, expecting error"
const ConsensusPostFinishDoc = "Tests consensus message processing after duty is finished"
const ConsensusDiffLengthSignersSignaturesDoc = "Tests consensus message processing with different length of signers and signatures, expecting error"
const ConsensusEmptySignatureDoc = "Tests consensus message processing with empty signature, expecting error"
const ConsensusNilSSVMessageDoc = "Tests consensus message processing with nil SSV message, expecting error"
const ConsensusInvalidSignatureDoc = "Tests consensus message processing with invalid signatures"
const ConsensusPastMessageDoc = "Tests consensus message processing with past messages"
const ConsensusPostDecidedDoc = "Tests consensus message processing after duty is already decided"
const ConsensusFutureDecidedNoInstanceDoc = "Tests consensus decided message processing with future messages when no instance exists"
const ConsensusFutureMessageDoc = "Tests consensus message processing with future messages"
const ConsensusInvalidDecidedValueDoc = "Tests consensus message processing with invalid decided values"
const ConsensusDecidedSlashableAttestationDoc = "Test that attempting to sign a slashable attestation results in an error"
const ConsensusFutureDecidedDoc = "Tests consensus decided message processing with future messages"

// Documentation for runner/construction tests
const RunnerConstructionManySharesDoc = "Test that only committee runner can be constructed with multiple shares"
const RunnerConstructionNoSharesDoc = "Test that all runners cannot be constructed without shares"
const RunnerConstructionOneShareDoc = "Test that all runners can be constructed with one share"

// Documentation for runner/duties/newduty tests
const NewDutyPostWrongDecidedDoc = "Tests new duty start after a wrong decided value, expecting error"
const NewDutyValidDoc = "Tests valid new duty start scenarios across different runner types"
const NewDutyFinishedDoc = "Tests new duty start after a previous duty has finished"
const NewDutyFirstHeightDoc = "Tests new duty start at first height"
const NewDutyNotDecidedDoc = "Tests new duty start when a previous duty has not been decided"
const NewDutyPostFutureDecidedDoc = "Tests new duty start after a future decided value, expecting error"
const NewDutyPostInvalidDecidedDoc = "Tests new duty start after an invalid decided value, expecting error"
const NewDutyConsensusNotStartedDoc = "Tests new duty start when consensus has not started"
const NewDutyPostDecidedDoc = "Tests new duty start after a previous duty has been decided"
const NewDutyDuplicateDutyFinishedDoc = "Tests new duty start with duplicate duty that has finished"
const NewDutyDuplicateDutyNotFinishedDoc = "Tests new duty start with duplicate duty that has not finished"

// Documentation for runner/duties/proposer tests
const ProposerBlindedReceivingNormalBlockDoc = "Tests full happy flow for a blinded proposer runner that accepts a normal block proposal"
const ProposerNormalReceivingBlindedBlockDoc = "Tests full happy flow for a normal proposer runner that accepts a blinded block proposal"
const ProposeBlindedBlockDecidedRegularDoc = "Tests proposing a blinded block but the decided block is a regular block"
const ProposeRegularBlockDecidedBlindedDoc = "Tests proposing a regular block but the decided block is a blinded block"

// Documentation for runner/duties/synccommitteeaggregator tests
const SyncCommitteeAggregatorProofAllAggregatorDoc = "Tests sync committee aggregator proof validation when all selection proofs are aggregators"
const SyncCommitteeAggregatorProofNoneAggregatorDoc = "Tests sync committee aggregator proof validation when none of the selection proofs are aggregators"
const SyncCommitteeAggregatorProofSomeAggregatorDoc = "Tests sync committee aggregator proof validation when some selection proofs are aggregators"

// Documentation for runner/postconsensus tests
const PostConsensusValidMsgDoc = "Tests valid post-consensus message processing with 4 operators"
const PostConsensusQuorumDoc = "Tests post-consensus quorum message processing across different runner types"
const PostConsensusInvalidMsgDoc = "Tests post-consensus message processing with invalid message, expecting error"
const PostConsensusDuplicateMsgDoc = "Tests post-consensus message processing with duplicate SignedPartialSignatureMessages from the same signer, expecting correct handling of duplicates"
const PostConsensusInvalidMsgSigDoc = "Tests post-consensus message processing with invalid message signature. No error is generated since the SignedPartialSignatureMessage.Signature is no longer checked"
const PostConsensusInvalidMsgSlotDoc = "Tests post-consensus message processing with invalid message slot, expecting error"
const PostConsensusInvalidSignedMsgDiffLengthDoc = "Tests post-consensus message processing with invalid signed message (different signature lengths)"
const PostConsensusInvalidSignedMsgEmptySignatureDoc = "Tests post-consensus message processing with invalid signed message (empty signature)"
const PostConsensusInvalidSignedMsgNoSignatureDoc = "Tests post-consensus message processing with invalid signed message (no signatures)"
const PostConsensusInvalidSignedMsgNoSignersDoc = "Tests post-consensus message processing with invalid signed message (no signers)"
const PostConsensusInvalidThenQuorumDoc = "Tests post-consensus quorum formation with an invalid message followed by a valid quorum, expecting error then successful termination"
const PostConsensusInvalidValidatorIndexDoc = "Tests post-consensus message processing with invalid validator indexes, expecting error"
const PostConsensusInvalidValidatorIndexQuorumDoc = "Tests post-consensus message processing with invalid validator indexes in quorum, expecting error"
const PostConsensusInconsistentBeaconSignerDoc = "Tests post-consensus message processing with inconsistent beacon signer, expecting error"
const PostConsensusValidMsg7OperatorsDoc = "Tests valid post-consensus message processing with 7 operators"
const PostConsensusValidMsg10OperatorsDoc = "Tests valid post-consensus message processing with 10 operators"
const PostConsensusValidMsg13OperatorsDoc = "Tests valid post-consensus message processing with 13 operators"
const PostConsensusQuorum7OperatorsDoc = "Tests post-consensus quorum message processing with 7 operators"
const PostConsensusQuorum10OperatorsDoc = "Tests post-consensus quorum message processing with 10 operators"
const PostConsensusQuorum13OperatorsDoc = "Tests post-consensus quorum message processing with 13 operators"
const PostConsensusTooFewRootsDoc = "Tests post-consensus message processing with too few roots, expecting error handling for missing roots."
const PostConsensusTooManyRootsDoc = "Tests post-consensus message processing with too many roots, expecting error handling for excess roots."
const PostConsensusUnorderedExpectedRootsDoc = "Tests post-consensus message processing with unordered expected roots, expecting correct handling regardless of order."
const PostConsensusPostQuorumDoc = "Tests post-consensus message processing after quorum is reached, expecting error"
const PostConsensusDuplicateMsgDiffRootThenQuorumDoc = "Tests post-consensus message processing where duplicate messages from the same signer have different roots, followed by a quorum, expecting correct error and recovery handling."
const PostConsensusDuplicateMsgDifferentRootsDoc = "Tests post-consensus message processing with duplicate SignedPartialSignatureMessages from the same signer but with different signing roots, expecting error"
const PostConsensusInconsistentOperatorSignerDoc = "Tests post-consensus message processing with inconsistent operator signer, expecting error"
const PostConsensusInvalidAndValidValidatorIndexesQuorumDoc = "Tests post-consensus message processing with a mix of invalid and valid validator indexes in quorum, expecting correct error and recovery handling"
const PostConsensusInvalidBeaconSignatureInQuorumDoc = "Tests post-consensus message processing with invalid beacon signature in quorum, expecting error"
const PostConsensusInvalidDecidedValueDoc = "Tests post-consensus message processing with invalid decided value, expecting error"
const PostConsensusInvalidExpectedRootDoc = "Tests post-consensus message processing with invalid expected root, expecting error"
const PostConsensusInvalidOperatorSignatureDoc = "Tests post-consensus message processing with invalid operator signature, expecting error"
const PostConsensusInvalidQuorumThenValidQuorumDoc = "Tests post-consensus message processing with an invalid quorum followed by a valid quorum, expecting error then successful termination"
const PostConsensusInvalidSignedMessageDifferentLengthDoc = "Tests post-consensus message processing with invalid signed message (different number of signers and signatures)"
const PostConsensusInvalidSignedMessageEmptySignatureDoc = "Tests post-consensus message processing with invalid signed message (empty signature)"
const PostConsensusMixedCommitteesDoc = "Tests post-consensus message processing with mixed committees, expecting correct handling of multiple committee scenarios"
const PostConsensusNilMsgDoc = "Tests post-consensus message processing with nil SignedSSVMessage, expecting error"
const PostConsensusNoRunningDutyDoc = "Tests post-consensus message processing when there is no running duty, expecting error"
const PostConsensusPartialInvalidRootQuorumThenValidQuorumDoc = "Tests post-consensus message processing where a partial invalid root quorum is followed by a valid quorum, expecting error and recovery handling."
const PostConsensusPartialInvalidSignatureQuorumThenValidQuorumDoc = "Tests post-consensus message processing where a partial invalid signature quorum is followed by a valid quorum, expecting error and recovery handling."
const PostConsensusPostFinishDoc = "Tests post-consensus message processing after duty is finished"
const PostConsensusPreDecidedDoc = "Tests post-consensus message processing before duty is decided, expecting error"
const PostConsensusUnknownSignerDoc = "Tests post-consensus message processing with unknown signer, expecting error"

// Documentation for runner/preconsensus tests
const PreConsensusDuplicateMsgDifferentRootsDoc = "Tests pre-consensus message processing with duplicate messages having different roots, expecting error"
const PreConsensusDuplicateMsgDoc = "Tests pre-consensus message processing with duplicate messages"
const PreConsensusInconsistentBeaconSignerDoc = "Tests pre-consensus message processing with inconsistent beacon signer, expecting error"
const PreConsensusInconsistentOperatorSignerDoc = "Tests pre-consensus message processing with inconsistent operator signer, expecting error"
const PreConsensusInvalidBeaconSignatureInQuorumDoc = "Tests pre-consensus message processing with invalid beacon signature in quorum, expecting error"
const PreConsensusInvalidExpectedRootDoc = "Tests pre-consensus message processing with invalid expected root, expecting error"
const PreConsensusInvalidMessageSignatureDoc = "Tests pre-consensus message processing with invalid message signature. No error is generated since the SignedPartialSignatureMessage.Signature is no longer checked"
const PreConsensusInvalidMessageSlotDoc = "Tests pre-consensus message processing with invalid message slot, expecting error"
const PreConsensusInvalidOperatorSignatureDoc = "Tests pre-consensus message processing with invalid operator signature, expecting error"
const PreConsensusInvalidQuorumThenValidQuorumDoc = "Tests pre-consensus message processing with an invalid quorum followed by a valid quorum, expecting error then successful termination"
const PreConsensusInvalidSignedMessageDifferentLengthDoc = "Tests pre-consensus message processing with invalid signed message (different number of signers and signatures)"
const PreConsensusInvalidSignedMessageEmptySignatureDoc = "Tests pre-consensus message processing with invalid signed message (empty signature)"
const PreConsensusInvalidSignedMessageNoSignatureDoc = "Tests pre-consensus message processing with invalid signed message (no signatures)"
const PreConsensusInvalidSignedMessageNoSignersDoc = "Tests pre-consensus message processing with invalid signed message (no signers)"
const PreConsensusInvalidSignedMessageDoc = "Tests pre-consensus message processing with invalid signed message, expecting error"
const PreConsensusInvalidThenQuorumDoc = "Tests pre-consensus message processing with an invalid message then a valid quorum, expecting error then successful termination"
const PreConsensusNilMsgDoc = "Tests pre-consensus message processing with nil SignedSSVMessage, expecting error"
const PreConsensusNoRunningDutyDoc = "Tests pre-consensus message processing when there is no running duty, expecting error"
const PreConsensusPostDecidedDoc = "Tests pre-consensus message processing after duty is decided"
const PreConsensusPostFinishDoc = "Tests pre-consensus message processing after duty is finished"
const PreConsensusPostQuorumDoc = "Tests pre-consensus message processing after quorum is reached"
const PreConsensusQuorum7OperatorsDoc = "Tests pre-consensus message processing with 7 operators"
const PreConsensusQuorum10OperatorsDoc = "Tests pre-consensus message processing with 10 operators"
const PreConsensusQuorum13OperatorsDoc = "Tests pre-consensus message processing with 13 operators"
const PreConsensusQuorumDoc = "Tests pre-consensus quorum message processing across different runner types"
const PreConsensusTooFewRootsDoc = "Tests pre-consensus message processing with too few roots, expecting error"
const PreConsensusTooManyRootsDoc = "Tests pre-consensus message processing with too many roots, expecting error"
const PreConsensusUnknownSignerDoc = "Tests pre-consensus message processing with unknown signer"
const PreConsensusUnorderedExpectedRootsDoc = "Tests pre-consensus message processing with unordered expected roots, expecting error"
const PreConsensusValidMsg7OperatorsDoc = "Tests pre-consensus message processing with 7 operators"
const PreConsensusValidMsg10OperatorsDoc = "Tests pre-consensus message processing with 10 operators"
const PreConsensusValidMsg13OperatorsDoc = "Tests pre-consensus message processing with 13 operators"
const PreConsensusValidMsgDoc = "Tests pre-consensus message processing across different runner types"

// Documentation for valcheckattestation tests
const ValCheckAttestationBeaconVoteDataNilDoc = "Tests attestation value check with nil attestation data"
const ValCheckAttestationFarFutureTargetDoc = "Tests attestation value check with target epoch too far in the future"
const ValCheckAttestationMajoritySlashableDoc = "Tests attestation value check with majority slashable attestation"
const ValCheckAttestationMinoritySlashableDoc = "Tests attestation value check with minority slashable attestation (source and target different from previous)"
const ValCheckAttestationSlashableDoc = "Tests attestation value check with slashable attestation"
const ValCheckAttestationSourceHigherThanTargetDoc = "Tests attestation value check with source epoch higher than target epoch"
const ValCheckAttestationValidNonSlashableSlotDoc = "Tests attestation value check with valid non-slashable slot"
const ValCheckAttestationValidDoc = "Tests attestation value check with valid attestation"

// Documentation for valcheckduty tests
const ValCheckDutyFarFutureDutySlotDoc = "Tests duty value check with duty slot too far in the future"
const ValCheckDutyWrongDutyTypeDoc = "Tests duty value check with wrong duty type"
const ValCheckDutyWrongValidatorIndexDoc = "Tests duty value check with wrong validator index across different roles"
const ValCheckDutyWrongValidatorPKDoc = "Tests duty value check with wrong validator public key across different roles"

// Documentation for valcheckproposer tests
const ValCheckProposerBlindedBlockDoc = "Tests proposer value check with blinded block data"
