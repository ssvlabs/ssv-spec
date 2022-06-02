package qbft

type Timer interface {
	// TimeoutForRound will reset running timer if exists and will start a new timer for a specific round
	TimeoutForRound(round Round)
}

//func uponRoundTimeout(state *State, Config IConfig) error {
//	state.Round++
//	roundChange, err := CreateRoundChange(state, Config, state.Round)
//	if err != nil {
//		return errors.Wrap(err, "could not generate round change msg")
//	}
//
//	if err := Config.GetNetwork().Broadcast(roundChange); err != nil {
//		return errors.Wrap(err, "failed to broadcast round change message")
//	}
//
//	return nil
//}
