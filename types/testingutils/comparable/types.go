package testingutilscomparable

import "fmt"

type Difference string

func Sprintf(format string, params ...interface{}) Difference {
	return Difference(fmt.Sprintf(format, params...))
}

var singleTab = "    "
var doubleTab = singleTab + singleTab

func Print(differences []Difference) {
	fmt.Printf("Compare Differences:\n")
	for _, diff := range differences {
		fmt.Printf("%s%s\n", singleTab, diff)
	}
}

func NestedCompare(prefix string, differences []Difference) Difference {
	if len(differences) > 0 {
		ret := fmt.Sprintf("%s", prefix)
		for _, diff := range differences {
			ret += fmt.Sprintf("\n%s%s", doubleTab, diff)
		}
		return Difference(ret)
	}
	return ""
}
