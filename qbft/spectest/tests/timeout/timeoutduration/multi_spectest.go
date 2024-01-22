package timeoutduration

import (
	"testing"
)

type MultiSpecTest struct {
	Name  string
	Tests []*TimeoutDurationTest
}

func (test *MultiSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

func (test *MultiSpecTest) TestName() string {
	return test.Name
}

func (test *MultiSpecTest) Run(t *testing.T) {
	for _, test := range test.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}
