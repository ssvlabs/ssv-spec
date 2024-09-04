package messageratetest

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
)

// SeveralBigCommittees tests the expected message rate for the case of several big committees, each with 500 validator
func SeveralBigCommittees() tests.SpecTest {

	return &MessageRateTest{
		Name: "several big committees",

		TestCases: []TestCase{
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 1000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 2_112.1248,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 5000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 10_560.6243,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 10000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 21_121.2487,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 20000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 42_242.4975,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 50_000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 105_606.2438,
			},
			TestCase{
				CommitteesConfig: []CommitteesWithValidators{
					CommitteesWithValidators{
						NumCommittees: 100_000,
						NumValidators: 500,
					},
				},
				ExpectedMessageRate: 211_212.4876,
			},
		},
	}
}
