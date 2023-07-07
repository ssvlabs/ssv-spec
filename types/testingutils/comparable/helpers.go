package comparable

import (
	"encoding/json"
	"fmt"
	spec2 "github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
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
// see https://github.com/bloxapp/ssv-spec/issues/178
func FixIssue178(input *types.ConsensusData, version spec2.DataVersion) *types.ConsensusData {
	byts, err := input.Encode()
	if err != nil {
		panic(err.Error())
	}
	ret := &types.ConsensusData{}
	if err := ret.Decode(byts); err != nil {
		panic(err.Error())
	}
	ret.Version = version

	return ret
}

// UnmarshalSSVStateComparison reads a json derived from 'test' and unmarshals it into 'targetState'
func UnmarshalSSVStateComparison(testName string, testType string, targetState types.Root) (types.Root, error) {
	scDir, err := GetSCDir(testType)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(scDir, fmt.Sprintf("%s.json", testName))
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, targetState)
	if err != nil {
		return nil, err
	}

	return targetState, nil
}

// GetSCDir returns the path to the state comparison folder for the given test type
func GetSCDir(testType string) (string, error) {
	basedir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get working directory for SC dir")
	}
	basedir = filepath.Join(strings.TrimSuffix(basedir, "main.go"), "state_comparison", testType)
	scDir := strings.NewReplacer(
		"*", "",
		".", "_").
		Replace(basedir)
	return scDir, nil
}
