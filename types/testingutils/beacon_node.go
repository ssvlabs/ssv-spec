package testingutils

import (
	"encoding/hex"
	altair "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/prysmaticlabs/go-bitfield"
)

var TestingAttestationData = &spec.AttestationData{
	Slot:            12,
	Index:           3,
	BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Source: &spec.Checkpoint{
		Epoch: 0,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
	Target: &spec.Checkpoint{
		Epoch: 1,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
}

var TestingBeaconBlock = &bellatrix.BeaconBlock{
	Slot:          12,
	ProposerIndex: 10,
	ParentRoot:    spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	StateRoot:     spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Body: &bellatrix.BeaconBlockBody{
		RANDAOReveal: spec.BLSSignature{},
		ETH1Data: &spec.ETH1Data{
			DepositRoot:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			DepositCount: 100,
			BlockHash:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Graffiti:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		ProposerSlashings: []*spec.ProposerSlashing{},
		AttesterSlashings: []*spec.AttesterSlashing{},
		Attestations: []*spec.Attestation{
			{
				AggregationBits: bitfield.NewBitlist(122),
				Data:            TestingAttestationData,
				Signature:       spec.BLSSignature{},
			},
		},
		Deposits:       []*spec.Deposit{},
		VoluntaryExits: []*spec.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: spec.BLSSignature{},
		},
		ExecutionPayload: &bellatrix.ExecutionPayload{
			ParentHash:    spec.Hash32{},
			FeeRecipient:  bellatrix.ExecutionAddress{},
			StateRoot:     spec.Hash32{},
			ReceiptsRoot:  spec.Hash32{},
			LogsBloom:     [256]byte{},
			PrevRandao:    [32]byte{},
			BlockNumber:   100,
			GasLimit:      1000000,
			GasUsed:       800000,
			Timestamp:     123456789,
			BaseFeePerGas: [32]byte{},
			BlockHash:     spec.Hash32{},
			Transactions:  []bellatrix.Transaction{},
		},
	},
}

var TestingAggregateAndProof = &spec.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  spec.BLSSignature{},
	Aggregate: &spec.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       spec.BLSSignature{},
		Data:            TestingAttestationData,
	},
}

const (
	TestingDutySlot       = 12
	TestingDutyEpoch      = 0
	TestingValidatorIndex = 1
)

var TestingSyncCommitteeBlockRoot = spec.Root{}

var TestingContributionProofRoots = func() []spec.Root {
	byts1, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565f")
	byts2, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb35653")
	byts3, _ := hex.DecodeString("81451c58b079c5af84ebe4b92900d3e9c5a346678cb6dc3c4b7eea2c9cb3565d")

	ret := make([]spec.Root, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := spec.Root{}
		copy(b[:], byts)
		ret = append(ret, b)
	}
	return ret
}()
var TestingContributionProofsSigned = func() []spec.BLSSignature {
	// signed with 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
	byts1, _ := hex.DecodeString("b18833bb7549ec33e8ac414ba002fd45bb094ca300bd24596f04a434a89beea462401da7c6b92fb3991bd17163eb603604a40e8dd6781266c990023446776ff42a9313df26a0a34184a590e57fa4003d610c2fa214db4e7dec468592010298bc")
	byts2, _ := hex.DecodeString("9094342c95146554df849dc20f7425fca692dacee7cb45258ddd264a8e5929861469fda3d1567b9521cba83188ffd61a0dbe6d7180c7a96f5810d18db305e9143772b766d368aa96d3751f98d0ce2db9f9e6f26325702088d87f0de500c67c68")
	byts3, _ := hex.DecodeString("a7f88ce43eff3aa8cdd2e3957c5bead4e21353fbecac6079a5398d03019bc45ff7c951785172deee70e9bc5abbc8ca6a0f0441e9d4cc9da74c31121357f7d7c7de9533f6f457da493e3314e22d554ab76613e469b050e246aff539a33807197c")

	ret := make([]spec.BLSSignature, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := spec.BLSSignature{}
		copy(b[:], byts)
		ret = append(ret, b)
	}
	return ret
}()

var TestingSyncCommitteeContributions = []*altair.SyncCommitteeContribution{
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 0,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         spec.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 1,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         spec.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 2,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         spec.BLSSignature{},
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

func blsSigFromHex(str string) spec.BLSSignature {
	byts, _ := hex.DecodeString(str)
	ret := spec.BLSSignature{}
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
func (bn *testingBeaconNode) GetAttestationData(slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error) {
	return TestingAttestationData, nil
}

// SubmitAttestation submit the attestation to the node
func (bn *testingBeaconNode) SubmitAttestation(attestation *spec.Attestation) error {
	return nil
}

// GetBeaconBlock returns beacon block by the given slot and committee index
func (bn *testingBeaconNode) GetBeaconBlock(slot spec.Slot, committeeIndex spec.CommitteeIndex, graffiti, randao []byte) (*bellatrix.BeaconBlock, error) {
	return TestingBeaconBlock, nil
}

// SubmitBeaconBlock submit the block to the node
func (bn *testingBeaconNode) SubmitBeaconBlock(block *bellatrix.SignedBeaconBlock) error {
	return nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
func (bn *testingBeaconNode) SubmitAggregateSelectionProof(slot spec.Slot, committeeIndex spec.CommitteeIndex, slotSig []byte) (*spec.AggregateAndProof, error) {
	return TestingAggregateAndProof, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *testingBeaconNode) SubmitSignedAggregateSelectionProof(msg *spec.SignedAggregateAndProof) error {
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *testingBeaconNode) GetSyncMessageBlockRoot() (spec.Root, error) {
	return TestingSyncCommitteeBlockRoot, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *testingBeaconNode) SubmitSyncMessage(msg *altair.SyncCommitteeMessage) error {
	return nil
}

// GetSyncSubcommitteeIndex returns sync committee indexes for aggregator
func (bn *testingBeaconNode) GetSyncSubcommitteeIndex(slot spec.Slot, pubKey spec.BLSPubKey) ([]uint64, error) {
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
func (bn *testingBeaconNode) GetSyncCommitteeContribution(slot spec.Slot, subnetID uint64, pubKey spec.BLSPubKey) (*altair.SyncCommitteeContribution, error) {
	return TestingSyncCommitteeContributions[subnetID], nil
}

// SubmitSignedContributionAndProof broadcasts to the network
func (bn *testingBeaconNode) SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error {
	return nil
}

func (bn *testingBeaconNode) DomainData(epoch spec.Epoch, domain spec.DomainType) (spec.Domain, error) {
	// epoch is used to calculate fork version, here we hard code it
	return types.ComputeETHDomain(domain, types.GenesisForkVersion, types.GenesisValidatorsRoot)
}
