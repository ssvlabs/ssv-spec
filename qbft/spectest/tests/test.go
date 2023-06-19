package tests

import (
	"github.com/bloxapp/ssv-spec/types"
	"testing"
)

type TestF func() SpecTest

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
	GetPostState() (types.Encoder, error)
}
