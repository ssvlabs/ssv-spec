package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

var (
	maliciousOperatorID uint32 = 2
)

func GetBlameSpecTest(testName string, data []byte) *FrostSpecTest {

	requestID := testingutils.GetRandRequestID()
	ks := testingutils.Testing4SharesSet()

	threshold := 3
	operators := []types.OperatorID{1, 2, 3, 4}

	initMessages := make(map[uint32][]*dkg.SignedMessage)
	initMsgBytes := testingutils.InitMessageDataBytes(
		operators,
		uint16(threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	)
	for _, operatorID := range operators {
		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: requestID,
				Data:       initMsgBytes,
			}),
		}
	}

	blameProtocolMessageBytes := []byte(data)
	blameSignedMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: requestID,
			Data:       blameProtocolMessageBytes,
		},
		Signer: types.OperatorID(maliciousOperatorID),
	}
	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(blameSignedMessage, ks.DKGOperators[2].ETHAddress)
	blameSignedMessage.Signature = sig

	return &FrostSpecTest{
		Name:   testName,
		Keyset: ks,

		RequestID: requestID,
		Threshold: uint64(threshold),
		Operators: operators,

		ExpectedOutcome: testingutils.TestOutcome{
			BlameOutcome: testingutils.TestBlameOutcome{
				Valid: true,
			},
		},
		ExpectedError: "",

		InputMessages: map[int]MessagesForNodes{
			0: initMessages,
			int(frost.Round1): {
				maliciousOperatorID: []*dkg.SignedMessage{blameSignedMessage},
			},
		},
	}
}

var validRound1MsgBytes = []byte(`{"round":2,"round1":{"Commitment":["k4MGdjUKb5vaeYp9RKAjfli3H5uZfH8L3GMOyA5S7wiLi4Y33svE9h7TEYYyWVyZ","gTzEEi8g/nnVPMNDVxweRg19h2pgbIvvvLW8uUCKTqWoN1ziDjIhgJEg8oFhm04x","tp1XvVz+tm2Q1VwoyX/f44xUCNwHeYFD/xhEUwHG9+oL7/foxIjMpmRM3p7PRYvu"],"ProofS":"OSFPgOln/HjbVNReoKpvScuPAqCiP9anwHZnBkuZEDI=","ProofR":"C3LlHscA6yGlCEpPImTcB2drK7UUZN6Teur7uF2OOPc=","Shares":{"1":"BAZR3YHRROd8C9Sdm86LUd3MokNA67Vm6PV9WGO5zYUvZBN3VzFSNi3LWFCzwLgl8s/5EbTW0mtIvTz4xrnGV51DXO9G/Dpm5aTboIly98aM/AOEQMj3OHmEhsUsWPM5ADTTf9E+FBuEns78J9q6ckR0BDCeZMg60BPJhV1hoPGO","3":"BC9VN12p/aB3n52j4A/9ULr7brl4oGvF7pq4In5vhF7fNV5MUx4IlsYRIOm5w6GiMOvhd3WbgtDFTENmZet3ALT4d2sDlXpjjeSSjNH0unlNUCFqxt7u31e4kKG+NLJwEMKzkD2G+8QyEUw47/SQmipQErNzDMMCucQrGzrR7EHK","4":"BMxCC/3HWsJMqCqknlC46IrMqO0EMTCUWkSkNcq/0cZuw+cG4TeTpAjA3JFE1AEEqyVRu+J1DmaIoYnFwfxs+SPLwoRYcGhqBsyXTJbvr89J8rQ6oHsFUG80RPtb1nIolsCch32nwR2cvmHsIrpyCugWWGkORxIk0Dq3RAroMRYT"}}}`)

func BlameTypeInvalidCommitment() *FrostSpecTest {
	return GetBlameSpecTest(
		"Blame Type Invalid Commitment - Happy Flow",
		makeInvalidForInvalidCommitment(validRound1MsgBytes),
	)
}

func BlameTypeInvalidScalar() *FrostSpecTest {
	return GetBlameSpecTest(
		"Blame Type Invalid Scalar - Happy Flow",
		makeInvalidForInvalidScalar(validRound1MsgBytes),
	)
}

func BlameTypeInvalidShare_FailedShareDecryption() *FrostSpecTest {
	return GetBlameSpecTest(
		"Blame Type Invalid Share (Unable to Decrypt) - Happy Flow",
		makeInvalidForFailedEcies(validRound1MsgBytes),
	)
}

