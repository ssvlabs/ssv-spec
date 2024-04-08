package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

type ClusterRunner struct {
	BaseRunner *BaseRunner
	beacon     BeaconNode
	network    Network
	signer     types.KeyManager
	valCheck   qbft.ProposedValueCheckF
}

func NewClusterRunner(beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot) *ClusterRunner {
	return &ClusterRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleAttester,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},
		beacon:  beacon,
		network: network,
		signer:  signer,
	}
}

func (cr ClusterRunner) StartNewDuty(duty types.Duty) error {
	return cr.BaseRunner.baseStartNewDuty(cr, duty)
}

func (cr ClusterRunner) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) Decode(data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetRoot() ([32]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetBaseRunner() *BaseRunner {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetBeaconNode() BeaconNode {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetValCheckF() qbft.ProposedValueCheckF {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetSigner() types.KeyManager {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) GetNetwork() Network {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) HasRunningDuty() bool {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) ProcessConsensus(msg *qbft.SignedMessage) error {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	//TODO implement me
	panic("implement me")
}

func (cr ClusterRunner) executeDuty(duty types.Duty) error {

	//TODO committeeIndex is 0, is this correct?
	attData, ver, err := cr.GetBeaconNode().GetAttestationData(duty.DutySlot(), 0)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	vote := types.BeaconVote{
		BlockRoot: attData.BeaconBlockRoot,
		Source:    attData.Source,
		Target:    attData.Target,
	}
	voteByts, err := vote.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal attestation data")
	}

	//TODO should duty be empty?
	input := &types.ConsensusData{
		Duty:    types.BeaconDuty{},
		Version: ver,
		DataSSZ: voteByts,
	}

	if err := cr.BaseRunner.decide(cr, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}
