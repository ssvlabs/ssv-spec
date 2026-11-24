// Package gloas holds the Gloas (ePBS / EIP-7732) beacon wire types the SSV spec references, per
// SIP #94 — the §4 block family (bid-only BeaconBlock with the EIP-8282 five-list
// ExecutionRequests), the §3 payload-attestation containers, the §5 proposer preferences, and the
// §6 blinded execution-payload envelope. They are defined against the consensus-specs snapshot
// pinned by SIP #94. go-eth2-client ships most of them too (v0.30.0+), with identical tags;
// generate.go says why the spec keeps its own and how their equality is pinned. DataVersionGloas
// is go-eth2-client's (version.go).
//
// The full (unblinded) ExecutionPayloadEnvelope and ExecutionPayload are deliberately not defined
// here: the spec signs the §6 envelope over the blinded form only (the signing root is identical),
// so the full-payload types remain node-side.
package gloas
