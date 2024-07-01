package ssv

import (
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CreateRunnerFn func(shareMap map[spec.ValidatorIndex]*types.Share) *CommitteeRunner

type Committee struct {
	Runners           map[spec.Slot]*CommitteeRunner
	CommitteeMember   types.CommitteeMember
	SignatureVerifier types.SignatureVerifier
	CreateRunnerFn    CreateRunnerFn
	Share             map[spec.ValidatorIndex]*types.Share
}

// NewCommittee creates a new cluster
func NewCommittee(
	committeeMember types.CommitteeMember,
	verifier types.SignatureVerifier,
	share map[spec.ValidatorIndex]*types.Share,
	createRunnerFn CreateRunnerFn,
) *Committee {
	c := &Committee{
		Runners:           make(map[spec.Slot]*CommitteeRunner),
		CommitteeMember:   committeeMember,
		SignatureVerifier: verifier,
		CreateRunnerFn:    createRunnerFn,
		Share:             share,
	}
	return c
}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty *types.CommitteeDuty) error {
	if len(duty.ValidatorDuties) == 0 {
		return errors.New("no beacon duties")
	}
	if _, exists := c.Runners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("CommitteeRunner for slot %d already exists", duty.Slot))
	}
	c.Runners[duty.Slot] = c.CreateRunnerFn(c.Share)
	return c.Runners[duty.Slot].StartNewDuty(duty, c.CommitteeMember.GetQuorum())
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := c.SignatureVerifier.Verify(signedSSVMessage, c.CommitteeMember.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage
	if err := c.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}

		if err := qbftMsg.Validate(); err != nil {
			return errors.Wrap(err, "invalid qbft Message")
		}

		runner, exists := c.Runners[spec.Slot(qbftMsg.Height)]
		if !exists {
			return errors.New("no runner found for message's slot")
		}
		return runner.ProcessConsensus(signedSSVMessage)
	case types.SSVPartialSignatureMsgType:
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return errors.New("PartialSignatureMessage has more than 1 signer")
		}

		if err := pSigMessages.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		if pSigMessages.Type == types.PostConsensusPartialSig {
			runner, exists := c.Runners[pSigMessages.Slot]
			if !exists {
				return errors.New("no runner found for message's slot")
			}
			return runner.ProcessPostConsensus(pSigMessages)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

func (c *Committee) validateMessage(msg *types.SSVMessage) error {
	if !(c.CommitteeMember.CommitteeID.MessageIDBelongs(msg.GetID())) {
		return errors.New("msg ID doesn't match committee ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
