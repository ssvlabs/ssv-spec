package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	ecies "github.com/ecies/go/v2"
)

type ProtocolConfig struct {
	identifier      dkg.RequestID
	threshold       uint32
	operatorID      types.OperatorID
	operators       []uint32
	operatorsOld    []uint32
	oldKeyGenOutput *dkg.KeyGenOutput
	sessionSK       *ecies.PrivateKey
}

func (c *ProtocolConfig) isResharing() bool {
	return len(c.operatorsOld) > 0
}

func (c *ProtocolConfig) inOldCommittee() bool {
	for _, id := range c.operatorsOld {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}

func (c *ProtocolConfig) inNewCommittee() bool {
	for _, id := range c.operators {
		if types.OperatorID(id) == c.operatorID {
			return true
		}
	}
	return false
}
