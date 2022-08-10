package testingutils

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/bloxapp/ssv-spec/types"
)

var TestingAttestationData = &phase0.AttestationData{
	Slot:            12,
	Index:           3,
	BeaconBlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
}
var TestingAttestationRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f") //[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}

var TestingBeaconBlock = &altair.BeaconBlock{
	Slot:          12,
	ProposerIndex: 10,
	ParentRoot:    phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	StateRoot:     phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Body: &altair.BeaconBlockBody{
		RANDAOReveal: phase0.BLSSignature{},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			DepositCount: 100,
			BlockHash:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Graffiti:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		ProposerSlashings: []*phase0.ProposerSlashing{},
		AttesterSlashings: []*phase0.AttesterSlashing{},
		Attestations: []*phase0.Attestation{
			{
				AggregationBits: bitfield.NewBitlist(122),
				Data:            TestingAttestationData,
				Signature:       phase0.BLSSignature{},
			},
		},
		Deposits:       []*phase0.Deposit{},
		VoluntaryExits: []*phase0.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: phase0.BLSSignature{},
		},
	},
}
var TestingBeaconBlockRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
var TestingRandaoRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")

var TestingAggregateAndProof = &phase0.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  phase0.BLSSignature{},
	Aggregate: &phase0.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       phase0.BLSSignature{},
		Data:            TestingAttestationData,
	},
}
var TestingSignedAggregateAndProofRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
var TestingSelectionProofRoot, _ = hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")

const (
	TestingDutySlot       = 12
	TestingValidatorIndex = 1
)

var TestingSyncCommitteeBlockRoot = phase0.Root{}

var TestingContributionProofRoots = func() []phase0.Root {
	byts1, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
	byts2, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb35653")
	byts3, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565d")

	ret := make([]phase0.Root, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := phase0.Root{}
		copy(b[:], byts)
		ret = append(ret, b)
	}
	return ret
}()
var TestingContributionProofsSigned = func() []phase0.BLSSignature {
	// signed with 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
	byts1, _ := hex.DecodeString("8cef237a0e3a1bba095e9534df220e0dccd3de740d50510c9b3cedb6a0c4ca5bd23b5ae672698260333fc47e532741c303efb98f0636a3515d615535c0e072ed470514c1fafda9335bc4919127a1fd3107b9990e3b857075e1f63a27bfd6b216")
	byts2, _ := hex.DecodeString("978f80611fb452449d413902487eec69a531bdab16ab51433582fdf9bb900d7b63de10fb048204591c06322ba3fa1cff0c83077ee0c17416dc718f6cca82c5c94115679646f5fa08410d794cd6974d562ffe522d70eb89f340064fe99bf7471d")
	byts3, _ := hex.DecodeString("addc5f3b2534b8c542b500c8a9e50fef49f18fbaafaf995a3ebe2a3d50836703f2d92cb5885d7b6cafdc6c09d889e6bc10d3470c51deb56d26996ea63f57a909a7a9bec0140ed79e683695b111ca36cbe4776d11ad50ad3d81e96290c064608c")

	ret := make([]phase0.BLSSignature, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := phase0.BLSSignature{}
		copy(b[:], byts)
		ret = append(ret, b)
	}
	return ret
}()
var TestingContributionRoots = func() [][]byte {
	byts1, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
	byts2, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565a")
	byts3, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb35656")
	return [][]byte{
		byts1, byts2, byts3,
	}
}()
var TestingSyncCommitteeContributions = []*altair.SyncCommitteeContribution{
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 0,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 1,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 2,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
}

var TestingAttesterDuty = &types.Duty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingProposerDuty = &types.Duty{
	Type:                    types.BNRoleProposer,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingAggregatorDuty = &types.Duty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingSyncCommitteeDuty = &types.Duty{
	Type:                    types.BNRoleSyncCommittee,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingSyncCommitteeContributionDuty = &types.Duty{
	Type:                    types.BNRoleSyncCommitteeContribution,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingUnknownDutyType = &types.Duty{
	Type:                    100,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingWrongDutyPK = &types.Duty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingWrongValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

func blsSigFromHex(str string) phase0.BLSSignature {
	byts, _ := hex.DecodeString(str)
	ret := phase0.BLSSignature{}
	copy(ret[:], byts)
	return ret
}

type testingBeaconNode struct {
}

func NewTestingBeaconNode() *testingBeaconNode {
	return &testingBeaconNode{}
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *testingBeaconNode) GetBeaconNetwork() types.BeaconNetwork {
	return types.NowTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *testingBeaconNode) GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	return TestingAttestationData, nil
}

// SubmitAttestation submit the attestation to the node
func (bn *testingBeaconNode) SubmitAttestation(attestation *phase0.Attestation) error {
	return nil
}

// GetBeaconBlock returns beacon block by the given slot and committee index
func (bn *testingBeaconNode) GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (*altair.BeaconBlock, error) {
	return TestingBeaconBlock, nil
}

// SubmitBeaconBlock submit the block to the node
func (bn *testingBeaconNode) SubmitBeaconBlock(block *altair.SignedBeaconBlock) error {
	return nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
func (bn *testingBeaconNode) SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, slotSig []byte) (*phase0.AggregateAndProof, error) {
	return TestingAggregateAndProof, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *testingBeaconNode) SubmitSignedAggregateSelectionProof(msg *phase0.SignedAggregateAndProof) error {
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *testingBeaconNode) GetSyncMessageBlockRoot() (phase0.Root, error) {
	return TestingSyncCommitteeBlockRoot, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *testingBeaconNode) SubmitSyncMessage(msg *altair.SyncCommitteeMessage) error {
	return nil
}

// GetSyncSubcommitteeIndex returns sync committee indexes for aggregator
func (bn *testingBeaconNode) GetSyncSubcommitteeIndex(slot phase0.Slot, pubKey phase0.BLSPubKey) ([]uint64, error) {
	// each subcommittee index correlates to TestingContributionProofRoots by index
	return []uint64{0, 1, 2}, nil
}

// IsSyncCommitteeAggregator returns tru if aggregator
func (bn *testingBeaconNode) IsSyncCommitteeAggregator(proof []byte) (bool, error) {
	return true, nil
}

// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
func (bn *testingBeaconNode) SyncCommitteeSubnetID(subCommitteeID uint64) (uint64, error) {
	// each subcommittee index correlates to TestingContributionProofRoots by index
	return subCommitteeID, nil
}

// GetSyncCommitteeContribution returns
func (bn *testingBeaconNode) GetSyncCommitteeContribution(slot phase0.Slot, subnetID uint64, pubKey phase0.BLSPubKey) (*altair.SyncCommitteeContribution, error) {
	return TestingSyncCommitteeContributions[subnetID], nil
}

// SubmitSignedContributionAndProof broadcasts to the network
func (bn *testingBeaconNode) SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error {
	return nil
}
