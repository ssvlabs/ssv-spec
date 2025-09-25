package consensusdataproposer

import (
	"fmt"
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type ProposerSpecTest struct {
	Name            string
	Type            string
	Documentation   string
	Blinded         bool
	DataCd          []byte
	DataBlk         []byte
	ExpectedBlkRoot types.ExpectedBlkRoot
	ExpectedCdRoot  types.ExpectedCdRoot
	ExpectedError   string
}

func (test *ProposerSpecTest) TestName() string {
	return test.Name
}

func (test *ProposerSpecTest) Run(t *testing.T) {
	// decode cd
	cd := &types.ValidatorConsensusData{}
	require.NoError(t, cd.Decode(test.DataCd))

	// blk data - GetBlockData now handles both blinded and regular blocks
	vBlk, hashRoot, err := cd.GetBlockData()
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
		return
	}
	require.NoError(t, err)
	require.NotNil(t, hashRoot)
	require.NotNil(t, vBlk)

	// Verify the block type matches the test expectation
	require.Equal(t, test.Blinded, vBlk.Blinded, "block blinded state mismatch")

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
	if test.Blinded {
		switch vBlk.Version {
		case spec.DataVersionCapella:
			require.NotNil(t, vBlk.CapellaBlinded)
			blkSSZ, err = vBlk.CapellaBlinded.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionDeneb:
			require.NotNil(t, vBlk.DenebBlinded)
			blkSSZ, err = vBlk.DenebBlinded.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionElectra:
			require.NotNil(t, vBlk.ElectraBlinded)
			blkSSZ, err = vBlk.ElectraBlinded.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionFulu:
			require.NotNil(t, vBlk.FuluBlinded)
			blkSSZ, err = vBlk.FuluBlinded.MarshalSSZ()
			require.NoError(t, err)
		default:
			require.Fail(t, fmt.Sprintf("unknown blinded block version %d", vBlk.Version))
		}
	} else {
		switch vBlk.Version {
		case spec.DataVersionCapella:
			require.NotNil(t, vBlk.Capella)
			blkSSZ, err = vBlk.Capella.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionDeneb:
			require.NotNil(t, vBlk.Deneb)
			blkSSZ, err = vBlk.Deneb.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionElectra:
			require.NotNil(t, vBlk.Electra)
			blkSSZ, err = vBlk.Electra.MarshalSSZ()
			require.NoError(t, err)
		case spec.DataVersionFulu:
			require.NotNil(t, vBlk.Fulu)
			blkSSZ, err = vBlk.Fulu.MarshalSSZ()
			require.NoError(t, err)
		default:
			require.Fail(t, fmt.Sprintf("unknown block version %d", vBlk.Version))
		}
	}
	require.EqualValues(t, test.DataBlk, blkSSZ)

	// compare cd roots
	cdRoot, err := cd.HashTreeRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedCdRoot, cdRoot)

	// compare cd data
	byts, err := cd.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.DataCd, byts)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewProposerSpecTest(name string, documentation string, blinded bool, dataCd []byte, dataBlk []byte, expectedBlkRoot [32]byte, expectedCdRoot [32]byte, expectedError string) *ProposerSpecTest {
	return &ProposerSpecTest{
		Name:            name,
		Type:            testdoc.ProposerSpecTestType,
		Documentation:   documentation,
		Blinded:         blinded,
		DataCd:          dataCd,
		DataBlk:         dataBlk,
		ExpectedBlkRoot: expectedBlkRoot,
		ExpectedCdRoot:  expectedCdRoot,
		ExpectedError:   expectedError,
	}
}
