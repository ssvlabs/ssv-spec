package queue

import (
	"testing"

	"github.com/bloxapp/ssv-spec/queue"
	"github.com/stretchr/testify/require"
)

type MessagePriorityPairTest struct {
	Name  string
	State *queue.State

	// Messages is the correctly-ordered slice of messages.
	A, B mockMessage

	// Prior determines whether A is expected to be prior to B.
	Prior bool
}

func (test *MessagePriorityPairTest) TestName() string {
	return test.Name
}

func (test *MessagePriorityPairTest) Run(t *testing.T) {
	a, err := queue.DecodeSSVMessage(test.A.toSSVMessage(test.State))
	require.NoError(t, err)
	b, err := queue.DecodeSSVMessage(test.A.toSSVMessage(test.State))
	require.NoError(t, err)

	prioritizer := queue.NewMessagePrioritizer(test.State)
	require.Equal(t, prioritizer.Prior(a, b), "incorrect message priority")
}
