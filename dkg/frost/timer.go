package frost

import (
	"context"
	"sync/atomic"
	"time"
)

type RoundTimeoutFunc func(ProtocolRound) time.Duration
type OnTimeoutFn func() error

var defaultTimeout int = 100
var defaultTimeDuration = time.Millisecond

// RoundTimeout returns the number of seconds until next timeout for a give round
func RoundTimeout(ProtocolRound) time.Duration {
	return time.Duration(defaultTimeout) * defaultTimeDuration
}

// RoundTimer helps to manage current instance rounds.
type RoundTimer struct {
	ctx context.Context
	// cancelCtx cancels the current context, will be called from Kill()
	cancelCtx context.CancelFunc
	// timer is the underlying time.Timer
	timer *time.Timer
	// result holds the result of the timer
	done OnTimeoutFn
	// round is the current round of the timer
	round        int64
	roundTimeout RoundTimeoutFunc
}

// New creates a new instance of RoundTimer.
func NewRoundTimer(pctx context.Context, done OnTimeoutFn) *RoundTimer {
	ctx, cancelCtx := context.WithCancel(pctx)
	return &RoundTimer{
		ctx:          ctx,
		cancelCtx:    cancelCtx,
		timer:        nil,
		done:         done,
		roundTimeout: RoundTimeout,
	}
}

// OnTimeout sets a function called on timeout.
func (t *RoundTimer) OnTimeout(done OnTimeoutFn) {
	t.done = done
}

// Round returns a round.
func (t *RoundTimer) Round() ProtocolRound {
	return ProtocolRound(atomic.LoadInt64(&t.round))
}

// TimeoutForRound times out for a given round.
func (t *RoundTimer) TimeoutForRound(round ProtocolRound) {
	atomic.StoreInt64(&t.round, int64(round))
	timeout := t.roundTimeout(round)
	// preparing the underlying timer
	timer := t.timer
	if timer == nil {
		timer = time.NewTimer(timeout)
	} else {
		timer.Stop()
		// draining the channel of existing timer
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(timeout)
	// spawns a new goroutine to listen to the timer
	go t.waitForRound(round, timer.C)
}

func (t *RoundTimer) waitForRound(round ProtocolRound, timeout <-chan time.Time) {
	ctx, cancel := context.WithCancel(t.ctx)
	defer cancel()
	done := t.done
	select {
	case <-ctx.Done():
	case <-timeout:
		if t.Round() == round {
			if done != nil {
				_ = done()

			}
		}
	}
}
