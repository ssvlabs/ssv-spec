package gloas

import (
	"encoding/json"
	"fmt"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockRootMatchesGoEth2Client pins the core invariant of the dynamic-ssz migration:
// ssv-spec's Gloas SignedBeaconBlock must produce the same hash_tree_root as go-eth2-client's
// spec/gloas type — the progressive-merkleization reference validated against consensus-specs on
// devnet-6. Decoding the same real block into both and comparing roots catches a future dynssz/tag
// drift that the round-trip and wire tests would miss: serialization can byte-match while HTR
// diverges, which is exactly the class of bug (wrong progressive root) this migration fixes.
//
// This exercises the progressive container merkleization (the block body) and, because an empty
// progressive list already hashes differently from an empty bounded list, a progressive-list tag drift on
// any body list too. TestExecutionPayloadBidRootMatchesGoEth2Client below adds a non-empty list.
func TestSignedBeaconBlockRootMatchesGoEth2Client(t *testing.T) {
	var ours SignedBeaconBlock
	require.NoError(t, ours.UnmarshalSSZ(TestingDevnet6SignedBeaconBlockSSZ), "decode into ssv-spec gloas type")
	ourRoot, err := ours.HashTreeRoot()
	require.NoError(t, err)

	var ref eth2gloas.SignedBeaconBlock
	require.NoError(t, ref.UnmarshalSSZ(TestingDevnet6SignedBeaconBlockSSZ), "decode into go-eth2-client spec/gloas type")
	refRoot, err := ref.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, refRoot, ourRoot,
		"ssv-spec Gloas SignedBeaconBlock hash_tree_root must match go-eth2-client spec/gloas (progressive merkleization)")
}

// TestExecutionPayloadBidRootMatchesGoEth2Client guards progressive-*list* merkleization: with a
// non-empty BlobKZGCommitments (a `ssz-type:"progressive-list"`), a progressive vs a binary
// MerkleizeWithMixin produce different roots, so this fails loudly if the tag/backend ever drifts.
// A real block with empty lists (the test above) cannot catch that.
func TestExecutionPayloadBidRootMatchesGoEth2Client(t *testing.T) {
	bid := &ExecutionPayloadBid{
		GasLimit:           30_000_000,
		BuilderIndex:       BuilderIndexSelfBuild,
		Slot:               66,
		Value:              12345,
		ExecutionPayment:   678,
		BlobKZGCommitments: []deneb.KZGCommitment{{0x01}, {0x02}, {0x03}, {0x04}, {0x05}},
	}
	ssz, err := bid.MarshalSSZ()
	require.NoError(t, err)

	var ours ExecutionPayloadBid
	require.NoError(t, ours.UnmarshalSSZ(ssz), "decode into ssv-spec gloas type")
	ourRoot, err := ours.HashTreeRoot()
	require.NoError(t, err)

	var ref eth2gloas.ExecutionPayloadBid
	require.NoError(t, ref.UnmarshalSSZ(ssz), "decode into go-eth2-client spec/gloas type")
	refRoot, err := ref.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, refRoot, ourRoot,
		"non-empty progressive-list root must match go-eth2-client spec/gloas")
}

// sszContainer is what the tests below need of a container, ssv-spec's or go-eth2-client's.
type sszContainer interface {
	MarshalSSZ() ([]byte, error)
	UnmarshalSSZ(data []byte) error
	HashTreeRoot() ([32]byte, error)
}

