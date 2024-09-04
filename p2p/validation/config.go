package validation

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Protocol configuration for the MessageValidator
type Config interface {
	ProposerForHeightAndRound(height qbft.Height, round qbft.Round) types.OperatorID
}
