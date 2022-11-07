package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests future msg invalid msg
func InvalidMsg() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		Height: 10,
		Round:  3,
	}, &qbft.Data{Root: inputData.Root})
	msg.Signature = nil
	msgEncoded, _ := msg.Encode()

	return &ControllerSyncSpecTest{
		Name: "future invalid msg",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: msgEncoded,
			},
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid future msg: invalid future msg: message signature is invalid",
	}
}
