package msgcontainer

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type SpecTest struct {
	Name                       string
	MsgsToAdd                  []*types.SignedPartialSignatureMessage
	PostMsgCount               int
	PostReconstructedSignature []string
	ExpectedErr                string
}

func (test *SpecTest) TestName() string {
	return test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	c := ssv.PartialSignatureContainer{}
	// add msgs
	for _, msg := range test.MsgsToAdd {
		c[msg.Signer] = msg
	}
	require.Len(t, c, test.PostMsgCount)

	msg := test.MsgsToAdd[0]
	// test signatures per signer match
	for _, msg := range msg.Message.Messages {
		require.Len(t, c.SignatureForRoot(msg.SigningRoot), test.PostMsgCount)
	}

	var lastErr error
	var sig []byte
	if len(test.PostReconstructedSignature) > 0 {
		for i, m := range msg.Message.Messages {
			sig, lastErr = c.ReconstructSignature(m.SigningRoot, testingutils.TestingValidatorPubKey[:])
			require.EqualValues(t, test.PostReconstructedSignature[i], hex.EncodeToString(sig))
		}
		if len(test.ExpectedErr) == 0 {
			require.NoError(t, lastErr)
		} else {
			require.EqualError(t, lastErr, test.ExpectedErr)
		}
	}
}
