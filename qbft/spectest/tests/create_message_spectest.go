package tests

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	CreateProposal    = "createProposal"
	CreatePrepare     = "CreatePrepare"
	CreateCommit      = "CreateCommit"
	CreateRoundChange = "CreateRoundChange"
)

type CreateMsgSpecTest struct {
	Name                                             string
	Value                                            []byte
	Round                                            qbft.Round
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
	case CreatePrepare:
		msg, lastErr = test.createPrepare()
	case CreateCommit:
		msg, lastErr = test.createCommit()
	case CreateRoundChange:
		msg, lastErr = test.createRoundChange()
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

func (test *CreateMsgSpecTest) createCommit() (*qbft.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	msgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	state := &qbft.State{
		Share: testingutils.TestingShare(ks),
		ID:    msgId,
	}
	config := testingutils.TestingConfig(ks)

	return qbft.CreateCommit(state, config, test.Value)
}

func (test *CreateMsgSpecTest) createPrepare() (*qbft.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	msgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	state := &qbft.State{
		Share: testingutils.TestingShare(ks),
		ID:    msgId,
	}
	config := testingutils.TestingConfig(ks)

	return qbft.CreatePrepare(state, config, test.Round, test.Value)
}

func (test *CreateMsgSpecTest) createProposal() (*qbft.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	msgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	state := &qbft.State{
		Share: testingutils.TestingShare(ks),
		ID:    msgId,
	}
	config := testingutils.TestingConfig(ks)

	prepareJustifications := make([]*qbft.SignedMessageHeader, 0)
	for _, rc := range test.PrepareJustifications {
		prepareHeader, err := rc.ToSignedMessageHeader()
		if err != nil {
			return nil, errors.Wrap(err, "could not convert signed msg to signed msg header")
		}
		prepareJustifications = append(prepareJustifications, prepareHeader)
	}

	return qbft.CreateProposal(state, config, test.Value, test.RoundChangeJustifications, prepareJustifications)
}

func (test *CreateMsgSpecTest) createRoundChange() (*qbft.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	msgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	state := &qbft.State{
		Share: testingutils.TestingShare(ks),
		ID:    msgId,
	}
	config := testingutils.TestingConfig(ks)

	if len(test.PrepareJustifications) > 0 {
		state.LastPreparedRound = test.PrepareJustifications[0].Message.Round
		state.LastPreparedValue = test.Value

		for _, msg := range test.PrepareJustifications {
			state.PrepareContainer.AddIfDoesntExist(msg)
		}
	}

	return qbft.CreateRoundChange(state, config, 1, test.Value)
}

func (test *CreateMsgSpecTest) TestName() string {
	return "qbft create message " + test.Name
}
