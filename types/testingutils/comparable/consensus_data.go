package testingutilscomparable

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
)

func comparePreConsensusJustifications(source, target []*types.SignedPartialSignatureMessage) []Difference {
	ret := make([]Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, Sprintf("Length", "source %d != target %d", len(source), len(target)))
		return ret
	}

	for i := range source {
		if diff := NestedCompare(fmt.Sprintf("Message %d", i), CompareSignedPartialSignatureMessage(source[i], target[i])); len(diff) > 0 {
			ret = append(ret, diff)
		}
	}

	return ret
}

func CompareConsensusData(source, target *types.ConsensusData) []Difference {
	ret := make([]Difference, 0)

	if (source == nil && target != nil) || (source != nil && target == nil) {
		ret = append(ret, Sprintf("ConsensusData nil?", "source %t != target %t", source == nil, target == nil))
		return ret
	}

	if diff := NestedCompare("Duty", CompareDuty(&source.Duty, &target.Duty)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if source.Version != target.Version {
		ret = append(ret, Sprintf("Version", "source %d != target %d", source.Version, target.Version))
	}

	if diff := NestedCompare("PreConsensusJustifications", comparePreConsensusJustifications(source.PreConsensusJustifications, target.PreConsensusJustifications)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if !bytes.Equal(source.DataSSZ, target.DataSSZ) {
		ret = append(ret, Sprintf("DataSSZ", "source %x != target %x", source.DataSSZ, target.DataSSZ))
	}

	return ret
}
