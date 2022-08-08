package messages

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgSpecTest struct {
	Name            string
	Messages        []*ssv.SignedPartialSignatureMessage
	EncodedMessages [][]byte
	ExpectedRoots   [][]byte
	ExpectedError   string
}

func (test *MsgSpecTest) TestName() string {
	return "msg " + test.Name
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error

	// test expected roots
	for i, byts := range test.EncodedMessages {
		m := &qbft.SignedMessage{}
		if err := m.Decode(byts); err != nil {
			lastErr = err
		}

		if len(test.ExpectedRoots) > 0 {
			r, err := m.GetRoot()
			if err != nil {
				lastErr = err
			}
			if !bytes.Equal(test.ExpectedRoots[i], r) {
				t.Fail()
			}
		}
	}

	// test encoding and validation
	for i, msg := range test.Messages {
		if err := msg.Validate(); err != nil {
			lastErr = err
		}

		if len(test.Messages) > 0 {
			r1, err := msg.Encode()
			if err != nil {
				lastErr = err
			}

			r2, err := test.Messages[i].Encode()
			if err != nil {
				lastErr = err
			}
			if !bytes.Equal(r2, r1) {
				t.Fail()
			}
		}
	}

	// check error
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}
