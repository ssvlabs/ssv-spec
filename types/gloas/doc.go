// Package gloas holds the Gloas (ePBS / EIP-7732) beacon wire types the SSV spec references, per
// SIP #94 — the §4 block family (bid-only BeaconBlock with the EIP-8282 five-list
// ExecutionRequests), the §3 payload-attestation containers, the §5 proposer preferences, and the
// §6 blinded execution-payload envelope. go-eth2-client does not ship Gloas yet, so these are
// vendored here against the consensus-specs snapshot pinned by SIP #94; once upstream ships
// (go-eth2-client PRs #269/#280), this package shrinks to whatever upstream still lacks and
// DataVersionGloas is reconciled with the upstream enum.
//
// The full (unblinded) ExecutionPayloadEnvelope and ExecutionPayload are deliberately not vendored:
// the spec models the §6 duty over the blinded form only (the signing root is identical), so the
// full-payload types remain node-side.
package gloas
