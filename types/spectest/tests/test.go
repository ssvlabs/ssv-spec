package tests

import "testing"

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}
