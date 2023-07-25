package messages

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgSpecTest struct {
	Name            string
	Messages        []*types.SignedPartialSignatureMessage
	EncodedMessages [][]byte
	ExpectedRoots   [][32]byte
	ExpectedError   string
}

func (test *MsgSpecTest) TestName() string {
	return "msg " + test.Name
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error

	for i, msg := range test.Messages {
		// test validation
		if err := msg.Validate(); err != nil {
			lastErr = err
		}

		// test encoding decoding
		if len(test.EncodedMessages) > 0 {
			byts, err := msg.Encode()
			require.NoError(t, err)
			require.EqualValues(t, test.EncodedMessages[i], byts)

			decoded := &types.SignedPartialSignatureMessage{}
			require.NoError(t, decoded.Decode(byts))
			decodedRoot, err := decoded.GetRoot()
			require.NoError(t, err)
			r, err := msg.GetRoot()
			require.NoError(t, err)
			require.EqualValues(t, r, decodedRoot)
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

func (tests *MsgSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}
