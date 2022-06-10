package dkg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type ProtocolOutput struct {
	Share       *bls.SecretKey
	ValidatorPK types.ValidatorPK
}

// Protocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type Protocol interface {
	Start() error
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *SignedMessage) (bool, *ProtocolOutput, error)
}

type Config struct {
	Network Network
}

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operators             []types.OperatorID
	WithdrawalCredentials []byte

	protocol Protocol
	config   *Config
}

func StartNewDKG(initMsg *Init, config *Config) (*Runner, error) {
	runner := &Runner{
		Operators:             initMsg.OperatorIDs,
		WithdrawalCredentials: initMsg.WithdrawalCredentials,

		protocol: NewSimpleDKG(config.Network),
		config:   config,
	}

	if err := runner.protocol.Start(); err != nil {
		return nil, errors.Wrap(err, "could not start dkg protocol")
	}

	return runner, nil
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, *SignedOutput, error) {
	// TODO - validate message

	finished, o, err := r.protocol.ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to process dkg msg")
	}

	if finished {
		ret, err := r.generateSignedOutput(o)
		if err != nil {
			return false, nil, errors.Wrap(err, "could not generate dkg SignedOutput")
		}
		return true, ret, nil
	}

	return false, nil, nil
}

func (r *Runner) generateSignedOutput(protocolOutput *ProtocolOutput) (*SignedOutput, error) {
	panic("implement")
}
