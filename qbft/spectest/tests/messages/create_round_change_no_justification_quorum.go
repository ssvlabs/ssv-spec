package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// CreateRoundChangeNoJustificationQuorum tests creating a round change msg that was previouly prepared
// but failed to extract a justification quorum (shouldn't happen).
// The result should be an unjustified round change.
func CreateRoundChangeNoJustificationQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := CreateRoundChangeNoJustificationQuorumSC()
	return &tests.CreateMsgSpecTest{
		CreateType:    tests.CreateRoundChange,
		Name:          "create round change no justification quorum",
		StateValue:    testingutils.TestingQBFTFullData,
		ExpectedState: sc.ExpectedState,
		PrepareJustifications: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		},
		ExpectedRoot: sc.Root(),
	}
}

func CreateRoundChangeNoJustificationQuorumSC() *comparable.StateComparison {
	expectedMsg := qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   0,
		Round:                    1,
		Identifier:               []byte{1, 2, 3, 4},
		Root:                     testingutils.TestingQBFTRootData,
		DataRound:                1,
		RoundChangeJustification: [][]byte{},
		PrepareJustification:     nil,
		FullData:                 testingutils.TestingQBFTFullData,
	}

	ks := testingutils.Testing4SharesSet()

	signedMsg := testingutils.SignQBFTMsg(ks.NetworkKeys[1], 1, &expectedMsg)

	return &comparable.StateComparison{ExpectedState: signedMsg}
}
