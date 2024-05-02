package messages

import (
	"github.com/pkg/errors"
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
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
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

	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	sig, err := config.GetShareSigner().SignRoot(&expectedMsg, types.QBFTSignatureType, config.GetSigningPubKey())
	if err != nil {
		panic(errors.Wrap(err, "unable to sign root for create_round_change_no_justification_quorum"))
	}
	signedMsg := &qbft.SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{1},
		Message:   expectedMsg,

		FullData: testingutils.TestingQBFTFullData,
	}
	return &comparable.StateComparison{ExpectedState: signedMsg}
}
