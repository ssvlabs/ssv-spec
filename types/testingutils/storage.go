package testingutils

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"strings"
)

type testingStorage struct {
	storage        map[string]*qbft.SignedMessage
	instancesState map[string]*qbft.State
	operators      map[types.OperatorID]*dkg.Operator
	keygenoupts    map[string]*dkg.KeyGenOutput
}

func NewTestingStorage() *testingStorage {
	ret := &testingStorage{
		storage:        make(map[string]*qbft.SignedMessage),
		instancesState: make(map[string]*qbft.State),
		operators:      make(map[types.OperatorID]*dkg.Operator),
		keygenoupts:    make(map[string]*dkg.KeyGenOutput),
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

func (s *testingStorage) SaveInstanceState(state *qbft.State) error {
	key := fmt.Sprintf("%s_%d", hex.EncodeToString(state.ID), state.Height)
	s.instancesState[key] = state
	return nil
}

func (s *testingStorage) GetInstanceState(identifier []byte, height qbft.Height) (*qbft.State, error) {
	key := fmt.Sprintf("%s_%d", hex.EncodeToString(identifier), height)
	return s.instancesState[key], nil
}

func (s *testingStorage) GetAlInstancesState(identifier []byte) ([]*qbft.State, error) {
	var res []*qbft.State
	for k, state := range s.instancesState {
		if strings.HasPrefix(k, hex.EncodeToString(identifier)) {
			res = append(res, state)
		}
	}
	return res, nil
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