// TestPopulatedContainersMatchGoEth2Client extends the two tests above to every container this package shares
// with go-eth2-client, each fully populated: every list is non-empty, the five ExecutionRequests lists
// included, and some are long enough to reach past the first subtrees of a progressive list. The devnet block
// leaves its lists empty, which checks each list's own kind but nothing about what goes into it: the element
// containers (the builder deposit and exit requests, the payload attestations) and the merkleization of a
// non-empty list are only exercised here. Each value's encoding must decode into go-eth2-client's type, encode
// back to the same bytes, and hash to the same root.
//
// Both sides are merkleized by the same library (dynamic-ssz), so this pins the containers' shape and tags,
// not the library.
func TestPopulatedContainersMatchGoEth2Client(t *testing.T) {
	block := populatedBeaconBlock()
	body := block.Body
	requests := body.ParentExecutionRequests
	payloadAttestation := body.PayloadAttestations[0]
	preferences := &ProposerPreferences{
		DependentRoot:  filledRoot(0x70),
		ProposalSlot:   66,
		ValidatorIndex: 7,
		FeeRecipient:   filledAddress(0x71),
		TargetGasLimit: 36_000_000,
	}
	auth := &BuilderRequestAuth{Data: filled(0x72, 300), Slot: 66}

	for _, tc := range []struct {
		name string
		ours sszContainer
		ref  sszContainer
	}{
		{"SignedBeaconBlock", &SignedBeaconBlock{Message: block, Signature: filledSignature(0x73)}, &eth2gloas.SignedBeaconBlock{}},
		{"BeaconBlock", block, &eth2gloas.BeaconBlock{}},
		{"BeaconBlockBody", body, &eth2gloas.BeaconBlockBody{}},
		{"SignedExecutionPayloadBid", body.SignedExecutionPayloadBid, &eth2gloas.SignedExecutionPayloadBid{}},
		{"ExecutionPayloadBid", body.SignedExecutionPayloadBid.Message, &eth2gloas.ExecutionPayloadBid{}},
		{"ExecutionRequests", requests, &eth2gloas.ExecutionRequests{}},
		{"BuilderDepositRequest", requests.BuilderDeposits[0], &eth2gloas.BuilderDepositRequest{}},
		{"BuilderExitRequest", requests.BuilderExits[0], &eth2gloas.BuilderExitRequest{}},
		{"PayloadAttestation", payloadAttestation, &eth2gloas.PayloadAttestation{}},
		{"PayloadAttestationData", payloadAttestation.Data, &eth2gloas.PayloadAttestationData{}},
		{
			"PayloadAttestationMessage",
			&PayloadAttestationMessage{ValidatorIndex: 7, Data: payloadAttestation.Data, Signature: filledSignature(0x74)},
			&eth2gloas.PayloadAttestationMessage{},
		},
		{"ProposerPreferences", preferences, &eth2gloas.ProposerPreferences{}},
		{
			"SignedProposerPreferences",
			&SignedProposerPreferences{Message: preferences, Signature: filledSignature(0x75)},
			&eth2gloas.SignedProposerPreferences{},
		},
		{"BuilderRequestAuth", auth, &eth2gloas.BuilderRequestAuth{}},
		{
			"SignedBuilderRequestAuth",
			&SignedBuilderRequestAuth{Message: auth, Signature: filledSignature(0x76)},
			&eth2gloas.SignedBuilderRequestAuth{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := tc.ours.MarshalSSZ()
			require.NoError(t, err)
			root, err := tc.ours.HashTreeRoot()
			require.NoError(t, err)

			require.NoError(t, tc.ref.UnmarshalSSZ(encoded), "decode into go-eth2-client spec/gloas type")
			refEncoded, err := tc.ref.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, encoded, refEncoded, "encoding must match go-eth2-client spec/gloas")
			refRoot, err := tc.ref.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, refRoot, root, "hash_tree_root must match go-eth2-client spec/gloas")
		})
	}
}

