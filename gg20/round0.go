package gg20

import (
	"github.com/bloxapp/ssv-spec/dkg/types"
	types2 "github.com/bloxapp/ssv-spec/gg20/types"
)

func (k *Keygen) r0Proceed() error {
	if k.Round != 0 {
		return ErrInvalidRound
	}

	msg := &types2.ParsedMessage{
		Header: &types.MessageHeader{
			MsgType:  k.HandleMessageType,
			Sender:   k.PartyI,
			Receiver: 0,
		},
		Body: &types2.KeygenMsgBody{
			Round1: &types2.Round1Msg{
				Commitment: k.GetCommitment(),
			},
		},
	}
	k.pushOutgoing(msg)
	k.Round = 1
	return nil
}

func (k *Keygen) r0CanProceed() error {
	if k.Round != 0 {
		return ErrInvalidRound
	}
	return nil
}
