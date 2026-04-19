package ssv

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type testRunner struct {
	startDuty   types.Duty
	startQuorum uint64
}

func (r *testRunner) Encode() ([]byte, error) { return nil, nil }
func (r *testRunner) Decode([]byte) error     { return nil }
func (r *testRunner) GetRoot() ([32]byte, error) {
	return [32]byte{}, nil
}
func (r *testRunner) GetBaseRunner() *BaseRunner               { return nil }
func (r *testRunner) GetBeaconNode() BeaconNode                { return nil }
func (r *testRunner) GetValCheckF() qbft.ProposedValueCheckF   { return nil }
func (r *testRunner) GetSigner() types.BeaconSigner            { return nil }
func (r *testRunner) GetOperatorSigner() *types.OperatorSigner { return nil }
func (r *testRunner) GetNetwork() Network                      { return nil }
func (r *testRunner) HasRunningDuty() bool                     { return false }
func (r *testRunner) ProcessPreConsensus(*types.PartialSignatureMessages) error {
	return nil
}
func (r *testRunner) ProcessConsensus(*types.SignedSSVMessage) error { return nil }
func (r *testRunner) ProcessPostConsensus(*types.PartialSignatureMessages) error {
	return nil
}
func (r *testRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, phase0.DomainType{}, nil
}
func (r *testRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, phase0.DomainType{}, nil
}
func (r *testRunner) executeDuty(types.Duty) error { return nil }
func (r *testRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	r.startDuty = duty
	r.startQuorum = quorum
	return nil
}

func validCommitteeMember() types.CommitteeMember {
	operatorIDs := []types.OperatorID{1, 2, 3, 4}
	return types.CommitteeMember{
		OperatorID:  1,
		CommitteeID: types.GetCommitteeID(operatorIDs),
		FaultyNodes: 1,
		Committee: []*types.Operator{
			{OperatorID: 1},
			{OperatorID: 2},
			{OperatorID: 3},
			{OperatorID: 4},
		},
		DomainType: types.JatoTestnet,
	}
}

func validShare(index phase0.ValidatorIndex) *types.Share {
	var validatorPubKey types.ValidatorPK
	validatorPubKey[0] = 1

	return &types.Share{
		ValidatorIndex:  index,
		ValidatorPubKey: validatorPubKey,
		Committee: []*types.ShareMember{
			{Signer: 1},
			{Signer: 2},
			{Signer: 3},
			{Signer: 4},
		},
		DomainType: types.JatoTestnet,
	}
}

func TestCommitteeStartDutyIgnoresMalformedForeignValidatorDuty(t *testing.T) {
	ownedIndex := phase0.ValidatorIndex(1)
	runner := &testRunner{}

	committee := NewCommittee(
		validCommitteeMember(),
		map[phase0.ValidatorIndex]*types.Share{
			ownedIndex: validShare(ownedIndex),
		},
		func(map[phase0.ValidatorIndex]*types.Share) Runner { return runner },
		func(map[phase0.ValidatorIndex]*types.Share) Runner { return runner },
	)

	var pubKey phase0.BLSPubKey
	pubKey[0] = 1

	err := committee.StartDuty(&types.CommitteeDuty{
		Slot: 10,
		ValidatorDuties: []*types.ValidatorDuty{
			{
				Type:                    types.BNRoleAttester,
				Slot:                    9,
				ValidatorIndex:          999,
				CommitteeLength:         4,
				ValidatorCommitteeIndex: 0,
			},
			{
				Type:                    types.BNRoleAttester,
				PubKey:                  pubKey,
				Slot:                    10,
				ValidatorIndex:          ownedIndex,
				CommitteeLength:         4,
				ValidatorCommitteeIndex: 0,
			},
		},
	})
	require.NoError(t, err)

	startedDuty, ok := runner.startDuty.(*types.CommitteeDuty)
	require.True(t, ok)
	require.Len(t, startedDuty.ValidatorDuties, 1)
	require.Equal(t, ownedIndex, startedDuty.ValidatorDuties[0].ValidatorIndex)
	require.Equal(t, committee.CommitteeMember.GetQuorum(), runner.startQuorum)
}

func TestCommitteeStartDutyRejectsInvalidOwnedAggregatorDuty(t *testing.T) {
	ownedIndex := phase0.ValidatorIndex(1)
	runner := &testRunner{}

	committee := NewCommittee(
		validCommitteeMember(),
		map[phase0.ValidatorIndex]*types.Share{
			ownedIndex: validShare(ownedIndex),
		},
		func(map[phase0.ValidatorIndex]*types.Share) Runner { return runner },
		func(map[phase0.ValidatorIndex]*types.Share) Runner { return runner },
	)

	var pubKey phase0.BLSPubKey
	pubKey[0] = 1

	err := committee.StartDuty(&types.AggregatorCommitteeDuty{
		Slot: 10,
		ValidatorDuties: []*types.ValidatorDuty{
			{
				Type:                    types.BNRoleAttester,
				PubKey:                  pubKey,
				Slot:                    10,
				ValidatorIndex:          ownedIndex,
				CommitteeLength:         4,
				ValidatorCommitteeIndex: 0,
			},
		},
	})
	require.EqualError(t, err, "invalid aggregator committee duty: invalid beacon role in validator duty")
}
