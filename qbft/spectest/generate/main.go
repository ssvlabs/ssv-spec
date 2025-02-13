package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/qbft/spectest"
)

//go:generate go run main.go

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

	// write small json files for each test
	// try to create directory if it doesn't exist
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("no caller info")
	}
	testsDir := filepath.Join(strings.TrimSuffix(basedir, "main.go"), "tests")
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
		name = "tests/" + name
		writeJson(name, byts)
	}

	// write large tests.json file
	byts, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	writeJson("tests", byts)

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
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	dir := filepath.Join(strings.TrimSuffix(basedir, "main.go"), "state_comparison")

	if err := os.RemoveAll(dir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(dir, 0700); err != nil {
		panic(err.Error())
	}
}

func clearTestsFolder() {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	dir := filepath.Join(strings.TrimSuffix(basedir, "main.go"), "tests")

	if err := os.RemoveAll(dir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(dir, 0700); err != nil {
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
	if err := os.WriteFile(file, byts, 0400); err != nil {
		panic(err.Error())
	}
}

func scDir(testType string) string {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")
	scDir := comparable2.GetSCDir(basedir, testType)
	return scDir
}

func writeJson(name string, data []byte) {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")

	// try to create directory if it doesn't exist
	_ = os.Mkdir(basedir, os.ModeDir)

	file := filepath.Join(basedir, name+".json")
	log.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0400); err != nil {
		panic(err.Error())
	}
}
