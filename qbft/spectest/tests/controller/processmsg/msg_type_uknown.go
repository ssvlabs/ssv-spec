package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/controller/futuremsg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgTypeUnknown tests a process msg with unknown msg type
func MsgTypeUnknown() *futuremsg.ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root})
	msgEncoded, _ := msg.Encode()

	return &futuremsg.ControllerSyncSpecTest{
		Name: "future msg unknown msg type",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.MsgType{0x9, 0x0, 0x0, 0x0}),
				Data: msgEncoded,
			},
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid msg: message type not supported",
	}
}
