package testingutils

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec/types"
)

func TestCapellaDenebProposerFixturesUseValidBLSBytes(t *testing.T) {
	types.InitBLS()

	assertCapellaBeaconBlockBLSFields(t, "capella", TestingBeaconBlockCapella.Body)
	assertDenebBeaconBlockBLSFields(t, "deneb", TestingBlockContentsDeneb.Block.Body)
}

func assertCapellaBeaconBlockBLSFields(t *testing.T, name string, body *capella.BeaconBlockBody) {
	t.Helper()

	assertBeaconBlockBLSFields(
		t,
		name,
		body.RANDAOReveal,
		body.ProposerSlashings,
		body.AttesterSlashings,
		body.Attestations,
		body.Deposits,
		body.VoluntaryExits,
		body.SyncAggregate,
		body.BLSToExecutionChanges,
	)
}

func assertDenebBeaconBlockBLSFields(t *testing.T, name string, body *deneb.BeaconBlockBody) {
	t.Helper()

	assertBeaconBlockBLSFields(
		t,
		name,
		body.RANDAOReveal,
		body.ProposerSlashings,
		body.AttesterSlashings,
		body.Attestations,
		body.Deposits,
		body.VoluntaryExits,
		body.SyncAggregate,
		body.BLSToExecutionChanges,
	)
}

func assertBeaconBlockBLSFields(
	t *testing.T,
	name string,
	randaoReveal phase0.BLSSignature,
	proposerSlashings []*phase0.ProposerSlashing,
	attesterSlashings []*phase0.AttesterSlashing,
	attestations []*phase0.Attestation,
	deposits []*phase0.Deposit,
	voluntaryExits []*phase0.SignedVoluntaryExit,
	syncAggregate *altair.SyncAggregate,
	blsToExecutionChanges []*capella.SignedBLSToExecutionChange,
) {
	t.Helper()

	assertValidBLSSignature(t, name+".randao_reveal", randaoReveal)

	for _, slashing := range proposerSlashings {
		assertValidBLSSignature(t, name+".proposer_slashings.signed_header_1.signature", slashing.SignedHeader1.Signature)
		assertValidBLSSignature(t, name+".proposer_slashings.signed_header_2.signature", slashing.SignedHeader2.Signature)
	}
	for _, slashing := range attesterSlashings {
		assertValidBLSSignature(t, name+".attester_slashings.attestation_1.signature", slashing.Attestation1.Signature)
		assertValidBLSSignature(t, name+".attester_slashings.attestation_2.signature", slashing.Attestation2.Signature)
	}
	for _, attestation := range attestations {
		assertValidBLSSignature(t, name+".attestations.signature", attestation.Signature)
	}
	for _, deposit := range deposits {
		assertValidBLSPubKey(t, name+".deposits.data.pubkey", deposit.Data.PublicKey)
		assertValidBLSSignature(t, name+".deposits.data.signature", deposit.Data.Signature)
	}
	for _, exit := range voluntaryExits {
		assertValidBLSSignature(t, name+".voluntary_exits.signature", exit.Signature)
	}
	assertValidBLSSignature(t, name+".sync_aggregate.sync_committee_signature", syncAggregate.SyncCommitteeSignature)
	for _, change := range blsToExecutionChanges {
		assertValidBLSPubKey(t, name+".bls_to_execution_changes.message.from_bls_pubkey", change.Message.FromBLSPubkey)
		assertValidBLSSignature(t, name+".bls_to_execution_changes.signature", change.Signature)
	}
}

func assertValidBLSPubKey(t *testing.T, name string, pubkey phase0.BLSPubKey) {
	t.Helper()

	var pk bls.PublicKey
	if err := pk.Deserialize(pubkey[:]); err != nil {
		t.Fatalf("%s is not a valid compressed BLS pubkey: %v", name, err)
	}
}

func assertValidBLSSignature(t *testing.T, name string, signature phase0.BLSSignature) {
	t.Helper()

	var sig bls.Sign
	if err := sig.Deserialize(signature[:]); err != nil {
		t.Fatalf("%s is not a valid compressed BLS signature: %v", name, err)
	}
}
