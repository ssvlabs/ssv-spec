package tests

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type RoundRobinSpecTest struct {
	Name      string
	Share     *types.CommitteeMember
	Heights   []qbft.Height
	Rounds    []qbft.Round
	Proposers []types.OperatorID
}

func (test *RoundRobinSpecTest) Run(t *testing.T) {
	require.True(t, len(test.Heights) > 0)
	for i, h := range test.Heights {
		r := test.Rounds[i]
		s := &qbft.State{
			Height:          h,
			Round:           r,
			CommitteeMember: test.Share,
		}

		require.EqualValues(t, test.Proposers[i], qbft.RoundRobinProposer(s, r))
	}
}

func (test *RoundRobinSpecTest) TestName() string {
	return "qbft round robin " + test.Name
}

func (test *RoundRobinSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}
