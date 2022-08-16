package testingutils

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type testingStorage struct {
	storage   map[string]*qbft.SignedMessage
	operators map[types.OperatorID]*dkg.Operator
}

func NewTestingStorage() *testingStorage {
	ret := &testingStorage{
		storage:   make(map[string]*qbft.SignedMessage),
		operators: map[types.OperatorID]*dkg.Operator{},
	}

	for i, s := range Testing13SharesSet().DKGOperators {
		ret.operators[i] = &dkg.Operator{
			OperatorID:       i,
			ETHAddress:       s.ETHAddress,
			EncryptionPubKey: &s.EncryptionKey.PublicKey,
		}
	}

	return ret
}

// SaveHighestDecided saves the Decided value as highest for a validator PK and role
func (s *testingStorage) SaveHighestDecided(signedMsg *qbft.SignedMessage) error {
	s.storage[hex.EncodeToString(signedMsg.Message.Identifier)] = signedMsg
	return nil
}

// GetHighestDecided returns highest decided if found, nil if didn't
func (s *testingStorage) GetHighestDecided(identifier []byte) (*qbft.SignedMessage, error) {
	return s.storage[hex.EncodeToString(identifier)], nil
}

// GetDKGOperator returns true and operator object if found by operator ID
func (s *testingStorage) GetDKGOperator(operatorID types.OperatorID) (bool, *dkg.Operator, error) {
	if ret, found := s.operators[operatorID]; found {
		return true, ret, nil
	}
	return false, nil, nil
}
