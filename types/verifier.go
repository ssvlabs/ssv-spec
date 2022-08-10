package types

import "github.com/ethereum/go-ethereum/common"

type DKGVerifier interface {
	VerifyDKGOutput(signature Signature, output Root, address common.Address) bool
}
