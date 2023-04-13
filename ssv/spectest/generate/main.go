package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/bloxapp/ssv-spec/ssv/spectest"
)

//go:generate go run main.go

func main() {
	all := map[string]spectest.SpecTest{}
	for _, t := range spectest.AllTests {
		n := reflect.TypeOf(t).String() + "_" + t.TestName()
		if all[n] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", n))
		}
		all[n] = t
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
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "ssv", "spectest", "generate")

	// try to create directory if it doesn't exist
	_ = os.Mkdir(path, os.ModeDir)

	file := filepath.Join(path, "tests.json")
	fmt.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0644); err != nil {
		panic(err.Error())
	}
}
