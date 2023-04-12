package testingutilscomparable

import (
	"testing"
)

func TestNestedCompare(t *testing.T) {
	d := make([]interface{}, 0)

	d = append(d, "1")
	d = append(d, []string{})

	d1 := Sprintf("diff 1", "---")
	d2 := NestedCompare("nested 2", []Difference{
		Sprintf("diff 2", "---"),
		d1,
	})
	d3 := NestedCompare("nested 3", []Difference{
		Sprintf("diff 3", "---"),
		d2,
		Sprintf("diff 3 end", "---"),
	})

	Print([]Difference{Sprintf("start", "---"), d3, Sprintf("end", "---")})
}
