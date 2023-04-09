package comparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareQBFTSignedMessage(source, target *qbft.SignedMessage) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if !bytes.Equal(source.Signature, target.Signature) {
		ret = append(ret, testingutilscomparable.Sprintf("Signature source/ target %x <---> %x", source.Signature, target.Signature))
	}

	if diff := testingutilscomparable.NestedCompare("Signers", CompareQBFTMessageSigners(source.Signers, target.Signers)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("Message", CompareQBFTMessage(source.Message, target.Message)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if !bytes.Equal(source.FullData, target.FullData) {
		ret = append(ret, testingutilscomparable.Sprintf("FullData source/ target %x <---> %x", source.FullData, target.FullData))
	}

	return ret
}

func CompareQBFTMessageSigners(source, target []types.OperatorID) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, testingutilscomparable.Sprintf("Committee length source/ target %d <---> %d", len(source), len(target)))
	}

	for i := range source {
		if source[i] != target[i] {
			ret = append(ret, testingutilscomparable.Sprintf("OperatorID source/ target %d <---> %d", source[i], target[i]))
		}
	}

	return ret
}

func CompareQBFTMessageJustifications(source, target [][]byte) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, testingutilscomparable.Sprintf("Committee length source/ target %d <---> %d", len(source), len(target)))
	}

	for i := range source {
		if !bytes.Equal(source[i], target[i]) {
			ret = append(ret, testingutilscomparable.Sprintf("Bytes %d source/ target %d <---> %d", i, source[i], target[i]))
		}
	}

	return ret
}

func CompareQBFTMessage(source, target qbft.Message) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if source.MsgType != target.MsgType {
		ret = append(ret, testingutilscomparable.Sprintf("MsgType source/ target %d <---> %d", source.MsgType, target.MsgType))
	}

	if source.Round != target.Round {
		ret = append(ret, testingutilscomparable.Sprintf("Round source/ target %d <---> %d", source.Round, target.Round))
	}

	if source.Height != target.Height {
		ret = append(ret, testingutilscomparable.Sprintf("Height source/ target %d <---> %d", source.Height, target.Height))
	}

	if !bytes.Equal(source.Identifier, target.Identifier) {
		ret = append(ret, testingutilscomparable.Sprintf("Identifier source/ target %x <---> %x", source.Identifier, target.Identifier))
	}

	if !bytes.Equal(source.Root[:], target.Root[:]) {
		ret = append(ret, testingutilscomparable.Sprintf("Root source/ target %x <---> %x", source.Root, target.Root))
	}

	if source.DataRound != target.DataRound {
		ret = append(ret, testingutilscomparable.Sprintf("DataRound source/ target %d <---> %d", source.DataRound, target.DataRound))
	}

	if diff := testingutilscomparable.NestedCompare("RoundChangeJustification", CompareQBFTMessageJustifications(source.RoundChangeJustification, target.RoundChangeJustification)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("PrepareJustification", CompareQBFTMessageJustifications(source.PrepareJustification, target.PrepareJustification)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	return ret
}
