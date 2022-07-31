# Network Configurations

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | JUL 22        |

This document contains information on the networking configurations used for `SSV.Network`.

## Node Config

### General

| ENV                  | YAML                 | Default Value          | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `NETWORK_PRIVATE_KEY`| `NetworkPrivateKey`  | -                      | No       | Key to use for libp2p/network identity |
| `NETWORK_ID`         | `p2p.NetworkID`      | `ssv-testnet-v2`       | No       | ID of the network that the node runs on|
| `TCP_PORT`           | `p2p.TcpPort`        | `13001`                | No       | TCP port to use as libp2p transport    |
| `UDP_PORT`           | `p2p.UdpPort`        | `12001`                | No       | UDP port to use for discovery          |
| `HOST_ADDRESS`       | `p2p.HostAddress`    | -                      | No       | External IP address                    |
| `HOST_DNS`           | `p2p.HostDNS`        | -                      | No       | External DNS address                   |
| `P2P_REQUEST_TIMEOUT`| `p2p.RequestTimeout` | `7s`                   | No       | Timeout for stream requests            |
| `P2P_MAX_BATCH_RESPONSE` |`p2p.MaxBatchResponse`| `25`               | No       | Maximum number of returned objects in a batch for sync protocols |

### Connectivity and Discovery

| ENV                  | YAML                 | Default Value          | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `BOOTNODES`          | `p2p.Bootnodes`      | Bootnode (Testnet)     | No       | ENR of network bootnodes               |
| `SUBNETS`            | `p2p.Subnets`        | -                      | No       | Subnets that this node will join upon start, represented as hex string|
| `P2P_MAX_PEERS`      | `p2p.MaxPeers`       | 60                     | No       | Maximum number connections to maintain |

### Logging and Tracing

| ENV                  | YAML                 | Default Value          | Required |  Description                             |
| ---                  | ---                  | ---                    | ---      | ---                                      |
| `P2P_LOG`            | `p2p.P2pLog`         | false                  | No       | Flag to turn on/off network debug logs, and other traces within network layer  |
| `PUBSUB_TRACE`       | `p2p.NetworkTrace`   | false                  | No       | Flag to turn on/off pubsub trace logs    |
| `DISCOVERY_TRACE`    | `p2p.NetworkTrace`   | false                  | No       | Flag to turn on/off discovery trace logs |

## Pubsub Parameters

| Parameter            |  Description | Default Value (libp2p) | Value  | Comments |
| ---                  | ---          | ---                    | ---    | ---      |
| `lastSeen` | Specifies how long a message will be remembered as seen | `2min` | `6.4min` | The value (6.4m) is the duration of one epoch as message become redundant afterwards |
| `peerOutboundQueueSize` | The size of the queue that is used for outbound messages | `32` | `256` | The value was increased to support high load of messages, using a larger value might cause high memory usage |
| `validationQueueSize`    | The size of the queue that is used for validation of incoming messages  | `32` | `256` | The value was increased to support high load of messages, using a larger value might cause high memory usage |
| `validateThrottle`    | The upper bound on the number of active validation goroutines across all topics | `2048` | `8192` | The value was decreased to avoid high resource usage and reduce overloading on the node |
| `msgID`    | A function that calculates an identifier for messages, to be used across pubsub components | `(msg) => msg.GetFrom() + msg.GetSeqno()` | `(msg) => hash(msg.GetData())` | We use a custom function that returns the content hash, so we won't process the same message multiple times
| `subFilter`       | Allows to control the topics that the node will subscribe to | `nil` | accept topics with the same fork if the node has interest in them | - |
| `floodPublish`       | Force peer's own messages to be published to all known peers for the topic | `false` | `false` | When turned on, this feature ensures reliability and protects from eclipse attacks. On the other hand it floods the network with duplicated message and therefore it was turned off |
| **Gossipsub**            |           |                     |     |       |
| `D`       | Sets the optimal degree for a GossipSub topic mesh. D should be set somewhere between Dlo and Dhi | `6` | `5`| Decreased as we want to reduce loading coming from pubsub |
| `Dlo`       | Sets the lower bound on the number of peers we keep in a GossipSub topic mesh | `5` | `3` | Decreased as we want to reduce loading coming from pubsub |
| `Dhi`       | Sets the upper bound on the number of peers we keep in a GossipSub topic mesh | `12` | `9` | Decreased as we want to reduce loading coming from pubsub |
| `HeartbeatInterval` | Controls the time between heartbeats | `1s` | `700ms` | - |
| `HistoryLength` | Controls the size of the message cache used for gossip | `5` | `100` | - |
| `HistoryGossip` | Controls how many cached message ids we will advertise in IHAVE gossip messages | `3` | `3` | - |
| `MaxIHaveLength` | Sets the maximum number of messages to include in an IHAVE message | `5000` | `1000` | Decreased value to avoid ihave floods |
| `MaxIHaveMessages` | Sets the maximum number of IHAVE messages to accept from a peer within a heartbeat | `10` | `32` | Increased as we want messages to be sent in batches, to reduce the amount of requests |

**TODO:** add scoring