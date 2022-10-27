package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateRoundChangePreviouslyPrepared tests creating a round change msg,previously prepared
func CreateRoundChangePreviouslyPrepared() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create round change previously prepared",
		Value:      []byte{1, 2, 3, 4},
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
			}),
		},
		ExpectedRoot: "d83735f8a99fabde910e4dff198163811be43ffbb1ffc62a80c75c44618270c1",
	}
}
