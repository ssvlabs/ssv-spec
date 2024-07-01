package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizePartialSignatureMessage                 = 144
	maxSizePartialSignatureMessages                = 144020
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
	for i := 0; i < 1000; i++ {
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
	return &StructureSizeTest{
		Name:                  "max PartialSignatureMessage",
		Object:                maxPartialSignatureMessage(),
		ExpectedEncodedLength: maxSizePartialSignatureMessage,
		IsMaxSize:             true,
	}
}

func MaxPartialSignatureMessages() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "max PartialSignatureMessages",
		Object:                maxPartialSignatureMessages(),
		ExpectedEncodedLength: maxSizePartialSignatureMessages,
		IsMaxSize:             true,
	}
}

func MaxPartialSignatureMessagesForPreConsensus() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "max PartialSignatureMessages for pre-consensus",
		Object:                maxPartialSignatureMessagesForPreConsensus(),
		ExpectedEncodedLength: maxSizePartialSignatureMessagesForPreConsensus,
		IsMaxSize:             false,
	}
}
