package share

import (
	"testing"
)

type SpecTest struct {
	Name string
	Data []byte
}

func (test *SpecTest) TestName() string {
	return test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	panic("implement")
}