// TestBlindedEnvelopeRootMatchesGoEth2ClientFullEnvelope pins the root equivalence that §6 envelope signing
// rests on (SIP #94 §6): the blinded envelope, which carries hash_tree_root(payload) and
// hash_tree_root(execution_requests) in place of the two fields, must hash to the same root as the full
// ExecutionPayloadEnvelope, or the cluster's signature would not verify for the envelope the beacon node
// publishes. The full types are go-eth2-client's, as this package does not define them.
//
// Both roots come from the same library here too. SIP #94 §6 also asks for a fixed expected root produced
// independently; this test is not that.
func TestBlindedEnvelopeRootMatchesGoEth2ClientFullEnvelope(t *testing.T) {
	full := populatedFullEnvelope(t)
	fullRoot, err := full.HashTreeRoot()
	require.NoError(t, err)
	payloadRoot, err := full.Payload.HashTreeRoot()
	require.NoError(t, err)
	requestsRoot, err := full.ExecutionRequests.HashTreeRoot()
	require.NoError(t, err)

	blinded := &BlindedExecutionPayloadEnvelope{
		PayloadRoot:           payloadRoot,
		ExecutionRequestsRoot: requestsRoot,
		BuilderIndex:          BuilderIndex(full.BuilderIndex),
		BeaconBlockRoot:       full.BeaconBlockRoot,
		ParentBeaconBlockRoot: full.ParentBeaconBlockRoot,
	}
	blindedRoot, err := blinded.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, fullRoot, blindedRoot,
		"blinded envelope hash_tree_root must match the full go-eth2-client ExecutionPayloadEnvelope's")
}

// filled returns n bytes counting up from seed, so that populated fields hold distinct data.
func filled(seed byte, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed + byte(i)
	}
	return b
}

func filledRoot(seed byte) (r phase0.Root) {
	copy(r[:], filled(seed, len(r)))
	return r
}

func filledHash(seed byte) (h phase0.Hash32) {
	copy(h[:], filled(seed, len(h)))
	return h
}

func filledSignature(seed byte) (s phase0.BLSSignature) {
	copy(s[:], filled(seed, len(s)))
	return s
}

func filledPubKey(seed byte) (p phase0.BLSPubKey) {
	copy(p[:], filled(seed, len(p)))
	return p
}

func filledAddress(seed byte) (a bellatrix.ExecutionAddress) {
	copy(a[:], filled(seed, len(a)))
	return a
}

