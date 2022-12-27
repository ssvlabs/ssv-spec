package testingutils

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

type testingStorage struct {
	operators       map[types.OperatorID]*dkg.Operator
	keygenoupts     map[string]*dkg.KeyGenOutput
	oldOperatorList map[string][]types.OperatorID
}

func NewTestingStorage() *testingStorage {
	ret := &testingStorage{
		operators:       make(map[types.OperatorID]*dkg.Operator),
		keygenoupts:     make(map[string]*dkg.KeyGenOutput),
		oldOperatorList: make(map[string][]types.OperatorID),
	}

	for i, s := range Testing13SharesSet().DKGOperators {
		ret.operators[i] = &dkg.Operator{
			OperatorID:       i,
			ETHAddress:       s.ETHAddress,
			EncryptionPubKey: &s.EncryptionKey.PublicKey,
		}
	}

	ret.oldOperatorList[TestingResharingKeySet().ValidatorPK.SerializeToHexStr()] = []types.OperatorID{1, 2, 3}
	return ret
}

// GetDKGOperator returns true and operator object if found by operator ID
func (s *testingStorage) GetDKGOperator(operatorID types.OperatorID) (bool, *dkg.Operator, error) {
	if ret, found := s.operators[operatorID]; found {
		return true, ret, nil
	}
	return false, nil, nil
}

func (s *testingStorage) SaveKeyGenOutput(output *dkg.KeyGenOutput) error {
	s.keygenoupts[hex.EncodeToString(output.ValidatorPK)] = output
	return nil
}

func (s *testingStorage) GetKeyGenOutput(pk types.ValidatorPK) (*dkg.KeyGenOutput, error) {
	return s.keygenoupts[hex.EncodeToString(pk)], nil
}

func (s *testingStorage) GetDKGOperators(pk types.ValidatorPK) ([]types.OperatorID, error) {
	return s.oldOperatorList[hex.EncodeToString(pk)], nil
}
