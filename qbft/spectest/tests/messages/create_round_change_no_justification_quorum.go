package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
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
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
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
	}

	encodedExpectedMsg, err := expectedMsg.Encode()
	if err != nil {
		panic(err)
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   [56]byte{1, 2, 3, 4},
		Data:    encodedExpectedMsg,
	}

	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)

	sig, err := config.OperatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		panic(err)
	}

	signedMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,

		FullData: testingutils.TestingQBFTFullData,
	}
	return &comparable.StateComparison{ExpectedState: signedMsg}
}