func BlameTypeInvalidShare_FailedValidationAgainstCommitment() *FrostSpecTest {

	requestID := testingutils.GetRandRequestID()
	ks := testingutils.Testing4SharesSet()

	threshold := 3
	operators := []types.OperatorID{1, 2, 3, 4}

	initMessages := make(map[uint32][]*dkg.SignedMessage)
	initMsgBytes := testingutils.InitMessageDataBytes(
		operators,
		uint16(threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	)
	for _, operatorID := range operators {
		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: requestID,
				Data:       initMsgBytes,
			}),
		}
	}

	pmData := `{"round":2,"round1":{"Commitment":["qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa","luG8uPgVeTmvRoE4MBPMpt/Vgp/oCAA9TBTzG07bRJ45L6Uo9uDGQ9gKFkt9+07n","kBb0Obbc/CYaBH/56rZtOUw6bz6DMFrbouRUNwh8lBfH8OjWg3NQwBIXE3Ir8lmn","p7CN+Aow0TPzKW0wFmbL3qNYuSnjLEo3Gtjapg71mDYn9+IGpmGtoqHAd1LUzxDA"],"ProofS":"YI5p9V013jgYX/78TJ3PSdI/5QEbZFKnRND0pTo6XAE=","ProofR":"MW4+TKI7AAf/q0OljiJiNLSkoAPXy4PzTXC2dFqhlAc=","Shares":{"1":"BLh2p4b+/slitPiMXooPEka+S6TqCcdQSB7Bzv1XTZNp0N5wpnI/jgA4qAwzg2YCbVdBzcG26FF5p/4FRHk1syDz0ljuJkv30ahpxt/bby1ItMnBKgy7p+zYOE9RkAlecpnowYohR3wj/Fxq/ln5gNRWDmMcWMePrflm5dpMCziY","3":"BMqwinOzpjBtLed3b/pCDuG2x9XQPzXlKIXHtR+8pK4R+qPbU6hB4Xgf/9D/b2PKs/jnH6XOKfLX7q1bC9DZkH59cmeeeAHFjy3YeObXyF3L7E1MX4NWHxmkjjWSLiH08M2MkQCtfrswWzIfOVT7YgFJSRRDy2sf94CA1WdbFAc3","4":"BKLu/3DyaZOIzXx6SkKrRhxojh30Y5uLXOqBGbt5hQiPZpbDtULBSdTr4XrUyXLdmi3JW+jiZihksHxHjZWgq7mdlhMspFYyEnqxmobMBuDxydmEa5VzdoCtsSrRc79kx9SkwMsNkmY52VhtkgwLlISCXSdmAQc8BIn3kT7xQcy1"}}}`

	blameProtocolMessageBytes := []byte(pmData)
	blameSignedMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: requestID,
			Data:       blameProtocolMessageBytes,
		},
		Signer: 2,
	}
	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(blameSignedMessage, ks.DKGOperators[2].ETHAddress)
	blameSignedMessage.Signature = sig

	return &FrostSpecTest{
		Name:   "Blame Type Invalid Share - Happy Flow",
		Keyset: ks,

		RequestID: requestID,
		Threshold: uint64(threshold),
		Operators: operators,

		ExpectedOutcome: testingutils.TestOutcome{
			BlameOutcome: testingutils.TestBlameOutcome{
				Valid: true,
			},
		},
		ExpectedError: "",

		InputMessages: map[int]MessagesForNodes{
			0: initMessages,
			2: {
				2: []*dkg.SignedMessage{blameSignedMessage},
			},
		},
	}
}

