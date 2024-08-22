package validation

import (
	"github.com/ssvlabs/ssv-spec/types"
)

// Performs data fetching from the Network
type NetworkDataFetcher interface {
	ValidDomain(domain []byte) bool
	CorrectTopic(committee []types.OperatorID, topic string) bool
	GetCommitteeInfo(msgID types.MessageID) *CommitteeInfo
	ExistingValidator(validatorPK types.ValidatorPK) bool
	ActiveValidator(validatorPK types.ValidatorPK) bool
	ValidatorLiquidated(validatorPK types.ValidatorPK) bool
	ExistingCommitteeID(committeeID types.CommitteeID) bool
}
