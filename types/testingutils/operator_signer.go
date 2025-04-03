package testingutils

import "github.com/ssvlabs/ssv-spec/types"

func NewOperatorSigner(keySet *TestKeySet, operatorID types.OperatorID) *types.OperatorSigner {
	return &types.OperatorSigner{
		SSVOperatorSK: keySet.OperatorKeys[operatorID],
		OperatorID:    operatorID,
	}
}
