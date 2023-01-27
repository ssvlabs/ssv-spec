package tests

import (
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/stretchr/testify/require"
)

type RoundRobinSpecTest struct {
	Name      string
	Share     *types.Share
	Heights   []alea.Height
	Rounds    []alea.Round
	Proposers []types.OperatorID
}

func (test *RoundRobinSpecTest) Run(t *testing.T) {
	require.True(t, len(test.Heights) > 0)
	for i, h := range test.Heights {
		r := test.Rounds[i]
		s := &alea.State{
			Height: h,
			Round:  r,
			Share:  test.Share,
		}

		require.EqualValues(t, test.Proposers[i], alea.RoundRobinProposer(s, r))
	}
}

func (test *RoundRobinSpecTest) TestName() string {
	return "alea round robin " + test.Name
}
