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
		maxSSVMessageFromQBFTMessage(),
		maxSizeSSVMessageFromQBFTMessage,
		true,
	)
}

func MaxSSVMessageFromPartialSignatureMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SSVMessage from PartialSignatureMessages",
		maxSSVMessageFromPartialSignatureMessages(),
		maxSizeSSVMessageFromPartialSignatureMessages,
		false,
	)
}