// populatedBeaconBlock returns a self-build block with every body field set and every list non-empty.
func populatedBeaconBlock() *BeaconBlock {
	attestationData := func(seed byte) *phase0.AttestationData {
		return &phase0.AttestationData{
			Slot:            65,
			Index:           1,
			BeaconBlockRoot: filledRoot(seed),
			Source:          &phase0.Checkpoint{Epoch: 1, Root: filledRoot(seed + 1)},
			Target:          &phase0.Checkpoint{Epoch: 2, Root: filledRoot(seed + 2)},
		}
	}
	header := func(seed byte) *phase0.SignedBeaconBlockHeader {
		return &phase0.SignedBeaconBlockHeader{
			Message: &phase0.BeaconBlockHeader{
				Slot:          64,
				ProposerIndex: 3,
				ParentRoot:    filledRoot(seed),
				StateRoot:     filledRoot(seed + 1),
				BodyRoot:      filledRoot(seed + 2),
			},
			Signature: filledSignature(seed + 3),
		}
	}
	indexedAttestation := func(seed byte) *eth2gloas.IndexedAttestation {
		return &eth2gloas.IndexedAttestation{
			AttestingIndices: []uint64{1, 2, 3, 4, 5, 6, 7},
			Data:             attestationData(seed),
			Signature:        filledSignature(seed + 3),
		}
	}
	attestation := func(seed byte) *eth2gloas.Attestation {
		aggregationBits := bitfield.NewBitlist(9)
		aggregationBits.SetBitAt(uint64(seed%9), true)
		committeeBits := bitfield.NewBitvector64()
		committeeBits.SetBitAt(uint64(seed%64), true)
		return &eth2gloas.Attestation{
			AggregationBits: aggregationBits,
			Data:            attestationData(seed),
			Signature:       filledSignature(seed + 3),
			CommitteeBits:   committeeBits,
		}
	}
	deposit := func(seed byte) *phase0.Deposit {
		proof := make([][]byte, 33)
		for i := range proof {
			proof[i] = filled(seed+byte(i), 32)
		}
		return &phase0.Deposit{
			Proof: proof,
			Data: &phase0.DepositData{
				PublicKey:             filledPubKey(seed),
				WithdrawalCredentials: filled(seed+1, 32),
				Amount:                32_000_000_000,
				Signature:             filledSignature(seed + 2),
			},
		}
	}
	payloadAttestation := func(seed byte, blobDataAvailable bool) *PayloadAttestation {
		aggregationBits := bitfield.NewBitvector512()
		aggregationBits.SetBitAt(uint64(seed), true)
		aggregationBits.SetBitAt(511, true)
		return &PayloadAttestation{
			AggregationBits: aggregationBits,
			Data: &PayloadAttestationData{
				BeaconBlockRoot:   filledRoot(seed),
				Slot:              65,
				PayloadPresent:    true,
				BlobDataAvailable: blobDataAvailable,
			},
			Signature: filledSignature(seed + 1),
		}
	}
	syncCommitteeBits := bitfield.NewBitvector512()
	syncCommitteeBits.SetBitAt(5, true)
	syncCommitteeBits.SetBitAt(500, true)

	return &BeaconBlock{
		Slot:          66,
		ProposerIndex: 7,
		ParentRoot:    filledRoot(0x01),
		StateRoot:     filledRoot(0x02),
		Body: &BeaconBlockBody{
			RANDAOReveal: filledSignature(0x03),
			ETH1Data:     &phase0.ETH1Data{DepositRoot: filledRoot(0x04), DepositCount: 5, BlockHash: filled(0x05, 32)},
			Graffiti:     filledRoot(0x06),
			ProposerSlashings: []*phase0.ProposerSlashing{
				{SignedHeader1: header(0x10), SignedHeader2: header(0x14)},
				{SignedHeader1: header(0x18), SignedHeader2: header(0x1c)},
			},
			AttesterSlashings: []*eth2gloas.AttesterSlashing{
				{Attestation1: indexedAttestation(0x20), Attestation2: indexedAttestation(0x24)},
			},
			Attestations: []*eth2gloas.Attestation{attestation(0x28), attestation(0x2c), attestation(0x30)},
			Deposits:     []*phase0.Deposit{deposit(0x34), deposit(0x38)},
			VoluntaryExits: []*phase0.SignedVoluntaryExit{
				{Message: &phase0.VoluntaryExit{Epoch: 3, ValidatorIndex: 4}, Signature: filledSignature(0x3c)},
				{Message: &phase0.VoluntaryExit{Epoch: 3, ValidatorIndex: 5}, Signature: filledSignature(0x3d)},
			},
			SyncAggregate: &altair.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: filledSignature(0x3e),
			},
			BLSToExecutionChanges: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     6,
						FromBLSPubkey:      filledPubKey(0x40),
						ToExecutionAddress: filledAddress(0x41),
					},
					Signature: filledSignature(0x42),
				},
			},
			SignedExecutionPayloadBid: &SignedExecutionPayloadBid{
				Message: &ExecutionPayloadBid{
					ParentBlockHash:       filledHash(0x50),
					ParentBlockRoot:       filledRoot(0x51),
					BlockHash:             filledHash(0x52),
					PrevRandao:            filledRoot(0x53),
					FeeRecipient:          filledAddress(0x54),
					GasLimit:              36_000_000,
					BuilderIndex:          BuilderIndexSelfBuild,
					Slot:                  66,
					Value:                 12345,
					ExecutionPayment:      678,
					BlobKZGCommitments:    []deneb.KZGCommitment{{0x01}, {0x02}, {0x03}, {0x04}, {0x05}, {0x06}},
					ExecutionRequestsRoot: filledRoot(0x55),
				},
				Signature: filledSignature(0x56),
			},
			PayloadAttestations:     []*PayloadAttestation{payloadAttestation(0x60, true), payloadAttestation(0x62, false)},
			ParentExecutionRequests: populatedExecutionRequests(),
		},
	}
}

