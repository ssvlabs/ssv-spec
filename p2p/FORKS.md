# Network Forking

|                        Contributors                         | Last Revision |
| :---------------------------------------------------------- | :------------ |
| [@amir-blox](https://github.com/amir-blox) & Matheus Franco | September 24  |

This document contains information on the changes in network domain for `SSV.Network` forks.


## Fork Genesis

### Validator topic mapping

Validator public key hash is used to determine the validator's subnet: \
`bloxstaking.ssv.{{hash(validatiorPubKey) % num_of_subnets}}`

### Number of subnets

Number of subnets for this fork is `128`.

### Message encoding

JSON is used for encoding/decoding of messages.

### Sync protocols ID

The following sync protocols are used with the
corresponding number of peers that are being called for a single request:

- `/ssv/sync/decided/last/0.0.1`, 10 peers
- `/ssv/sync/decided/history/0.0.1`, 10 peers
- `/ssv/sync/round/0.0.1`, 5 peers

### User agent

User Agent contains the node version and type
`SSV-Node/v0.x.x`


## Alan Fork

### Topic mapping

For the Alan fork, operators register to topics according to the committees it belongs to.

To exchange message for a validator duty, the operator subscribes to the topic

`committee_id % num_of_subnets`

where `committee_id` is the [`CommitteeID`](./../types/docs/README.md#operatorid-and-committeeid) \
of the committee that the validator is registered to.

### Committee Duty

The Attestation and Sync Committee duties are, now, unified in the `CommitteeDuty` type.
This optimization allows a committee to run only one consensus execution, as explained in the
[SSV protocol file](./../ssv/README.md).

### Message encoding

Messages are encoded using [`ssz`](https://github.com/ferranbt/fastssz).

### MessageID

The `MessageID` type, now, may hold a `CommitteeID` depending on the duty type. Check its [documentation](./../types/docs/README.md#messageid) for more details.

### Sync protocols

No sync protocols are used any longer.