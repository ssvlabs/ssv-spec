package ssv

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"sort"
)

type Cluster struct {
	Runners           map[spec.Slot]*ClusterRunner
	Network           Network
	Beacon            BeaconNode
	Operator          types.Operator
	Share             *types.Share
	Signer            types.BeaconSigner
	OperatorSigner    types.OperatorSigner
	SignatureVerifier types.SignatureVerifier
	CreateRunnerFn    func() *ClusterRunner
}

// NewCluster creates a new cluster
func NewCluster(
	network Network,
	beacon BeaconNode,
	operator types.Operator,
	share *types.Share,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	verifier types.SignatureVerifier,
	createRunnerFn func() *ClusterRunner,
) *Cluster {
	return &Cluster{
		Runners:           make(map[spec.Slot]*ClusterRunner),
		Network:           network,
		Beacon:            beacon,
		Operator:          operator,
		Share:             share,
		Signer:            signer,
		OperatorSigner:    operatorSigner,
		SignatureVerifier: verifier,
		CreateRunnerFn:    createRunnerFn,
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
func (c *Cluster) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	msg, err := signedSSVMessage.VerifyAndDecodeData(c.SignatureVerifier, c.Operator.Committee)
	if err != nil {
		return err
	}

	if err := c.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}
		runner := c.Runners[spec.Slot(signedMsg.Message.Height)]
		// TODO: check if runner is nil
		return runner.ProcessConsensus(signedMsg)
	case types.SSVPartialSignatureMsgType:
		signedMsg := &types.SignedPartialSignatureMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		if signedMsg.Message.Type == types.PostConsensusPartialSig {
			runner := c.Runners[signedMsg.Message.Slot]
			// TODO: check if runner is nil
			return runner.ProcessPostConsensus(signedMsg)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

func (c *Cluster) validateMessage(msg *types.SSVMessage) error {
	if !c.Operator.ClusterID.MessageIDBelongs(msg.GetID()) {
		return errors.New("Message ID does not match cluster IF")
	}
	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}

type ClusterID [32]byte

func (cid ClusterID) MessageIDBelongs(msgID types.MessageID) bool {
	id := msgID.GetSenderID()[16:]
	return bytes.Equal(id, cid[:])
}

// Return a 32 bytes ID for the cluster of operators
func getClusterID(committee []types.OperatorID) ClusterID {
	// sort
	sort.Slice(committee, func(i, j int) bool {
		return committee[i] < committee[j]
	})
	// Convert to bytes
	bytes := make([]byte, len(committee)*4)
	for i, v := range committee {
		binary.LittleEndian.PutUint32(bytes[i*4:], uint32(v))
	}
	// Hash
	return sha256.Sum256(bytes)
}