// populatedExecutionRequests returns Gloas execution requests with all five lists non-empty; the deposits are
// enough to spill into a progressive list's third subtree.
func populatedExecutionRequests() *ExecutionRequests {
	requests := &ExecutionRequests{
		Withdrawals: []*electra.WithdrawalRequest{
			{SourceAddress: filledAddress(0x90), ValidatorPubkey: filledPubKey(0x91), Amount: 1},
			{SourceAddress: filledAddress(0x92), ValidatorPubkey: filledPubKey(0x93), Amount: 2},
		},
		Consolidations: []*electra.ConsolidationRequest{
			{SourceAddress: filledAddress(0x94), SourcePubkey: filledPubKey(0x95), TargetPubkey: filledPubKey(0x96)},
		},
		BuilderDeposits: []*BuilderDepositRequest{
			{Pubkey: filledPubKey(0x97), WithdrawalCredentials: filledRoot(0x98), Amount: 3, Signature: filledSignature(0x99)},
			{Pubkey: filledPubKey(0x9a), WithdrawalCredentials: filledRoot(0x9b), Amount: 4, Signature: filledSignature(0x9c)},
		},
		BuilderExits: []*BuilderExitRequest{
			{SourceAddress: filledAddress(0x9d), Pubkey: filledPubKey(0x9e)},
		},
	}
	for i := byte(0); i < 6; i++ {
		requests.Deposits = append(requests.Deposits, &electra.DepositRequest{
			Pubkey:                filledPubKey(0x80 + i),
			WithdrawalCredentials: filled(0x81+i, 32),
			Amount:                32_000_000_000,
			Signature:             filledSignature(0x82 + i),
			Index:                 uint64(i),
		})
	}
	return requests
}

// populatedFullEnvelope returns go-eth2-client's full ExecutionPayloadEnvelope with a populated payload and the
// populated requests above. The payload is built through its JSON form: its base fee is a type from a module
// this one does not import directly.
func populatedFullEnvelope(t *testing.T) *eth2gloas.ExecutionPayloadEnvelope {
	hex := func(b []byte) string { return fmt.Sprintf("%#x", b) }
	payloadJSON, err := json.Marshal(map[string]any{
		"parent_hash":      hex(filled(0xa0, 32)),
		"fee_recipient":    hex(filled(0xa1, 20)),
		"state_root":       hex(filled(0xa2, 32)),
		"receipts_root":    hex(filled(0xa3, 32)),
		"logs_bloom":       hex(filled(0xa4, 256)),
		"prev_randao":      hex(filled(0xa5, 32)),
		"block_number":     "100",
		"gas_limit":        "36000000",
		"gas_used":         "21000",
		"timestamp":        "1700000000",
		"extra_data":       hex(filled(0xa6, 8)),
		"base_fee_per_gas": "7",
		"block_hash":       hex(filled(0xa7, 32)),
		"transactions":     []string{hex(filled(0xa8, 40)), hex(filled(0xa9, 3)), hex(filled(0xaa, 120))},
		"withdrawals": []map[string]string{
			{"index": "1", "validator_index": "2", "address": hex(filled(0xab, 20)), "amount": "3"},
			{"index": "4", "validator_index": "5", "address": hex(filled(0xac, 20)), "amount": "6"},
		},
		"blob_gas_used":     "131072",
		"excess_blob_gas":   "262144",
		"block_access_list": hex(filled(0xad, 50)),
		"slot_number":       "66",
	})
	require.NoError(t, err)
	payload := &eth2gloas.ExecutionPayload{}
	require.NoError(t, json.Unmarshal(payloadJSON, payload))

	// The requests are this package's populated ones, moved over through their encoding.
	encodedRequests, err := populatedExecutionRequests().MarshalSSZ()
	require.NoError(t, err)
	requests := &eth2gloas.ExecutionRequests{}
	require.NoError(t, requests.UnmarshalSSZ(encodedRequests))

	return &eth2gloas.ExecutionPayloadEnvelope{
		Payload:               payload,
		ExecutionRequests:     requests,
		BuilderIndex:          eth2gloas.BuilderIndex(BuilderIndexSelfBuild),
		BeaconBlockRoot:       filledRoot(0xb0),
		ParentBeaconBlockRoot: filledRoot(0xb1),
	}
}
