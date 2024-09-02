# Message Validation

> [!NOTE]
> For implementation details, access the module's [code documentation]().


Through the [Extended Validators](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#extended-validators) feature, GossipSub allows the application to define a validation function to be applied to a message before it reaches the application.

The output of the validation may be one of the following types:
- `ValidationAccept`: for a message which is considered valid and can be forwarded to the application. 
- `ValidationReject`: for a message that is considered malicious regarding the protocol rules. This message is dropped and the peer is punished according to the scoring parametrs.
- `ValidationIgnore`: for a message that is either irrelevant or doesn't contribute for the state progression of the application.

In case the validation function returns `ValidationIgnore` or `ValidationReject`, the message is dropped so that the node doesn't overlay bad messages to the network.

## Rules

For our case, we defined a list of rules that a message should follow, which can be divided into the following categories:
- PubSub: regards the validity of the pubsub message.
- Syntax: regards the validity of the composition of the message.
- Semantic: regards the validity of the meaning that a message's attribute may have.
- QBFT logic: regards the validity of a message or a sequence of messages considering the QBFT rules.
- Duty logic: regards the validity of a message or a sequence of messages considering the duty rules.


### PubSub

| Verification                     | Error                     | Classification | Explanation                                                |
|----------------------------------|---------------------------|----------------|------------------------------------------------------------|
| Empty pubsub.Message.Data        | ErrPubSubMessageHasNoData | Reject         | pubsub.Message.Data must not be empty.                     |
| Big pubsub.Message.Data          | ErrPubSubDataTooBig       | Ignore         | pubsub.Message.Data must be below a size limit             |
| Can't decode pubsub.Message.Data | ErrMalformedPubSubMessage | Reject         | pubsub.Message.Data must be decodable to SignedSSVMessage. |

### Syntax

| Verification                    | Error                                      | Classification | Explanation                                                               |
|---------------------------------|--------------------------------------------|----------------|---------------------------------------------------------------------------|
| Nil SignedSSVMessage            | ErrNilSignedSSVMessage                     | Reject         | SignedSSVMessage can't be nil.                                            |
| No signers                      | ErrNoSigners                               | Reject         | Len(SignedSSVMessage.OperatorIDs) must be >= 0.                           |
| No signatures                   | ErrNoSignatures                            | Reject         | Len(SignedSSVMessage.Signatures) must be >= 0.                            |
| Wrong signature size            | ErrWrongRSASignatureSize                   | Reject         | $\forall i$ Len(SignedSSVMessage.Signatures[i]) must be of expected size. |
| Sorted signers                  | ErrSignersNotSorted                        | Reject         | SignedSSVMessage.OperatorIDs must be sorted.                              |
| Non-zero signers                | ErrZeroSigner                              | Reject         | $\forall i$ SignedSSVMessage.OperatorIDs[i] must not be 0.                |
| Unique signers                  | ErrDuplicatedSigner                        | Reject         | SignedSSVMessage.OperatorIDs must have unique values.                     |
| Len(signers) = Len(signatures)  | ErrSignersAndSignaturesWithDifferentLength | Reject         | The number of signers and signatures must match                           |
| Nil SSVMessage                  | ErrNilSSVMessage                           | Reject         | SignedSSVMessage.SSVMessage can't be nil.                                 |
| Empty SSVMessage.Data           | ErrEmptyData                               | Reject         | SSVMessage.Data can't be empty.                                           |
| Big SSVMessage.Data             | ErrSSVMessageDataTooBig                    | Ignore         | SSVMessage.Data must be below a size limit.                               |
| Can't decode SSVMessage.Data    | ErrUndecodableData                         | Reject         | SSVMessage.Data must be decodable.                                        |
| Can't decode QBFT Justification | ErrUndecodableData                         | Reject         | A qbft.Message must have decodable justifications.                        |

### Semantics General Rules

| Verification         | Error                     | Classification | Explanation                                                           |
|----------------------|---------------------------|----------------|-----------------------------------------------------------------------|
| Signers in committee | ErrSignerNotInCommittee   | Reject         | Signers must belong to validator's or CommitteeID's committee.        |
| Different Domain     | ErrWrongDomain            | Ignore         | MsgID.Domain is different than self domain.                           |
| Invalid Role         | ErrInvalidRole            | Reject         | MsgID.Role is not known.                                              |
| Validator exists     | ErrUnknownValidator       | Ignore         | If MsgID.SenderID is a validator, it must exist.                      |
| Active Validator ID  | ErrValidatorNotAttesting  | Ignore         | If MsgID.SenderID is a validator, it must be active active validator. |
| Validator Liquidated | ErrValidatorLiquidated    | Ignore         | If MsgID.SenderID is a validator, it must not be liquidated.          |
| CommitteeID exists   | ErrNonExistentCommitteeID | Ignore         | If MsgID.SenderID is a committee, it must exist.                      |
| Wrong topic          | ErrIncorrectTopic         | Ignore         | The message should be sent in the correct topic                       |
| DKG Message          | ErrDKGMessage             | Reject         | MsgType can't be of DKG message.                                      |
| Unknown MsgType      | ErrUnknownSSVMessageType  | Reject         | MsgType is not known.                                                 |


### Consensus

#### Semantics


|           Verification            |              Error               | Classification |                         Explanation                         |
| --------------------------------- | -------------------------------- | -------------- | ----------------------------------------------------------- |
| Non-decided with multiple signers | ErrNonDecidedWithMultipleSigners | Reject         | Non-decided message must have one signer.                   |
| Decided with enough signers       | ErrDecidedNotEnoughSigners       | Reject         | Decided message must have at least a quorum of signers      |
| Prepare or commit with full data  | ErrPrepareOrCommitWithFullData   | Reject         | Prepare or Commit messages must not have FullData.          |
| Invalid full data hash            | ErrInvalidHash                   | Reject         | If there's a FullData field, Message.Root must be the hash. |
| Unknown MsgType                   | ErrUnknownQBFTMessageType        | Reject         | Message.MsgType must be known.                              |
| Zero round                        | ErrZeroRound                     | Reject         | Message.Round must not be zero.                             |
| Mismatched identifier             | ErrMismatchedIdentifier          | Reject         | Message.Identifier must match SSVMessage.Identifier         |


#### QBFT Logic

| Verification                            | Error                                  | Classification | Explanation                                                                                                           |
|-----------------------------------------|----------------------------------------|----------------|-----------------------------------------------------------------------------------------------------------------------|
| Not Leader                              | ErrSignerNotLeader                     | Reject         | Signer is not leader for round.                                                                                       |
| Different decided with same signers     | ErrDecidedWithSameSigners              | Ignore         | Two decided messages for duty with the same signers.                                                                  |
| Double Proposal with different FullData | ErrDuplicatedProposalWithDifferentData | Reject         | Signer already sent a different proposal for round.                                                                   |
| Double Proposal                         | ErrDuplicatedMessage                   | Reject         | Signer already sent a proposal for round.                                                                             |
| Double Prepare                          | ErrDuplicatedMessage                   | Reject         | Signer already sent a Prepare for round.                                                                              |
| Double Commit                           | ErrDuplicatedMessage                   | Reject         | Signer already sent a Commit for round.                                                                               |
| Double Round-Change                     | ErrDuplicatedMessage                   | Reject         | Signer already sent a Round-Change for round.                                                                         |
| Round in time-spread                    | ErrRoundNotAllowedInTimeSpread         | Ignore         | Round must be in the estimated round interval, considering the received the current time. Allows a spread of 1 round. |
| Already advanced round                  | ErrRoundAlreadyAdvanced                | Ignore         | Signer is already in a future round.                                                                                  |
| Unexpected Round-Change justification   | ErrUnexpectedRoundChangeJustifications | Reject         | Round-Change justification can only be for a Proposal or Round-Change.                                                |
| Unexpected Prepare justification        | ErrUnexpectedPrepareJustifications     | Reject         | Preapre justification can only be for Proposal.                                                                       |


#### Duty Logic

| Verification               | Error                         | Classification | Explanation                                                                    |
|----------------------------|-------------------------------|----------------|--------------------------------------------------------------------------------|
| Invalid role for consensus | ErrUnexpectedConsensusMessage | Reject         | SSVMessage.MsgID.Role must not be ValidatorRegistration or VoluntaryExit.      |
| Valid round for role       | ErrRoundTooHigh               | Reject         | For committee and aggregation, round can go up to 12. Else, it can go up to 6. |

### Partial Signatures

#### Semantics

| Verification               | Error                               | Classification | Explanation                                                                                                                                                                                                                                                                                                                                                                                                    |
|----------------------------|-------------------------------------|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Multiple signers           | ErrPartialSigMultipleSigners        | Reject         | Must have only 1 signer.                                                                                                                                                                                                                                                                                                                                                                                       |
| Unexpected FullData        | ErrFullDataNotInConsensusMessage    | Reject         | Must not have FullData.                                                                                                                                                                                                                                                                                                                                                                                        |
| Unknown type               | ErrInvalidPartialSignatureType      | Reject         | Type not known.                                                                                                                                                                                                                                                                                                                                                                                                |
| Wrong type for role        | ErrPartialSignatureTypeRoleMismatch | Reject         | Type must match role:<br> PostConsensusPartialSig for Committee,<br> RandaoPartialSig or PostConsensusPartialSig for Proposer,<br> SelectionProofPartialSig or PostConsensusPartialSig for Aggregator,<br> SelectionProofPartialSig or PostConsensusPartialSig for Sync committee contribution,<br> ValidatorRegistrationPartialSig for Validator Registration,<br> VoluntaryExitPartialSig for Voluntary Exit |
| No PartialSignatureMessage | ErrNoPartialSignatureMessages       | Reject         | Message must have at least one PartialSignatureMessage.                                                                                                                                                                                                                                                                                                                                                        |
| Wrong BLS Signature Size   | ErrWrongBLSSignatureSize            | Reject         | $\forall i$ PartialSignatureMessages.Message[i].Signature must have the correct length.                                                                                                                                                                                                                                                                                                                        |
| Inconsistent signer        | ErrInconsistentSigners              | Reject         | $\forall i$ PartialSignatureMessages.Message[i].Signer must be the same as the<br> SignedSSVMessage.OperatorIDs[i].                                                                                                                                                                                                                                                                                            |
| Validtor's index mismatch  | ErrValidatorIndexMismatch           | Ignore         | Except for the committee role, $\forall i$ PartialSignatureMessages.Message[i].ValidatorIndex must<br> belong to SSVMessage.SenderID(). For the committee role, we can't assume synchrony on the committee's validators sets.                                                                                                                                                                                  |


#### Duty Logic

|         Verification         |          Error           | Classification |                                                                                             Explanation                                                                                              |
| ---------------------------- | ------------------------ | -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Invalid signature type count | ErrInvalidPartialSignatureTypeCount                | Reject         | It's allow only:<br> 1 PostConsensusPartialSig, for Committee duty,<br> 1 RandaoPartialSig and 1 PostConsensusPartialSig for Proposer,<br> 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Aggregator,<br> 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Sync committee contribution,<br> 1 ValidatorRegistrationPartialSig for Validator Registration,<br> 1 VoluntaryExitPartialSig for Voluntary Exit. |
| Too many partial signatures  | ErrTooManyPartialSignatureMessages                | Reject         | For the committee role, it's allowed $min(2*V, V + $ SYNC_COMMITTEE_SIZE $)$ <br> where $V$ is the number of committee's validatos.<br> For sync committee contribution, it's allowed 13.<br> Else, only 1.  |
| Triple validator index       | ErrTripleValidatorIndexInPartialSignatures                | Reject         | A validator index can not be associated to more than 2 signatures. |


### General duty Logic

| Verification              | Error                                     | Classification | Explanation                                                                                                                                                                                          |
|---------------------------|-------------------------------------------|----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Already advanced slot     | ErrSlotAlreadyAdvanced                    | Ignore         | (Non-committee roles) Signer already advanced to later slot.                                                                                                                                         |
| Slot not in time for role | ErrEarlySlotMessage or ErrLateSlotMessage | Ignore         | Current time must be between duty's starting time and<br> +34 (committee and aggregator) or +3 (else) slots.                                                                                         |
| No beacon duty            | ErrNoDuty                                 | Ignore         | If Proposal or Sync committee contribution duty, check if duty exists with beacon node.                                                                                                              |
| Too many duties per epoch | ErrTooManyDutiesPerEpoch                  | Ignore         | If role is either aggregator, voluntary exit and validator registration,<br> it's allowed 2 duties per epoch. Else if committee,<br> 2*V (if no validator is doing sync committee).<br> Else accept. |

### Signature

The last check performed is the signature verification. This is the last check since it's, by far, the most expensive one.

|  Verification   |          Error           | Classification |        Explanation        |
| --------------- | ------------------------ | -------------- | ------------------------- |
| Valid signature | ErrSignatureVerification | Reject         | Signature must be correct |