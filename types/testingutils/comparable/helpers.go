package comparable

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"encoding/base64"
	"encoding/hex"

	spec2 "github.com/attestantio/go-eth2-client/spec"
	ssz "github.com/ferranbt/fastssz"
	"github.com/google/go-cmp/cmp"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/utils"
	"github.com/stretchr/testify/require"
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
func UnmarshalStateComparison[T types.Root](basedir string, testName string, testType string, targetState T) (T, error) {
	var nilT T
	basedir = filepath.Join(basedir, "generate")
	scDir := GetSCDir(basedir, testType)
	path := filepath.Join(scDir, fmt.Sprintf("%s.json", testName))

	byteValue, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nilT, err
	}

	err = utils.UnmarshalJSONWithHex(byteValue, targetState)
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
	byteValue, err := os.ReadFile(filepath.Clean(path))
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

	// Convert hex strings to byte arrays for SigningRoot
	convertSigningRootToBytes(testMap)
	convertSigningRootToBytes(expectedTestMap)

	diff := cmp.Diff(testMap, expectedTestMap)
	if diff != "" {
		t.Errorf("%s inputs changed. %v", testName, diff)
	}
}

// convertSigningRootToBytes recursively converts hex strings to byte arrays for SigningRoot fields
func convertSigningRootToBytes(m map[string]interface{}) {
	for k, v := range m {
		switch vv := v.(type) {
		case map[string]interface{}:
			convertSigningRootToBytes(vv)
		case []interface{}:
			// Check if this is a Signatures array
			if k == "Signatures" {
				// Convert each signature to base64 string
				for i, item := range vv {
					if str, ok := item.(string); ok {
						// Remove 0x prefix if present
						hexStr := strings.TrimPrefix(str, "0x")
						bytes, err := hex.DecodeString(hexStr)
						if err == nil {
							// Convert to base64
							vv[i] = base64.StdEncoding.EncodeToString(bytes)
						}
					}
				}
			} else if k == "MessageIDs" {
				// Convert each MessageID to float64 array
				for i, item := range vv {
					if str, ok := item.(string); ok {
						// Remove 0x prefix if present
						hexStr := strings.TrimPrefix(str, "0x")
						bytes, err := hex.DecodeString(hexStr)
						if err == nil {
							// Convert bytes to array of interface{}
							anyArray := make([]interface{}, len(bytes))
							for j, b := range bytes {
								anyArray[j] = float64(b)
							}
							vv[i] = anyArray
						}
					}
				}
			} else {
				// For other arrays, recursively process each item
				for _, item := range vv {
					if m, ok := item.(map[string]interface{}); ok {
						convertSigningRootToBytes(m)
					}
				}
			}
		case string:
			if k == "SigningRoot" || k == "ExpectedBlkRoot" || k == "ExpectedCdRoot" || k == "ExpectedRoot" || k == "MsgID" || k == "CommitteeID" || k == "DomainType" || k == "ForkVersion" || k == "Value" {
				// Remove 0x prefix if present
				hexStr := vv
				hexStr = strings.TrimPrefix(hexStr, "0x")
				bytes, err := hex.DecodeString(hexStr)
				if err == nil {
					// Convert bytes to array of float64
					floatArray := make([]float64, len(bytes))
					for i, b := range bytes {
						floatArray[i] = float64(b)
					}
					// Convert []float64 to []any to match the expected type
					anyArray := make([]interface{}, len(floatArray))
					for i, f := range floatArray {
						anyArray[i] = f
					}
					m[k] = anyArray
				}
			} else if k == "SSVOperatorPubKey" || k == "PartialSignature" || k == "FullData" || k == "StateValue" {
				// Remove 0x prefix if present
				hexStr := vv
				hexStr = strings.TrimPrefix(hexStr, "0x")
				bytes, err := hex.DecodeString(hexStr)
				if err == nil {
					// Convert to base64
					m[k] = base64.StdEncoding.EncodeToString(bytes)
				}
			}
		}
	}
}
