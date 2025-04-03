package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/qbft/spectest"
)

//go:generate go run main.go

var testsDir = "tests"
var stateComparisonDir = "state_comparison"

func main() {
	clearStateComparisonFolder()
	clearTestsFolder()

	all := map[string]tests.SpecTest{}
	for _, testF := range spectest.AllTests {
		test := testF()
		n := reflect.TypeOf(test).String() + "_" + test.TestName()
		if all[n] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", n))
		}
		all[n] = test
	}

	log.Printf("found %d tests\n", len(all))
	if len(all) != len(spectest.AllTests) {
		log.Fatalf("did not generate all tests\n")
	}

	if err := os.MkdirAll(testsDir, 0700); err != nil && !os.IsExist(err) {
		panic(err.Error())
	}
	for name, test := range all {
		byts, err := json.MarshalIndent(test, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ReplaceAll(name, "*", "")
		name = filepath.Join(testsDir, name)
		writeJson(name, byts)
	}

	// write large tests.json file
	byts, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	writeJson(testsDir, byts)

	// write state comparison json files
	for _, testF := range spectest.AllTests {
		test := testF()
		// generate post state comparison
		post, err := test.GetPostState()
		if err != nil {
			panic(errors.Wrapf(err, "failed to get post state for test: %s", test.TestName()).Error())
		}
		writeJsonStateComparison(test.TestName(), reflect.TypeOf(test).String(), post)
	}
}

func clearStateComparisonFolder() {
	if err := os.RemoveAll(stateComparisonDir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(stateComparisonDir, 0700); err != nil {
		panic(err.Error())
	}
}

func clearTestsFolder() {
	if err := os.RemoveAll(testsDir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(testsDir, 0700); err != nil {
		panic(err.Error())
	}
}

func writeJsonStateComparison(name, testType string, post interface{}) {
	if post == nil { // If nil, test not supporting post state comparison yet
		log.Printf("skipping state comparison json, not supported: %s\n", name)
		return
	}
	log.Printf("writing state comparison json: %s\n", name)

	byts, err := json.MarshalIndent(post, "", "		")
	if err != nil {
		panic(err.Error())
	}
	scDir := scDir(testType)

	// try to create directory if it doesn't exist
	if err := os.MkdirAll(scDir, 0700); err != nil && !os.IsExist(err) {
		panic(err.Error())
	}

	file := filepath.Join(scDir, fmt.Sprintf("%s.json", name))
	log.Printf("writing state comparison json: %s\n", file)
	if err := os.WriteFile(file, byts, 0664); err != nil {
		panic(err.Error())
	}
}

func scDir(testType string) string {
	return comparable2.GetSCDir(".", testType)
}

func writeJson(name string, data []byte) {
	file := name + ".json"
	log.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0664); err != nil {
		panic(err.Error())
	}
}
