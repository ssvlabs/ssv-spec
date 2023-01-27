
# Alea-BFT

## Introduction
This is a spec implementation for the Alea-BFT protocol, following the [paper](https://arxiv.org/abs/2202.02071) specification.

## Important note on message processing
The spec only deals with message process logic but it's also very important the way controller.ProcessMsg is called.
Message queueing and retry are important as there is no guarantee as to when a message is delivered.

## TODO
- [X] Implement new Alea-BFT messages
- [X] Adjust controller.go and instance.go
- [X] Implement the Alea-BFT protocol abstracting ABA and VCBC primitives
- [X] Implement ABA
- [X] Implement VCBC
- [X] Message tests
- [X] Happy Flow tests
- [ ] Proposal tests
- [ ] VCBC tests
- [ ] ABA tests
- [ ] FillGap and Filler tests
- [ ] Instance tests
- [ ] Controller tests