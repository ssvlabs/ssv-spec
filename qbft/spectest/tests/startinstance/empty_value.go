package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance empty value",
		testdoc.StartInstanceEmptyValueDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{},
			},
		},
		nil,
		types.QBFTValueInvalidErrorCode,
		nil,
		nil,
	)
}
