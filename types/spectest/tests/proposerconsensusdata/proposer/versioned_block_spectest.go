package consensusdataproposer

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type ProposerSpecTest struct {
	Name              string
	Type              string
	Documentation     string
	DataCd            []byte
	DataBlk           []byte
	ExpectedBlkRoot   types.ExpectedBlkRoot
	ExpectedCdRoot    types.ExpectedCdRoot
	ExpectedErrorCode int
}

func (test *ProposerSpecTest) TestName() string {
	return test.Name
}

func (test *ProposerSpecTest) Run(t *testing.T) {
	// decode cd
	cd := &types.ProposerConsensusData{}
	require.NoError(t, cd.Decode(test.DataCd))

	// blk data
	proposalData, err := cd.GetBlockData()
	tests.AssertErrorCode(t, test.ExpectedErrorCode, err)
	if err != nil {
		return
	}
	require.NotNil(t, proposalData)
	require.NotNil(t, proposalData.Block)

	// compare block roots
	blkRoot, err := proposalData.Block.HashTreeRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedBlkRoot, blkRoot)

	// compare blk data
	blkSSZ, err := proposalData.Block.MarshalSSZ()
	require.NoError(t, err)
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

func NewProposerSpecTest(name string, documentation string, dataCd []byte, dataBlk []byte, expectedBlkRoot [32]byte, expectedCdRoot [32]byte, expectedErrorCode int) *ProposerSpecTest {
	return &ProposerSpecTest{
		Name:              name,
		Type:              testdoc.ProposerSpecTestType,
		Documentation:     documentation,
		DataCd:            dataCd,
		DataBlk:           dataBlk,
		ExpectedBlkRoot:   expectedBlkRoot,
		ExpectedCdRoot:    expectedCdRoot,
		ExpectedErrorCode: expectedErrorCode,
	}
}
