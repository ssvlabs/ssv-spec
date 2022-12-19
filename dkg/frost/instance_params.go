package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

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
