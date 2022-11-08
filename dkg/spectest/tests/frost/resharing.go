package frost

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func Resharing() *FrostSpecTest {

	requestID := testingutils.GetRandRequestID()
	ks := testingutils.Testing13SharesSet()

	threshold := 3
	operators := []types.OperatorID{5, 6, 7, 8}
	operatorsOld := []types.OperatorID{1, 2, 3} //4}

	vk, _ := hex.DecodeString("8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812")
	reshareMsgBytes := testingutils.ReshareMessageDataBytes(
		operators,
		uint16(threshold),
		vk,
	)

	initMessages := make(map[uint32][]*dkg.SignedMessage)
	for _, operatorID := range append(operators, operatorsOld...) {
		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
				MsgType:    dkg.ReshareMsgType,
				Identifier: requestID,
				Data:       reshareMsgBytes,
			}),
		}
	}

	spectest := &FrostSpecTest{
		Name:   "Simple Resharing",
		Keyset: ks,

		RequestID: requestID,
		Threshold: uint64(threshold),
		Operators: operators,

		IsResharing:  true,
		OperatorsOld: operatorsOld,
		OldKeygenOutcomes: testingutils.TestOutcome{
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

		ExpectedOutcome: testingutils.TestOutcome{
			KeygenOutcome: testingutils.TestKeygenOutcome{
				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
				Share: map[uint32]string{
					5: "52046f0837c928ea5d5bbc893b90f3cd75a07a9d25092e2fbb0129825100c3be",
					6: "0f1d824d53df922ca8c15d639c802f84463a78cf69ef57e0b1cbb8b95cd1f458",
					7: "213989136198ba32e82eb8e449a843b7fb6a52007ba72794212d25d135a84679",
					8: "146adc07375723b4e869f703396758634172622d5a32414b092570cadb83ba20",
				},
				OperatorPubKeys: map[uint32]string{
					5: "81e52afe4656f4544715cc2a37724c939afa8462d57549ba242681b52c80d8ac7e6b259d03ba37ce688aeca5e1a346b3",
					6: "ac6d0b0ba2f3f581f520c59049c6dfb98ce12d87a3ee9ccc00b9e0ef13153b036c777a946d9ec78409a047d92ce942e7",
					7: "8d9b4d117564b4852ee7d060626e27bc93ec5dddde0fbbbe053aed7e54b0772b334ad74149fb1c6d3f1ff3d5b4d87fc8",
					8: "b11cb28641e5d6440e214d45abfc6a2158cbf163312144609e08236fee95aa096a61a0d70b4401d8daf4af69a1cca9ad",
				},
			},
		},
		ExpectedError: "",

		InputMessages: map[int]MessagesForNodes{
			0: initMessages,
		},
	}
	return spectest
}
