# Types

This packages defines the basic building blocks for the [SSV](./../ssv/README.md) and [QBFT](./../qbft/README.md) modules.

## DomainType

`DomainType` is a sequence of 4 bytes that uniquely identifies the network:
- The first two bytes are yet empty and lack meaning.
- The third byte is the network ID (e.g. Mainnet, Primus, Shifu, Jato)
- The fourth byte is the current fork number of the network.

## BeaconNetwork

`BeaconNetwork` is a type that allows accessing information about the blockchain, e.g. such as the current slot, the number of seconds per slot, and the fork version.

## BeaconRole and RunnerRole

The `BeaconRole` is an enumeration regarding the beacon duty types including: Attester, Aggregator, Proposer, Sync Committee, and Sync Committee Contribution. It also encompass the two extra duties: Validator Registration, and Voluntary Exit.

The `RunnerRole` is similar to the `BeaconRole` enumeration but it regards the types of [runners](./../ssv/docs/documentation.md#runner). Its difference to the `BeaconRole` is that it lacks the Attester and Sync Committee values while it has an extra `CommitteeRole`. This is due to the [CommitteeDuty optimization](./../ssv/README.md#the-committee-duty).

## BeaconSigner and OperatorSigner

`BeaconSigner` is an interface that provides methods for creating beacon signatures. Apart from signing, it servers as a slashing protection by providing a method to check if an attestation is slashable.

The `OperatorSigner` is a similar structure but performs signing for network messages.

## OperatorID and CommitteeID

`OperatorID` is a `uint64` type that uniquely identifies an operator in the network.

The `CommitteeID` is a sequence of 32 bytes that identifies a committee (a cluster of operators). It is the result of applying a hash function to the string that represents the sorted list of `OperatorID`s that belong to the committee.

## MessageID

`MessagID` is an identifier that is composed by:
- a [`DomainType`](#domaintype) that uniquely identifies the SSV network.
- a [`RunnerRole`](#beaconrole-and-runnerrole) that uniquely identifies the [runner](./../ssv/docs/documentation.md#runner) type.
- a `DutyExecutor` that is either a validator public key (if the `RunnerRole` is not the `CommitteeRole`) or a [`CommitteeID`](#operatorid-and-committeeid).

## Share / ShareMember and CommitteeMember / Operator

The `Share` represents a shared validator. Namely, it holds:
- the `ValidatorIndex`, its public key, recipient address and graffiti.
- a list of `ShareMembers` which consists of an `OperatorID` and its public share key.

The `CommitteeMember` is a similar structure, but it contains information about a committee. Namely, it holds:
- the [`CommitteeID`](#operatorid-and-committeeid)
- the number of faulty nodes
- a list of `Operators` which consists of an `OperatorID` and its network public key.

## SignedSSVMessage and SSVMessage

Messages are exchanged in the network using the `SignedSSVMessage` structure:
```go
type SignedSSVMessage struct {
	Signatures  [][]byte
	OperatorIDs []OperatorID
	SSVMessage  *SSVMessage
	FullData []byte
}
```
It contains the list of signers, signatures, a `SSVMessage`, and a `FullData` field used for proposed data in QBFT. The signers and signatures lists have more than one element only for Decided messages.

The `SSVMesage` identifies the type of the message (consensus or partial signature) and includes the encoded message along with a [`MessageID`](#messageid).
```go
type SSVMessage struct {
	MsgType MsgType
	MsgID   MessageID
	Data []byte
}
```

If the `MsgType` is of consensus type, then the `Data` should an encoded `qbft.Message` object. Else (if it's of partial signature type), it should be an encoded [`PartialSignatureMessages`](#partialsignaturemessages) object.

## PartialSignatureMessage(s)

The `PartialSignatureMessages` is the message used for the partial signature phases.

```go
type PartialSignatureMessages struct {
	Type     PartialSigMsgType
	Slot     phase0.Slot
	Messages []*PartialSignatureMessage
}
```

It contains the signature type (e.g. RANDAO signature, post-consensus signature), the duty's slot, and a list of `PartialSignatureMessage`, which contains the partial signature, the signing root, the sender (`OperatorID`), and the validator associated to the duty.

```go
type PartialSignatureMessage struct {
	PartialSignature Signature
	SigningRoot      [32]byte
	Signer           OperatorID
	ValidatorIndex   phase0.ValidatorIndex
}
```

Multiple `PartialSignatureMessage` may be required because:
- the Sync Committee Contribution duty may encompass multiple committees, and a signature is required for each duty.
- the `CommitteeDuty` encompass multiple duties and a signature is required for each.

## Duty, ValidatorDuty and CommitteeDuty

`Duty` is an interface to represent a duty that will be executed by a `Runner`. There are two implementations:
- `ValidatorDuty`: represents an Ethereum duty, including the duty type, the validator assigned to it, the slot, and more duty details.
- `CommitteeDuty`: contains a collection of `ValidatorDuty` objects. It's used for the `CommitteeRunner` which performs many Attestation and Sync Committee duties at once.

## ValidatorConsensusData and BeaconVote

But what objects do the operators propose in the consensus protocol? It can be a `ValidatorConsensusData` or a `BeaconDuty`.

The `ValidatorConsensusData` holds the associated duty, a blockchain version (e.g. as Phase0, Altair, Bellatrix, Capella, and Deneb), and an encoded object related to the duty (e.g. an AttestationData, a BeaconBlock). This structure is used for all except the `CommitteeDuty`.

The `BeaconVote` contains a blockchain head vote (LMD GHOST vote) and a source and target votes (Casper FFG). It's used for the `CommitteeDuty`.