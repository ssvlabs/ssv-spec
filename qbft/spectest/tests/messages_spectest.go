package tests

import (
	"testing"

	"encoding/hex"
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgSpecTest tests encoding and decoding of a msg
type MsgSpecTest struct {
	Name            string
	Messages        []*types.SignedSSVMessage
	EncodedMessages [][]byte
	ExpectedRoots   [][32]byte
	ExpectedError   string
	PrivateKeys     *PrivateKeyInfo `json:"PrivateKeys,omitempty"`
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

// SetPrivateKeys populates the PrivateKeys field with keys from the given TestKeySet
func (test *MsgSpecTest) SetPrivateKeys(ks *testingutils.TestKeySet) {
	privateKeyInfo := &PrivateKeyInfo{
		ValidatorSK:  hex.EncodeToString(ks.ValidatorSK.Serialize()),
		Shares:       make(map[types.OperatorID]string),
		OperatorKeys: make(map[types.OperatorID]string),
	}

	// Add share keys
	for operatorID, shareSK := range ks.Shares {
		privateKeyInfo.Shares[operatorID] = hex.EncodeToString(shareSK.Serialize())
	}

	// Add operator keys (RSA private keys used for signing)
	for operatorID, operatorKey := range ks.OperatorKeys {
		// For RSA keys, we'll include the modulus and exponent
		privateKeyInfo.OperatorKeys[operatorID] = fmt.Sprintf("N:%s,E:%d",
			operatorKey.N.String(), operatorKey.E)
	}

	test.PrivateKeys = privateKeyInfo
}
