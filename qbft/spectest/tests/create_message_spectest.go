package tests

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
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

type TestSigner struct {
	OperatorID types.OperatorID
	OperatorSK *rsa.PrivateKey
}

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

	// consts for CreateMsgSpecTest
	CommitteeMember *types.CommitteeMember `json:"CommitteeMember,omitempty"`
	Identifier      []byte                 `json:"Identifier,omitempty"`
	OperatorID      types.OperatorID       `json:"OperatorID,omitempty"`
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
	
	// Print the actual root for test updates
	fmt.Printf("Test: %s, Actual root: %x\n", test.TestName(), r)
	
	// Create output string to write to both stdout and file
	var output string
	
	// Debug: Print all message fields for comparison
	output += fmt.Sprintf("\n========== Test: %s ==========\n", test.Name)
	output += fmt.Sprintf("Root: %s\n", hex.EncodeToString(r[:]))
	output += fmt.Sprintf("ExpectedRoot: %s\n", test.ExpectedRoot)
	output += fmt.Sprintf("Round: %d\n", test.Round)
	output += fmt.Sprintf("Value (hash): %s\n", hex.EncodeToString(test.Value[:]))
	output += fmt.Sprintf("StateValue (full data): %s\n", hex.EncodeToString(test.StateValue))
	
	// Print the entire message as hex
	msgBytes, _ := msg.MarshalSSZ()
	output += fmt.Sprintf("Full SignedSSVMessage (hex): %s\n", hex.EncodeToString(msgBytes))
	
	// Print SSVMessage fields
	output += fmt.Sprintf("SSVMessage:\n")
	output += fmt.Sprintf("  MsgType: %d\n", msg.SSVMessage.MsgType)
	output += fmt.Sprintf("  MsgID: %s\n", hex.EncodeToString(msg.SSVMessage.MsgID[:]))
	output += fmt.Sprintf("  Data: %s\n", hex.EncodeToString(msg.SSVMessage.Data))
	
	// Decode and print QBFT message
	qbftMsg := &qbft.Message{}
	if err := qbftMsg.Decode(msg.SSVMessage.Data); err == nil {
		output += fmt.Sprintf("QBFTMessage:\n")
		output += fmt.Sprintf("  MsgType: %d\n", qbftMsg.MsgType)
		output += fmt.Sprintf("  Height: %d\n", qbftMsg.Height)
		output += fmt.Sprintf("  Round: %d\n", qbftMsg.Round)
		output += fmt.Sprintf("  Identifier: %s\n", hex.EncodeToString(qbftMsg.Identifier))
		output += fmt.Sprintf("  Root: %s\n", hex.EncodeToString(qbftMsg.Root[:]))
		output += fmt.Sprintf("  DataRound: %d\n", qbftMsg.DataRound)
		
		// Justifications
		if len(qbftMsg.RoundChangeJustification) > 0 {
			output += fmt.Sprintf("  RoundChangeJustifications: %d messages\n", len(qbftMsg.RoundChangeJustification))
			for i, rcJust := range qbftMsg.RoundChangeJustification {
				output += fmt.Sprintf("    RC[%d] (hex): %s\n", i, hex.EncodeToString(rcJust))
				// Try to decode the justification message
				rcMsg := &types.SignedSSVMessage{}
				if err := rcMsg.UnmarshalSSZ(rcJust); err == nil {
					output += fmt.Sprintf("    RC[%d] Decoded SignedSSVMessage:\n", i)
					output += fmt.Sprintf("      OperatorIDs: %v\n", rcMsg.OperatorIDs)
					output += fmt.Sprintf("      Signatures: %d signatures\n", len(rcMsg.Signatures))
					for j, sig := range rcMsg.Signatures {
						output += fmt.Sprintf("        Signature[%d]: %s\n", j, hex.EncodeToString(sig))
					}
					output += fmt.Sprintf("      SSVMessage.MsgType: %d\n", rcMsg.SSVMessage.MsgType)
					output += fmt.Sprintf("      SSVMessage.MsgID: %s\n", hex.EncodeToString(rcMsg.SSVMessage.MsgID[:]))
					output += fmt.Sprintf("      SSVMessage.Data: %s\n", hex.EncodeToString(rcMsg.SSVMessage.Data))
					if rcMsg.FullData != nil && len(rcMsg.FullData) > 0 {
						output += fmt.Sprintf("      FullData: %s\n", hex.EncodeToString(rcMsg.FullData))
					}
					
					// Decode the inner QBFT message
					innerQbft := &qbft.Message{}
					if err := innerQbft.Decode(rcMsg.SSVMessage.Data); err == nil {
						output += fmt.Sprintf("      Inner QBFTMessage:\n")
						output += fmt.Sprintf("        MsgType: %d\n", innerQbft.MsgType)
						output += fmt.Sprintf("        Height: %d\n", innerQbft.Height)
						output += fmt.Sprintf("        Round: %d\n", innerQbft.Round)
						output += fmt.Sprintf("        Identifier: %s\n", hex.EncodeToString(innerQbft.Identifier))
						output += fmt.Sprintf("        Root: %s\n", hex.EncodeToString(innerQbft.Root[:]))
						output += fmt.Sprintf("        DataRound: %d\n", innerQbft.DataRound)
						if len(innerQbft.RoundChangeJustification) > 0 {
							output += fmt.Sprintf("        Has %d RoundChangeJustifications\n", len(innerQbft.RoundChangeJustification))
						}
						if len(innerQbft.PrepareJustification) > 0 {
							output += fmt.Sprintf("        Has %d PrepareJustifications\n", len(innerQbft.PrepareJustification))
						}
					}
				}
			}
		}
		if len(qbftMsg.PrepareJustification) > 0 {
			output += fmt.Sprintf("  PrepareJustifications: %d messages\n", len(qbftMsg.PrepareJustification))
			for i, prepJust := range qbftMsg.PrepareJustification {
				output += fmt.Sprintf("    Prepare[%d] (hex): %s\n", i, hex.EncodeToString(prepJust))
				// Try to decode the justification message
				prepMsg := &types.SignedSSVMessage{}
				if err := prepMsg.UnmarshalSSZ(prepJust); err == nil {
					output += fmt.Sprintf("    Prepare[%d] Decoded SignedSSVMessage:\n", i)
					output += fmt.Sprintf("      OperatorIDs: %v\n", prepMsg.OperatorIDs)
					output += fmt.Sprintf("      Signatures: %d signatures\n", len(prepMsg.Signatures))
					for j, sig := range prepMsg.Signatures {
						output += fmt.Sprintf("        Signature[%d]: %s\n", j, hex.EncodeToString(sig))
					}
					output += fmt.Sprintf("      SSVMessage.MsgType: %d\n", prepMsg.SSVMessage.MsgType)
					output += fmt.Sprintf("      SSVMessage.MsgID: %s\n", hex.EncodeToString(prepMsg.SSVMessage.MsgID[:]))
					output += fmt.Sprintf("      SSVMessage.Data: %s\n", hex.EncodeToString(prepMsg.SSVMessage.Data))
					if prepMsg.FullData != nil && len(prepMsg.FullData) > 0 {
						output += fmt.Sprintf("      FullData: %s\n", hex.EncodeToString(prepMsg.FullData))
					}
					
					// Decode the inner QBFT message
					innerQbft := &qbft.Message{}
					if err := innerQbft.Decode(prepMsg.SSVMessage.Data); err == nil {
						output += fmt.Sprintf("      Inner QBFTMessage:\n")
						output += fmt.Sprintf("        MsgType: %d\n", innerQbft.MsgType)
						output += fmt.Sprintf("        Height: %d\n", innerQbft.Height)
						output += fmt.Sprintf("        Round: %d\n", innerQbft.Round)
						output += fmt.Sprintf("        Identifier: %s\n", hex.EncodeToString(innerQbft.Identifier))
						output += fmt.Sprintf("        Root: %s\n", hex.EncodeToString(innerQbft.Root[:]))
						output += fmt.Sprintf("        DataRound: %d\n", innerQbft.DataRound)
						if len(innerQbft.RoundChangeJustification) > 0 {
							output += fmt.Sprintf("        Has %d RoundChangeJustifications\n", len(innerQbft.RoundChangeJustification))
						}
						if len(innerQbft.PrepareJustification) > 0 {
							output += fmt.Sprintf("        Has %d PrepareJustifications\n", len(innerQbft.PrepareJustification))
						}
					}
				}
			}
		}
	} else {
		output += fmt.Sprintf("Failed to decode QBFT message: %v\n", err)
	}
	
	// Print signatures
	output += fmt.Sprintf("Signatures:\n")
	for i, sig := range msg.Signatures {
		output += fmt.Sprintf("  Signature[%d]: %s\n", i, hex.EncodeToString(sig))
	}
	
	output += fmt.Sprintf("OperatorIDs: %v\n", msg.OperatorIDs)
	
	if msg.FullData != nil {
		output += fmt.Sprintf("FullData: %s\n", hex.EncodeToString(msg.FullData))
	}
	
	// Print to stdout
	fmt.Print(output)
	
	// Also append to file
	file, err := os.OpenFile("/tmp/go_messages_full.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer file.Close()
		file.WriteString(output)
	}

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

	qbftMsgValidate := &qbft.Message{}
	err = qbftMsgValidate.Decode(msg.SSVMessage.Data)
	require.NoError(t, err)

	err = qbftMsgValidate.Validate()
	require.NoError(t, err)

	// remove consts for state comparison
	test.CommitteeMember = nil
	test.Identifier = nil
	test.OperatorID = 0

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

	// Use StateValue (full data) instead of Value (which is already a hash)
	return qbft.CreateProposal(state, signer, test.StateValue, testingutils.ToProcessingMessages(test.
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

	// Use StateValue (full data) instead of Value (which is already a hash)
	return qbft.CreateRoundChange(state, signer, qbft.FirstRound, test.StateValue)
}

func (test *CreateMsgSpecTest) TestName() string {
	return "qbft create message " + test.Name
}

func (test *CreateMsgSpecTest) GetPostState() (interface{}, error) {
	test.CommitteeMember = nil
	test.Identifier = nil
	test.OperatorID = 0
	test.PrivateKeys = nil

	return test, nil
}

func NewCreateMsgSpecTest(name string, documentation string, value [32]byte, stateValue []byte, round qbft.Round, roundChangeJustifications []*types.SignedSSVMessage, prepareJustifications []*types.SignedSSVMessage, createType string, expectedRoot string, expectedState types.Root, expectedError string, ks *testingutils.TestKeySet) *CreateMsgSpecTest {
	committeeMember := &types.CommitteeMember{}
	operatorID := types.OperatorID(0)
	if ks != nil {
		committeeMember = testingutils.TestingCommitteeMember(ks)
		operatorID = testingutils.TestingOperatorSigner(ks).OperatorID
	}

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

		// consts for CreateMsgSpecTest
		CommitteeMember: committeeMember,
		Identifier:      testingutils.TestingIdentifier,
		OperatorID:      operatorID,
	}
}
