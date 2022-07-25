package keygen

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg/algorithms/vss"
	"github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func (k *Keygen) r2Proceed() error {
	if k.Round != 2 {
		return ErrInvalidRound
	}
	ids := make([]*bls.Fr, k.PartyCount)
	for i, id := range k.Committee {
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, id)
		fe := new(bls.Fr)
		fe.SetLittleEndian(bytes)
		ids[i] = fe
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
		receiver := k.Committee[i]
		if i+1 != int(k.PartyI) {
			msg := &ParsedMessage{
				Header: &types.MessageHeader{
					SessionId: k.SessionID,
					MsgType:   k.HandleMessageType,
					Sender:    k.PartyI,
					Receiver:  receiver,
				},
				Body: &KeygenMsgBody{
					Round3: &Round3Msg{
						Share: share.Share.Serialize(),
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
	for _, id := range k.Committee {
		r1Msg := k.Round1Msgs[id]
		r2Msg := k.Round2Msgs[id]
		if r1Msg == nil || r2Msg == nil || r1Msg.Body.Round1 == nil || r2Msg.Body.Round2 == nil {
			return ErrExpectMessage
		}
		if !k.VerifyCommitment(*r1Msg.Body.Round1, *r2Msg.Body.Round2, r2Msg.Header.Sender) {
			// TODO: Handle blame?
			return fmt.Errorf("decomm doesn't match comm for party %d", r2Msg.Header.Sender)
		}
	}
	return nil
}
