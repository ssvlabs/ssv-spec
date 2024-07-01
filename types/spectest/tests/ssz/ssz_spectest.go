package ssz

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
)

type SSZSpecTest struct {
	Name          string
	Data          []byte
	ExpectedRoot  [32]byte
	ExpectedError string
}

func (test *SSZSpecTest) TestName() string {
	return test.Name
}

func (test *SSZSpecTest) Run(t *testing.T) {
	cd := &types.ValidatorConsensusData{}
	require.NoError(t, cd.Decode(test.Data))

	vBlk, _, err := cd.GetBlockData()
	require.NoError(t, err)
	require.NotNil(t, vBlk)

	var withdrawals types.SSZWithdrawals
	switch vBlk.Version {
	case spec.DataVersionCapella:
		withdrawals = vBlk.Capella.Body.ExecutionPayload.Withdrawals
	}

	root, err := withdrawals.HashTreeRoot()
	require.NoError(t, err)
	require.NotNil(t, root)
	require.EqualValues(t, test.ExpectedRoot, root)
}
