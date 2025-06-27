package tests

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type RoundRobinSpecTest struct {
	Name          string
	Type          string
	Documentation string
	Share         *types.CommitteeMember
	Heights       []qbft.Height
	Rounds        []qbft.Round
	Proposers     []types.OperatorID
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

func NewRoundRobinSpecTest(name string, documentation string, share *types.CommitteeMember, heights []qbft.Height, rounds []qbft.Round, proposers []types.OperatorID) *RoundRobinSpecTest {
	return &RoundRobinSpecTest{
		Name:          name,
		Type:          "Round robin proposer: validation of fair proposer selection algorithm",
		Documentation: documentation,
		Share:         share,
		Heights:       heights,
		Rounds:        rounds,
		Proposers:     proposers,
	}
}
