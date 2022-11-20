package testingutils

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"strings"
)

type TestingStorage struct {
	storage        map[string]*qbft.SignedMessage
	instancesState map[string]*qbft.State
	operators      map[types.OperatorID]*dkg.Operator
	keygenoupts    map[string]*dkg.KeyGenOutput
}

func NewTestingStorage() *TestingStorage {
	ret := &TestingStorage{
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
func (s *TestingStorage) SaveHighestDecided(signedMsg *qbft.SignedMessage) error {
	s.storage[hex.EncodeToString(signedMsg.Message.Identifier)] = signedMsg
	return nil
}

// GetHighestDecided returns highest decided if found, nil if didn't
func (s *TestingStorage) GetHighestDecided(identifier []byte) (*qbft.SignedMessage, error) {
	return s.storage[hex.EncodeToString(identifier)], nil
}

func (s *TestingStorage) SaveInstanceState(state *qbft.State) error {
	key := fmt.Sprintf("%s_%d", hex.EncodeToString(state.ID), state.Height)

	copiedState := &qbft.State{}
	*copiedState = *state
	s.instancesState[key] = copiedState
	return nil
}

func (s *TestingStorage) GetInstanceState(identifier []byte, height qbft.Height) (*qbft.State, error) {
	key := fmt.Sprintf("%s_%d", hex.EncodeToString(identifier), height)
	state := s.instancesState[key]
	if state == nil {
		return state, nil
	}
	// in order to mock storage without same pointer
	copiedState := &qbft.State{}
	*copiedState = *state
	return copiedState, nil
}

func (s *TestingStorage) GetAllInstancesState(identifier []byte) ([]*qbft.State, error) {
	var res []*qbft.State
	for k, state := range s.instancesState {
		if strings.HasPrefix(k, hex.EncodeToString(identifier)) {
			res = append(res, state)
		}
	}
	return res, nil
}

// GetDKGOperator returns true and operator object if found by operator ID
func (s *TestingStorage) GetDKGOperator(operatorID types.OperatorID) (bool, *dkg.Operator, error) {
	if ret, found := s.operators[operatorID]; found {
		return true, ret, nil
	}
	return false, nil, nil
}

func (s *TestingStorage) SaveKeyGenOutput(output *dkg.KeyGenOutput) error {
	s.keygenoupts[hex.EncodeToString(output.ValidatorPK)] = output
	return nil
}

func (s *TestingStorage) GetKeyGenOutput(pk types.ValidatorPK) (*dkg.KeyGenOutput, error) {
	return s.keygenoupts[hex.EncodeToString(pk)], nil
}
