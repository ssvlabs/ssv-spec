package keygen

import "github.com/bloxapp/ssv-spec/dkg/base"

func (k *Keygen) r0Proceed() error {
	if k.Round != 0 {
		return ErrInvalidRound
	}

	msg := &ParsedMessage{
		Header: &base.MessageHeader{
			MsgType:  k.HandleMessageType,
			Sender:   k.PartyI,
			Receiver: 0,
		},
		Body: &KeygenMsgBody{
			Round1: &Round1Msg{
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
