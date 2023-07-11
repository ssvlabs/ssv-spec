
# QBFT

## Introduction
This is a spec implementation for the QBFT protocol, following [formal verification spec](https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification) / [github repo](https://github.com/ConsenSys/qbft-formal-spec-and-verification).

To understand the implementation and its comparison to the H. Moniz [IBFT paper](https://arxiv.org/pdf/2002.03613.pdf), please follow this [documentation](docs/README.md).

## Important note on message processing
The spec only deals with message process logic but it's also very important the way controller.ProcessMsg is called.
Message queueing and retry are important as there is no guarantee as to when a message is delivered.
Examples:
* A proposal message can be delivered after its respective prepare
* A next round message can be delivered before the timer hits timeout
* A late commit message can decide the instance even if it started the next round
* A message can fail to process because it's "too early" or "too late"

Because of the above, there is a need to order and queue messages based on their round and type so to not lose message and make the protocol round change.


## TODO
- [X] Support 4,7,10,13 committee sizes
- [X] Message encoding and validation spec tests
- [X] proposal/ prepare/ commit spec tests
- [X] round change spec tests
- [//] Unified test suite, compatible with the formal verification spec
- [//] Align according to spec and [Roberto's comments](./roberto_comments)
- [ ] Remove round check from upon commit as it can be for any round?
- [ ] Use data hashes instead of full data in msgs to save space in justifications


## Additional documents

- to understand the class structure of the qbft package, please check the [class diagram](docs/CLASSDIAGRAM.md).
- a little informal [text](docs/QBFT_LITTLE_EXPLANATION.md) on how qbft works.
- a [sequence diagram](docs/SEQUENCE_DIAGRAM.md) representing how the controller and an instance interact.
- a [Gantt chart](docs/GANTT_CHART.md) of a successful normal case scenario of QBFT.

