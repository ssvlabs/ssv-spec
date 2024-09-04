# Pubsub Scoring

| Contributors                               | Status | Last Revision |
|:-------------------------------------------|:-------|:--------------|
| [@amir-blox](https://github.com/amir-blox) & Mathus Franco | Updated  | SEP 24        |

This document contains information on the pubsub scoring strategies and configurations used in `SSV.Network`.

## Overview

`GossipSub v1.1` introduced the pubsub [scoring](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#peer-scoring) feature.
With this feature, the node maintains a single score for each individual peer. The score is used to regulate what types of communication may be performed with the corresponding peer (e.g. graylisting a peer that sends too many malicious messages). The score for each peer is computed locally taking into consideration its observed behaviour and it's not shared in the network.

To use the scoring feature, we need to define several scoring parameters, as defined in its [documentation](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#the-score-function). To compute the parameters, we based our calculations in a [script for ETH 2](https://gist.github.com/blacktemplar/5c1862cb3f0e32a1a7fb0b25e79e6e2c)
since it has shown robustness in scale.

For a quick summary, the scoring is composed of topic scores ($S_1, S_2, S_3, S_{3b},$ and $S_4$) and general scores ($S_5, S_6,$ and $S_7$) and it's defined as

$$S(P) = max\{TopicScoreCap, \sum_{i}TopicScore(P,t_i)\} + S_5(P) + S_6(P) + S_7(P)$$

with

$$TopicScore(P,t) = w_{t} \times (S_1(P,t) + S_2(P,t) + S_3(P,t) + S_{3b}(P,t) + S_4(P,t))$$

Each different subscore ($S_1$, ..., $S_7$) serves for a different purpose (e.g. penalizing or rewarding depending on the peer behaviour) and each holds its own counter ($P_1$, ..., $P_7$). Not to classify a peer as good or bad for eternity, the counters decays according to their specific rates every `DecayInterval` interval.

Additional sources of information are available in libp2p documentation:
* [Score function](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#the-score-function)
* [Parameters overview](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#overview-of-new-parameters)
* [Score parameters Go file](https://github.com/libp2p/go-libp2p-pubsub/blob/b8a6a868adce87101b3f61cb1b1a644db627c59f/score_params.go)

<br />

## Peer Score Thresholds

Scoring thresholds are used by libp2p to determine whether a peer should be removed from topic's mesh,
penalized or even ignored if the score drops too low. \
See [this section](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#score-thresholds)
for more details regarding the different thresholds.

It's used the following scoring thresholds:

| Name | Description | Default Value (libp2p) | Default Value (SSV) | Default Value (ETH2) |
|---   | ---  | ---                    | ---                 | ---                  |
| `GossipThreshold` | below this threshold, no gossip is emitted to or received from the peer | `-4000`   | `-4000` | `-4000` |
| `PublishThreshold`| below this threshold, no self-published messages are propagated to the peer | `-8000`   | `-8000` | `-8000` |
| `GraylistThreshold`| below this threshold, all RPCs from/to the peer are ignored | `-16000`| `-16000` | `-16000` |
| `AcceptPXThreshold`| above this threshold, We accept supplied peers from the peer | `100`      | `100` | `100` |
| `OpportunisticGraftThreshold`| if the median mesh is below this threshold, opportunistic grafting is triggered | `5`  | `5` | `5` |

<br />

## Peer Score Params

[`PeerScoreParams`](https://github.com/libp2p/go-libp2p-pubsub/blob/b8a6a868adce87101b3f61cb1b1a644db627c59f/score_params.go#L66) defines a set of general parameters that are independent of specific topic behaviours.

| Parameter |  Value | Description |
| --- | --- | --- |
| `TopicScoreCap`| 32.72 | Cap for the sum of topic scores |
| `DecayInterval` | `32 * (12 * Time.Second)` (one epoch) | Time until applying the decay function |
| `DecayToZero` | 0.01 | Below this value, a counter is set to 0 |
| `RetainScore` | `100 * 32 * (12 * Time.Second)` (100 epochs) | Time we save a peer's score |
| `SeenMsgTTL` | `550 * HeartbeatInterval` (500 heartbeats) | Time until forgetting about a message |

### P5

| Parameter | Value |
| --- | --- |
| `AppSpecificWeight` | 0 |

$P_5$ is computed as follows:
$$S_5(P) = w_5 * P_5$$

Notes:
- P5 is not currently being used. So, the weight is set to 0.

### P6
| Parameter | Value |
| --- | --- |
| `IPColocationFactorWhitelist` | 0 |
| `IPColocationFactorWeight` | -32.72 (-`TopicScoreCap`) |
| `IPColocationFactorThreshold` | 10 |

$S_6$ is computed as follows:
$$S_6(P) = w_6 * (max\{0, P_6 - th_6\})^2$$

Notes:
- P6 penalizes a peer if it shares the same IP with too many other peers. It's useful against Sybil attacks.
- The P6 counter is always the total number of peers with the peer's IP and it doesn't decay.
- Whitelist is the list of IPs not to be penalized. As default, this list is empty.
- Threshold is the value after which the peer starts to get penalized. In our case, we allow at most 10.
- The weight is defined so that, after a peer exceeds the threshold, its total score is at most 0.


### P7
| Parameter | Value |
| --- | --- |
| `BehaviourPenaltyDecay` | 0.6309573444801932 |
| `BehaviourPenaltyWeight` | -8.986961427779512 |
| `BehaviourPenaltyThreshold` | 6 |

$S_7$ is computed as follows:
$$S_7(P) = w_7 * (max\{0, P_7 - th_7\})^2$$

Notes:
- P7 applies penalizations due to misbehaviours according to [GossipSub router's rules](https://github.com/libp2p/go-libp2p-pubsub/blob/b8a6a868adce87101b3f61cb1b1a644db627c59f/score_params.go#L93).
- The decay is set to that 1 is transformed into 0.01 after 10 epochs (considering the `DecayInterval`).
- Threshold is the counter value after which the peer starts to get penalized. In our case, it's defined as 6.
- The weight is set so that *a peer may keep indefinitely doing at most 10 behaviour penalties per epoch and, at infinity, it will reach `GossipThreshold`*.

If the counter keeps getting incremented by $r$ and decayed by $d$, it will be, at infinity,
$$P_7 = r + r * \times d + r * \times d^2 + ... = \frac{r}{1-d}$$

So, a peer doing 10 misbehaviours per epoch will have, at infinity,
$$P_7 = \frac{10}{1 - 0.6309573444801932} = 27.097138638119553$$

Thus, we want that
$$GossipThreshold = w_7 * (P_7 - th_7)^2$$
$$w_7 = \frac{GossipThreshold}{(P_7 - th_7)^2} = -8.986961427779512$$

<br />

## Topic Score Params

[`TopicScoreParams`](https://github.com/libp2p/go-libp2p-pubsub/blob/b8a6a868adce87101b3f61cb1b1a644db627c59f/score_params.go#L117) enables us to score the peer depending on more specific topic behaviours.

Different topics may have different parameters values. Here, we define a common set of formulas that are applied for all topics.

| Parameter |  Value |
| --- | --- |
| `TopicWeight`| 0.03125 |

Notes:
- We define the total topics weight to be $4.0$. Since we have $128$ topics, each topic has a weight of $\frac{4.0}{128} = 0.03125$.

### P1

| Parameter | Value |
| --- | --- |
| `TimeInMeshWeight` | 0.0333... |
| `TimeInMeshQuantum` | 12 |
| `TimeInMeshCap` | 300 |

$S_1$ is computed as follows:
$$S_1(P) = w_1 * P_1$$

with $P_1$ being capped by $Cap_1$ (i.e. $P_1\leftarrow max\{P_1, Cap_1\}$).

Notes:
- P1 rewards a peer for staying a long period in the node's mesh.
- The $P_1$ counter is computed by the time it's in the mesh divided by `TimeInMeshQuantum`. The $P_1$ is not decayed.
- We defined the quantum to be a slot duration (12 seconds).
- We defined the $Cap_1$ to be the counter value after one hour. So, $Cap_1 = \frac{60 * 60}{12} = 300$.
- We defined the weight to be the maximum value for P1, 10, divided by the cap. So, $w_1 = \frac{10}{Cap_1} = 0.0333...$

### P2

| Parameter | Value |
| --- | --- |
| `FirstMessageDeliveriesDecay` | 0.3162277660168379 |
| `FirstMessageDeliveriesWeight` | Check below |
| `FirstMessageDeliveriesCap` | Check below |


$S_2$ is computed as follows:
$$S_2(P) = w_2 * P_2$$

with $P_2$ being capped by $Cap_2$ (i.e. $P_2\leftarrow max\{P_2, Cap_2\}$).

Notes:
- P2 rewards a peer that sends valid messages that the node has never seen before.
- The $P_2$ counter is incremented whenever a never-seen valid message is received from the peer.
- The decay is defined so that $P_2$ decays from 1 to 0.01 in 4 epochs.
- The $Cap_2$ is defined as the counter value that a peer gets when sending the double of its expectation.

So, let's say that the expected message rate for a topic is $m$ and the mesh size is $D$. The peer is expected to send, in average, $m/D$ never-seen valid messages.
If it keeps sending twice this value, $2m/D$, the counter will be, at infinity,
$$P_2 = 2m/D + 2m/D \times d + 2m/D \times d^2 + ...$$

where $d$ is the decay. Thus, we have that
$$P_2 = \frac{2m/D}{1-d}$$

- The weight is defined as the maximum value for P2, 80, divided by the cap. So $w_2 = \frac{80}{Cap_2}$.

Here, we present no numerical values for $Cap_2$ and $w_2$ since they depend on the topic's message rate. Check the [message_rate_estimation.go](./messagerate/message_rate_estimation.go) file to see how this can be computed.

### P3 and P3b

The P3 and P3b scores are used to penalize peers that doesn't fullfil an expectation of message delivery in the topic. For now, these scores are unused (i.e. the weights are set to 0).

### P4

| Parameter | Value |
| --- | --- |
| `InvalidMessageDeliveriesDecay` | 0.954992586021436 |
| `InvalidMessageDeliveriesWeight` | -1280 |

$S_4$ is computed as follows:
$$S_4(P) = w_4 * P_4^2$$

Notes:
- P4 is used to penalize peers that send malicious messages. A message are considered malicious if its validation output is `ValidationReject`.
- The decay is set to that the $P_4$ counter decays from 1 to 0.01 in 100 epochs.
- The weight is defined so that, if the counter reaches 20, its score goes to `GraylistThreshold`. Thus

$$GraylistThreshold = w_t \times w_4 \times P_4^2$$
$$w_4 = \frac{GraylistThreshold}{w_t \times P_4^2} = \frac{-16000}{0.03125 \times 20^2} = -1280$$

where $w_t$ is the topic's weight.
