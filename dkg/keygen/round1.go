package keygen

func (k *Keygen) r1Proceed() error {
	if k.Round != 1 {
		return ErrInvalidRound
	}
	msg := &Message{
		Sender: k.PartyI,
		Body: MessageBody{
			Round2: &Round2Msg{
				DeCommmitment: k.GetDecommitment(),
				BlindFactor:   k.BlindFactor[:],
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
	for _, message := range k.Round1Msgs {
		if message == nil {
			return ErrExpectMessage
		}
	}
	return nil
}
