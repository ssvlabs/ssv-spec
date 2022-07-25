package keygen

import "github.com/bloxapp/ssv-spec/dkg/types"

func (k *Keygen) r1Proceed() error {
	if k.Round != 1 {
		return ErrInvalidRound
	}
	msg := &ParsedMessage{
		Header: &types.MessageHeader{
			SessionId: k.SessionID,
			MsgType:   k.HandleMessageType,
			Sender:    k.PartyI,
			Receiver:  0,
		},
		Body: &KeygenMsgBody{
			Round2: &Round2Msg{
				Decommitment: k.GetDecommitment(),
				BlindFactor:  k.BlindFactor[:],
			},
		},
	}
	k.pushOutgoing(msg)
	k.Round = 2
	return nil
}

func (k *Keygen) r1CanProceed() error {
	if k.Round != 1 {
		return ErrInvalidRound
	}
	for _, id := range k.Committee {
		if msg := k.Round1Msgs[id]; msg == nil {
			return ErrExpectMessage
		}
	}
	return nil
}
