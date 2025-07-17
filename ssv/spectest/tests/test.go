package tests

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
	GetPostState() (interface{}, error)
}

type TestF func() SpecTest

type PrivateKeyInfo struct {
	ValidatorSK  string
	Shares       map[types.OperatorID]string
	OperatorKeys map[types.OperatorID]string
}
