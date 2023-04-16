package main

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/bloxapp/ssv-spec/qbft/spectest"
)

//go:generate go run main.go

func main() {
	all := map[string]tests.SpecTest{}
	for _, testF := range spectest.AllTests {
		test := testF()
		n := reflect.TypeOf(test).String() + "_" + test.TestName()
		if all[n] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", n))
		}
		all[n] = test
	}

	byts, err := json.Marshal(all)
	if err != nil {
		panic(err.Error())
	}

	if len(all) != len(spectest.AllTests) {
		panic("did not generate all tests\n")
	}

	fmt.Printf("found %d tests\n", len(all))
	writeJson(byts)
}

func writeJson(data []byte) {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")

	// try to create directory if it doesn't exist
	_ = os.Mkdir(basedir, os.ModeDir)

	file := filepath.Join(basedir, "tests.json")
	fmt.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0644); err != nil {
		panic(err.Error())
	}
}
