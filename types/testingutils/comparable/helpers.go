package comparable

import (
	"encoding/json"
	"fmt"
	spec2 "github.com/attestantio/go-eth2-client/spec"
	ssz "github.com/ferranbt/fastssz"
	"github.com/google/go-cmp/cmp"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func NoErrorEncoding(obj ssz.Marshaler) []byte {
	ret, err := obj.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	return ret
}

// FixIssue178 fixes consensus data fields which are nil instead of empty slice
// If we change the fields in ssv_msgs.go it will break a lot of roots, we're slowly fixing them
// SHOULD BE REMOVED once all tests are fixes
// see https://github.com/ssvlabs/ssv-spec/issues/178
func FixIssue178(input *types.ValidatorConsensusData, version spec2.DataVersion) *types.ValidatorConsensusData {
	byts, err := input.Encode()
	if err != nil {
		panic(err.Error())
	}
	ret := &types.ValidatorConsensusData{}
	if err := ret.Decode(byts); err != nil {
		panic(err.Error())
	}
	ret.Version = version
	return ret
}

// UnmarshalStateComparison reads a json derived from 'testName' and unmarshals it into 'targetState'
func UnmarshalStateComparison[T types.Root](basedir string, testName string, testType string, targetState T) (T,
	error) {
	var nilT T
	basedir = filepath.Join(basedir, "generate")
	scDir := GetSCDir(basedir, testType)
	path := filepath.Join(scDir, fmt.Sprintf("%s.json", testName))
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nilT, err
	}

	err = json.Unmarshal(byteValue, targetState)
	if err != nil {
		return nilT, err
	}

	return targetState, nil
}

// readStateComparison reads a json derived from 'testName' and unmarshals it into a json object
func readStateComparison(basedir string, testName string, testType string) (map[string]interface{}, error) {
	basedir = filepath.Join(basedir, "generate")
	scDir := GetSCDir(basedir, testType)
	path := filepath.Join(scDir, fmt.Sprintf("%s.json", testName))
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetExpectedStateFromScFile(testName string, testType string) (map[string]interface{}, error) {
	basedir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	expectedState, err := readStateComparison(basedir, testName, testType)
	if err != nil {
		return nil, err
	}
	return expectedState, nil
}

// GetSCDir returns the path to the state comparison folder for the given test type
func GetSCDir(basedir string, testType string) string {
	testType = strings.NewReplacer(
		"*", "",
		".", "_").
		Replace(testType)
	return filepath.Join(basedir, "state_comparison", testType)
}

// CompareWithJson compares the given test with the expected state from the state comparison folder
func CompareWithJson(t *testing.T, test any, testName string, testType string) {
	// marshal test into json
	byts, err := json.Marshal(test)
	require.NoError(t, err)
	//unmarshal json into map
	var testMap map[string]interface{}
	err = json.Unmarshal(byts, &testMap)
	require.NoError(t, err)

	expectedTestMap, err := GetExpectedStateFromScFile(testName, testType)
	require.NoError(t, err)

	diff := cmp.Diff(testMap, expectedTestMap)
	if diff != "" {
		t.Errorf("%s inputs changed. %v", testName, diff)
	}
}
