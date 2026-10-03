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
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/p2p/spectest"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
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
		name := reflect.TypeOf(test).String() + "_" + test.TestName()
		if all[name] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", name))
		}
		all[name] = test
	}

	log.Printf("found %d tests\n", len(all))
	if len(all) != len(spectest.AllTests) {
		log.Fatalf("did not generate all tests\n")
	}

	if err := os.MkdirAll(testsDir, 0o700); err != nil && !os.IsExist(err) {
		panic(err.Error())
	}
	for name, test := range all {
		byts, err := json.MarshalIndent(test, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		writeJSON(filepath.Join(testsDir, sanitize(name)), byts)
	}

	byts, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	writeJSON(testsDir, byts)

	for _, testF := range spectest.AllTests {
		test := testF()
		post, err := test.GetPostState()
		if err != nil {
			panic(errors.Wrapf(err, "failed to get post state for test: %s", test.TestName()).Error())
		}
		writeJSONStateComparison(test, post)
	}
}

func clearStateComparisonFolder() {
	if err := os.RemoveAll(stateComparisonDir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(stateComparisonDir, 0o700); err != nil {
		panic(err.Error())
	}
}

func clearTestsFolder() {
	if err := os.RemoveAll(testsDir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(testsDir, 0o700); err != nil {
		panic(err.Error())
	}
}

func writeJSONStateComparison(test tests.SpecTest, post interface{}) {
	if post == nil {
		log.Printf("skipping state comparison json, not supported: %s\n", test.TestName())
		return
	}

	byts, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	scDir := comparable2.GetSCDir(".", reflect.TypeOf(test).String())
	if err := os.MkdirAll(scDir, 0o700); err != nil && !os.IsExist(err) {
		panic(err.Error())
	}

	file := filepath.Join(scDir, fmt.Sprintf("%s.json", sanitize(test.TestName())))
	log.Printf("writing state comparison json: %s\n", file)
	if err := os.WriteFile(file, byts, 0o664); err != nil {
		panic(err.Error())
	}
}

func writeJSON(name string, data []byte) {
	file := name + ".json"
	log.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0o664); err != nil {
		panic(err.Error())
	}
}

func sanitize(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	return strings.ReplaceAll(name, "*", "")
}
