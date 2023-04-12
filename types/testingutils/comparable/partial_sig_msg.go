package testingutilscomparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
)

func comparePartialSignatureMessageSlice(source, target []*types.PartialSignatureMessage) []Difference {
	ret := make([]Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, Sprintf("Length", "source %d != target %d", len(source), len(target)))
		return ret
	}

	for i := range source {
		if !bytes.Equal(source[i].PartialSignature, target[i].PartialSignature) {
			ret = append(ret, Sprintf("PartialSignature", "source %x != target %x", source[i].PartialSignature, target[i].PartialSignature))
		}

		if !bytes.Equal(source[i].SigningRoot[:], target[i].SigningRoot[:]) {
			ret = append(ret, Sprintf("SigningRoot", "source %x != target %x", source[i].SigningRoot, target[i].SigningRoot))
		}

		if source[i].Signer != target[i].Signer {
			ret = append(ret, Sprintf("Signer", "source %d != target %d", source[i].Signer, target[i].Signer))
		}
	}

	return ret
}

func comparePartialSignatureMessages(source, target *types.PartialSignatureMessages) []Difference {
	ret := make([]Difference, 0)

	if source.Type != target.Type {
		ret = append(ret, Sprintf("Type", "source %d != target %d", source.Type, target.Type))
	}

	if source.Slot != target.Slot {
		ret = append(ret, Sprintf("Slot", "source %d != target %d", source.Slot, target.Slot))
	}

	if diff := NestedCompare("Messages", comparePartialSignatureMessageSlice(source.Messages, target.Messages)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	return ret
}

func CompareSignedPartialSignatureMessage(source, target *types.SignedPartialSignatureMessage) []Difference {
	ret := make([]Difference, 0)

	if diff := NestedCompare("Message", comparePartialSignatureMessages(&source.Message, &target.Message)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if !bytes.Equal(source.Signature, target.Signature) {
		ret = append(ret, Sprintf("Signature", "source %x != target %x", source.Signature, target.Signature))
	}

	if source.Signer != target.Signer {
		ret = append(ret, Sprintf("Signer", "source %d != target %d", source.Signer, target.Signer))
	}

	return ret
}
