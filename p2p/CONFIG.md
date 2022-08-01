# Network Configurations

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | JUL 22        |

This document contains information on the networking configurations used for `SSV.Network`.

## Node Config

### General

Below is a set of general configurations for an `SSV` node:

| ENV                  | YAML                 | Default Value (SSV)    | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `NETWORK_PRIVATE_KEY`| `NetworkPrivateKey`  | -                      | No       | Key to use for libp2p/network identity |
| `TCP_PORT`           | `p2p.TcpPort`        | `13001`                | No       | TCP port to use as libp2p transport    |
| `UDP_PORT`           | `p2p.UdpPort`        | `12001`                | No       | UDP port to use for discovery          |
| `HOST_ADDRESS`       | `p2p.HostAddress`    | -                      | No       | External IP address                    |
| `HOST_DNS`           | `p2p.HostDNS`        | -                      | No       | External DNS address                   |
| `P2P_REQUEST_TIMEOUT`| `p2p.RequestTimeout` | `7s`                   | No       | Timeout for stream requests            |
| `P2P_MAX_BATCH_RESPONSE` |`p2p.MaxBatchResponse`| `25`               | No       | Maximum number of returned objects in a batch for sync protocols |

### Connectivity and Discovery

The following set of configurations affects connectivity and discovery:

| ENV                  | YAML                 | Default Value (SSV)    | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `NETWORK_ID`         | `p2p.NetworkID`      | `ssv-testnet-v2`       | No       | ID of the network that the node runs on|
| `BOOTNODES`          | `p2p.Bootnodes`      | Bootnode (Testnet)     | No       | ENR of network bootnodes               |
| `SUBNETS`            | `p2p.Subnets`        | -                      | No       | Subnets that this node will join upon start, represented as hex string|
| `P2P_MAX_PEERS`      | `p2p.MaxPeers`       | 60                     | No       | Maximum number connections to maintain |

### Logging and Tracing

The following set of configurations affects logging/tracing w.r.t network components.

Note that in order to have a detail log, the node operator must set log level (`LOG_LEVEL`) to `debug`.

| ENV                  | YAML                 | Default Value (SSV)    | Required |  Description                             |
| ---                  | ---                  | ---                    | ---      | ---                                      |
| `P2P_LOG`            | `p2p.P2pLog`         | false                  | No       | Flag to turn on/off network debug logs, and other traces within network layer  |
| `PUBSUB_TRACE`       | `p2p.NetworkTrace`   | false                  | No       | Flag to turn on/off pubsub trace logs    |
| `DISCOVERY_TRACE`    | `p2p.NetworkTrace`   | false                  | No       | Flag to turn on/off discovery trace logs |

## Pubsub Parameters

The following parameters are used when initializing and interaction with the pubsub router.

