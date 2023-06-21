package tests

import "testing"

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
	GetPostState() (interface{}, error)
}

type TestF func() SpecTest
