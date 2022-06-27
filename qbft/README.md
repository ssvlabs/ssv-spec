
# QBFT

## Introduction
This is a spec implementation for the QBFT protocol, following [formal verification spec](https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification) / [github repo](https://github.com/ConsenSys/qbft-formal-spec-and-verification).

## Important note on message processing
The spec only deals with message process logic but it's also very important the way controller.ProcessMsg is called.
Message queueing and retry are important as there is no guarantee as to when a message is delivered.
Examples:
* A proposal message can be delivered after its respective prepare
* A next round message can be delivered before the timer hits timeout
* A late commit message can decide the instance even if it started the next round
* A message can fail to process because it's "too early" or "too late"

Because of the above, there is a need to order and queue messages based on their round and type so to not lose message and make the protocol round change.

## QBFT syncing
QBFT requires an instance to be decided before moving to a new instance which means that 2f+1 nodes need to be on the same instance and round in order to come to a consensus.  
To ensure a node is alway in sync with the current state of its peers QBFT requires syncing in some cases.  
There are 2 sync types: highest decided and highest round change.

[Sync protocol](./p2p/README.md#protocols)

### Highest Decided
Highest decided sync is when a node asks some of its peers to send their highest known decided message (by height number). The node receives multiple decieded messages, validates them and saves the highest one. The next instance it starts will have heights known height + 1.


Triggers:
- Node Start
- f+1 round changes from higher instances

### Highest Round Change
For SSV starting a new QBFT instance is time sensitive, there could be a situation in which a node starts an instance late in which case other nodes could be in later rounds waiting for long timeouts.  
For that reason a way to fetch the latest round change message they've sent on demand is crucial for a timely consensus.

Triggers:
- Timeout on first round without receiving any messages from other peers

## Decided Messages
Decided message is an aggregate commit message with >= quorum threshold unique signatures of the QBFT committee.  
Decided messages are sent to let other peers know that an instance ended even if they didn't receive/ process the QBFT messages.

There could be many other instances sending decided messages and many permutations of the same decided message (different order/ signer ids) which could cloug the P2P network, for that reason <b>decided messages are broadcasted (by default) only by the leader of the decided round</b>.  
This rule is not strictly enforced, other peers could broadcast decided messages as well.

## TODO
- [X] Support 4,7,10,13 committee sizes
- [X] Message encoding and validation spec tests
- [//] proposal/ prepare/ commit spec tests
- [//] round change spec tests
- [ ] Unified test suite, compatible with the formal verification spec
- [ ] Align according to spec and [Roberto's comments](./roberto_comments)
- [ ] Remove round check from upon commit as it can be for any round?
- [ ] RoundChange spec tests
- [ ] Should remove NextProposalData from RoundChangeData? What use it has?
- [ ] Timeout trigger
- [ ] Round-robin leader protocol