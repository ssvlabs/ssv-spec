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

type Committee struct {
	Runners                 map[spec.Slot]*CommitteeRunner
	Operator                types.Operator
	SignatureVerifier       types.SignatureVerifier
	CreateRunnerFn          func() *CommitteeRunner
	HighestAttestingSlotMap map[types.ValidatorPK]spec.Slot
}

// NewCommittee creates a new cluster
func NewCommittee(
	network Network,
	beacon BeaconNode,
	operator types.Operator,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	verifier types.SignatureVerifier,
	createRunnerFn func() *CommitteeRunner,
) *Committee {
	return &Committee{
		Runners:           make(map[spec.Slot]*CommitteeRunner),
		Operator:          operator,
		SignatureVerifier: verifier,
		CreateRunnerFn:    createRunnerFn,
	}

}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty *types.CommitteeDuty) error {
	// do we need slot?
	if _, exists := c.Runners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("CommitteeRunner for slot %d already exists", duty.Slot))
	}
	c.Runners[duty.Slot] = c.CreateRunnerFn()
	validatorToStopMap := make(map[spec.Slot]types.ValidatorPK)
	duty, validatorToStopMap, c.HighestAttestingSlotMap = FilterCommitteeDuty(duty, c.HighestAttestingSlotMap)
	c.StopDuties(validatorToStopMap)
	return c.Runners[duty.Slot].StartNewDuty(duty)
}

func (c *Committee) StopDuties(validatorToStopMap map[spec.Slot]types.ValidatorPK) {
	for slot, validator := range validatorToStopMap {
		c.Runners[slot].StopDuty(validator)
	}
}

// FilterCommitteeDuty filters the committee duty. It returns the new duty, the validators to stop and the highest attesting slot map
func FilterCommitteeDuty(duty *types.CommitteeDuty, slotMap map[types.ValidatorPK]spec.Slot) (
	*types.CommitteeDuty,
	map[spec.Slot]types.ValidatorPK,
	map[types.ValidatorPK]spec.Slot) {
	validatorsToStop := make(map[spec.Slot]types.ValidatorPK)

	for i, beaconDuty := range duty.BeaconDuties {
		validatorPK := types.ValidatorPK(beaconDuty.PubKey)
		if slotMap[validatorPK] < beaconDuty.Slot {
			validatorsToStop[beaconDuty.Slot] = validatorPK
			slotMap[validatorPK] = beaconDuty.Slot
		} else { // else don't run duty with old slot
			duty.BeaconDuties[i] = nil
		}
	}
	return duty, validatorsToStop, slotMap
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := c.SignatureVerifier.Verify(signedSSVMessage, c.Operator.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}
		runner := c.Runners[spec.Slot(qbftMsg.Height)]
		// TODO: check if runner is nil
		return runner.ProcessConsensus(signedSSVMessage)
	case types.SSVPartialSignatureMsgType:
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}
		if pSigMessages.Type == types.PostConsensusPartialSig {
			runner := c.Runners[pSigMessages.Slot]
			// TODO: check if runner is nil
			return runner.ProcessPostConsensus(pSigMessages)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

func (c *Committee) validateMessage(msg *types.SSVMessage) error {
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
