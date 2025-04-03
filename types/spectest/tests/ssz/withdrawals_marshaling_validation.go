package ssz

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SSZWithdrawalsMarshaling tests a valid capella withdrawals marshaling
func SSZWithdrawalsMarshaling() *SSZSpecTest {
	withdrawals := types.SSZWithdrawals(testingutils.TestingBeaconBlockCapella.Body.ExecutionPayload.Withdrawals)

	root, err := withdrawals.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &SSZSpecTest{
		Name:         "ssz withdrawals marshalling",
		Data:         testingutils.TestProposerConsensusDataBytsV(spec.DataVersionCapella),
		ExpectedRoot: root,
	}
}
