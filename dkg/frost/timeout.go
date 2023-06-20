package frost

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
)

func (fr *Instance) UponRoundTimeout() error {
	if fr.state.GetCurrentRound() != fr.state.roundTimer.Round() {
		return nil
	}
	prevRound := fr.state.GetCurrentRound()
	fr.state.SetCurrentRound(common.Timeout)

	msg := &ProtocolMsg{
		Round: common.Timeout,
		TimeoutMessage: &TimeoutMessage{
			Round: prevRound,
		},
	}
	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return err
	}
	return fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}

func (fr *Instance) ProcessTimeoutMessage() (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {
	return true, &dkg.ProtocolOutcome{}, nil
}
