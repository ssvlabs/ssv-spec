package ssvmsg

import (
	comparable2 "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
	reflect2 "reflect"
	"testing"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type SSVMessageTest struct {
	Name               string
	MessageIDs         []types.MessageID
	BelongsToValidator bool
}

func (test *SSVMessageTest) TestName() string {
	return "ssvmessage " + test.Name
}

func (test *SSVMessageTest) Run(t *testing.T) {

	ks := testingutils.Testing4SharesSet()

	// For each message ID, test if MessageIDBelongs returns the expected value
	for _, msgID := range test.MessageIDs {
		belongs := testingutils.TestingShare(ks).ValidatorPubKey.MessageIDBelongs(msgID)
		require.Equal(t, test.BelongsToValidator, belongs)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}
