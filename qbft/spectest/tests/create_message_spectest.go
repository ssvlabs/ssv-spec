package tests

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	CreateProposal = "createProposal"
)

type CreateMsgSpecTest struct {
	Name                                             string
	Value                                            []byte
	RoundChangeJustifications, PrepareJustifications []*qbft.SignedMessage
	CreateType                                       string
	ExpectedRoot                                     string
	ExpectedError                                    string
}

func (test *CreateMsgSpecTest) Run(t *testing.T) {
	var msg *qbft.SignedMessage
	var lastErr error
	switch test.CreateType {
	case CreateProposal:
		msg, lastErr = test.createProposal()
	default:
		t.Fail()
	}

	r, err := msg.GetRoot()
	if err != nil {
		lastErr = err
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	require.EqualValues(t, test.ExpectedRoot, hex.EncodeToString(r))
}

func (test *CreateMsgSpecTest) createProposal() (*qbft.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &qbft.State{
		Share: testingutils.TestingShare(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfig(ks)

	return qbft.CreateProposal(state, config, test.Value, test.RoundChangeJustifications, test.PrepareJustifications)
}

func (test *CreateMsgSpecTest) TestName() string {
	return test.Name
}
