
# QBFT

## Introduction
This is a spec implementation for the QBFT protocol, following [formal verification spec](https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification).

## TODO
- [X] Support 4,7,10,13 committee sizes
- [ ] Unified test suite, compatible with the formal verification spec
- [ ] Align according to spec and [Roberto's comments](./roberto_comments)
- [ ] Remove round check from upon commit as it can be for any round?
- [ ] RoundChange spec tests