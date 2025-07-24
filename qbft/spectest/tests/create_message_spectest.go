package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
	"github.com/stretchr/testify/require"
)

const (
	CreateProposal    = "createProposal"
	CreatePrepare     = "CreatePrepare"
	CreateCommit      = "CreateCommit"
	CreateRoundChange = "CreateRoundChange"
)

type CreateMsgSpecTest struct {
	Name          string
	Type          string
	Documentation string
	// ISSUE 217: rename to root
	Value qbft.Value
	// ISSUE 217: rename to value
	StateValue                                       []byte
	Round                                            qbft.Round
	RoundChangeJustifications, PrepareJustifications []*types.SignedSSVMessage
	CreateType                                       string
	ExpectedRoot                                     string
	ExpectedState                                    types.Root `json:"-"` // Field is ignored by encoding/json"
	ExpectedError                                    string
	PrivateKeys                                      *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for CreateMsgSpecTest
// This is a workaround to handle the ExpectedRoot field which is a string as it conflicts with EncodingTest.ExpectedRoot which is a [32]byte
func (test *CreateMsgSpecTest) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a raw map to extract ExpectedRoot
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract ExpectedRoot as string before hex processing
	var expectedRoot string
	if er, ok := raw["ExpectedRoot"].(string); ok {
		expectedRoot = er
	}

	// Remove ExpectedRoot from raw map so it doesn't get hex processed
	delete(raw, "ExpectedRoot")

	// Marshal the remaining data back to JSON
	remainingData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	// Create a temporary struct without ExpectedRoot for hex processing
	type CreateMsgSpecTestWithoutExpectedRoot struct {
		Name                                             string
		Type                                             string
		Documentation                                    string
		Value                                            qbft.Value
		StateValue                                       []byte
		Round                                            qbft.Round
		RoundChangeJustifications, PrepareJustifications []*types.SignedSSVMessage
		CreateType                                       string
		ExpectedState                                    types.Root `json:"-"`
		ExpectedError                                    string
	}

	temp := &CreateMsgSpecTestWithoutExpectedRoot{}

	if err := json.Unmarshal(remainingData, temp); err != nil {
		return err
	}

	// Copy all fields from temp to test
	test.Name = temp.Name
	test.Type = temp.Type
	test.Documentation = temp.Documentation
	test.Value = temp.Value
	test.StateValue = temp.StateValue
	test.Round = temp.Round
	test.RoundChangeJustifications = temp.RoundChangeJustifications
	test.PrepareJustifications = temp.PrepareJustifications
	test.CreateType = temp.CreateType
	test.ExpectedState = temp.ExpectedState
	test.ExpectedError = temp.ExpectedError

	// Set ExpectedRoot as string
	test.ExpectedRoot = expectedRoot

	return nil
}

func (test *CreateMsgSpecTest) Run(t *testing.T) {
	var msg *types.SignedSSVMessage
	var err error
	switch test.CreateType {
	case CreateProposal:
		msg, err = test.createProposal()
	case CreatePrepare:
		msg, err = test.createPrepare()
	case CreateCommit:
		msg, err = test.createCommit()
	case CreateRoundChange:
		msg, err = test.createRoundChange()
	default:
		t.Fail()
	}

	if err != nil && len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
		return
	}
	require.NoError(t, err)

	if test.Round == qbft.NoRound {
		require.Fail(t, "qbft round is invalid")
	}

	r, err := msg.GetRoot()
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
		return
	}
	require.NoError(t, err)

	if test.ExpectedRoot != hex.EncodeToString(r[:]) {
		fmt.Printf("expected: %v\n", test.ExpectedRoot)
		fmt.Printf("actual: %v\n", hex.EncodeToString(r[:]))
		// diff := typescomparable.PrintDiff(test.ExpectedState, msg)
		require.Fail(t, "post state not equal", "")
	}
	require.EqualValues(t, test.ExpectedRoot, hex.EncodeToString(r[:]))

	// Validate message
	err = msg.Validate()
	require.NoError(t, err)

	qbftMsg := &qbft.Message{}
	err = qbftMsg.Decode(msg.SSVMessage.Data)
	require.NoError(t, err)

	err = qbftMsg.Validate()
	require.NoError(t, err)

	typescomparable.CompareWithJson(t, test, test.TestName(), reflect.TypeOf(test).String())
}

