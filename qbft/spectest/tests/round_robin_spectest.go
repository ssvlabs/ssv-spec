package tests

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
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

		operators := make([]types.OperatorID, 0)
		for _, operator := range test.Share.Committee {
			operators = append(operators, operator.OperatorID)
		}

		require.EqualValues(t, test.Proposers[i], qbft.RoundRobinProposer(h, r, operators))
	}
}

func (test *RoundRobinSpecTest) TestName() string {
	return "qbft round robin " + test.Name
}

func (test *RoundRobinSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}
