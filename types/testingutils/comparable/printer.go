package comparable

import (
	"encoding/json"
	"fmt"
)

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
