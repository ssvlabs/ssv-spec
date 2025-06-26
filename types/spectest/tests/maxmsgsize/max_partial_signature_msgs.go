package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizePartialSignatureMessage                 = 144
	maxSizePartialSignatureMessages                = 217748
	maxSizePartialSignatureMessagesForPreConsensus = 1892
)

func maxPartialSignatureMessage() *types.PartialSignatureMessage {

	signature := [96]byte{1}

	return &types.PartialSignatureMessage{
		PartialSignature: signature[:],
		SigningRoot:      [32]byte{1},
		Signer:           1,
		ValidatorIndex:   1,
	}
}

func maxPartialSignatureMessages() *types.PartialSignatureMessages {

	msgs := make([]*types.PartialSignatureMessage, 0)
	for i := 0; i < 1512; i++ {
		msgs = append(msgs, maxPartialSignatureMessage())
	}

	return &types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     1,
		Messages: msgs,
	}
}

func maxPartialSignatureMessagesForPreConsensus() *types.PartialSignatureMessages {

	msgs := make([]*types.PartialSignatureMessage, 0)
	for i := 0; i < 13; i++ {
		msgs = append(msgs, maxPartialSignatureMessage())
	}

	return &types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     1,
		Messages: msgs,
	}
}

func MaxPartialSignatureMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"max PartialSignatureMessage",
		"Test the maximum size of a single partial signature message",
		maxPartialSignatureMessage(),
		maxSizePartialSignatureMessage,
		true,
	)
}

func MaxPartialSignatureMessages() *StructureSizeTest {
	return NewStructureSizeTest(
		"max PartialSignatureMessages",
		"Test the maximum size of partial signature messages collection",
		maxPartialSignatureMessages(),
		maxSizePartialSignatureMessages,
		true,
	)
}

func MaxPartialSignatureMessagesForPreConsensus() *StructureSizeTest {
	return NewStructureSizeTest(
		"max PartialSignatureMessages for pre-consensus",
		"Test the maximum size of partial signature messages for pre-consensus phase",
		maxPartialSignatureMessagesForPreConsensus(),
		maxSizePartialSignatureMessagesForPreConsensus,
		false,
	)
}
