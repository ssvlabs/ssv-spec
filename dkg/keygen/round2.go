package keygen

import (
	"errors"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg/vss"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func (k *Keygen) r2Proceed() error {
	if k.Round != 2 {
		return ErrInvalidRound
	}
	ids := make([]*bls.Fr, k.PartyCount)
	for i := 0; i < int(k.PartyCount); i++ {
		id := new(bls.Fr)
		id.SetInt64(int64(i + 1))
		ids[i] = id
	}
	commitments, allShares, err := vss.Create(k.Coefficients, ids)
	if err != nil {
		return err
	}
	if len(commitments) != len(k.Coefficients) {
		return errors.New("invalid length of coefficients commitments")
	}
	if len(allShares) != int(k.PartyCount) {
		return errors.New("invalid length of shares")
	}
	yiComms := make([][]byte, len(k.Coefficients))

	for i, commitment := range commitments {
		yiComms[i] = commitment.Serialize()
	}

	for i, share := range allShares {
		receiver := uint16(i+1)
		if i+1 != int(k.PartyI) {
			msg := &Message{
				Sender:   k.PartyI,
				Receiver: &receiver,
				Body: MessageBody{
					Round3: &Round3Msg{
						Commitments: yiComms,
						Share:       share.Share.Serialize(),
					},
				},
			}
			k.pushOutgoing(msg)
		} else {
			k.ownShare = share.Share
		}
	}

	k.Round = 3
	return nil
}

func (k *Keygen) r2CanProceed() error {
	if k.Round != 2 {
		return ErrInvalidRound
	}
	for i, r2Msg := range k.Round2Msgs {
		r1Msg := k.Round1Msgs[i]
		if r1Msg == nil || r2Msg == nil || r1Msg.Body.Round1 == nil || r2Msg.Body.Round2 == nil {
			return errors.New("expected message not found")
		}
		if !VerifyYiCommitment(*r1Msg.Body.Round1, *r2Msg.Body.Round2, r2Msg.Sender) {
			// TODO: Handle blame?
			return fmt.Errorf("decomm doesn't match comm for party %d", r2Msg.Sender)
		}

	}
	return nil
}
