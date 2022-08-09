# Network Forking

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | JUL 22        |

This document contains information on the changes in network domain for `SSV.network` forks.


### validator topic mapping

**genesis** - Validator public key hash is used to determine the validator's subnet: \
`bloxstaking.ssv.{{hash(validatiorPubKey) % num_of_subnets}}`

### number of subnets

**genesis** - number of subnets for this fork is `128`.

### msgID function

msgID function to use by libp2p's `pubsub.Router` for calculating `msg_id`.

**genesis** - a content based msgID function is used, see [msg-id](#message-id): \
`hash(msg-data)[:20]`

<br />

### message encoding

**genesis** - JSON is used for encoding/decoding of messages.

<br />

**TBD** [SSZ](https://github.com/ethereum/consensus-specs/blob/v0.11.1/ssz/simple-serialize.md)
will be used to encode/decode network messages.
It is efficient with traversing on fields, and is the standard encoding in ETH 2.0.

<br />

### sync protocols ID

**genesis** - the following sync protocols are used with the
corresponding number of peers that are being called for a single request:

- `/ssv/sync/decided/last/0.0.1`, 10 peers
- `/ssv/sync/decided/history/0.0.1`, 10 peers
- `/ssv/sync/round/0.0.1`, 5 peers

<br />

### user agent

**genesis** - User Agent contains the node version and type
`SSV-Node/v0.x.x`



