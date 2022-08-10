package tests

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	init := &dkg.Init{
		OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
		Threshold:             3,
		WithdrawalCredentials: testingutils.TestingWithdrawalCredentials,
		Fork:                  testingutils.TestingForkVersion,
	}
	initBytes, _ := init.Encode()
	root := testingutils.DespositDataSigningRoot(ks, init)

	return &MsgProcessingSpecTest{
		Name: "happy flow",
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: identifier,
				Data:       initBytes,
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       nil, // GLNOTE: Dummy message simulating the KeyGenProtocol to complete
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(2, root, ks.Shares[2]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(3, root, ks.Shares[3]),
			}),

			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       testingutils.SignedOutputBytes(identifier, 2, root, ks.DKGOperators[2].ETHAddress, ks.Shares[2], ks.ValidatorSK),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       testingutils.SignedOutputBytes(identifier, 3, root, ks.DKGOperators[3].ETHAddress, ks.Shares[3], ks.ValidatorSK),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(1, root, ks.Shares[1]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       testingutils.SignedOutputBytes(identifier, 1, root, ks.DKGOperators[1].ETHAddress, ks.Shares[1], ks.ValidatorSK),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			1: testingutils.SignedOutputObject(identifier, 1, root, ks.DKGOperators[1].ETHAddress, ks.Shares[1], ks.ValidatorSK),
			2: testingutils.SignedOutputObject(identifier, 2, root, ks.DKGOperators[2].ETHAddress, ks.Shares[2], ks.ValidatorSK),
			3: testingutils.SignedOutputObject(identifier, 3, root, ks.DKGOperators[3].ETHAddress, ks.Shares[3], ks.ValidatorSK),
		},
		KeySet:        testingutils.Testing4SharesSet(),
		ExpectedError: "",
	}
}
