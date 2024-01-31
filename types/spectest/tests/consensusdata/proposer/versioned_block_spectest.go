package consensusdataproposer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/types"
)

type ProposerSpecTest struct {
	Name            string
	Blinded         bool
	DataCd          []byte
	DataBlk         []byte
	ExpectedBlkRoot [32]byte
	ExpectedCdRoot  [32]byte
	ExpectedError   string
}

func (test *ProposerSpecTest) TestName() string {
	return test.Name
}

// ISSUE 219: make multi proposer spec test for all (upcoming) versions

func (test *ProposerSpecTest) Run(t *testing.T) {
	// decode cd
	cd := &types.ConsensusData{}
	require.NoError(t, cd.Decode(test.DataCd))

	if test.Blinded {
		// blk data
		vBlk, hashRoot, err := cd.GetBlindedBlockData()
		if len(test.ExpectedError) != 0 {
			require.EqualError(t, err, test.ExpectedError)
			return
		}
		require.NoError(t, err)
		require.NotNil(t, hashRoot)
		require.NotNil(t, vBlk)

		// compare block roots
		blkRoot, err := vBlk.Root()
		require.NoError(t, err)
		require.NotNil(t, blkRoot)

		root, err := hashRoot.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)

		require.EqualValues(t, blkRoot, root)
		require.EqualValues(t, test.ExpectedBlkRoot, blkRoot)

		// compare blk data
		var blkSSZ []byte
		switch vBlk.Version {
		case spec.DataVersionBellatrix:
			require.NotNil(t, vBlk.Bellatrix)
			blkSSZ, err = vBlk.Bellatrix.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionCapella:
			require.NotNil(t, vBlk.Capella)
			blkSSZ, err = vBlk.Capella.MarshalSSZ()
			require.NoError(t, err)
		default:
			require.Failf(t, "unknown blinded block version %s", vBlk.Version.String())
		}
		require.EqualValues(t, test.DataBlk, blkSSZ)

	} else {
		// blk data
		vBlk, hashRoot, err := cd.GetBlockData()
		if len(test.ExpectedError) != 0 {
			require.EqualError(t, err, test.ExpectedError)
			return
		}

		require.NoError(t, err)
		require.NotNil(t, hashRoot)
		require.NotNil(t, vBlk)

		// compare block roots
		blkRoot, err := vBlk.Root()
		require.NoError(t, err)
		require.NotNil(t, blkRoot)

		root, err := hashRoot.HashTreeRoot()
		require.NoError(t, err)
		require.NotNil(t, root)

		require.EqualValues(t, blkRoot, root)
		require.EqualValues(t, test.ExpectedBlkRoot, blkRoot)

		// compare blk data
		var blkSSZ []byte
		switch vBlk.Version {
		case spec.DataVersionBellatrix:
			require.NotNil(t, vBlk.Bellatrix)
			blkSSZ, err = vBlk.Bellatrix.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionCapella:
			require.NotNil(t, vBlk.Capella)
			blkSSZ, err = vBlk.Capella.MarshalSSZ()
			require.NoError(t, err)
		default:
			require.Failf(t, "unknown block version %s", vBlk.Version.String())
		}
		require.EqualValues(t, test.DataBlk, blkSSZ)
	}

	// compare cd roots
	cdRoot, err := cd.HashTreeRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedCdRoot, cdRoot)

	// compare cd data
	byts, err := cd.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.DataCd, byts)
}
