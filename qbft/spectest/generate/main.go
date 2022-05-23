package main

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft/spectest"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

func main() {
	all := map[string]spectest.SpecTest{}
	for _, t := range spectest.AllTests {
		all[reflect.TypeOf(t).String()+"_"+t.TestName()] = t
	}

	byts, _ := json.Marshal(all)
	fmt.Printf("found %d tests\n", len(all))
	writeJson(byts)
}

func writeJson(data []byte) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "qbft", "spectest", "generate")
	fileName := "tests.json"
	fullPath := path + "/" + fileName

	fmt.Printf("writing spec tests json to: %s\n", fullPath)
	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
		panic(err.Error())
	}
}
