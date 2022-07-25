package testutils

import (
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

type MockStorage struct {
	keySet testingutils.TestKeySet
}

func (m MockStorage) GetDKGOperator(operatorID types.OperatorID) (bool, *dkgtypes.Operator, error) {
	found := m.keySet.DKGOperators[operatorID]

	o := &dkgtypes.Operator{
		OperatorID:       operatorID,
		ETHAddress:       found.ETHAddress,
		EncryptionPubKey: &found.EncryptionKey.PublicKey,
	}

	return true, o, nil
}

func newMockStorage(keySet testingutils.TestKeySet) *MockStorage {
	return &MockStorage{
		keySet: keySet,
	}
}
