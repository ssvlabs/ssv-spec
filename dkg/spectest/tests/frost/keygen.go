package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func Keygen() *FrostSpecTest {

	requestID := testingutils.GetRandRequestID()
	ks := testingutils.Testing4SharesSet()

	threshold := 3
	operators := []types.OperatorID{1, 2, 3, 4}
	initMsgBytes := testingutils.InitMessageDataBytes(
		operators,
		uint16(threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	)

	initMessages := make(map[uint32][]*dkg.SignedMessage)
	for _, operatorID := range operators {
		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: requestID,
				Data:       initMsgBytes,
			}),
		}
	}

	return &FrostSpecTest{
		Name:   "Simple Keygen",
		Keyset: ks,

		RequestID: requestID,
		Threshold: uint64(threshold),
		Operators: operators,

		ExpectedOutcome: testingutils.TestOutcome{
			KeygenOutcome: testingutils.TestKeygenOutcome{
				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
				Share: map[uint32]string{
					1: "5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912",
					2: "533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7",
					3: "442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186",
					4: "2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f",
				},
				OperatorPubKeys: map[uint32]string{
					1: "add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa",
					2: "8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5",
					3: "b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051",
					4: "a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037",
				},
			},
		},
		ExpectedError: "",

		InputMessages: map[int]MessagesForNodes{
			0: initMessages,
		},
	}
}
