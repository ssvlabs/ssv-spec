package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	maxSizeSignedSSVMessageFromQBFTWithNoJustification  = 3700
	maxSizeSignedSSVMessageFromQBFTWith1Justification   = 51852
	maxSizeSignedSSVMessageFromQBFTWith2Justification   = 9114816
	maxSizeSignedSSVMessageFromPartialSignatureMessages = 221316
	maxSizeFullConsensusData                            = 8388836
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

func MaxSignedSSVMessageFromQBFTMessageWithNoJustification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SignedSSVMessage from qbftMessage with no justification",
		testdoc.StructureSizeTestMaxSignedSSVMessageFromQBFTWithNoJustificationDoc,
		maxSignedSSVMessageFromObject(maxQbftMessageNoJustification()),
		maxSizeSignedSSVMessageFromQBFTWithNoJustification,
		false,
	)
}

func MaxSignedSSVMessageFromQBFTMessageWith1Justification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SignedSSVMessage from qbftMessage with 1 justification",
		testdoc.StructureSizeTestMaxSignedSSVMessageFromQBFTWith1JustificationDoc,
		maxSignedSSVMessageFromObject(maxQbftMessageWith1Justification()),
		maxSizeSignedSSVMessageFromQBFTWith1Justification,
		false,
	)
}

func MaxSignedSSVMessageFromQBFTMessageWith2Justification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SignedSSVMessage from qbftMessage with 2 justifications",
		testdoc.StructureSizeTestMaxSignedSSVMessageFromQBFTWith2JustificationDoc,
		maxSignedSSVMessageWithFullDataFromObject(maxQbftMessageWith2Justification()),
		maxSizeSignedSSVMessageFromQBFTWith2Justification,
		true,
	)
}

func MaxSignedSSVMessageFromPartialSignatureMessages() *StructureSizeTest {
	return NewStructureSizeTest(
		"max SignedSSVMessage from PartialSignatureMessages",
		testdoc.StructureSizeTestMaxSignedSSVMessageFromPartialSignatureMessagesDoc,
		maxSignedSSVMessageFromObject(maxPartialSignatureMessages()),
		maxSizeSignedSSVMessageFromPartialSignatureMessages,
		false,
	)
}
