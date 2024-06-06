package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizeSignedSSVMessageFromQBFTWithNoJustification  = 3700
	maxSizeSignedSSVMessageFromQBFTWith1Justification   = 51852
	maxSizeSignedSSVMessageFromQBFTWith2Justification   = 4945164
	maxSizeSignedSSVMessageFromPartialSignatureMessages = 147588
	maxSizeFullData                                     = 4219184
)

func maxSignedSSVMessageFromObject(obj types.Encoder) *types.SignedSSVMessage {

	objBytes, err := obj.Encode()
	if err != nil {
		panic(err)
	}

	signatures := [][]byte{}
	for i := 0; i < 13; i++ {
		sig := [256]byte{1}
		signatures = append(signatures, sig[:])
	}
	signers := [13]types.OperatorID{1}

	return &types.SignedSSVMessage{
		Signatures:  signatures,
		OperatorIDs: signers[:],
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   [56]byte{1},
			Data:    objBytes,
		},
	}
}

func maxSignedSSVMessageWithFullDataFromObject(obj types.Encoder) *types.SignedSSVMessage {
	msg := maxSignedSSVMessageFromObject(obj)
	msg.FullData = maxFullData()
	return msg
}

func SignedSSVMessageFromQBFTMessageWithNoJustification() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "max SignedSSVMessage from qbftMessage with no justification",
		Object:                maxSignedSSVMessageFromObject(maxQbftMessageNoJustification()),
		ExpectedEncodedLength: maxSizeSignedSSVMessageFromQBFTWithNoJustification,
		IsMaxSizeForType:      false,
	}
}

func SignedSSVMessageFromQBFTMessageWith1Justification() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "max SignedSSVMessage from qbftMessage with 1 justification",
		Object:                maxSignedSSVMessageFromObject(maxQbftMessageWith1Justification()),
		ExpectedEncodedLength: maxSizeSignedSSVMessageFromQBFTWith1Justification,
		IsMaxSizeForType:      false,
	}
}

func SignedSSVMessageFromQBFTMessageWith2Justification() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "max SignedSSVMessage from qbftMessage with 2 justifications",
		Object:                maxSignedSSVMessageWithFullDataFromObject(maxQbftMessageWith2Justification()),
		ExpectedEncodedLength: maxSizeSignedSSVMessageFromQBFTWith2Justification,
		IsMaxSizeForType:      false,
	}
}

func SignedSSVMessageFromPartialSignatureMessages() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "max SignedSSVMessage from PartialSignatureMessages",
		Object:                maxSignedSSVMessageFromObject(maxPartialSignatureMessages()),
		ExpectedEncodedLength: maxSizeSignedSSVMessageFromPartialSignatureMessages,
		IsMaxSizeForType:      false,
	}
}
