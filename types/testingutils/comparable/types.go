package testingutilscomparable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
)

var nestingStart = "{"
var nestingEnd = "}"
var noTab = ""
var singleTab = "    "

type Difference map[string]interface{}

func (diff Difference) Empty() bool {
	return len(diff) == 0
}

func Sprintf(key, format string, params ...interface{}) Difference {
	return map[string]interface{}{
		key: fmt.Sprintf(format, params...),
	}
}

func Print(differences []Difference) {
	diff := NestedCompare("Compare Differences", differences)

	byts, err := json.Marshal(diff)
	if err != nil {
		panic(err.Error())
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, byts, "", "    "); err != nil {
		panic(err.Error())
	}
	fmt.Printf("%s\n", prettyJSON.String())
}

func NestedCompare(prefix string, differences []Difference) Difference {
	if len(differences) > 0 {
		nested := make(map[string]interface{})
		for _, diff := range differences {
			if len(nested) == 0 {
				nested = diff
			} else {
				maps.Copy(nested, diff)
			}
		}
		return map[string]interface{}{
			prefix: nested,
		}
	}
	return map[string]interface{}{}
}

func PrintDiff(source, target interface{}) {
	fmt.Printf("\n\n############ Struct Diff Report ############\n\n" +
		"Instructions:\n" +
		"   1) The below is a json dump of the compared objects\n" +
		"   2) They should perfectly match, but they don't\n" +
		"   3) Go to https://www.jsondiff.com/ to compare between the structs\n" +
		"   4) Paste the 'Source' json on the left\n" +
		"   5) Paste the 'Target' json on the right\n" +
		"   6) Hit 'Compare'\n" +
		"   7) Look for differences, source should match the target. If it doesn't find where it doesn't and fix it!\n\n")
	byts, _ := json.Marshal(source)
	fmt.Printf("   Source: \n"+
		"      %s\n", string(byts))
	byts, _ = json.Marshal(target)
	fmt.Printf("   Target: \n"+
		"      %s\n", string(byts))
	fmt.Printf("\n############################################\n\n")
}
