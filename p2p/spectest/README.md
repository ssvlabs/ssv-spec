# P2P Spectests

This package contains spec tests for the P2P validation surface.

## Layout

- `all_tests.go`: registry of all P2P spec tests.
- `run_test.go`: Go test entrypoint.
- `tests/`: concrete spec test types and grouped case files.
- `testdoc/`: centralized test type and case documentation strings.
- `generate/`: JSON generator for test vectors and future state-comparison outputs.

## Current scope

The current suite covers pubsub message validation in:

- malformed payload handling
- signed SSV message type filtering
- consensus message validation
- pre-consensus partial-signature validation
- post-consensus partial-signature validation

## Running

```bash
go test ./p2p/...
```

## Generating

```bash
go generate ./p2p/spectest/generate
```
