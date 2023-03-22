package ssz

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SSZWithdrawalsMarshaling tests a valid capella withdrawals marshaling
func SSZWithdrawalsMarshaling() *SSZSpecTest {

	withdrawals := types.SSZWithdrawals(testingutils.Withdrawals)

	root, err := withdrawals.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &SSZSpecTest{
		Name:         "ssz withdrawals marshalling",
		Data:         testingutils.TestProposerConsensusDataByts,
		ExpectedRoot: root,
	}
}
