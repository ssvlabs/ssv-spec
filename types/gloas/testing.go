package gloas

import (
	_ "embed"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// TestingParentRoot and TestingStateRoot are the fixture block's parent and state roots: non-zero and distinct,
// so a §6 envelope whose parent_beacon_block_root is not taken from the block's parent root gets another root
// (SIP #94 §6).
var (
	TestingParentRoot = phase0.Root{0x70, 0x61, 0x72, 0x65, 0x6e, 0x74}
	TestingStateRoot  = phase0.Root{0x73, 0x74, 0x61, 0x74, 0x65}
)

// TestingBeaconBlock returns a minimal self-build Gloas BeaconBlock for the slot, with the required
// fixed-size body fields populated so it round-trips through SSZ. The bid is for the block's slot and parent,
// as the proposer value check requires (SIP #94 §4), and commits the execution-requests root of an empty
// ExecutionRequests; the §6 blinded envelope derived from the decided block copies that same root straight
// from the bid (SIP #94 §6), so its post-consensus signing root is well-defined. For use in tests.
func TestingBeaconBlock(slot phase0.Slot) *BeaconBlock {
	requestsRoot, err := (&ExecutionRequests{}).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	return &BeaconBlock{
		Slot: slot,
		// == testingutils.TestingValidatorIndex (gloas can't import testingutils): must match the duty's
		// validator so the proposer-index value check (SIP #94 §4) accepts the fixture block.
		ProposerIndex: 1,
		ParentRoot:    TestingParentRoot,
		StateRoot:     TestingStateRoot,
		Body: &BeaconBlockBody{
			ETH1Data:      &phase0.ETH1Data{BlockHash: make([]byte, 32)},
			SyncAggregate: &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &SignedExecutionPayloadBid{Message: &ExecutionPayloadBid{
				ParentBlockRoot:       TestingParentRoot,
				BuilderIndex:          BuilderIndexSelfBuild,
				Slot:                  slot,
				ExecutionRequestsRoot: requestsRoot,
			}},
			ParentExecutionRequests: &ExecutionRequests{},
		},
	}
}

// TestingBlockLists is the number of elements in each of a block's limited lists (see ValidateListLimits), for
// TestingBeaconBlockWithLists.
type TestingBlockLists struct {
	ProposerSlashings, AttesterSlashings, Attestations, Deposits, VoluntaryExits, BLSToExecutionChanges int
	PayloadAttestations                                                                                 int
	WithdrawalRequests, ConsolidationRequests, BuilderDepositRequests, BuilderExitRequests              int
	BlobKZGCommitments                                                                                  int
}

// TestingBeaconBlockWithLists returns TestingBeaconBlock for the slot with the given number of elements in each
// of its limited lists. The elements are empty but complete, so the block encodes. For use in tests.
func TestingBeaconBlockWithLists(slot phase0.Slot, lists TestingBlockLists) *BeaconBlock {
	attestationData := func() *phase0.AttestationData {
		return &phase0.AttestationData{Source: &phase0.Checkpoint{}, Target: &phase0.Checkpoint{}}
	}
	header := func() *phase0.SignedBeaconBlockHeader {
		return &phase0.SignedBeaconBlockHeader{Message: &phase0.BeaconBlockHeader{}}
	}

	block := TestingBeaconBlock(slot)
	body := block.Body
	for i := 0; i < lists.ProposerSlashings; i++ {
		body.ProposerSlashings = append(body.ProposerSlashings, &phase0.ProposerSlashing{SignedHeader1: header(), SignedHeader2: header()})
	}
	for i := 0; i < lists.AttesterSlashings; i++ {
		body.AttesterSlashings = append(body.AttesterSlashings, &eth2gloas.AttesterSlashing{
			Attestation1: &eth2gloas.IndexedAttestation{Data: attestationData()},
			Attestation2: &eth2gloas.IndexedAttestation{Data: attestationData()},
		})
	}
	for i := 0; i < lists.Attestations; i++ {
		body.Attestations = append(body.Attestations, &eth2gloas.Attestation{
			AggregationBits: bitfield.NewBitlist(1),
			Data:            attestationData(),
			CommitteeBits:   bitfield.NewBitvector64(),
		})
	}
	for i := 0; i < lists.Deposits; i++ {
		proof := make([][]byte, 33)
		for j := range proof {
			proof[j] = make([]byte, 32)
		}
		body.Deposits = append(body.Deposits, &phase0.Deposit{
			Proof: proof,
			Data:  &phase0.DepositData{WithdrawalCredentials: make([]byte, 32)},
		})
	}
	for i := 0; i < lists.VoluntaryExits; i++ {
		body.VoluntaryExits = append(body.VoluntaryExits, &phase0.SignedVoluntaryExit{Message: &phase0.VoluntaryExit{}})
	}
	for i := 0; i < lists.BLSToExecutionChanges; i++ {
		body.BLSToExecutionChanges = append(body.BLSToExecutionChanges, &capella.SignedBLSToExecutionChange{Message: &capella.BLSToExecutionChange{}})
	}
	for i := 0; i < lists.PayloadAttestations; i++ {
		body.PayloadAttestations = append(body.PayloadAttestations, &PayloadAttestation{
			AggregationBits: bitfield.NewBitvector512(),
			Data:            &PayloadAttestationData{},
		})
	}

	requests := body.ParentExecutionRequests
	for i := 0; i < lists.WithdrawalRequests; i++ {
		requests.Withdrawals = append(requests.Withdrawals, &electra.WithdrawalRequest{})
	}
	for i := 0; i < lists.ConsolidationRequests; i++ {
		requests.Consolidations = append(requests.Consolidations, &electra.ConsolidationRequest{})
	}
	for i := 0; i < lists.BuilderDepositRequests; i++ {
		requests.BuilderDeposits = append(requests.BuilderDeposits, &BuilderDepositRequest{})
	}
	for i := 0; i < lists.BuilderExitRequests; i++ {
		requests.BuilderExits = append(requests.BuilderExits, &BuilderExitRequest{})
	}

	bid := body.SignedExecutionPayloadBid.Message
	bid.BlobKZGCommitments = make([]deneb.KZGCommitment, lists.BlobKZGCommitments)
	return block
}

// TestingExternalBuilderIndex is a non-self-build BuilderIndex for the §4 external-bid path — any value
// other than BuilderIndexSelfBuild marks the bid as won by an external builder.
const TestingExternalBuilderIndex = BuilderIndex(3)

// TestingBeaconBlockExternalBuild returns the minimal Gloas block of TestingBeaconBlock but with an
// external-builder bid (BuilderIndex != BuilderIndexSelfBuild): the §4 path where produce yields a bare
// block and no envelope, so the decided value's payload_root is zero. For use in tests.
func TestingBeaconBlockExternalBuild(slot phase0.Slot) *BeaconBlock {
	blk := TestingBeaconBlock(slot)
	blk.Body.SignedExecutionPayloadBid.Message.BuilderIndex = TestingExternalBuilderIndex
	return blk
}

// TestingDevnet6SignedBeaconBlockSSZ is a real Gloas SignedBeaconBlock (slot 66) captured from lighthouse v8.2.0
// on glamsterdam-devnet-6; its ParentExecutionRequests carries the EIP-8282 five-list ExecutionRequests (all
// empty). For use in tests.
//
//go:embed testdata/devnet6_gloas_block.ssz
var TestingDevnet6SignedBeaconBlockSSZ []byte