func BlameTypeInconsistentMessage() *FrostSpecTest {

	requestID := testingutils.GetRandRequestID()
	ks := testingutils.Testing4SharesSet()

	threshold := 3
	operators := []types.OperatorID{1, 2, 3, 4}

	initMessages := make(map[uint32][]*dkg.SignedMessage)
	initMsgBytes := testingutils.InitMessageDataBytes(
		operators,
		uint16(threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	)
	for _, operatorID := range operators {
		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: requestID,
				Data:       initMsgBytes,
			}),
		}
	}

	pmData1 := `{"round":2,"round1":{"Commitment":["r4pOd119gLXzt06xwvmXudIYrEFHl7ZyT7yXDMz3Wt/CmK+KkPRem6nq4ov5Sf3q","gQh8Bd8lJmokT9zzFUK/javWp8z8VOIp5R/kCyXxCoYqpICyOwmg0XVYZMLwIj/Q","qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa"],"ProofS":"RY6MjnapCPtt6cR9YXWdbdd3Me8BSJCrlTJpX9y5bL8=","ProofR":"NARDVvKTH6pSjiYtMwUqiZSYv1lKNVk6deDB9FcddkA=","Shares":{"1":"BOSDxrY8bmwO+WdwYs/TgDC+viXCYQqldNoEmSutOHrljBoIGmS9KKmxbAYEpdTtk+ahyyOnG0lHn3WTrN9PJEeYE6QpcGgrRkUWOq/RSwHQX50R00iUmCnXH5B3WVUdyTTAOzkvenfWrq6W+uVQ4Vu00k590W9xCbBvtGM4UXJ+","3":"BOlPoCeJDaUsr3bRVGPlU0JZ1OgPm8StbA93DYyEaL5e5Y7PNzEyCnrDPWVoVqnNbPk6GikWHoGd/sOJCB4l7fBiyd0H0H6Ypwz44MFhEu8qgBWxFGeG730HZKv4+6mj048Tfkj1l+tTHdqI8O3GjzwWD51UOl1aIV68swslQqeL","4":"BHZa/+riJkzhM7PFkIFzhlUkqX2K3P1iQZO1wJRTmyvPuqqYnAc0KsbkSnDSq7GTwwA5L+jtle3Y4NlxVFH5lq9RYntwNnDyRliDwzxis8xlRDQtrnAFfIySw+rDJa7clxWUTavMjHeEawDWYv9MIKbPId0AwrlXMRb7pycDMoAW"}}}`

	pmData2 := `{"round":2,"round1":{"Commitment":["r4pOd119gLXzt06xwvmXudIYrEFHl7ZyT7yXDMz3Wt/CmK+KkPRem6nq4ov5Sf3q","gQh8Bd8lJmokT9zzFUK/javWp8z8VOIp5R/kCyXxCoYqpICyOwmg0XVYZMLwIj/Q","qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa"],"ProofS":"RY6MjnapCPtt6cR9YXWdbdd3Me8BSJCrlTJpX9y5bL8=","ProofR":"NARDVvKTH6pSjiYtMwUqiZSYv1lKNVk6deDB9FcddkA=","Shares":{"1":"BOSDxrY8bmwO+WdwYs/TgDC+viXCYQqldNoEmSutOHrljBoIGmS9KKmxbAYEpdTtk+ahyyOnG0lHn3WTrN9PJEeYE6QpcGgrRkUWOq/RSwHQX50R00iUmCnXH5B3WVUdyTTAOzkvenfWrq6W+uVQ4Vu00k590W9xCbBvtGM4UXJ+","3":"BOlPoCeJDaUsr3bRVGPlU0JZ1OgPm8StbA93DYyEaL5e5Y7PNzEyCnrDPWVoVqnNbPk6GikWHoGd/sOJCB4l7fBiyd0H0H6Ypwz44MFhEu8qgBWxFGeG730HZKv4+6mj048Tfkj1l+tTHdqI8O3GjzwWD51UOl1aIV68swslQqeL","4":"BHZa/+riJkzhM7PFkIFzhlUkqX2K3P1iQZO1wJRTmyvPuqqYnAc0KsbkSnDSq7GTwwA5L+jtle3Y4NlxVFH5lq9RYntwNnDyRliDwzxis8xlRDQtrnAFfIySw+rDJa7clxWUTavMjHeEawDWYv9MIKbPId0AwrlXMRb7pycDMoWA"}}}` // root changed

	data1SignedMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: requestID,
			Data:       []byte(pmData1),
		},
		Signer: 2,
	}
	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(data1SignedMessage, ks.DKGOperators[2].ETHAddress)
	data1SignedMessage.Signature = sig

	data2SignedMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: requestID,
			Data:       []byte(pmData2),
		},
		Signer: 2,
	}
	sig, _ = testingutils.NewTestingKeyManager().SignDKGOutput(data2SignedMessage, ks.DKGOperators[2].ETHAddress)
	data2SignedMessage.Signature = sig

	return &FrostSpecTest{
		Name:   "Blame Type Inconsisstent Message - Happy Flow",
		Keyset: ks,

		RequestID: requestID,
		Threshold: uint64(threshold),
		Operators: operators,

		ExpectedOutcome: testingutils.TestOutcome{
			BlameOutcome: testingutils.TestBlameOutcome{
				Valid: true,
			},
		},
		ExpectedError: "",

		InputMessages: map[int]MessagesForNodes{
			0: initMessages,
			2: {
				2: []*dkg.SignedMessage{data1SignedMessage, data2SignedMessage},
			},
		},
	}
}

func makeInvalidForFailedEcies(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.Shares[maliciousOperatorID] = []byte("rubbish-value")
	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInvalidScalar(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.ProofR = []byte("rubbish-value")
	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInvalidCommitment(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.Commitment[len(protocolMessage.Round1Message.Commitment)-1] = []byte("rubbish-value")
	d, _ := protocolMessage.Encode()
	return d
}
