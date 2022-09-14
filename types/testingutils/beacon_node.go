package testingutils

import (
	"encoding/hex"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/prysmaticlabs/go-bitfield"
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

var TestingBeaconBlock = &bellatrix.BeaconBlock{
	Slot:          12,
	ProposerIndex: 10,
	ParentRoot:    phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	StateRoot:     phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Body: &bellatrix.BeaconBlockBody{
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
		ExecutionPayload: &bellatrix.ExecutionPayload{
			ParentHash:    phase0.Hash32{},
			FeeRecipient:  bellatrix.ExecutionAddress{},
			StateRoot:     phase0.Hash32{},
			ReceiptsRoot:  phase0.Hash32{},
			LogsBloom:     [256]byte{},
			PrevRandao:    [32]byte{},
			BlockNumber:   100,
			GasLimit:      1000000,
			GasUsed:       800000,
			Timestamp:     123456789,
			BaseFeePerGas: [32]byte{},
			BlockHash:     phase0.Hash32{},
			Transactions:  []bellatrix.Transaction{},
		},
	},
}

var TestingAggregateAndProof = &phase0.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  phase0.BLSSignature{},
	Aggregate: &phase0.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       phase0.BLSSignature{},
		Data:            TestingAttestationData,
	},
}

const (
	TestingDutySlot       = 12
	TestingDutySlot2      = 50
	TestingDutyEpoch      = 0
	TestingDutyEpoch2     = 1
	TestingValidatorIndex = 1

	UnknownDutyType = 100
)

var TestingSyncCommitteeBlockRoot = phase0.Root{}

var TestingContributionProofIndexes = []uint64{0, 1, 2}
var TestingContributionProofsSigned = func() []phase0.BLSSignature {
	// signed with 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
	byts1, _ := hex.DecodeString("b18833bb7549ec33e8ac414ba002fd45bb094ca300bd24596f04a434a89beea462401da7c6b92fb3991bd17163eb603604a40e8dd6781266c990023446776ff42a9313df26a0a34184a590e57fa4003d610c2fa214db4e7dec468592010298bc")
	byts2, _ := hex.DecodeString("9094342c95146554df849dc20f7425fca692dacee7cb45258ddd264a8e5929861469fda3d1567b9521cba83188ffd61a0dbe6d7180c7a96f5810d18db305e9143772b766d368aa96d3751f98d0ce2db9f9e6f26325702088d87f0de500c67c68")
	byts3, _ := hex.DecodeString("a7f88ce43eff3aa8cdd2e3957c5bead4e21353fbecac6079a5398d03019bc45ff7c951785172deee70e9bc5abbc8ca6a0f0441e9d4cc9da74c31121357f7d7c7de9533f6f457da493e3314e22d554ab76613e469b050e246aff539a33807197c")

	ret := make([]phase0.BLSSignature, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := phase0.BLSSignature{}
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
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

// TestingProposerDutyNextEpoch testing for a second duty start
var TestingProposerDutyNextEpoch = &types.Duty{
	Type:                    types.BNRoleProposer,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot2,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingAggregatorDuty = &types.Duty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

// TestingAggregatorDutyNextEpoch testing for a second duty start
var TestingAggregatorDutyNextEpoch = &types.Duty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot2,
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

// TestingSyncCommitteeContributionNexEpochDuty testing for a second duty start
var TestingSyncCommitteeContributionNexEpochDuty = &types.Duty{
	Type:                    types.BNRoleSyncCommitteeContribution,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot2,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingUnknownDutyType = &types.Duty{
	Type:                    UnknownDutyType,
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

//func blsSigFromHex(str string) phase0.BLSSignature {
//	byts, _ := hex.DecodeString(str)
//	ret := phase0.BLSSignature{}
//	copy(ret[:], byts)
//	return ret
//}

type TestingBeaconNode struct {
	syncCommitteeAggregatorRoots map[string]bool
}

func NewTestingBeaconNode() *TestingBeaconNode {
	return &TestingBeaconNode{}
}

// SetSyncCommitteeAggregatorRootHexes FOR TESTING ONLY!! sets which sync committee aggregator roots will return true for aggregator
func (bn *TestingBeaconNode) SetSyncCommitteeAggregatorRootHexes(roots map[string]bool) {
	bn.syncCommitteeAggregatorRoots = roots
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *TestingBeaconNode) GetBeaconNetwork() types.BeaconNetwork {
	return types.NowTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *TestingBeaconNode) GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	return TestingAttestationData, nil
}

// SubmitAttestation submit the attestation to the node
func (bn *TestingBeaconNode) SubmitAttestation(attestation *phase0.Attestation) error {
	return nil
}

// GetBeaconBlock returns beacon block by the given slot and committee index
func (bn *TestingBeaconNode) GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (*bellatrix.BeaconBlock, error) {
	return TestingBeaconBlock, nil
}

// SubmitBeaconBlock submit the block to the node
func (bn *TestingBeaconNode) SubmitBeaconBlock(block *bellatrix.SignedBeaconBlock) error {
	return nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
func (bn *TestingBeaconNode) SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, slotSig []byte) (*phase0.AggregateAndProof, error) {
	return TestingAggregateAndProof, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *TestingBeaconNode) SubmitSignedAggregateSelectionProof(msg *phase0.SignedAggregateAndProof) error {
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *TestingBeaconNode) GetSyncMessageBlockRoot() (phase0.Root, error) {
	return TestingSyncCommitteeBlockRoot, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *TestingBeaconNode) SubmitSyncMessage(msg *altair.SyncCommitteeMessage) error {
	return nil
}

// GetSyncSubcommitteeIndex returns sync committee indexes for aggregator
func (bn *TestingBeaconNode) GetSyncSubcommitteeIndex(slot phase0.Slot, pubKey phase0.BLSPubKey) ([]uint64, error) {
	// each subcommittee index correlates to TestingContributionProofRoots by index
	return TestingContributionProofIndexes, nil
}

// IsSyncCommitteeAggregator returns tru if aggregator
func (bn *TestingBeaconNode) IsSyncCommitteeAggregator(proof []byte) (bool, error) {
	if len(bn.syncCommitteeAggregatorRoots) != 0 {
		if val, found := bn.syncCommitteeAggregatorRoots[hex.EncodeToString(proof)]; found {
			return val, nil
		}
		return false, nil
	}
	return true, nil
}

// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
func (bn *TestingBeaconNode) SyncCommitteeSubnetID(subCommitteeID uint64) (uint64, error) {
	// each subcommittee index correlates to TestingContributionProofRoots by index
	return subCommitteeID, nil
}

// GetSyncCommitteeContribution returns
func (bn *TestingBeaconNode) GetSyncCommitteeContribution(slot phase0.Slot, subnetID uint64, pubKey phase0.BLSPubKey) (*altair.SyncCommitteeContribution, error) {
	return TestingSyncCommitteeContributions[subnetID], nil
}

// SubmitSignedContributionAndProof broadcasts to the network
func (bn *TestingBeaconNode) SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error {
	return nil
}

func (bn *TestingBeaconNode) DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error) {
	// epoch is used to calculate fork version, here we hard code it
	return types.ComputeETHDomain(domain, types.GenesisForkVersion, types.GenesisValidatorsRoot)
}
