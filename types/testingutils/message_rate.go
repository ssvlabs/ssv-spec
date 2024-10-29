package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/messagerate"
	"github.com/ssvlabs/ssv-spec/types"
)

var ToValidatorIndex = func(lst []int) []phase0.ValidatorIndex {
	validatorIndexList := make([]phase0.ValidatorIndex, 0)
	for _, valIdx := range lst {
		validatorIndexList = append(validatorIndexList, phase0.ValidatorIndex(valIdx))
	}
	return validatorIndexList
}

var TestingDisjointCommittees = func(numCommittees int, numValidatorsPerCommittee int) []*messagerate.Committee {
	committees := make([]*messagerate.Committee, 0)
	for i := 0; i < numCommittees; i++ {
		committees = append(committees, &messagerate.Committee{
			Operators:  []types.OperatorID{1, 2, 3, uint64(4 + i)},
			Validators: ToValidatorIndex(ValidatorIndexList(numValidatorsPerCommittee)),
		})
	}
	return committees
}
