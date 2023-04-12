package comparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/qbft"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareQBFTState(source, target *qbft.State) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if shareDiff := testingutilscomparable.NestedCompare("Share", testingutilscomparable.CompareShare(source.Share, target.Share)); len(shareDiff) > 0 {
		ret = append(ret, shareDiff)
	}

	if !bytes.Equal(source.ID, target.ID) {
		ret = append(ret, testingutilscomparable.Sprintf("ID", "source %x != target %x", source.ID, target.ID))
	}

	if source.Round != target.Round {
		ret = append(ret, testingutilscomparable.Sprintf("Round", "source %d != target %d", source.Round, target.Round))
	}

	if source.Height != target.Height {
		ret = append(ret, testingutilscomparable.Sprintf("Height", "source %d != target %d", source.Height, target.Height))
	}

	if source.LastPreparedRound != target.LastPreparedRound {
		ret = append(ret, testingutilscomparable.Sprintf("LastPreparedRound", "source %d != target %d", source.LastPreparedRound, target.LastPreparedRound))
	}

	if !bytes.Equal(source.LastPreparedValue, target.LastPreparedValue) {
		ret = append(ret, testingutilscomparable.Sprintf("LastPreparedValue", "source %x != target %x", source.LastPreparedValue, target.LastPreparedValue))
	}

	if diff := testingutilscomparable.NestedCompare("ProposalAcceptedForCurrentRound", CompareQBFTSignedMessage(source.ProposalAcceptedForCurrentRound, target.ProposalAcceptedForCurrentRound)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if source.Decided != target.Decided {
		ret = append(ret, testingutilscomparable.Sprintf("Decided", "source %t != target %t", source.Decided, target.Decided))
	}

	if !bytes.Equal(source.DecidedValue, target.DecidedValue) {
		ret = append(ret, testingutilscomparable.Sprintf("DecidedValue", "source %x != target %x", source.DecidedValue, target.DecidedValue))
	}

	if diff := testingutilscomparable.NestedCompare("ProposeContainer", CompareQBFTMessageContainer(source.ProposeContainer, target.ProposeContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("PrepareContainer", CompareQBFTMessageContainer(source.PrepareContainer, target.PrepareContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("CommitContainer", CompareQBFTMessageContainer(source.CommitContainer, target.CommitContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("RoundChangeContainer", CompareQBFTMessageContainer(source.RoundChangeContainer, target.RoundChangeContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	r1, err := source.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	r2, err := target.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	if !bytes.Equal(r1[:], r2[:]) {
		ret = append(ret, testingutilscomparable.Sprintf("Root", "source %x != target %x", r1, r2))
	}

	return ret
}
