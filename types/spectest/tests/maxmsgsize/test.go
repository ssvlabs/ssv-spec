package maxmsgsize

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
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

		Object interface{}
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

	var getDecoder = func() *json.Decoder {
		decoder := json.NewDecoder(strings.NewReader(string(byts)))
		decoder.DisallowUnknownFields()
		return decoder
	}

	var possibleObjects = []types.Encoder{
		&qbft.Message{},
		&types.PartialSignatureMessage{},
		&types.PartialSignatureMessages{},
		&types.SSVMessage{},
		&types.SignedSSVMessage{},
		&types.ValidatorConsensusData{},
		&types.BeaconVote{},
	}

	for _, obj := range possibleObjects {
		err := getDecoder().Decode(&obj)
		if err == nil {
			t.Object = obj
			return nil
		}
	}

	panic("unknown type")
}
