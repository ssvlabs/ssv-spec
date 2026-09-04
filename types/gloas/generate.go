package gloas

// ssv-spec's canonical Gloas (ePBS) containers. The spec defines the types every client (Anchor, the
// SSV node) must match; go-eth2-client's fork now carries the same containers with identical tags and
// the node aliases those, but ssv-spec keeps its own so the spec does not depend on a client library —
// their equality is pinned by beacon_block_root_test.go. Aliasing the fork here is a possible follow-up.
//
// The types use progressive containers/lists (EIP-7688/7916) that fastssz cannot merkleize, so
// pk910/dynamic-ssz generates the encoders. Regenerate with `go generate`. Hand-written type files are
// snake_case.go; the generated encoders are flatcase *_ssz.go — do not edit those. dynssz-gen overwrites
// them in place (no `rm` first): it type-checks the package to load it, and the hand-written Encode/Decode
// wrappers reference the generated MarshalSSZ, so deleting the files first would break generation — a stale
// encoder for a removed type fails the build rather than slipping through.
//go:generate go tool dynssz-gen -config generate.yaml
