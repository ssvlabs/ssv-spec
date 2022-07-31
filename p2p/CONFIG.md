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

Note that we compare default values with both libp2p and ETH2.0 nodes. For more information regards ETH2.0 configuration please refer to [their p2p spec](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#the-gossip-domain-gossipsub).

| Parameter            |  Description | Default Value (libp2p) | Default Value (SSV) | Default Value (ETH2.0) | Comments |
| ---                  | ---          | ---                    | ---                 | ---                    | ---      |
| `lastSeen` | Specifies how long a message will be remembered as seen | `2min` | `6.4min` | `6.4min` | The value (6.4m) is the duration of one epoch as message become redundant afterwards |
| `peerOutboundQueueSize` | The size of the queue that is used for outbound messages | `32` | `256` | `600` | The value was increased to support high load of messages, using a larger value might cause high memory usage |
| `validationQueueSize`| The size of the queue that is used for validation of incoming messages  | `32` | `256` | `600` | The value was increased to support high load of messages, using a larger value might cause high memory usage |
| `validateThrottle`   | The upper bound on the number of active validation goroutines across all topics | `8192` | `4096` | N/A | The value was decreased to avoid high resource usage and reduce overloading on the node |
| `msgID`    | A function that calculates an identifier for messages, to be used across pubsub components | `(msg) => msg.GetFrom() + msg.GetSeqno()` | `(msg) => sha256(msg.GetData())[:20]` | `(msg) => sha256(MESSAGE_DOMAIN_VALID_SNAPPY + snappy_decompress(message.data))[:20]` | We use a custom function that returns the content hash, so we won't process the same message multiple times
| `subscriptionFilter` | Allows to control the topics that the node will subscribe to | `nil` | Accept topics with the same fork if the node has interest in them, if we didn't reach subscriptions limit (`129`) | accept topics with supported fork digest if the node has interest in them and we didn't reach subscriptions limit (`200`) | - |
| `floodPublish`       | Force peer's own messages to be published to all known peers for the topic | `false` | `false` | `false` | When turned on, this feature ensures reliability and protects from eclipse attacks. On the other hand it floods the network with duplicated message and therefore it was turned off |
| `signaturePolicy`    | The mode of operation for producing and verifying message signatures in the pubsub router level | `StrictSign` | `StrictSign` | `StrictNoSign` | **TODO** |
| `msgValidator`       | A function that is invoked by pubsub for incoming messages before they are being processed | `None` | Decodes the message and validate that it was sent on the right topic | A more complete validation, according to message type (according to topic). More details can be [found here](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#global-topics) | **TODO** |
| **Gossipsub**  | | | | |
| `D`       | Sets the optimal degree for a GossipSub topic mesh. D should be set somewhere between Dlo and Dhi | `6` | `5`| `8` | Decreased as we want to reduce loading coming from pubsub |
| `Dlo`       | Sets the lower bound on the number of peers we keep in a GossipSub topic mesh | `5` | `3` | `6` | Decreased as we want to reduce loading coming from pubsub |
| `Dhi`       | Sets the upper bound on the number of peers we keep in a GossipSub topic mesh | `12` | `9` | `12` | Decreased as we want to reduce loading coming from pubsub |
| `HeartbeatInterval` | Controls the time between heartbeats | `1s` | `700ms` | `700ms` | As we work intensively with pubsub, we had to decrease the heartbeet interval to reduce the latency created by a higher heartbeat |
| `HistoryLength` | Controls the size of the message cache used for gossip (`IWANT` responses) | `5` | `100` | `6` | The value was increased to reduce the number of `IWANT` messages in the network |
| `HistoryGossip` | Controls how many cached message ids we will advertise in IHAVE gossip messages | `3` | `3` | `3` | - |
| `MaxIHaveLength` | Sets the maximum number of messages to include in an IHAVE message | `5000` | `1000` | N/A | Decreased value to avoid ihave floods |
| `MaxIHaveMessages` | Sets the maximum number of IHAVE messages to accept from a peer within a heartbeat | `10` | `32` | N/A | Increased as we want messages to be sent in batches, to reduce the amount of requests |
| **Scoring** (**TODO**) | | | | |