Note that we compare default values with both libp2p defaults and ETH2.0 nodes. 
For more information regards ETH2.0 configuration please refer to 
[their p2p spec](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#the-gossip-domain-gossipsub).

### Last Seen

Specifies how long a message will be remembered as seen.

**Default Value (libp2p):** `2min` \
**Default Value (SSV):** `6.4min` \
**Default Value (ETH2):** `6.4min`

Was increased to `6.4m` which is the duration of one epoch, as messages become redundant afterwards.

### Outbound Queue Size

The size of the queue that is used for outbound messages.

**Default Value (libp2p):** `32` \
**Default Value (SSV):** `256` \
**Default Value (ETH2):** `600`

The value was increased to support higher load of messages.
SSV uses a smaller number than ETH2.0 due to high memory usage, 
which is caused by the size of raw JSON messages.

**NOTE** This parameter will be considered once encoding is changed to `SSZ` and compression is applied (`snappy` or `s2`). 

### Validation Queue Size

The size of the queue that is used for validation of incoming messages.

**Default Value (libp2p):** `32` \
**Default Value (SSV):** `256` \
**Default Value (ETH2):** `600`

The value was increased to support higher load of messages.
SSV uses a smaller number than ETH2.0 due to high memory usage,
which is caused by the size of raw JSON messages.

**NOTE** This parameter will be considered once encoding is changed to `SSZ` and compression is applied (`snappy` or `s2`).

### Msg ID

A function that calculates an identifier for messages, to be used across pubsub components.

**Default Value (libp2p):** 
```go
func MsgID(msg) string {
	return msg.GetFrom() + msg.GetSeqno()
}
```
**Default Value (SSV):**
```go
func MsgID(msg) string {
    return string(sha256(msg.GetData())[:20])
}
```
**Default Value (ETH2):**
```go
func MsgID(msg) string {
    h := sha256(MESSAGE_DOMAIN_VALID_SNAPPY + snappy_decompress(message.data))
    return string(h[:20])
}
```

SSV uses a custom function that returns the content hash, so we won't process the same message multiple times.


### Subscription Filter

Allows to control the topics that the node will subscribe to.

**Default Value (libp2p):** Accept all topics \
**Default Value (SSV):** Accept topics with the same fork if the node has interest in them, 
if we didn't reach subscriptions limit (`129`) \
**Default Value (ETH2):** Accept topics with supported fork digest if the node has interest in them,
and it didn't reach subscriptions limit (`200`)


### Flood Publish

Force peer's own messages to be published to all known peers for the topic.

**Default Value (libp2p):** `false` \
**Default Value (SSV):** `false` \
**Default Value (ETH2):** `false`

When turned on, this feature ensures reliability and protects from eclipse attacks. 
On the other hand it floods the network with duplicated message and therefore it was turned off.


### Signature Policy

The mode of operation for producing and verifying message signatures in the pubsub router level.

**Default Value (libp2p):** `StrictSign` \
**Default Value (SSV):** `StrictSign` \
**Default Value (ETH2):** `StrictNoSign`

**TODO**


### Topic Message Validator

A function that is invoked by pubsub for incoming messages before they are being processed.

**Default Value (libp2p):** None \
**Default Value (SSV):** Decodes the message and validate that it was sent on the right topic \
**Default Value (ETH2):** A more complete validation, according to message type (according to topic). \
More details can be [found here](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#global-topics)

**TODO**


### Gossipsub: D

Sets the optimal degree for a GossipSub topic mesh. D should be set somewhere between Dlo and Dhi.

**Default Value (libp2p):** `6` \
**Default Value (SSV):** `8` \
**Default Value (ETH2):** `8`

### Gossipsub: Dlo

Sets the lower bound on the number of peers we keep in a GossipSub topic mesh.

**Default Value (libp2p):** `5` \
**Default Value (SSV):** `6` \
**Default Value (ETH2):** `6`

### Gossipsub: Dhi

Sets the upper bound on the number of peers we keep in a GossipSub topic mesh.

**Default Value (libp2p):** `5` \
**Default Value (SSV):** `6` \
**Default Value (ETH2):** `6`

### Gossipsub: HeartbeatInterval

Controls the time between heartbeats, which are used across pubsub components to align on timing.

**Default Value (libp2p):** `1s` \
**Default Value (SSV):** `700ms` \
**Default Value (ETH2):** `700ms`

As SSV nodes work intensively with pubsub, we had to decrease the heartbeat interval 
to reduce the latency created by a higher heartbeat.

### Gossipsub: HistoryLength

Controls the size of the message cache used for gossip (`IWANT` responses).

**Default Value (libp2p):** `5` \
**Default Value (SSV):** `6` \
**Default Value (ETH2):** `6`

The value was increased to reduce the number of `IWANT` messages in the network.

### Gossipsub: HistoryGossip

Controls how many cached message ids we will advertise in IHAVE gossip messages.

**Default Value (libp2p):** `3` \
**Default Value (SSV):** `4` \
**Default Value (ETH2):** `3`

**TODO**

### Gossipsub: MaxIHaveLength

Sets the maximum number of messages to include in an IHAVE message.

**Default Value (libp2p):** `5000` \
**Default Value (SSV):** `1500` \
**Default Value (ETH2):** `5000`

Decreased value to avoid ihave floods.
**TODO**

### Gossipsub: MaxIHaveMessages

Sets the maximum number of IHAVE messages to accept from a peer within a heartbeat.

**Default Value (libp2p):** `10` \
**Default Value (SSV):** `32` \
**Default Value (ETH2):** `10`

Increased as we want messages to be sent in batches, to reduce the amount of requests.

<br />

**TODO: Scoring Params**
