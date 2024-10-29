package validation

import (
	"time"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/maxmsgsize"
)

// Time spreads
const (
	LateMessageMargin     = time.Second * 3       // The duration past a message's TTL in which it is still considered valid
	ClockErrorTolerance   = time.Millisecond * 50 // The maximum amount of clock error we expect to see between nodes
	AllowedRoundsInFuture = 1
	LateSlotAllowance     = 2
)

// Ethereum values
const (
	SyncCommitteeSize = 512
)

// Allowed sizes
const (
	RsaSignatureSize     = 256
	PartialSignatureSize = 96
	MaxSignatures        = 13
	MaxSSVMessageData    = maxmsgsize.MaxSizeSSVMessageFromQBFTMessage
	MaxEncodedMsgSize    = maxmsgsize.MaxSizeSignedSSVMessageFromQBFTWith2Justification
)
