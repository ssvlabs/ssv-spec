package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	expectedSizePrepareQBFTMessage     = 132
	expectedSizeCommitQBFTMessage      = 132
	expectedSizeRoundChangeQBFTMessage = 1596
	expectedSizeProposalQBFTMessage    = 7452
)

func expectedPrepare() *qbft.Message {
	msgID := [56]byte{1}

	qbftMsg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     1,
		Round:      1,
		Identifier: msgID[:],
		Root:       [32]byte{},
	}

	return qbftMsg
}

func expectedCommit() *qbft.Message {
	msgID := [56]byte{1}

	qbftMsg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     1,
		Round:      1,
		Identifier: msgID[:],
		Root:       [32]byte{},
	}

	return qbftMsg
}

func expectedRoundChange(quorum int) *qbft.Message {
	msgID := [56]byte{1}

	rcJustification := make([]*types.SignedSSVMessage, 0)

	for i := 0; i < quorum; i++ {
		rcJustification = append(rcJustification, expectedSignedSSVMessageFromObject(expectedPrepare(), 1))
	}

	rcJustificationBytes, err := qbft.MarshalJustifications(rcJustification)
	if err != nil {
		panic(err)
	}

	qbftMsg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   1,
		Round:                    2,
		Identifier:               msgID[:],
		Root:                     [32]byte{},
		DataRound:                1,
		RoundChangeJustification: rcJustificationBytes,
	}

	return qbftMsg
}

func expectedProposal(quorum int) *qbft.Message {
	msgID := [56]byte{1}

	rcJustification := make([]*types.SignedSSVMessage, 0)

	for i := 0; i < quorum; i++ {
		rcJustification = append(rcJustification, expectedSignedSSVMessageFromObject(expectedRoundChange(quorum), 1))
	}

	rcJustificationBytes, err := qbft.MarshalJustifications(rcJustification)
	if err != nil {
		panic(err)
	}

	prepareJustification := make([]*types.SignedSSVMessage, 0)

	for i := 0; i < quorum; i++ {
		prepareJustification = append(prepareJustification, expectedSignedSSVMessageFromObject(expectedPrepare(), 1))
	}

	prepareJustificationBytes, err := qbft.MarshalJustifications(prepareJustification)
	if err != nil {
		panic(err)
	}

	qbftMsg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   1,
		Round:                    2,
		Identifier:               msgID[:],
		Root:                     [32]byte{},
		DataRound:                1,
		RoundChangeJustification: rcJustificationBytes,
		PrepareJustification:     prepareJustificationBytes,
	}

	return qbftMsg
}

func ExpectedPrepareQBFTMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected prepare qbftMessage",
		"Test the expected size of a prepare QBFT message with no justifications",
		expectedPrepare(),
		expectedSizePrepareQBFTMessage,
		false,
	)
}

func ExpectedCommitQBFTMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected commit qbftMessage",
		"Test the expected size of a commit QBFT message with no justifications",
		expectedCommit(),
		expectedSizeCommitQBFTMessage,
		false,
	)
}

func ExpectedRoundChangeQBFTMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected round change qbftMessage",
		"Test the expected size of a round change QBFT message with 3 round change justifications",
		expectedRoundChange(3),
		expectedSizeRoundChangeQBFTMessage,
		false,
	)
}

func ExpectedProposalQBFTMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected proposal qbftMessage",
		"Test the expected size of a proposal QBFT message with 3 round change and 3 prepare justifications",
		expectedProposal(3),
		expectedSizeProposalQBFTMessage,
		false,
	)
}
