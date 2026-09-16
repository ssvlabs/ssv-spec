package testingutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

func TestBeaconBlockRoot(t *testing.T) {
	for _, v := range SupportedBlockVersions {
		// Gloas (ePBS §4) has no api.VersionedProposal form and no blinded variant: the block is
		// bid-only, so the "blinded" fixture is byte-identical to the regular one.
		if v == gloas.DataVersionGloas {
			require.EqualValues(t, TestingBeaconBlockBytesV(v), TestingBlindedBeaconBlockBytesV(v),
				fmt.Sprintf("%s, blinded fixture should equal the bid-only block", v.String()))
			continue
		}
		r1, _ := TestingBeaconBlockV(v).Root()
		r2, _ := TestingBlindedBeaconBlockV(v).Root()
		require.EqualValues(t, r1, r2, fmt.Sprintf("%s, hash root should be equal for both BeaconBlock and BlindedBeaconBlock", v.String()))
	}
}
