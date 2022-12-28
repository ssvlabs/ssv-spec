package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

// InstanceParams contains properties needed to start the protocol like requestID,
// operatorID, threshold, operator list etc.
type InstanceParams struct {
	// unique identifier for this instance of keygen or resharing
	identifier dkg.RequestID
	// minimum number of operators require to sign a message
	threshold uint32
	// operator ID for which this instance is running
	operatorID types.OperatorID
	// list of all the operators in the keygen. New committee in case of resharing
	operators []uint32
	// list of operators from old committee in case of resharing. Atleast
	// t (from old keygen) number of operators required to proceed with resharing.
	operatorsOld []uint32
	// Keygen output from old keygen instance from this operator in old committee
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
