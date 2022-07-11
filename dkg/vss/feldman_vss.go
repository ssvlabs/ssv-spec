package vss

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type (
	Share struct {
		Threshold int
		ID        *bls.Fr
		Share     *bls.Fr
	}
	Commitments  []*bls.PublicKey
	Shares       []*Share
	Coefficients []bls.Fr
)

var (
	ErrNumSharesBelowThreshold = fmt.Errorf("not enough shares to satisfy the threshold")
)

// Returns an array ai for f(x) = a0 + a1*x^1 + ... + a{n-1}*x^{n-1}
// The first coefficient a0 will be used as ui in Feldman VSS
func CreatePolynomial(length int) Coefficients {
	v := make([]bls.Fr, length)
	for i := 0; i < length; i++ {
		fe := bls.Fr{}
		fe.SetByCSPRNG()
		v[i] = fe
	}
	return v
}

// Returns a new array of secret shares created by Shamir's Secret Sharing Algorithm,
// requiring a minimum number of shares to recreate, of length shares, from the input secret
//
func Create(poly Coefficients, indexes []*bls.Fr) (Commitments, Shares, error) {
	if poly == nil || indexes == nil {
		return nil, nil, fmt.Errorf("polynomial or indexes == nil: %v %v", poly, indexes)
	}
	if len(poly) < 2 {
		return nil, nil, errors.New("invalid polynomial")
	}

	threshold := len(poly) - 1
	ids, err := CheckIndexes(indexes)
	if err != nil {
		return nil, nil, err
	}

	num := len(indexes)
	if num < len(poly)-1 {
		return nil, nil, ErrNumSharesBelowThreshold
	}

	v := make(Commitments, len(poly))
	for i, ai := range poly {
		v[i] = bls.CastToSecretKey(&ai).GetPublicKey()
	}

	shares := make(Shares, num)
	for i := 0; i < num; i++ {
		share := new(bls.Fr)
		bls.FrEvaluatePolynomial(share, poly, ids[i])
		shares[i] = &Share{Threshold: threshold, ID: ids[i], Share: share}
	}
	return v, shares, nil
}

// Check share ids of Shamir's Secret Sharing, return error if duplicate or 0 value found
func CheckIndexes(indexes []*bls.Fr) ([]*bls.Fr, error) {

	visited := make(map[string]struct{})
	for _, v := range indexes {
		if v.IsZero() {
			return nil, errors.New("party index should not be 0")
		}
		vModStr := hex.EncodeToString(v.Serialize())
		if _, ok := visited[vModStr]; ok {
			return nil, fmt.Errorf("duplicate indexes %s", vModStr)
		}
		visited[vModStr] = struct{}{}
	}
	return indexes, nil
}

func (share *Share) Verify(threshold int, vs Commitments) bool {
	if share.Threshold != threshold || vs == nil {
		return false
	}

	xi := new(bls.Fr)
	xi.SetInt64(int64(1))
	v := &bls.PublicKey{}
	v.Deserialize(vs[0].Serialize())

	for j := 1; j <= threshold; j++ {
		vi := new(bls.G1)
		vi.Deserialize(vs[j].Serialize())
		bls.FrMul(xi, xi, share.ID)
		bls.G1Mul(vi, vi, xi)
		v.Add(bls.CastToPublicKey(vi))
	}

	return bls.CastToSecretKey(share.Share).GetPublicKey().IsEqual(v)
}

func (shares Shares) ReConstruct() (secret *bls.Fr, err error) {
	if shares != nil && shares[0].Threshold > len(shares) {
		return nil, ErrNumSharesBelowThreshold
	}
	//
	//// x coords
	xs := make([]bls.Fr, shares[0].Threshold)
	indices := make([]bls.Fr, shares[0].Threshold)
	for i, share := range shares {
		if i >= shares[0].Threshold {
			break
		}
		err := indices[i].Deserialize(share.ID.Serialize())
		if err != nil {
			return nil, err
		}
		err = xs[i].Deserialize(share.Share.Serialize())
		if err != nil {
			return nil, err
		}

	}
	secret = new(bls.Fr)
	err = bls.FrLagrangeInterpolation(secret, indices, xs)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
