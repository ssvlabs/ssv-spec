package testingutils

import (
	"github.com/ssvlabs/ssv-spec/types"
)

// Verifies a list of SignedSSVMessage using the operators list
func VerifyListOfSignedSSVMessages(msgs []*types.SignedSSVMessage, operators []*types.Operator) error {
	for _, msg := range msgs {
		err := types.Verify(msg, operators)
		if err != nil {
			return err
		}
	}
	return nil
}
