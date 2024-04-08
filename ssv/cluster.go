package ssv

import (
	"fmt"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type Cluster struct {
	Runners        map[spec.Slot]*ClusterRunner
	Network        Network
	Beacon         BeaconNode
	Signer         types.KeyManager
	CreateRunnerFn func() *ClusterRunner
}

func NewCluster(network Network, beacon BeaconNode, signer types.KeyManager, createRunnerFn func() *ClusterRunner) *Cluster {
	return &Cluster{
		Runners:        make(map[spec.Slot]*ClusterRunner),
		Network:        network,
		Beacon:         beacon,
		Signer:         signer,
		CreateRunnerFn: createRunnerFn,
	}
}

// StartDuty starts a new duty for the given slot
func (c *Cluster) StartDuty(duty *types.ClusterDuty) error {
	// do we need slot?
	if _, exists := c.Runners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("ClusterRunner for slot %d already exists", duty.Slot))
	}
	c.Runners[duty.Slot] = c.CreateRunnerFn()
	return c.Runners[duty.Slot].StartNewDuty(duty)
}

// ProcessMessage processes Network Message of all types
func (c *Cluster) ProcessMessage(msg *types.SSVMessage) error {
	//TODO validate message
	//dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.GetID())
	//if dutyRunner == nil {
	//	return errors.Errorf("could not get duty runner for msg ID")
	//}
	//
	//if err := v.validateMessage(dutyRunner, msg); err != nil {
	//	return errors.Wrap(err, "Message invalid")
	//}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}
		return c.Runners[spec.Slot(signedMsg.Message.Height)].ProcessConsensus(signedMsg)
	case types.SSVPartialSignatureMsgType:
		signedMsg := &types.SignedPartialSignatureMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		if signedMsg.Message.Type == types.PostConsensusPartialSig {
			return c.Runners[signedMsg.Message.Slot].ProcessPostConsensus(signedMsg)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

func validateMessage(msg *types.SSVMessage) error {
	panic("implement me")
}
