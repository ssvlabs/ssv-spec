package ssz

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SSZWithdrawalsMarshaling tests a valid capella withdrawals marshaling
func SSZWithdrawalsMarshaling() *SSZSpecTest {
	withdrawals := types.SSZWithdrawals(testingutils.TestingBeaconBlockCapella.Body.ExecutionPayload.Withdrawals)

	root, err := withdrawals.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewSSZSpecTest(
		"ssz withdrawals marshalling",
		testdoc.SSZSpecTestWithdrawalsMarshalingDoc,
		testingutils.TestProposerConsensusDataBytsV(spec.DataVersionCapella),
		root,
		0,
	)
}
