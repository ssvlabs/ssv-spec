# Network Forking

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | JUL 22        |

This document contains information on the changes in network domain for `SSV.network` forks.


## Fork Genesis

### validator topic mapping

Validator public key hash is used to determine the validator's subnet: \
`bloxstaking.ssv.{{hash(validatiorPubKey) % num_of_subnets}}`

### number of subnets

Number of subnets for this fork is `128`.

### msgID function

msgID function to use by libp2p's `pubsub.Router` for calculating `msg_id`.

a content based msgID function is used, see [msg-id](#message-id): \
`hash(msg-data)[:20]`

### message encoding

JSON is used for encoding/decoding of messages.

### sync protocols ID

The following sync protocols are used with the
corresponding number of peers that are being called for a single request:

- `/ssv/sync/decided/last/0.0.1`, 10 peers
- `/ssv/sync/decided/history/0.0.1`, 10 peers
- `/ssv/sync/round/0.0.1`, 5 peers

### user agent

User Agent contains the node version and type
`SSV-Node/v0.x.x`



