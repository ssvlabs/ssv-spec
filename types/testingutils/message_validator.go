package testingutils

import (
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

func TestingMessageValidator() *validation.MessageValidator {
	return validation.NewMessageValidator(
		NewTestingBeaconNode().GetBeaconNetwork(),
		NewTestingNetworkDataFetcher(),
		NewTestingDutyFetcher(),
		NewMessageValidatorConfig(),
	)
}

type MessageValidatorConfig struct {
}

func NewMessageValidatorConfig() *MessageValidatorConfig {
	return &MessageValidatorConfig{}
}

func (c *MessageValidatorConfig) ProposerForHeightAndRound(height qbft.Height, round qbft.Round) types.OperatorID {
	return types.OperatorID(1)
}
