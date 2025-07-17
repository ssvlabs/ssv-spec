package tests

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
)

type TestF func() SpecTest

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
	GetPostState() (interface{}, error)
}

// PrivateKeyInfo contains the private keys used for signature generation
type PrivateKeyInfo struct {
	ValidatorSK  string                      `json:"ValidatorSK"`  // BLS secret key
	Shares       map[types.OperatorID]string `json:"Shares"`       // BLS share keys
	OperatorKeys map[types.OperatorID]string `json:"OperatorKeys"` // RSA private keys
}
