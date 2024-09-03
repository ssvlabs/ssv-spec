# Network Configurations

|                        Contributors                         | Last Revision |
| :---------------------------------------------------------- | :------------ |
| [@amir-blox](https://github.com/amir-blox) & Matheus Franco | September 24  |

These are the configurations for the frameworks used in the `SSV.Network`.

> [!NOTE]
> Scoring configuration is detailed in a [separated document](./SCORING.md).

## Node Config

### General

Below is a set of general configurations for an `SSV` node.

| ENV                  | YAML                 | Default Value (SSV)    | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `NETWORK_PRIVATE_KEY`| `NetworkPrivateKey`  | -                      | No       | Key to use for libp2p/network identity |
| `TCP_PORT`           | `p2p.TcpPort`        | `13001`                | No       | TCP port to use as libp2p transport    |
| `UDP_PORT`           | `p2p.UdpPort`        | `12001`                | No       | UDP port to use for discovery          |
| `HOST_ADDRESS`       | `p2p.HostAddress`    | -                      | No       | External IP address                    |
| `HOST_DNS`           | `p2p.HostDNS`        | -                      | No       | External DNS address                   |
| `P2P_REQUEST_TIMEOUT`| `p2p.RequestTimeout` | `7s`                   | No       | Timeout for stream requests            |
| `P2P_MAX_BATCH_RESPONSE` |`p2p.MaxBatchResponse`| `25`               | No       | Maximum number of returned objects in a batch for sync protocols |

### Discovery

The following set of configurations is used for discovery.

| ENV                  | YAML                 | Default Value (SSV)    | Required |  Description                           |
| ---                  | ---                  | ---                    | ---      | ---                                    |
| `NETWORK_ID`         | `p2p.NetworkID`      | `ssv-testnet-v2`       | No       | ID of the network that the node runs on|
| `BOOTNODES`          | `p2p.Bootnodes`      | Bootnode (Testnet)     | No       | ENR of network bootnodes               |
| `SUBNETS`            | `p2p.Subnets`        | -                      | No       | Subnets that this node will join upon start, represented as hex string|

## GossipSub Parameters

To understand the GossipSub parameters, visit the [v1.0](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.0.md) and [v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#overview-of-new-parameters) documentations.

We compare the parameter values with both libp2p and ETH2.0 default values. \
For more information on libp2p defaults, see
[gossipsub-v1.0 spec > parameters](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.0.md#parameters). \
For more information on ETH2.0 configuration, please refer to
[their p2p spec](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#the-gossip-domain-gossipsub).

| Parameter                | Default Value (libp2p) | Default Value (SSV) | Default Value (ETH2) |
|--------------------------|------------------------|---------------------|----------------------|
| `FloodPublish`           | `false`                | `false`             | `false`              |
| `MessageSignaturePolicy` | `StrictNoSign`         | `StrictSign`        | `StrictNoSign`       |
| `OutboundQueueSize`      | `32`                   | `512`               | `600`                |
| `ValidateQueueSize`      | `32`                   | `512`               | `600`                |
| `ValidateThrottle`       | `8192`                 | `8192`              | `8192`               |
| `D`                      | `6`                    | `8`                 | `8`                  |
| `Dlo`                    | `5`                    | `6`                 | `6`                  |
| `Dhi`                    | `12`                   | `12`                | `12`                 |
| `Dout`                   | `2`                    | `4`                 | `2`                  |
| `HeartbeatInterval`      | `1s`                   | `700ms`             | `700ms`              |
| `HistoryLength`          | `5`                    | `6`                 | `6`                  |
| `HistoryGossip`          | `3`                    | `4`                 | `3`                  |
| `MaxIHaveLength`         | `5000`                 | `1500`              | `5000`               |
| `MaxIHaveMessages`       | `10`                   | `32`                | `10`                 |
| `SeenMsgTTL`             | `2min`                 | `6.4min`            | `6.4min`             |


### The MsgID function

To detect duplicated messages, GossipSub uses a `MsgID` function that computes the identifier of an upcoming message.
If the identifier has already been seen, the router drops the message.

The framework allows the application to define such function.
If not defined, by default, it will use the (`sender`,`msg_seq`(sequence number)) tuple.

Since we don't enable the sequence number field, we define a message's identifier by the hash of its content.

```go
func MsgID(msg) string {
    return string(sha256(msg.GetData())[:20])
}
```

See [pubsub spec > message identification](https://github.com/libp2p/specs/blob/master/pubsub/README.md#message-identification) for more details.
