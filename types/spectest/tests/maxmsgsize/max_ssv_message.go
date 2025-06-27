package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizeSSVMessageFromQBFTMessage              = 722480
	maxSizeSSVMessageFromPartialSignatureMessages = 217816
)

func maxSSVMessageFromData(data []byte) *types.SSVMessage {
	return &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   [56]byte{1},
		Data:    data,
	}
}

func maxSSVMessageFromQBFTMessage() *types.SSVMessage {

	msg := maxQbftMessageWith2Justification()
	msgBytes, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	return maxSSVMessageFromData(msgBytes)
}

func maxSSVMessageFromPartialSignatureMessages() *types.SSVMessage {

	msg := maxPartialSignatureMessages()
	msgBytes, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	return maxSSVMessageFromData(msgBytes)
}

func MaxSSVMessageFromQBFTMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SSVMessage from qbftMessage",
		"Test the maximum size of an SSVMessage containing a QBFT message with 2 justifications",
		maxSSVMessageFromQBFTMessage(),
		maxSizeSSVMessageFromQBFTMessage,
		true,
	)
}

func MaxSSVMessageFromPartialSignatureMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SSVMessage from PartialSignatureMessages",
		"Test the maximum size of an SSVMessage containing partial signature messages",
		maxSSVMessageFromPartialSignatureMessages(),
		maxSizeSSVMessageFromPartialSignatureMessages,
		false,
	)
}
