package comparable

import (
	"encoding/json"
	"fmt"
	spec2 "github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"os"
	"path/filepath"
	"reflect"
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

// UnmarshalSSVStateComparison unmarshals the json rperesenting the test's post state to the targetState struct
func UnmarshalSSVStateComparison(test tests.SpecTest, targetState types.Root) (types.Root, error) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate", "state_comparison", reflect.TypeOf(test).String(),
		fmt.Sprintf("%s.json", test.TestName()))
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
