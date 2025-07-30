package tests

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgSpecTest tests encoding and decoding of a msg
type MsgSpecTest struct {
	Name            string
	Type            string
	Documentation   string
	Messages        []*types.SignedSSVMessage
	EncodedMessages [][]byte
	ExpectedRoots   [][32]byte
	ExpectedError   string
	PrivateKeys     *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error

	for i, msg := range test.Messages {
		if err := msg.Validate(); err != nil {
			lastErr = err
			continue
		}

		qbftMessage := &qbft.Message{}
		require.NoError(t, qbftMessage.Decode(msg.SSVMessage.Data))
		if err := qbftMessage.Validate(); err != nil {
			lastErr = err
			continue
		}

		if len(test.EncodedMessages) > 0 {
			byts, err := msg.Encode()
			require.NoError(t, err)
			require.EqualValues(t, test.EncodedMessages[i], byts)
		}

		if len(test.ExpectedRoots) > 0 {
			r, err := msg.GetRoot()
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

func (test *MsgSpecTest) GetPostState() (interface{}, error) {
	// remove private keys
	test.PrivateKeys = nil

	return test, nil
}

func NewMsgSpecTest(name string, documentation string, messages []*types.SignedSSVMessage, encodedMessages [][]byte, expectedRoots [][32]byte, expectedError string, ks *testingutils.TestKeySet) *MsgSpecTest {
	return &MsgSpecTest{
		Name:            name,
		Type:            testdoc.MsgSpecTestType,
		Documentation:   documentation,
		Messages:        messages,
		EncodedMessages: encodedMessages,
		ExpectedRoots:   expectedRoots,
		ExpectedError:   expectedError,
		PrivateKeys:     testingutils.BuildPrivateKeyInfo(ks),
	}
}
