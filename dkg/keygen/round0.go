package keygen

func (k *Keygen) r0Proceed() error {
	if k.Round != 0 {
		return ErrInvalidRound
	}

	msg := &Message{
		Sender: k.PartyI,
		Body: MessageBody{
			Round1: &Round1Msg{
				Commitment: k.GetYiCommitment(),
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
