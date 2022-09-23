# Pubsub Scoring

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | SEP 22        |

This document contains information on the pubsub scoring strategies and configurations used in `SSV.Network`.

**NOTE** that the work on scoring is in progress. 
fine-tuning scoring params might introduce updates to this document.  

<br />

## Overview

`gossipsub v1.1` introduced pubsub [scoring](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#peer-scoring),
the idea is that a node maintains a score for each connected peer. 
The score helps to identify and ignore or prune malicious, badly connected or just slow peers 
(e.g. small machines that don't keep up).

Additional sources of information are available in libp2p documentation:
* [Score function](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#the-score-function)
* [Parameters overview](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#overview-of-new-parameters)

Scoring parameters calculations are based on [ETH 2](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c),
as it is a similar solution that has already shown robustness in scale.

<br />

## Peer Score Thresholds

Scoring thresholds are used by libp2p to determine whether a peer should be removed from topic's mesh,
penalized or even ignored if the score drops too low. \
See [this section](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#score-thresholds)
for more details regards the different thresholds.

`SSV.Network` uses the following scoring thresholds:

| Name | Description | Default Value (libp2p) | Default Value (SSV) | Default Value (ETH2) |
|---   | ---  | ---                    | ---                 | ---                  |
| `GossipThreshold` | below this threshold, gossip propagation is supressed | `-4000`   | `-4000` | `-4000` |
| `PublishThreshold`| below this threshold, message publishing is supressed | `-8000`   | `-8000` | `-8000` |
| `GraylistThreshold`| below this threshold, message processing is supressed altogether | `-16000`| `-16000` | `-16000` |
| `AcceptPXThreshold`| below this threshold, Peer Exchange will be ignored | `100`      | `100` | `100` |
| `OpportunisticGraftThreshold`| the median mesh score threshold before triggering opportunistic grafting | `5`  | `5` | `5` |

<br />

## Peer Score Params

Peer score params defines a set of general parameters to apply to all peers, 
regardless of what topic/s they are connected and sending messages.

As mentioned above, this solution is based on ETH2 
(see [function](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c#file-generate-scoring-params-py-L144))
for calculating peer score params.

```json
{
  "TopicScoreCap": 56.25, // = maxPositiveScore / 4.0
  "AppSpecificWeight": 1, // not used for now
  "IPColocationFactorWeight": -56.25, // = -topicScoreCap
  "IPColocationFactorThreshold": 10,
  "IPColocationFactorWhitelist": [],
  "BehaviourPenaltyWeight": -47.74333136903669, // = gossipThreshold / decayConvergence(rate = 10.0) ^ 2
  "BehaviourPenaltyThreshold": 10,
  "BehaviourPenaltyDecay": 0.6309573444801931,  // = decay oneEpoch*10
  "DecayInterval": 384000000000, // = oneEpoch
  "DecayToZero": 0.01,
  "RetainScore": 3840000000000, // = oneEpoch * 10
  "SeenMsgTTL": 385000000000
}
```

#### Input Arguments

The following arguments can be provided as an input to peer score params procedure:

- one epoch duration should set the time offset for scoring (`12 seconds * 32 slots`)
- msgID cache TTL
- IP colocation threshold the amount of peers allowed from the same IP (SSV default `10`)
- IP colocation weight helps to determine the penalty for multiple peers of the same IP
- IP whitelist to enable specific peers

<br />

## Topic Score Params

Topic score params enables to define some more specific values in the topic level, 
ensuring that each topic has its own set of config according to the topic nature.

As mentioned above, our solution for calculating score params is based on ETH2
[procedure](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c#file-generate-scoring-params-py-L174),
with some minor adaptations to ssv.

Topic score params accepts multiple types of arguments:

**Network level arguments**

* `ActiveValidators` is the amount of validators in the network
* `Subnets` is the number of subnets in the network
  * `Groups` is the amount of groups used in the network. **TBD** pending network topology 
* `OneEpochDuration` is used as a time-frame length to control scoring in a dynamic way
* `TotalTopicsWeight` is the weight of all the available topics in the network (decided + subnets)

**Topic level arguments**

* `TopicWeight` is the weight of the topic
* `ExpectedMsgRate` is the expected rate for the topic (in 12sec): \
`activeValidators / float64(slotsPerEpoch)`
* `InvalidMsgDecayTime` defines the decay for invalid messages (P4),
passing a zero value disables scoring of message validation.
* `FirstMsgDecayTime` defines the decay time for first message deliveries (P2)
* `MeshMsgDecayTime` defines the decay time for mesh message deliveries (P3)
  * `MeshMsgCapFactor` defines the factor to use to apply on the mesh message deliveries cap
  * `MeshMsgActivationTime` defines time in the mesh before penalties are being applied
* `D` is the gossip degree of the topic

Topic level arguments might change between different topics, currently there are two supported topics:
* **decided topic** for broadcasting decided messages across the network
* **subnet topic** for broadcasting consensus messages (of multiple validators)

### Decided Topic Params

Decided topic is subscribed by all peers, used for broadcasting decided messages across the network.
Params are based on `beacon aggregate and proof` topic (ETH2) due to the high rate of messages.

The following topic level arguments are used for constructing decided topic score params:

| Arg | Value | Comments |
| --- | -------- | --- |
|TotalTopicsWeight| `4.0 + 0.5 = 4.5` | The total weight consists of all subnets weight (4.0) and decided topic weight |
|TopicWeight| `0.5` | weights more than the subnet topics as its the most busiest topic |
|ExpectedMsgRate| `activeValidators / float64(slotsPerEpoch)` | rate per 12s |
|FirstMsgDecayTime| `1` | |
|MeshMsgDecayTime| `16` | |
|MeshMsgCapFactor| `32.0` | |
|MeshMsgActivationTime| One epoch duration | |

**NOTE** might be updated once we finalize this topic, after fine-tuning. 

Examples:

<details>
  <summary><b>1k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.5,
  "ExpectedMsgRate": 31.25,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 1,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 32,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.5,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 5.0687999999999995,
  "FirstMessageDeliveriesDecay": 0.01,
  "FirstMessageDeliveriesCap": 7.891414141414142,
  "MeshMessageDeliveriesWeight": -0.12514111392630922,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 59.96616130565809,
  "MeshMessageDeliveriesThreshold": 1.8739425408018153,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -0.12514111392630922,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

<details>
  <summary><b>10k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.5,
  "ExpectedMsgRate": 312.5,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 1,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 32,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.5,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 0.50688,
  "FirstMessageDeliveriesDecay": 0.01,
  "FirstMessageDeliveriesCap": 78.91414141414141,
  "MeshMessageDeliveriesWeight": -0.01222081190686613,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 191.89171617810592,
  "MeshMessageDeliveriesThreshold": 5.99661613056581,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -0.01222081190686613,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

<details>
  <summary><b>51k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.5,
  "ExpectedMsgRate": 1593.75,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 1,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 32,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.5,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 0.09938823529411765,
  "FirstMessageDeliveriesDecay": 0.01,
  "FirstMessageDeliveriesCap": 402.4621212121212,
  "MeshMessageDeliveriesWeight": -0.01222081190686613,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 191.89171617810592,
  "MeshMessageDeliveriesThreshold": 5.99661613056581,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -0.01222081190686613,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

[beacon aggregate and proof topic params for 51K validators](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c#file-51000-toml-L63)

### Subnet Topic Params

Subnet topics are used for exchanging consensus messages of multiple validator committees.
Params are based on `attestation subnets` (ETH2) due to their similar nature/orientation.

The following topic level arguments are used for constructing subnet topic score params:

| Arg | Value | Comments |
| --- | -------- | --- |
|TotalTopicsWeight| `4.0 + 0.5 = 4.5` (subnets + decided) | The total weight consists of all subnets weight (4.0) and decided topic weight |
|TopicWeight| `4.0 / subnetsCount` | all the subnets currently share the same weight, **TBD** higher weight for busier subnets |
|ExpectedMsgRate| `(activeValidators / subnetsCount) * msgsPerEpoch / slotsPerEpoch` | rate per 12s |
|FirstMsgDecayTime| `8` | |
|MeshMsgDecayTime| `16` | |
|MeshMsgCapFactor| `16.0` | |
|MeshMsgActivationTime| One epoch duration | |

**NOTE** might be updated once we finalize this topic, after fine-tuning

Examples:

<details>
  <summary><b>1k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.03125,
  "ExpectedMsgRate": 2.197265625,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 8,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 16,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.03125,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 31.869332124805872,
  "FirstMessageDeliveriesDecay": 0.5623413251903491,
  "FirstMessageDeliveriesCap": 1.255125141730395,
  "MeshMessageDeliveriesWeight": -450,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 2.1081853584020425,
  "MeshMessageDeliveriesThreshold": 0.13176158490012765,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -450,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

<details>
  <summary><b>10k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.03125,
  "ExpectedMsgRate": 21.97265625,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 8,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 16,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.03125,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 3.186933212480587,
  "FirstMessageDeliveriesDecay": 0.5623413251903491,
  "FirstMessageDeliveriesCap": 12.55125141730395,
  "MeshMessageDeliveriesWeight": -16.199996132888096,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 21.081853584020426,
  "MeshMessageDeliveriesThreshold": 1.3176158490012766,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -16.199996132888096,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

<details>
  <summary><b>51k validators</b></summary>

Topic arguments:
```json
{
  "TopicWeight": 0.03125,
  "ExpectedMsgRate": 112.060546875,
  "InvalidMsgDecayTime": 0,
  "FirstMsgDecayTime": 8,
  "MeshMsgDecayTime": 16,
  "MeshMsgCapFactor": 16,
  "MeshMsgActivationTime": 384000000000,
  "D": 8
}
```

Topic score params:
```json
{
  "TopicWeight": 0.03125,
  "TimeInMeshWeight": 0.03333333333333333,
  "TimeInMeshQuantum": 12000000000,
  "TimeInMeshCap": 300,
  "FirstMessageDeliveriesWeight": 0.624888865192272,
  "FirstMessageDeliveriesDecay": 0.5623413251903491,
  "FirstMessageDeliveriesCap": 64.01138222825014,
  "MeshMessageDeliveriesWeight": -0.7821319620394324,
  "MeshMessageDeliveriesDecay": 0.7498942093324558,
  "MeshMessageDeliveriesCap": 95.94585808905296,
  "MeshMessageDeliveriesThreshold": 5.99661613056581,
  "MeshMessageDeliveriesWindow": 2000000000,
  "MeshMessageDeliveriesActivation": 384000000000,
  "MeshFailurePenaltyWeight": -0.7821319620394324,
  "MeshFailurePenaltyDecay": 0.7498942093324558,
  "InvalidMessageDeliveriesWeight": 0,
  "InvalidMessageDeliveriesDecay": 0.1
}
```
</details>

[beacon attestations subnet topic params for 51K validators](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c#file-51000-toml-L87)
