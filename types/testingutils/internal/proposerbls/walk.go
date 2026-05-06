// Package proposerbls walks the BLS-typed fields of Capella and Deneb proposer
// block bodies. The regression test in types/testingutils and the regen script
// in types/testingutils/scripts both consume this so the field inventory stays
// in one place.
package proposerbls

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// Visitor receives every BLS-typed field in document order. Each callback is
// invoked with a JSON-path label (for diagnostics) and the raw bytes (48 for
// Pubkey, 96 for Signature).
type Visitor struct {
	Pubkey    func(label string, b []byte)
	Signature func(label string, b []byte)
}

func WalkCapella(body *capella.BeaconBlockBody, v Visitor) {
	v.Signature("randao_reveal", body.RANDAOReveal[:])
	for i, ps := range body.ProposerSlashings {
		v.Signature(fmt.Sprintf("proposer_slashings[%d].signed_header_1.signature", i), ps.SignedHeader1.Signature[:])
		v.Signature(fmt.Sprintf("proposer_slashings[%d].signed_header_2.signature", i), ps.SignedHeader2.Signature[:])
	}
	for i, as := range body.AttesterSlashings {
		v.Signature(fmt.Sprintf("attester_slashings[%d].attestation_1.signature", i), as.Attestation1.Signature[:])
		v.Signature(fmt.Sprintf("attester_slashings[%d].attestation_2.signature", i), as.Attestation2.Signature[:])
	}
	for i, a := range body.Attestations {
		v.Signature(fmt.Sprintf("attestations[%d].signature", i), a.Signature[:])
	}
	for i, d := range body.Deposits {
		v.Pubkey(fmt.Sprintf("deposits[%d].data.pubkey", i), d.Data.PublicKey[:])
		v.Signature(fmt.Sprintf("deposits[%d].data.signature", i), d.Data.Signature[:])
	}
	for i, ve := range body.VoluntaryExits {
		v.Signature(fmt.Sprintf("voluntary_exits[%d].signature", i), ve.Signature[:])
	}
	v.Signature("sync_aggregate.sync_committee_signature", body.SyncAggregate.SyncCommitteeSignature[:])
	for i, btc := range body.BLSToExecutionChanges {
		v.Pubkey(fmt.Sprintf("bls_to_execution_changes[%d].message.from_bls_pubkey", i), btc.Message.FromBLSPubkey[:])
		v.Signature(fmt.Sprintf("bls_to_execution_changes[%d].signature", i), btc.Signature[:])
	}
}

// WalkDeneb mirrors WalkCapella; the BLS-typed field set is unchanged from
// Capella to Deneb. Kept as a separate function because the body types are
// distinct go-eth2-client structs with no shared interface.
func WalkDeneb(body *deneb.BeaconBlockBody, v Visitor) {
	v.Signature("randao_reveal", body.RANDAOReveal[:])
	for i, ps := range body.ProposerSlashings {
		v.Signature(fmt.Sprintf("proposer_slashings[%d].signed_header_1.signature", i), ps.SignedHeader1.Signature[:])
		v.Signature(fmt.Sprintf("proposer_slashings[%d].signed_header_2.signature", i), ps.SignedHeader2.Signature[:])
	}
	for i, as := range body.AttesterSlashings {
		v.Signature(fmt.Sprintf("attester_slashings[%d].attestation_1.signature", i), as.Attestation1.Signature[:])
		v.Signature(fmt.Sprintf("attester_slashings[%d].attestation_2.signature", i), as.Attestation2.Signature[:])
	}
	for i, a := range body.Attestations {
		v.Signature(fmt.Sprintf("attestations[%d].signature", i), a.Signature[:])
	}
	for i, d := range body.Deposits {
		v.Pubkey(fmt.Sprintf("deposits[%d].data.pubkey", i), d.Data.PublicKey[:])
		v.Signature(fmt.Sprintf("deposits[%d].data.signature", i), d.Data.Signature[:])
	}
	for i, ve := range body.VoluntaryExits {
		v.Signature(fmt.Sprintf("voluntary_exits[%d].signature", i), ve.Signature[:])
	}
	v.Signature("sync_aggregate.sync_committee_signature", body.SyncAggregate.SyncCommitteeSignature[:])
	for i, btc := range body.BLSToExecutionChanges {
		v.Pubkey(fmt.Sprintf("bls_to_execution_changes[%d].message.from_bls_pubkey", i), btc.Message.FromBLSPubkey[:])
		v.Signature(fmt.Sprintf("bls_to_execution_changes[%d].signature", i), btc.Signature[:])
	}
}
