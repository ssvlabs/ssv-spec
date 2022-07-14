package keygen

import (
	"errors"
	"github.com/bloxapp/ssv-spec/dkg/dlog"
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
		temp.Deserialize(r2Msg.Body.Round2.DeCommitment[0])
		pk.Add(temp)
	}
	vkVec := make([][]byte, int(k.PartyCount))
	for i, r4Msg := range k.Round4Msgs {
		vkVec[i] = r4Msg.Body.Round4.PubKey
	}
	k.Output = &LocalKeyShare{
		Index:           k.PartyI,
		Threshold:       uint32(len(k.Coefficients) - 1),
		ShareCount:      k.PartyCount,
		PublicKey:       pk.Serialize(),
		SecretShare:     k.skI.Serialize(),
		SharePublicKeys: vkVec,
	}
	k.Round = 5
	return nil
}

func (k *Keygen) r4CanProceed() error {
	if k.Round != 4 {
		return ErrInvalidRound
	}
	for _, r4Msg := range k.Round4Msgs {
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
