package keygen

import (
	"errors"
	"github.com/bloxapp/ssv-spec/dkg/algorithms/dlog"
	"github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/gogo/protobuf/sortkeys"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var (
	ErrInvalidProof = errors.New("invalid proof")
)

func (k *Keygen) r4Proceed() error {
	if k.Round != 4 {
		return ErrInvalidRound
	}

	pk := new(bls.PublicKey)
	for _, r2Msg := range k.Round2Msgs {
		temp := new(bls.PublicKey)
		temp.Deserialize(r2Msg.Body.Round2.Decommitment[0])
		pk.Add(temp)
	}
	var vkVec [][]byte
	sortkeys.Uint64s(k.Committee)
	for _, id := range k.Committee {
		r4Msg := k.Round4Msgs[id]
		vkVec = append(vkVec, r4Msg.Body.Round4.PubKey)
	}
	k.Output = &types.LocalKeyShare{
		Index:           k.PartyI,
		Threshold:       uint64(len(k.Coefficients) - 1),
		PublicKey:       pk.Serialize(),
		SecretShare:     k.skI.Serialize(),
		Committee:       k.Committee,
		SharePublicKeys: vkVec,
	}
	k.Round = 5
	return nil
}

func (k *Keygen) r4CanProceed() error {
	if k.Round != 4 {
		return ErrInvalidRound
	}
	for _, id := range k.Committee {
		r4Msg := k.Round4Msgs[id]
		if r4Msg == nil || r4Msg.Body.Round4 == nil {
			return ErrExpectMessage
		}
		proof := &dlog.Proof{
			Commitment: new(bls.PublicKey),
			PubKey:     new(bls.PublicKey),
			Response:   new(bls.Fr),
		}
		proof.Commitment.Deserialize(r4Msg.Body.Round4.Commitment)
		proof.PubKey.Deserialize(r4Msg.Body.Round4.PubKey)
		proof.Response.Deserialize(r4Msg.Body.Round4.ChallengeResponse)
		if !proof.Verify() {
			return ErrInvalidProof
		}
	}
	return nil
}
