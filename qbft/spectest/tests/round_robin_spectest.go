package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type RoundRobinSpecTest struct {
	Name      string
	Share     *types.Share
	Heights   []qbft.Height
	Rounds    []qbft.Round
	Proposers []types.OperatorID
}

func (test *RoundRobinSpecTest) Run(t *testing.T) {
	require.True(t, len(test.Heights) > 0)
	for i, h := range test.Heights {
		r := test.Rounds[i]
		s := &qbft.State{
			Height: h,
			Round:  r,
			Share:  test.Share,
		}

		require.EqualValues(t, test.Proposers[i], qbft.RoundRobinProposer(s, r))
	}
}

func (test *RoundRobinSpecTest) TestName() string {
	return "qbft round robin " + test.Name
}
