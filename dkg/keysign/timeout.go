package keysign

import "github.com/bloxapp/ssv-spec/dkg"

func (sp *Instance) UponRoundTimeout() error {
	// if sp.state.GetCurrentRound() != sp.state.roundTimer.Round() {
	// 	return nil
	// }
	// prevRound := sp.state.GetCurrentRound()
	// sp.state.SetCurrentRound(Timeout)

	// msg := &ProtocolMsg{
	// 	Round: Timeout,
	// 	TimeoutMessage: &TimeoutMessage{
	// 		Round: prevRound,
	// 	},
	// }
	// bcastMsg, err := sp.saveSignedMsg(msg)
	// if err != nil {
	// 	return err
	// }
	// return sp.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
	return nil
}

func (sp *Instance) ProcessTimeoutMessage() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {
	return true, &dkg.ProtocolOutcome{}, nil
}
