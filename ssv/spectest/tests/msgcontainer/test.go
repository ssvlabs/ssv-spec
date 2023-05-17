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
		require.Len(t, c.SignatureForRoot(msg.SigningRoot), len(test.MsgsToAdd))
	}

	if len(test.PostReconstructedSignature) > 0 {
		for i, m := range msg.Message.Messages {
			sig, err := c.ReconstructSignature(m.SigningRoot, testingutils.TestingValidatorPubKey[:])
			require.NoError(t, err)
			require.EqualValues(t, test.PostReconstructedSignature[i], hex.EncodeToString(sig))
		}
	}
}
