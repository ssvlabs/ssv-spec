package proposer

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SevenOperators tests round-robin proposer selection for 7 member committee
func SevenOperators() tests.SpecTest {
	var p types.OperatorID
	heights := make([]qbft.Height, 0)
	rounds := make([]qbft.Round, 0)
	proposers := make([]types.OperatorID, 0)
	for h := qbft.FirstHeight; h < 100; h++ {
		p = types.OperatorID(h%7) + 1
		for r := qbft.FirstRound; r < 100; r++ {
			heights = append(heights, h)
			rounds = append(rounds, r)
			proposers = append(proposers, p)

			p++
			if p == 8 {
				p = 1
			}
		}
	}

	//fmt.Printf("h:%v\nr:%v\np:%v\n", heights, rounds, proposers)

	return &tests.RoundRobinSpecTest{
		Name:      "7 member committee",
		Share:     testingutils.TestingCommitteeMember(testingutils.Testing7SharesSet()),
		Heights:   heights,
		Rounds:    rounds,
		Proposers: proposers,
	}
}
