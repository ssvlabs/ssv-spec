package timeout

import "github.com/bloxapp/ssv-spec/qbft/spectest"

// ErrorOnBroadcast tests calling UponRoundTimeout and having a broadcast error, should still change state regardless
func ErrorOnBroadcast() spectest.SpecTest {
	panic("implement")
}
