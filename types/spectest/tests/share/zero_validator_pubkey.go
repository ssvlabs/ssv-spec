package share

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ZeroValidatorPubKey tests a share with an all-zero validator public key
func ZeroValidatorPubKey() *ValidationTest {
	share := testingutils.TestingShare(testingutils.Testing4SharesSet(), testingutils.TestingValidatorIndex)
	share.ValidatorPubKey = types.ValidatorPK{}

	return NewValidationTest(
		"zero validator pubkey",
		testdoc.ShareValidationZeroValidatorPubKeyDoc,
		[]*types.Share{share},
		types.InvalidShareErrorCode,
	)
}
