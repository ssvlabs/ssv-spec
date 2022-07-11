package dlog

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

type Proof struct {
	Commitment *bls.PublicKey
	PubKey     *bls.PublicKey
	Response   *bls.Fr
}

type Knowledge struct {
	SecretKey    *bls.SecretKey
	RandomNumber *bls.Fr
}

func (k *Knowledge) GetCommitment() *bls.PublicKey {
	r := bls.CastToSecretKey(k.RandomNumber)
	return r.GetPublicKey()
}

func (k *Knowledge) GetChallenge() *bls.Fr {
	generator := new(bls.PublicKey)
	bls.GetGeneratorOfPublicKey(generator)

	var challengeBytes []byte
	challengeBytes = append(challengeBytes, k.GetCommitment().Serialize()...)
	challengeBytes = append(challengeBytes, generator.Serialize()...)
	challengeBytes = append(challengeBytes, k.SecretKey.GetPublicKey().Serialize()...)
	challenge := new(bls.Fr)
	challenge.SetHashOf(challengeBytes)
	return challenge
}

func (k *Knowledge) Prove() Proof {
	response, cx := new(bls.Fr), new(bls.Fr)
	bls.FrMul(cx, k.GetChallenge(), bls.CastFromSecretKey(k.SecretKey))
	bls.FrAdd(response, k.RandomNumber, cx)

	return Proof{
		Commitment: k.GetCommitment(),
		PubKey:     k.SecretKey.GetPublicKey(),
		Response:   response,
	}
}

func (p *Proof) Verify() bool {
	generator := new(bls.PublicKey)
	bls.GetGeneratorOfPublicKey(generator)

	var challengeBytes []byte
	challengeBytes = append(challengeBytes, p.Commitment.Serialize()...)
	challengeBytes = append(challengeBytes, generator.Serialize()...)
	challengeBytes = append(challengeBytes, p.PubKey.Serialize()...)
	challenge := new(bls.Fr)
	challenge.SetHashOf(challengeBytes)
	cY := new(bls.G1)
	bls.G1Mul(cY, bls.CastFromPublicKey(p.PubKey), challenge)

	// s = r + cx
	// sG = Commitment + cY
	pkResponse := bls.CastToSecretKey(p.Response).GetPublicKey()
	expectedPkResponse := new(bls.PublicKey)
	expectedPkResponse.Deserialize(p.Commitment.Serialize())
	expectedPkResponse.Add(bls.CastToPublicKey(cY))
	return pkResponse.IsEqual(expectedPkResponse)
}