func (test *CreateMsgSpecTest) createCommit() (*types.SignedSSVMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &qbft.State{
		CommitteeMember: testingutils.TestingCommitteeMember(ks),
		ID:              testingutils.TestingIdentifier,
		Round:           test.Round,
	}
	signer := testingutils.TestingOperatorSigner(ks)

	return qbft.CreateCommit(state, signer, test.Value)
}

func (test *CreateMsgSpecTest) createPrepare() (*types.SignedSSVMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &qbft.State{
		CommitteeMember: testingutils.TestingCommitteeMember(ks),
		ID:              testingutils.TestingIdentifier,
		Round:           test.Round,
	}
	signer := testingutils.TestingOperatorSigner(ks)

	return qbft.CreatePrepare(state, signer, test.Round, test.Value)
}

func (test *CreateMsgSpecTest) createProposal() (*types.SignedSSVMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &qbft.State{
		CommitteeMember: testingutils.TestingCommitteeMember(ks),
		ID:              testingutils.TestingIdentifier,
		Round:           test.Round,
	}
	signer := testingutils.TestingOperatorSigner(ks)

	return qbft.CreateProposal(state, signer, test.Value[:], testingutils.ToProcessingMessages(test.
		RoundChangeJustifications), testingutils.ToProcessingMessages(test.PrepareJustifications))
}

func (test *CreateMsgSpecTest) createRoundChange() (*types.SignedSSVMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &qbft.State{
		CommitteeMember:  testingutils.TestingCommitteeMember(ks),
		ID:               testingutils.TestingIdentifier,
		PrepareContainer: qbft.NewMsgContainer(),
		Round:            test.Round,
	}
	signer := testingutils.TestingOperatorSigner(ks)

	if len(test.PrepareJustifications) > 0 {
		prepareMsg, err := qbft.DecodeMessage(test.PrepareJustifications[0].SSVMessage.Data)
		if err != nil {
			return nil, err
		}
		state.LastPreparedRound = prepareMsg.Round
		state.LastPreparedValue = test.StateValue

		for _, msg := range test.PrepareJustifications {
			_, err := state.PrepareContainer.AddFirstMsgForSignerAndRound(testingutils.ToProcessingMessage(msg))
			if err != nil {
				return nil, errors.Wrap(err, "could not add first message for signer")
			}
		}
	}

	return qbft.CreateRoundChange(state, signer, qbft.FirstRound, test.Value[:])
}

func (test *CreateMsgSpecTest) TestName() string {
	return "qbft create message " + test.Name
}

func (test *CreateMsgSpecTest) GetPostState() (interface{}, error) {
	// remove private keys
	test.PrivateKeys = nil

	return test, nil
}

func NewCreateMsgSpecTest(name string, documentation string, value [32]byte, stateValue []byte, round qbft.Round, roundChangeJustifications []*types.SignedSSVMessage, prepareJustifications []*types.SignedSSVMessage, createType string, expectedRoot string, expectedState types.Root, expectedError string, ks *testingutils.TestKeySet) *CreateMsgSpecTest {
	return &CreateMsgSpecTest{
		Name:                      name,
		Type:                      testdoc.CreateMsgSpecTestType,
		Documentation:             documentation,
		Value:                     value,
		StateValue:                stateValue,
		Round:                     round,
		RoundChangeJustifications: roundChangeJustifications,
		PrepareJustifications:     prepareJustifications,
		CreateType:                createType,
		ExpectedRoot:              expectedRoot,
		ExpectedState:             expectedState,
		ExpectedError:             expectedError,
		PrivateKeys:               testingutils.BuildPrivateKeyInfo(ks),
	}
}
