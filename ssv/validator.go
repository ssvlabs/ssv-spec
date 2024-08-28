package ssv

import (
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// Validator represents an SSV ETH consensus validator Share assigned, coordinates duty execution and more.
// Every validator has a validatorID which is validator's public key.
// Each validator has multiple DutyRunners, for each duty type.
type Validator struct {
	DutyRunners     DutyRunners
	CommitteeMember *types.CommitteeMember
	Share           *types.Share
}

func NewValidator(
	committeeMember *types.CommitteeMember,
	share *types.Share,
	runners map[types.RunnerRole]Runner,
) *Validator {
	return &Validator{
		DutyRunners:     runners,
		Share:           share,
		CommitteeMember: committeeMember,
	}
}

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty types.Duty) error {
	role := duty.RunnerRole()
	dutyRunner := v.DutyRunners[role]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", role.String())
	}
	return dutyRunner.StartNewDuty(duty, v.CommitteeMember.GetQuorum())
}

// ProcessMessage processes a message of all types
func (v *Validator) ProcessMessage(msg *types.SignedSSVMessage) error {
	// Get runner for message
	dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.SSVMessage.MsgID)
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID")
	}
	// Process message
	return RunnerProcessMessage(dutyRunner, msg)
}

// Returns true if message is intended to the validator according to its MessageID
func (v *Validator) isMessageForValidator(msg *types.SignedSSVMessage) bool {
	return v.Share.ValidatorPubKey.MessageIDBelongs(msg.SSVMessage.MsgID)
}
