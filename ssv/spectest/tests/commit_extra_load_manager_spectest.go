package tests

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/stretchr/testify/require"
)

type CommitExtraLoadManagerTest struct {
	Name                         string
	CommitExtraLoadManagerF      qbft.CommitExtraLoadManagerF
	FullDataToCreate             [][]byte
	ExpectedCommitExtraLoadRoots [][32]byte
	ValidateSignedMessages       []*qbft.SignedMessage
	ValidateFullData             [][]byte
	ExpectedError                string
}

func (test *CommitExtraLoadManagerTest) TestName() string {
	return "commitextraloadmanager " + test.Name
}

func (test *CommitExtraLoadManagerTest) Run(t *testing.T) {

	// Create manager
	manager := test.CommitExtraLoadManagerF()

	var lastErr error

	// Create CommitExtraLoads
	for idx, fullData := range test.FullDataToCreate {
		commitExtraLoad, err := manager.Create(fullData)
		if err == nil {
			actualRoot, err := commitExtraLoad.GetRoot()
			if err != nil {
				panic(err.Error())
			}
			if !bytes.Equal(actualRoot[:], test.ExpectedCommitExtraLoadRoots[idx][:]) {
				panic(fmt.Sprintf("Roots not equal. Expected: %v Actual: %v", actualRoot, test.ExpectedCommitExtraLoadRoots[idx]))
			}
		} else {
			lastErr = err
		}
	}

	// Validate and process messages
	for idx, msg := range test.ValidateSignedMessages {
		err := manager.Validate(msg, test.ValidateFullData[idx])
		if err == nil {
			err = manager.Process(msg)
		}
		if err != nil {
			lastErr = err
		}
	}

	// Check expected error
	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *CommitExtraLoadManagerTest) GetPostState() (interface{}, error) {
	return nil, nil
}
