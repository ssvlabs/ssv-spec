package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type IConfig interface {
	// GetSigner returns a Signer instance
	GetSigner() types.DKGSigner
	// GetNetwork returns a p2p Network instance
	GetNetwork() dkg.Network
	// GetStorage returns a Storage instance
	GetStorage() dkg.Storage
}

type Config struct {
	network dkg.Network
	signer  types.DKGSigner
	storage dkg.Storage
}

// GetSigner returns a Signer instance
func (c *Config) GetSigner() types.DKGSigner {
	return c.signer
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() dkg.Network {
	return c.network
}

// GetStorage returns a Storage instance
func (c *Config) GetStorage() dkg.Storage {
	return c.storage
}

// ProtocolRound is enum for all the rounds in the protocol
type ProtocolRound int

const (
	Uninitialized ProtocolRound = iota
	Preparation
	Round1
	Round2
	KeygenOutput
	Blame
)

var rounds = []ProtocolRound{
	Uninitialized,
	Preparation,
	Round1,
	Round2,
	KeygenOutput,
	Blame,
}

// State tracks protocol's current round, stores messages in MsgContainer, stores
// session key and operator's secret shares
type State struct {
	currentRound   ProtocolRound
	participant    *frost.DkgParticipant
	sessionSK      *ecies.PrivateKey
	msgContainer   IMsgContainer
	operatorShares map[uint32]*bls.SecretKey
}

func initState() *State {
	return &State{
		currentRound:   Uninitialized,
		msgContainer:   newMsgContainer(),
		operatorShares: make(map[uint32]*bls.SecretKey),
	}
}

func (state *State) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
	msg, err := state.msgContainer.GetPreparationMsg(operatorID)
	if err != nil {
		return nil, errors.Wrapf(err, "no session pk found for the operator")
	}
	sessionPK, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
	if err != nil {
		return nil, err
	}
	return ecies.Encrypt(sessionPK, data)
}

// InstanceParams contains properties needed to start the protocol like requestID,
// operatorID, threshold, operator list etc.
type InstanceParams struct {
	identifier      dkg.RequestID
	threshold       uint32
	operatorID      types.OperatorID
	operators       []uint32
	operatorsOld    []uint32
	oldKeyGenOutput *dkg.KeyGenOutput
}

func (c *InstanceParams) isResharing() bool {
	return len(c.operatorsOld) > 0
}

func (c *InstanceParams) inOldCommittee() bool {
	for _, id := range c.operatorsOld {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}

func (c *InstanceParams) inNewCommittee() bool {
	for _, id := range c.operators {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}
