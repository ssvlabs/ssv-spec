package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

// Holds information about a committee including its validators and operators
type CommitteeInfo struct {
	Operators   []*types.Operator
	Validators  []phase0.ValidatorIndex
	CommitteeID types.CommitteeID
}

// Returns the list of operatorIDs that belong to the committee
func (ci *CommitteeInfo) OperatorIDs() []types.OperatorID {
	operatorIDs := make([]types.OperatorID, 0)
	for _, operator := range ci.Operators {
		operatorIDs = append(operatorIDs, operator.OperatorID)
	}
	return operatorIDs
}
