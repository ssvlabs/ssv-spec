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

// State maintains value for current round, messages, sessions key and
// operator shares. these properties will change over the runtime of the protocol
type State struct {
	currentRound   ProtocolRound
	participant    *frost.DkgParticipant
	sessionSK      *ecies.PrivateKey
	msgContainer   IMsgContainer
	operatorShares map[uint32]*bls.SecretKey
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
