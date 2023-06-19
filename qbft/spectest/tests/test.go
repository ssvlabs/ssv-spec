package tests

import (
	"testing"
)

type TestF func() SpecTest

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
	GetPostState() (interface{}, error)
}
