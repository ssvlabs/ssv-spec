package tests

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type MultiMsgProcessingSpecTest struct {
	Name        string
	Tests       []*MsgProcessingSpecTest
	PrivateKeys *PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (tests *MultiMsgProcessingSpecTest) TestName() string {
	return tests.Name
}

func (tests *MultiMsgProcessingSpecTest) Run(t *testing.T) {
	tests.overrideStateComparison(t)

	for _, test := range tests.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.RunAsPartOfMultiTest(t)
		})
	}
}

// overrideStateComparison overrides the post state comparison for all tests in the multi test
func (tests *MultiMsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	testsName := strings.ReplaceAll(tests.TestName(), " ", "_")
	for _, test := range tests.Tests {
		path := filepath.Join(testsName, test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(tests).String())
	}
}

func (tests *MultiMsgProcessingSpecTest) GetPostState() (interface{}, error) {
	ret := make(map[string]types.Root, len(tests.Tests))
	for _, test := range tests.Tests {
		_, _, err := test.runPreTesting()
		if err != nil && test.ExpectedError != err.Error() {
			return nil, err
		}
		ret[test.Name] = test.Runner
	}
	return ret, nil
}

func (tests *MultiMsgProcessingSpecTest) SetPrivateKeys(ks *testingutils.TestKeySet) {
	privateKeyInfo := &PrivateKeyInfo{
		ValidatorSK:  hex.EncodeToString(ks.ValidatorSK.Serialize()),
		Shares:       make(map[types.OperatorID]string),
		OperatorKeys: make(map[types.OperatorID]string),
	}

	// Add share keys
	for operatorID, shareSK := range ks.Shares {
		privateKeyInfo.Shares[operatorID] = hex.EncodeToString(shareSK.Serialize())
	}

	// Add operator keys (RSA private keys used for signing)
	for operatorID, operatorKey := range ks.OperatorKeys {
		privateKeyInfo.OperatorKeys[operatorID] = fmt.Sprintf("N:%s,E:%d",
			operatorKey.N.String(), operatorKey.E)
	}

	tests.PrivateKeys = privateKeyInfo
}
