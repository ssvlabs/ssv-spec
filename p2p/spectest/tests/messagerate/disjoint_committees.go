package messageratetest

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
)

// DisjointCommittees tests the expected message rate for the maximum case of disjoint committees, each with 1 validator
func DisjointCommittees() tests.SpecTest {

	return &MessageRateTest{
		Name: "disjoint committees",
		TestCases: []TestCase{
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 41.40675,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 5_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 207.03375,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 10_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 414.0675,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 20_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 828.135,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 50_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 2_070.3375,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 100_000,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 4_140.675,
			},
		},
	}
}
