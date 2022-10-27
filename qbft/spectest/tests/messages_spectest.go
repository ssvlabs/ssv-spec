package tests

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

// MsgSpecTest tests encoding and decoding of a msg
type MsgSpecTest struct {
	Name            string
	Messages        []*types.Message
	EncodedMessages [][]byte
	ExpectedRoots   [][]byte
	ExpectedError   string
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error

	for i, msg := range test.Messages {
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			lastErr = err
		}
		if err := signedMsg.Validate(msg.GetID().GetMsgType()); err != nil {
			lastErr = err
			continue
		}

		switch msg.GetID().GetMsgType() {
		case types.ConsensusProposeMsgType:

		case types.ConsensusPrepareMsgType:

		case types.ConsensusCommitMsgType:

		case types.ConsensusRoundChangeMsgType:
			if signedMsg.Message.Prepared() {
				if len(signedMsg.RoundChangeJustifications) == 0 {
					lastErr = errors.New("round change justification invalid")
				}
			}
		}

		if len(test.EncodedMessages) > 0 {
			byts, err := signedMsg.Encode()
			require.NoError(t, err)
			require.EqualValues(t, test.EncodedMessages[i], byts)
		}

		if len(test.ExpectedRoots) > 0 {
			r, err := signedMsg.GetRoot()
			require.NoError(t, err)
			require.EqualValues(t, test.ExpectedRoots[i], r)
		}
	}

	// check error
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *MsgSpecTest) TestName() string {
	return "qbft message " + test.Name
}
