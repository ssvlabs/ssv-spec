package comparable

import (
	"encoding/json"
	"fmt"
)

func PrintDiff(source, target interface{}) string {
	byts1, _ := json.Marshal(source)
	byts2, _ := json.Marshal(target)
	return fmt.Sprintf("\n\n############ Struct Diff Report ############\n\n"+
		"Instructions:\n"+
		"   1) The below is a json dump of the compared objects\n"+
		"   2) They should perfectly match, but they don't\n"+
		"   3) Go to https://www.jsondiff.com/ to compare between the structs\n"+
		"   4) Paste the 'Source' json on the left\n"+
		"   5) Paste the 'Target' json on the right\n"+
		"   6) Hit 'Compare'\n"+
		"   7) Look for differences, source should match the target. If it doesn't find where it doesn't and fix it!\n\n"+
		"   Source: \n%s\n"+
		"   Target: \n%s\n"+
		"\n############################################\n\n", string(byts1), string(byts2))
}
