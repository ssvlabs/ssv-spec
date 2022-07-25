package vss_test

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/gg20/algorithms/vss"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)
}

func TestCheckIndexesDup(t *testing.T) {
	indexes := make([]*bls.Fr, 0)
	for i := 0; i < 1000; i++ {
		fe := new(bls.Fr)
		fe.SetByCSPRNG()
		indexes = append(indexes, fe)
	}
	_, e := vss.CheckIndexes(indexes)
	require.NoError(t, e)

	indexes = append(indexes, indexes[99])
	_, e = vss.CheckIndexes(indexes)
	require.Error(t, e)
}

func TestCheckIndexesZero(t *testing.T) {
	indexes := make([]*bls.Fr, 0)
	for i := 0; i < 1000; i++ {
		fe := new(bls.Fr)
		fe.SetByCSPRNG()
		indexes = append(indexes, fe)
	}
	_, e := vss.CheckIndexes(indexes)
	require.NoError(t, e)
}

func TestCreate(t *testing.T) {
	num, threshold := 5, 3

	poly := vss.CreatePolynomial(threshold + 1)

	ids := make([]*bls.Fr, 0)
	for i := 0; i < num; i++ {
		id := new(bls.Fr)
		id.SetByCSPRNG()
		ids = append(ids, id)
	}

	vs, _, err := vss.Create(poly, ids)
	require.Nil(t, err)

	require.Equal(t, threshold+1, len(vs))

	require.Equal(t, threshold+1, len(vs))

	// ensure that each vs has two points on the curve
	for _, v := range vs {
		require.False(t, v.IsZero())
		require.True(t, v.IsValidOrder())
	}
}

func TestVerify(t *testing.T) {
	num, threshold := 4, 2

	poly := vss.CreatePolynomial(threshold + 1)

	ids := make([]*bls.Fr, 0)
	for i := 0; i < num; i++ {
		id := new(bls.Fr)
		id.SetByCSPRNG()
		ids = append(ids, id)
	}

	vs, shares, err := vss.Create(poly, ids)
	require.NoError(t, err)

	var vsStr, shareStr []string
	for _, v := range vs {
		vsStr = append(vsStr, v.SerializeToHexStr())
	}
	for _, share := range shares {
		shareStr = append(shareStr, share.ID.GetString(10)+": "+hex.EncodeToString(share.Share.Serialize()))
	}

	for i := num - 1; i >= 0; i-- {
		require.True(t, shares[i].Verify(threshold, vs))
	}
}

func TestReconstruct(t *testing.T) {
	num, threshold := 5, 3

	poly := vss.CreatePolynomial(threshold + 1)

	ids := make([]*bls.Fr, 0)
	for i := 0; i < num; i++ {
		id := new(bls.Fr)
		id.SetByCSPRNG()
		ids = append(ids, id)
	}

	_, shares, err := vss.Create(poly, ids)
	require.NoError(t, err)

	secret2, err2 := shares[:threshold-1].ReConstruct()
	require.Error(t, err2) // not enough shares to satisfy the threshold
	require.Nil(t, secret2)

	secret3, err3 := shares[:threshold].ReConstruct()
	require.NoError(t, err3)
	require.False(t, secret3.IsZero())
	require.True(t, secret3.IsValid())

	secret4, err4 := shares[:num].ReConstruct()
	require.NoError(t, err4)
	require.False(t, secret4.IsZero())
	require.True(t, secret4.IsValid())
}
