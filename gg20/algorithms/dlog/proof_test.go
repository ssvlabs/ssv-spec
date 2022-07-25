package dlog_test

import (
	"github.com/bloxapp/ssv-spec/gg20/algorithms/dlog"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)
}

func TestProof(t *testing.T) {
	secret := new(bls.SecretKey)
	secret.SetByCSPRNG()
	r := new(bls.Fr)
	r.SetByCSPRNG()
	knowledge := dlog.Knowledge{
		SecretKey:    secret,
		RandomNumber: r,
	}
	proof := knowledge.Prove()
	verified := proof.Verify()
	require.True(t, verified)

	// Tweaked response should fail
	one := new(bls.Fr)
	one.SetInt64(int64(1))
	bls.FrAdd(proof.Response, proof.Response, one)
	verified = proof.Verify()
	require.False(t, verified)
}
