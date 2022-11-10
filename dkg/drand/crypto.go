package drand

import (
	"github.com/drand/kyber"
	"github.com/drand/kyber/share"
	"github.com/drand/kyber/share/dkg"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func pubKeysFromPolyEvaluation(suite dkg.Suite, commitments []kyber.Point, evalPoints []uint32) ([]*bls.PublicKey, error) {
	exp := share.NewPubPoly(suite, suite.Point().Base(), commitments)
	ret := make([]*bls.PublicKey, 0)
	for _, x := range evalPoints {
		pubshare := exp.Eval(int(x))
		byts, err := pubshare.V.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal pub share")
		}

		pk := &bls.PublicKey{}
		if err := pk.Deserialize(byts); err != nil {
			return nil, errors.Wrap(err, "could not deserialize pubkey")
		}
		ret = append(ret, pk)
	}
	return ret, nil
}
