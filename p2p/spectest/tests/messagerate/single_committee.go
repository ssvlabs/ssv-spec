package messageratetest

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
)

// SingleCommittee tests the expected message rate for the minimum case of an unique committee
func SingleCommittee() tests.SpecTest {

	return &MessageRateTest{
		Name: "single committee",

		TestCases: []TestCase{
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 1,
					},
				},
				ExpectedMessageRate: 0.041406,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 10,
					},
				},
				ExpectedMessageRate: 0.361920,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 50,
					},
				},
				ExpectedMessageRate: 1.087111,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 100,
					},
				},
				ExpectedMessageRate: 1.372784,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 250,
					},
				},
				ExpectedMessageRate: 1.680669,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 2.112124,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 750,
					},
				},
				ExpectedMessageRate: 2.543187,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 1_000,
					},
				},
				ExpectedMessageRate: 2.974250,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 5_000,
					},
				},
				ExpectedMessageRate: 9.87125,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1,
						NumValidators: 10_000,
					},
				},
				ExpectedMessageRate: 18.4925,
			},
		},
	}
}
