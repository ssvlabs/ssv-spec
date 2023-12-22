package ssz

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/types"
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
	cd := &types.ConsensusData{}
	require.NoError(t, cd.Decode(test.Data))

	vBlk, _, err := cd.GetBlockData()
	require.NoError(t, err)
	require.NotNil(t, vBlk)

	switch vBlk.Version {

	case spec.DataVersionCapella:

		// Test BeaconBlock root
		root, err := vBlk.Capella.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)
		require.EqualValues(t, test.ExpectedRoot, root)

		// test SSZWithdrawals
		var withdrawals types.SSZWithdrawals
		withdrawals = types.SSZWithdrawals(vBlk.Capella.Body.ExecutionPayload.Withdrawals)
		root, err = withdrawals.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)

	case spec.DataVersionDeneb:

		// Test BeaconBlock root
		root, err := vBlk.Deneb.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)
		require.EqualValues(t, test.ExpectedRoot, root)

		// test SSZWithdrawals
		var withdrawals types.SSZWithdrawals
		withdrawals = types.SSZWithdrawals(vBlk.Deneb.Body.ExecutionPayload.Withdrawals)
		root, err = withdrawals.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)

		// test SSZBlobKZGCommitments
		var blobKZGCommitments types.SSZBlobZGCommitments
		blobKZGCommitments = types.SSZBlobZGCommitments(vBlk.Deneb.Body.BlobKZGCommitments)
		root, err = blobKZGCommitments.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)
	}
}
