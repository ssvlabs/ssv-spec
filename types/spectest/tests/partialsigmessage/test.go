package partialsigmessage

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type MsgSpecTest struct {
	Name              string
	Type              string
	Documentation     string
	Messages          []*types.PartialSignatureMessages
	EncodedMessages   [][]byte
	ExpectedRoots     [][32]byte
	ExpectedErrorCode int
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

			decoded := &types.PartialSignatureMessages{}
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
	tests.AssertErrorCode(t, test.ExpectedErrorCode, lastErr)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func (tests *MsgSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

func NewMsgSpecTest(name, documentation string, messages []*types.PartialSignatureMessages, encodedMessages [][]byte, expectedRoots [][32]byte, expectedErrorCode int) *MsgSpecTest {
	return &MsgSpecTest{
		Name:              name,
		Type:              testdoc.MsgSpecTestType,
		Documentation:     documentation,
		Messages:          messages,
		EncodedMessages:   encodedMessages,
		ExpectedRoots:     expectedRoots,
		ExpectedErrorCode: expectedErrorCode,
	}
}
