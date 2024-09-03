# SSV Specifications - Networking

|                        Contributors                         | Last Revision |
| :---------------------------------------------------------- | :------------ |
| [@amir-blox](https://github.com/amir-blox) & Matheus Franco | September 24 |

This document contains the networking specification for `SSV.Network`.

## Overview

- [SSV Specifications - Networking](#ssv-specifications---networking)
  - [Overview](#overview)
  - [Fundamentals](#fundamentals)
    - [Stack](#stack)
    - [Types of Nodes](#types-of-nodes)
    - [Identity](#identity)
  - [GossipSub](#gossipsub)
    - [Topic subscription](#topic-subscription)
    - [Messages on the wire](#messages-on-the-wire)
    - [Scoring](#scoring)
    - [Message Validation](#message-validation)
  - [Discovery](#discovery)
    - [Handshake Protocol](#handshake-protocol)
    - [Bootnodes](#bootnodes)
    - [ENR](#enr)
    - [Subnets Discovery](#subnets-discovery)
  - [Security](#security)
  - [Forks](#forks)

## Fundamentals

### Stack

`SSV.Network` is a permissionless P2P network that consists of operator nodes.

The discovery layer, used to discover peers, is built upon the [Discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md) protocol.

The networking layer, used to exchange messages, is built upon the [GossipSub](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md) protocol.

Both frameworks are used by multiple decentralized projects including ETH 2.0.

The configuration used for each framework is described in [configuration document](./CONFIG.md).

### Types of Nodes

There are several types of ndodes in the network:

`Operator` is responsible for executing validators duties. \
It holds registry data and the validators consensus data.

`Exporter` is responsible for collecting and exporting information from the network. \
It collects registry data and consensus data (decided messages) of all the validators in the network.

`Bootnode` is a public peer which is responsible for helping new peers to find other peers in the network.
It has a stable ENR that is provided with default configuration, so other peers could join the network easily.

### Identity

Identity in the network is based on two types of keys:

`Network Key` is used to create the network identity [`peer.ID`](https://docs.libp2p.io/concepts/peer-id).
It will be used by all network peers to set up a secured connection. \
Unless provided, it will be generated and stored locally for future use. \
It can be revoked in case it is compromised.

`Operator Key` is used for decryption of share's keys.

<br />

## GossipSub

To exchange protocol messages, the operators use the [GossipSub](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md) protocol, which is a P2P gossip protocol that follows the publish/subscribe architecture.

### Topic subscription

Operators subscribe to topics according to the committees that it belongs to. This helps balancing the load of messages that the operator listens to. For a given committee, the operator should subscribe to the topic computed as

`committee_id % num_of_subnets`

where [`CommitteeID`](./../types/docs/README.md#operatorid-and-committeeid) is an identifier for the committee (in practice, it's a hash of the list of operators' IDs). Notice that deterministic mapping is ensured as long as the number of subnets doesn't change.


### Messages on the wire

All the messages that are being transmitted over the wire are wrapped with libp2p's [pubsub RPC structure](https://github.com/libp2p/go-libp2p-pubsub/blob/master/pubsub.go#L244).

Inside this structure, it's added an encoded [`SignedSSVMessage`](./../types/docs/README.md#signedssvmessage-and-ssvmessage) object used for the [SSV protocol](./../ssv/README.md).

### Scoring

TODO.

### Message Validation

TODO.

<br />

## Discovery

For discovery, `SSV.Network` uses the [Discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md) protocol for finding other participants in the network.

DiscV5 uses UDP communication and a DHT to store node records (`ENR`) of discovered peers.
It allows walking randomly on the nodes in the table, and act according to application needs.

In SSV, new nodes are filtered by score, that is calculated from past behavior and properties (`ENR` entries).
If the score is above a certain threshold, the node tries to connect and handshake with the new node.

As discv5 is standalone (i.e. not depends on libp2p), the communication is encrypted and authenticated using session keys,
established in a separate [handshake process](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md#sessions).


### Handshake Protocol

`/ssv/info/0.0.1`

The handshake protocol allows peers to connect by exchanging `NodeInfo` data. \
It ensures that the other node is on the same SSV network.

[go-libp2p-core/record.Record](https://github.com/libp2p/go-libp2p-core/blob/master/record/record.go) and
[go-libp2p-core/record.Envelope](https://github.com/libp2p/go-libp2p-core/blob/master/record/envelope.go)
are common utilities that libp2p provides for these cases.
Therefore, they are used for signing and verifying `NodeInfo` messages in SSV.


```protobuf
syntax = "proto3";
import "gogo.proto";

// NodeInfo contains node's information
message NodeInfo {
  // network_id is the id of the node's network
  bytes network_id   = 1 [(gogoproto.nullable) = false];
  // fork_version is the current fork used by the nodee
  bytes fork_version = 2 [(gogoproto.nullable) = true];
  // metadata of the node
  Metadata metadata  = 3;
}

// Metadata holds node's general information
message Metadata {
  // operator_id of the node
  string operator_id     = 1;
  // subnets is an hex string that represents subnets with interest
  string subnets         = 2;
  // node_version is the current ssv-node version
  string node_version    = 3;
  // execution_node is the eth1 node used by the node
  string execution_node  = 4;
  // consensus_node is the eth2 node used by the node
  string consensus_node = 5;
}
```

### Bootnodes

The discovery protocol counts with bootnodes (a peer that have a publica nd static ENR) to enable new peers to join the network. \
For the sake of flexibility,
bootnode/s ENR values are configurable and can be changed on demand by operators. \
A bootnode doesn't start a libp2p host for TCP communication,
its role ends once a new peer finds existing peers in the network.

### ENR

The [Ethereum Node Records](https://github.com/ethereum/devp2p/blob/master/enr.md) (ENR) is a structured data that holds peer information.
Records contain a signature, a sequence number (for republishing record) and arbitrary key/value pairs.

`ENR` structure in `SSV.Network` consists of the following key/value pairs:

| Key         | Description                                                    |
|:------------|:---------------------------------------------------------------|
| `id`        | name of identity scheme, e.g. "v4"                             |
| `secp256k1` | compressed secp256k1 public key of the network key, 33 bytes   |
| `ip`        | IPv4 address, 4 bytes                                          |
| `tcp`       | TCP port, big endian integer                                   |
| `udp`       | UDP port, big endian integer                                   |
| `type`      | node type, integer; 1 (operator), 2 (exporter), 3 (bootnode)   |
| `forkv`     | fork version, integer                                          |
| `subnets`   | bitlist, 0 for irrelevant and 1 for assigned subnet            |


### Subnets Discovery

As `ENR` has a size limit (`< 300` bytes),
discv5 won't support multiple key/value pairs for storing subnets of operators,
which could have made it easier to find nodes with common subnets.

Instead, an array of flags is used to
represent the assignment of subnets for an operator, similar to how it is implemented in Ethereum 2.0
[phase 0](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#attestation-subnet-bitfield).

[Discovery v5.2](https://github.com/ethereum/devp2p/milestone/3) will introduce the
[Topic Index](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#the-topic-index)
feature which helps to lookup nodes by their advertised topics. In SSV these topics would be the operator's subnets. \

For more information:
- [DiscV5 Theory > Topic Advertisement](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md#topic-advertisement)
- [discv5: topic index design thread](https://github.com/ethereum/devp2p/issues/136)

See [Consensus specs > phase 0 > p2p interface](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#why-are-we-using-discv5-and-not-libp2p-kademlia-dht)
for details on why discv5 was chosen over libp2p Kad DHT in Ethereum.

<br />

## Security

As mentioned above, `GossipSub v1.1` comes with a set of tools for protecting the network from bad peers.
E.g., msg validation will protect from malicious messages, while scoring will cause the bad peers to get pruned. \
[Gossipsub v1.1 Evaluation Report](https://gateway.ipfs.io/ipfs/QmRAFP5DBnvNjdYSbWhEhVRJJDFCLpPyvew5GwCCB4VxM4)
describes some potential attacks and how they are mitigated.

Connection gating does IP limiting and protects against known bad peers.
The gater is invoked in an early stage, before the other components processes
the request to avoid redundant resources allocations.

In addition, the discovery system is naturally a good candidate for security problems. \
DiscV5 spec specifies potential vulnerabilities in their system and how they were (or will be) mitigated,
see [DiscV5 Rationale > security-goals](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#security-goals).
The major ones include routing table pollution, traffic redirection, spamming or replayed messages.
Peers with bad scores will be filtered during discovery, ensuring that attacking peers are
known and ignored all over the system.

<br />

## Forks

Future network forks will follow the general forks mechanism and design in SSV.

See the [network forks file](./FORKS.md) for more information.

<br />
