package queue

import (
	"testing"

	"github.com/bloxapp/ssv-spec/queue"
	"github.com/stretchr/testify/require"
)

type MessagePriorityPairTest struct {
	Name  string
	State *queue.State

	// Messag A is expected to be prior to message B.
	A, B mockMessage
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
	require.True(t, prioritizer.Prior(a, b), "incorrect message priority")
}
