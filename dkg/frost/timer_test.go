package frost

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/stretchr/testify/require"
)

func TestRoundTimer_TimeoutForRound(t *testing.T) {
	t.Run("TimeoutForRound", func(t *testing.T) {
		count := int32(0)
		onTimeout := func() error {
			atomic.AddInt32(&count, 1)
			return nil
		}
		timer := NewRoundTimer(context.Background(), onTimeout)
		timer.roundTimeout = func(round common.ProtocolRound) time.Duration {
			return 1100 * time.Millisecond
		}
		timer.StartRoundTimeoutTimer(common.ProtocolRound(1))
		require.Equal(t, int32(0), atomic.LoadInt32(&count))
		<-time.After(timer.roundTimeout(common.ProtocolRound(1)) + time.Millisecond*10)
		require.Equal(t, int32(1), atomic.LoadInt32(&count))
	})

	t.Run("timeout round before elapsed", func(t *testing.T) {
		count := int32(0)
		onTimeout := func() error {
			atomic.AddInt32(&count, 1)
			return nil
		}
		timer := NewRoundTimer(context.Background(), onTimeout)
		timer.roundTimeout = func(round common.ProtocolRound) time.Duration {
			return 1100 * time.Millisecond
		}

		timer.StartRoundTimeoutTimer(common.ProtocolRound(1))
		<-time.After(timer.roundTimeout(common.ProtocolRound(1)) / 2)
		timer.StartRoundTimeoutTimer(common.ProtocolRound(2)) // reset before elapsed
		require.Equal(t, int32(0), atomic.LoadInt32(&count))
		<-time.After(timer.roundTimeout(common.ProtocolRound(2)) + time.Millisecond*10)
		require.Equal(t, int32(1), atomic.LoadInt32(&count))
	})
}
