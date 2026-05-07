package testingutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeaconBlockRoot(t *testing.T) {
	for _, v := range SupportedBlockVersions {
		r1, _ := TestingBeaconBlockV(v).Root()
		r2, _ := TestingBlindedBeaconBlockV(v).Root()
		require.EqualValues(t, r1, r2, fmt.Sprintf("%s, hash root should be equal for both BeaconBlock and BlindedBeaconBlock", v.String()))
	}
}
