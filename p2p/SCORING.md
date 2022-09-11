# Scoring

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) | DRAFT  | JUL 22        |

This document contains information on the scoring strategies and configurations used in `SSV.Network`.

<br />

## Overview

`gossipsub v1.1` introduced pubsub [scoring](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#peer-scoring),
the idea is that a node maintains a score for each connected peer. 
The score helps to identify and ignore or prune malicious, badly connected and slow peers.

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
  "TopicScoreCap": 56.25, // maxPositiveScore / 4.0
  "AppSpecificWeight": 1, // not used for now
  "IPColocationFactorWeight": -56.25, // -topicScoreCap
  "IPColocationFactorThreshold": 10,
  "IPColocationFactorWhitelist": [],
  "BehaviourPenaltyWeight": -47.74333136903669, // gossipThreshold / decayConvergence(rate = 10.0) ^ 2
  "BehaviourPenaltyThreshold": 10,
  "BehaviourPenaltyDecay": 0.6309573444801931,  // decay oneEpoch*10
  "DecayInterval": 384000000000, // oneEpoch
  "DecayToZero": 0.01,
  "RetainScore": 3840000000000, // oneEpoch * 10
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

As mentioned above, this solution is based on ETH2
(see [function](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c#file-generate-scoring-params-py-L174))
for calculating topic score params.

### Decided Topic Params

Decided topic is based on `aggregation topic` (ETH2) due to the high rate of messages.

### Subnet Topic Params

Subnet topics are based on `attestation subnets` (ETH2) due to their similar nature/orientation.
