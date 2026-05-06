package testingutils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils/internal/proposerbls"
)

func TestBeaconBlockRoot(t *testing.T) {
	for _, v := range SupportedBlockVersions {
		r1, _ := TestingBeaconBlockV(v).Root()
		r2, _ := TestingBlindedBeaconBlockV(v).Root()
		require.EqualValues(t, r1, r2, fmt.Sprintf("%s, hash root should be equal for both BeaconBlock and BlindedBeaconBlock", v.String()))
	}
}

func TestProposerFixturesBLSDecodable(t *testing.T) {
	types.InitBLS()

	// Electra/Fulu fixtures use real Pectra-devnet-6 bytes and are already
	// strict-decodable; not re-tested here.
	visitor := func(t *testing.T) proposerbls.Visitor {
		return proposerbls.Visitor{
			Pubkey:    func(label string, b []byte) { assertBLSPublicKey(t, label, b) },
			Signature: func(label string, b []byte) { assertBLSSignature(t, label, b) },
		}
	}
	t.Run("capella", func(t *testing.T) {
		proposerbls.WalkCapella(TestingBeaconBlockCapella.Body, visitor(t))
	})
	t.Run("deneb", func(t *testing.T) {
		proposerbls.WalkDeneb(TestingBlockContentsDeneb.Block.Body, visitor(t))
	})
}

func assertBLSPublicKey(t *testing.T, label string, b []byte) {
	t.Helper()
	buf := bytes.Clone(b)
	var pk bls.PublicKey
	require.NoErrorf(t, pk.Deserialize(buf), "%s: %#x is not a valid compressed BLS public key", label, buf)
}

func assertBLSSignature(t *testing.T, label string, b []byte) {
	t.Helper()
	buf := bytes.Clone(b)
	var sig bls.Sign
	require.NoErrorf(t, sig.Deserialize(buf), "%s: %#x is not a valid compressed BLS signature", label, buf)
}
