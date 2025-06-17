package maxmsgsize

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/utils"
	"github.com/stretchr/testify/require"
)

type StructureSizeTest struct {
	Name                  string
	Object                types.Encoder
	ExpectedEncodedLength int
	IsMaxSize             bool
}

func (test *StructureSizeTest) TestName() string {
	return "structure size test " + test.Name
}

func (test *StructureSizeTest) Run(t *testing.T) {
	if test.Object == nil {
		t.Fatal("no input object")
	}

	// Check if object respects SSZ tags sizes
	checkSSZTags(t, getReflectValueForObject(test.Object), test.IsMaxSize)

	// Check expected size
	encodedObject, err := test.Object.Encode()
	require.NoError(t, err)
	require.Equal(t, test.ExpectedEncodedLength, len(encodedObject))
}

// Custom JSON unmarshaller for StructureSizeTest since json can't unmarshal the types.Encoder interface
func (t *StructureSizeTest) UnmarshalJSON(data []byte) error {
	// Define alias with a decodable Object field
	type Alias struct {
		Name                  string
		ExpectedEncodedLength int
		IsMaxSize             bool
		Object                interface{}
	}

	// Unmarshal alias
	aliasObj := &Alias{}
	if err := json.Unmarshal(data, &aliasObj); err != nil {
		return err
	}
	t.Name = aliasObj.Name
	t.ExpectedEncodedLength = aliasObj.ExpectedEncodedLength
	t.IsMaxSize = aliasObj.IsMaxSize

	// Treat Object field with appropriate decoder
	byts, err := json.Marshal(aliasObj.Object)
	if err != nil {
		return err
	}

	// First try to determine the type from the JSON
	var objMap map[string]interface{}
	if err := json.Unmarshal(byts, &objMap); err != nil {
		return err
	}

	// Try to determine the type based on the fields present
	var correctType types.Encoder
	switch {
	case objMap["MsgType"] != nil && objMap["MsgID"] != nil:
		correctType = &types.SSVMessage{}
	case objMap["Signatures"] != nil && objMap["OperatorIDs"] != nil:
		correctType = &types.SignedSSVMessage{}
	case objMap["SigningRoot"] != nil && objMap["Signer"] != nil && objMap["ValidatorIndex"] != nil:
		correctType = &types.PartialSignatureMessage{}
	case objMap["Type"] != nil && objMap["Slot"] != nil && objMap["Messages"] != nil:
		correctType = &types.PartialSignatureMessages{}
	case objMap["Round"] != nil && objMap["Height"] != nil:
		correctType = &qbft.Message{}
	case objMap["Duty"] != nil && objMap["DataSSZ"] != nil:
		correctType = &types.ValidatorConsensusData{}
	case objMap["BlockRoot"] != nil && objMap["Source"] != nil && objMap["Target"] != nil:
		correctType = &types.BeaconVote{}
	default:
		return fmt.Errorf("could not determine object type from JSON")
	}

	// Try to unmarshal with hex handling
	if err := utils.UnmarshalJSONWithHex(byts, correctType); err != nil {
		return err
	}
	t.Object = correctType
	return nil
}
