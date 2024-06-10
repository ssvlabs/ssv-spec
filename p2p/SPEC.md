# SSV Specifications - Networking

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | SEP  22        |

This document contains the networking specification for `SSV.Network`.

## Overview

- [Fundamentals](#fundamentals)
    - [Stack](#stack)
    - [Transport](#transport)
    - [Messaging](#messaging)
    - [Network Peers](#network-peers)
    - [Identity](#identity)
    - [Network Discovery](#network-discovery)
    - [Peer Scoring](#peer-scoring)
- [Wire](#wire)
    - [Consensus](#consensus-protocol)
    - [Sync](#sync-protocols)
    - [Handshake](#handshake-protocol)
- [Network Layer ](#networking)
    - [PubSub](#pubsub)
    - [PubSub Scoring](#pubsub-scoring)
    - [Message Validation](#pubsub-validation)
    - [Discovery](#discovery)
    - [Subnets](#subnets)
    - [Peers Connectivity](#peers-connectivity)
    - [Peers Balancing](#peers-balancing)
    - [Connection Gating](#connection-gating)
    - [Security](#security)
    - [Configurations](#configurations)
    - [Forks](#forks)
    - [Relayers](#relayers)

## Fundamentals

### Stack

`SSV.Network` is a permission-less P2P network, consists of operator nodes that are grouped in multiple subnets,
signing validators' duties after reaching to consensus for each duty.

The networking layer is built with [Libp2p](https://libp2p.io/),
a modular framework for P2P networking that is used by multiple decentralized projects, including ETH 2.0.

### Transport

Network peers must support the following transports:
- `TCP` is used by libp2p for setting up communication channels between peers.
  default port: `12001`
- `UDP` is used for discovery by [discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md).
  default port: `13001`

[go-libp2p-noise](https://github.com/libp2p/go-libp2p-noise)
is used to secure transport, for more details see [noise protocol](https://noiseprotocol.org/noise.html)
and [libp2p spec](https://github.com/libp2p/specs/blob/master/noise/README.md).

Multiplexing of protocols over channels is achieved using [yamux](https://github.com/libp2p/go-libp2p-yamux) protocol.


### Messaging

Messages in the network are formatted with `protobuf` (NOTE: `v0` messages are encoded/decoded with JSON),
and being transported p2p with one of the following methods:

**Streams**

Libp2p allows to create a bidirectional stream between two peers and implement the corresponding wire messaging protocol.

[Streams](https://ipfs.io/ipfs/QmVqNrDfr2dxzQUo4VN3zhG4NV78uYFmRpgSktWDc2eeh2/specs/7-properties/#71-communication-model---streams)
are used in the network for direct messages between peers.

**PubSub**

GossipSub ([v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md))
is the pubsub protocol used in `SSV.Network`

The main purpose is for broadcasting messages to a group (AKA subnet) of nodes. \
In addition, the machinery helps to determine liveliness and maintain peers scoring.


### Network Peers

There are several types of nodes in the network:

`Operator` is responsible for executing validators duties. \
It holds registry data and the validators consensus data.

`Exporter` is responsible for collecting and exporting information from the network. \
It collects registry data and consensus data (decided messages) of all the validators in the network.

`Bootnode` is a public peer which is responsible for helping new peers to find other peers in the network.
It has a stable ENR that is provided with default configuration, so other peers could join the network easily.


### Identity

Identity in the network is based on two types of keys:

`Network Key` is used to create network/[libp2p identity](https://docs.libp2p.io/concepts/peer-id) (`peer.ID`),
will be used by all network peers to set up a secured connection. \
Unless provided, the key will be generated and stored locally for future use,
and can be revoked in case it was compromised.

`Operator Key` is used for decryption of share's keys that are used for signing/verifying consensus messages and duties. \
Exporter and Bootnode does not hold this key.


### Network Discovery

[discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md)
is used in `SSV.Network` as the discovery component.

More information is available in [Discovery section](#discovery)


### Peer Scoring

Peer scoring in `SSV.Network` is how we protect the network from bad peers,
by scoring them according to a predefined set of scores.

For more info please refer to [Pubsub Scoring](#pubsub-scoring) section


------


## Wire

All the messages that are being transmitted over the network are wrapped with the following structure:

<details>
  <summary><b>protobuf</b></summary>

  ```protobuf
syntax = "proto3";
import "gogo.proto";

// SignedMessage holds a message and it's corresponding signature
message SSVMessage {
  // type of the message
  MsgType msg_type = 1 [(gogoproto.nullable) = false];
  // id of the message
  bytes msg_id     = 2 [(gogoproto.nullable) = false];
  // message data (encoded)
  bytes data      = 3 [(gogoproto.nullable) = false];
}

// MsgType is an enum that represents the type of message 
enum MsgType {
  // consensus/QBFT messages
  Consensus              = 0;
  // partial signatures sent post consensus
  Signature = 1;
}
```

</details>

Note that all pubsub messages in the network are wrapped with libp2p's message structure 
([see pubsub RPC](https://github.com/libp2p/specs/blob/master/pubsub/README.md#the-rpc)).

## Consensus Protocol

`IBFT`/`QBFT` consensus protocol is used to govern `SSV` network.
`QBFT` ensures that consensus can be reached by a committee of `n`
operator nodes while tolerating a certain amount of `f` faulty nodes as defined by `n ≥ 3f + 1`.

As part of the algorithm, nodes are exchanging messages with other nodes in the committee. \
Once the committee reaches consensus, the nodes will publish the decided message across the network.

Consensus messages are being sent in the network over pubsub topics (see [subnets](#subnets))

More information regarding the protocol can be found in [iBFT annotated paper (By Blox)](/ibft/IBFT.md)

### Message Structure

Messages structure and more information can be found in 
[ssv-spec/types](https://github.com/ssvlabs/ssv-spec/blob/main/types/messages.go).

---

## Handshake Protocol

`/ssv/info/0.0.1`

The handshake protocol allows peers to identify, by exchanging signed information AKA `NodeInfo`. \
It must be performed for every connection, and check that the other node is on the same ssv network.

[go-libp2p-core/record.Record](https://github.com/libp2p/go-libp2p-core/blob/master/record/record.go) and
[go-libp2p-core/record.Envelope](https://github.com/libp2p/go-libp2p-core/blob/master/record/envelope.go)
are a common utilities that libp2p provides for these cases,
therefore they are used for signing and verification of `NodeInfo` messages in SSV. 

<details>
  <summary><b>protobuf</b></summary>

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
</details>

<br />
<br />

## Networking


### Pubsub

The main purpose is for broadcasting messages to a group (AKA subnet) of nodes, 
using a gossip approach to avoid the overhead of maintaining multiple direct connections. \
In addition, the following are achieved as well:

- subscriptions metadata helps to get liveliness information of nodes
- pubsub scoring enables to prune bad/malicious peers based on network behavior and application-specific rules

The following sections details on how pubsub is used in `SSV.network`.
In addition, parameters configuration is described [here](./CONFIG.md#pubsub-parameters).


### Pubsub Scoring

`gossipsub v1.1` introduced pubsub [scoring](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#peer-scoring),
the idea is that each individual peer maintains a score for other peers.
The score is locally computed by each individual peer based on observed behaviour and is not shared.

**NOTE** that [topic scores](./SCORING.md#topic-score-params), [peer scores](./SCORING.md#peer-score-params) 
and [thresholds](./SCORING.md#peer-score-thresholds) are detailed in the sibling scoring spec document. 


### Pubsub Validation

Message validation is applied on the topic level.
Each incoming message will be validated to avoid relaying bad messages,
which could affect peer score.

[Extended Validators](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#extended-validators)
allows the application to aid in the gossipsub peer-scoring scheme.
We utilize `ACCEPT`, `REJECT`, and `IGNORE` as the result of the validation.

#### Basic Validation

This validation pipeline is the baseline and will be applied for all pubsub topics. It consists of the following rules:

- `ACCEPT` message from my peer
- `REJECT` empty message
- `REJECT` message with corrupted or invalid top-level structure

**NOTE** any message will be decoded only once as part of the basic validation.

**NOTE** As of the time this spec was written, additional validation is performed by QBFT components in an async way.
As messages might pass the base validation but fail at a later point, signing policy of pubsub is turned on 
to ensure authenticity of pubsub message senders. Once a more complete validation is added, we will reduce
the signing policy as it becomes redundant.

<br />

### Subnets

Consensus messages are being sent in the network over a pubsub topic.

A subnet of peers consists of
operators that are responsible for multiple committees,
reusing the same topic to communicate on behalf of multiple validators. \
Operator nodes will validate and store highest decided messages and last change round messages
of all the committees in the subnets they participate.

`Decided` topic is used for propagation of decided messages across all the nodes in the network,
that should store the last decided message of each committee in the network.
Having such redundancy of decided messages helps to maintain the multiple states (per validator) across the network.

**Committee Mapping**

The committee ID is a hash of the operators IDs.
This helps to distribute committees across subnets in a balanced, distributed way:

`committee_id % num_of_subnets`

Deterministic mapping is ensured as long as the number of subnets doesn't change,
therefore it's a fixed number.

<br />

### Discovery

[discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md)
is a system for finding other participants in a peer-to-peer network,
it is used in `SSV.network` to complement discovery.

DiscV5 works on top of UDP, it uses a DHT to store node records (`ENR`) of discovered peers.
It allows walking randomly on the nodes in the table, and act according to application needs.

In SSV, new nodes are filtered by score, that is calculated from past behavior and properties (`ENR` entries).
If the score is above threshold, the node tries to connect and handshake with the new node.

As discv5 is standalone (i.e. not depends on libp2p), the communication is encrypted and authenticated using session keys,
established in a separate [handshake process](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md#sessions).

**Bootnode**

A peer that have a public, static ENR to enable new peers to join the network. For the sake of flexibility,
bootnode/s ENR values are configurable and can be changed on demand by operators. \
Bootnode doesn't start a libp2p host for TCP communication,
its role ends once a new peer finds existing peers in the network.

#### ENR

[Ethereum Node Records](https://github.com/ethereum/devp2p/blob/master/enr.md) is a format that holds peer information.
Records contain a signature, sequence (for republishing record) and arbitrary key/value pairs.

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


#### Subnets Discovery

As `ENR` has a size limit (`< 300` bytes),
discv5 won't support multiple key/value pairs for storing subnets of operators,
which could have made it easier to find nodes with common subnets.

Instead, an array of flags is used,
representing the assignment of subnets for an operator.
Similar to how it implemented in Ethereum 2.0
[phase 0](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#attestation-subnet-bitfield).

[Discovery v5.2](https://github.com/ethereum/devp2p/milestone/3) will introduce
[Topic Index](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#the-topic-index)
which helps to lookup nodes by their advertised topics. In SSV these topics would be the operator's subnets. \
For more information:
- [DiscV5 Theory > Topic Advertisement](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md#topic-advertisement)
- [discv5: topic index design thread](https://github.com/ethereum/devp2p/issues/136)

See [Consensus specs > phase 0 > p2p interface](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/p2p-interface.md#why-are-we-using-discv5-and-not-libp2p-kademlia-dht)
for details on why discv5 was chosen over libp2p Kad DHT in Ethereum.

<br />

### Peers Connectivity

In a fully-connected network, where each peer is connected to all other peers in the network,
running nodes will consume many resources to process all network related tasks e.g. parsing, peers management etc.

Limiting connected peers count should reduce the network throughput and resource consumption, 
while the node still receives messages via gossiping.

#### Peers Balancing

Once reached to peer limit, the node will stop looking for new nodes,
and won't accept incoming connections from relevant peers. \
Instead, there is a peer balancing procedure that ensures (over time) that connections are being replaced, 
and we'll try form a balanced set of peers across all the subnets of our interest.

Tagging is done using `go-libp2p-core/connmgr.ConnManager` interface to protect / unprotect peers. \
Once peers were tagged, we trim the peers that are unprotected

Peers balancing includes the following steps:

**NOTE:** in this process we rely on information that was sent by the peer,
and on the internal state of connected peers (provided by libp2p's host and pubsub router).

1. continue if we reached peers limit in the node level, or stop otherwise.
2. tag best `n` peers where `n = maxPeers - 1`
  1. calculate scores for subnets:
    1. subnet w/o peers - `2`
    2. subnet w/ less than the minimum (<= 2) - `1`
    3. subnet w/ overflow of peers (>= 5) - `-1`
  2. calculate peers scores according to their subnets,
     by a counting the subnets scores and giving bonus score for peers with multiple shared subnets.
  3. **TBD** pubsub scoring will be taken into account (once added)
3. trim untagged peers

#### Connection Gating

Connection Gating allows safeguarding against bad/pruned peers that tries to reconnect multiple times.
Inbound and outbound connections are intercepted and being checked before other components process the connection.

IP limiting is also a procedure that will be applied by the connection gater in future versions.

See libp2p's [ConnectionGater](https://github.com/libp2p/go-libp2p-core/blob/master/connmgr/gater.go)
interface for more info.

<br />

### Security

As mentioned above, `gossipsub v1.1` comes with a set of tools for protecting the network from bad peers,
e.g. msg validation will protect from malicious messages, while scoring will cause the bad peers to get pruned. \
[Gossipsub v1.1 Evaluation Report](https://gateway.ipfs.io/ipfs/QmRAFP5DBnvNjdYSbWhEhVRJJDFCLpPyvew5GwCCB4VxM4)
describes some potential attacks and how they are mitigated.

Connection gating does IP limiting and protects against known, bad peers.
The gater is invoked in an early stage, before the other components processes
the request to avoid redundant resources allocations.

In addition, the discovery system is naturally a good candidate for security problems. \
DiscV5 specs specifies potential vulnerabilities in their system and how they were (or will be) mitigated,
see [DiscV5 Rationale > security-goals](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#security-goals).
The major ones includes routing table pollution, traffic redirection, spamming or replayed messages.
Peers with bad scores will be filtered during discovery, ensuring that attacking peers are
known and ignored all over the system.

<br />

### Configurations

Configuration set used in `SSV.network` is described in [this document](./CONFIG.md).

<br />

### Forks

Future network forks will follow the general forks mechanism and design in SSV,
where some procedures are called in the context of current fork version.

See [network forks](./FORKS.md) for more information.

<br />
