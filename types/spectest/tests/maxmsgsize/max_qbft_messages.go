package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	MaxSizeQBFTMessageWithNoJustification = 132
	MaxSizeQBFTMessageWith1Justification  = 48284
	MaxSizeQBFTMessageWith2Justification  = 722412
)

func maxFullData() []byte {
	fullData := [MaxSizeFullConsensusData]byte{}
	return fullData[:]
}

func maxQbftMessageNoJustification() *qbft.Message {
	msgID := [56]byte{1}

	qbftMsg := &qbft.Message{
		MsgType:                  qbft.PrepareMsgType,
		Height:                   1,
		Round:                    1,
		Identifier:               msgID[:],
		Root:                     [32]byte{},
		DataRound:                1,
		RoundChangeJustification: make([][]byte, 0),
		PrepareJustification:     make([][]byte, 0),
	}

	return qbftMsg
}

func maxQbftMessageWith1Justification() *qbft.Message {

	justification := make([]*types.SignedSSVMessage, 0)
	for i := 0; i < 13; i++ {
		justification = append(justification, maxSignedSSVMessageFromObject(maxQbftMessageNoJustification()))
	}

	justificationBytes, err := qbft.MarshalJustifications(justification)
	if err != nil {
		panic(err)
	}

	qbftMsg := maxQbftMessageNoJustification()
	qbftMsg.RoundChangeJustification = justificationBytes

	return qbftMsg
}

func maxQbftMessageWith2Justification() *qbft.Message {

	justification1 := make([]*types.SignedSSVMessage, 0)
	for i := 0; i < 13; i++ {
		justification1 = append(justification1, maxSignedSSVMessageFromObject(maxQbftMessageNoJustification()))
	}

	justification1Bytes, err := qbft.MarshalJustifications(justification1)
	if err != nil {
		panic(err)
	}

	justification2 := make([]*types.SignedSSVMessage, 0)
	for i := 0; i < 13; i++ {
		justification2 = append(justification2, maxSignedSSVMessageFromObject(maxQbftMessageWith1Justification()))
	}

	justification2Bytes, err := qbft.MarshalJustifications(justification2)
	if err != nil {
		panic(err)
	}

	qbftMsg := maxQbftMessageNoJustification()
	qbftMsg.PrepareJustification = justification1Bytes
	qbftMsg.RoundChangeJustification = justification2Bytes

	return qbftMsg
}

func MaxQBFTMessageWithNoJustification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max qbftMessage with no justification",
		testdoc.StructureSizeTestMaxQBFTMessageWithNoJustificationDoc,
		maxQbftMessageNoJustification(),
		MaxSizeQBFTMessageWithNoJustification,
		false,
	)
}

func MaxQBFTMessageWith1Justification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max qbftMessage with 1 justification",
		testdoc.StructureSizeTestMaxQBFTMessageWith1JustificationDoc,
		maxQbftMessageWith1Justification(),
		MaxSizeQBFTMessageWith1Justification,
		false,
	)
}

func MaxQBFTMessageWith2Justification() *StructureSizeTest {
	return NewStructureSizeTest(
		"max qbftMessage with 2 justifications",
		testdoc.StructureSizeTestMaxQBFTMessageWith2JustificationDoc,
		maxQbftMessageWith2Justification(),
		MaxSizeQBFTMessageWith2Justification,
		true,
	)
}
