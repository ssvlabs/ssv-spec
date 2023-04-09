package testingutilscomparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
)

func CompareShare(source, target *types.Share) []Difference {
	ret := make([]Difference, 0)

	if source.OperatorID != target.OperatorID {
		ret = append(ret, Sprintf("OperatorID source/ target %d <---> %d", source.OperatorID, target.OperatorID))
	}

	if !bytes.Equal(source.ValidatorPubKey, target.ValidatorPubKey) {
		ret = append(ret, Sprintf("ValidatorPubKey source/ target %x <---> %x", source.ValidatorPubKey, target.ValidatorPubKey))
	}

	if !bytes.Equal(source.SharePubKey, target.SharePubKey) {
		ret = append(ret, Sprintf("SharePubKey source/ target %x <---> %x", source.SharePubKey, target.SharePubKey))
	}

	if shareDiff := NestedCompare("Committee", CompareCommittee(source.Committee, target.Committee)); len(shareDiff) > 0 {
		ret = append(ret, shareDiff)
	}

	if source.Quorum != target.Quorum {
		ret = append(ret, Sprintf("Quorum source/ target %d <---> %d", source.Quorum, target.Quorum))
	}

	if source.PartialQuorum != target.PartialQuorum {
		ret = append(ret, Sprintf("PartialQuorum source/ target %d <---> %d", source.PartialQuorum, target.PartialQuorum))
	}

	if !bytes.Equal(source.DomainType[:], target.DomainType[:]) {
		ret = append(ret, Sprintf("DomainType source/ target %x <---> %x", source.DomainType, target.DomainType))
	}

	if !bytes.Equal(source.FeeRecipientAddress[:], target.FeeRecipientAddress[:]) {
		ret = append(ret, Sprintf("FeeRecipientAddress source/ target %x <---> %x", source.FeeRecipientAddress, target.FeeRecipientAddress))
	}

	if !bytes.Equal(source.Graffiti, target.Graffiti) {
		ret = append(ret, Sprintf("Graffiti source/ target %x <---> %x", source.Graffiti, target.Graffiti))
	}

	return ret
}
