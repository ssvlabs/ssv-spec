package validation

import (
	"testing"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsub_pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SpecTest to test the message validation checks
type MessageValidationTest struct {
	Name          string
	Messages      [][]byte
	Topic         string
	ReceivedAt    time.Time
	ExpectedError string
}

func (test *MessageValidationTest) TestName() string {
	return "messagevalidation " + test.Name
}

func (test *MessageValidationTest) Run(t *testing.T) {

	peerID := peer.ID("peerID")
	messageValidator := testingutils.TestingMessageValidator()
	topic := test.GetTopic()

	var lastErr error

	for _, msg := range test.Messages {

		receivedAt := test.GetReceivedAt(msg, messageValidator)

		// Construct pubsub message from the data
		pubsubMessage := &pubsub.Message{
			ReceivedFrom: peerID,
			Message: &pubsub_pb.Message{
				Data:  msg,
				Topic: &topic,
			},
		}

		// Process message
		err := messageValidator.HandlePubsubMessage(pubsubMessage, receivedAt)
		if err != nil {
			lastErr = err
		}

		// Update state if no error
		if err == nil {
			if err := messageValidator.UpdateState(peerID, pubsubMessage); err != nil {
				lastErr = err
			}
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *MessageValidationTest) GetPostState() (interface{}, error) {
	return nil, nil
}

// Returns the topic to be used. If not set, use default value
func (test *MessageValidationTest) GetTopic() string {
	topic := testingutils.TestingTopic
	if len(test.Topic) > 0 {
		topic = test.Topic
	}
	return topic
}

// Returns the ReceivedAt time to be used. If not set, use messgae's slot starting time
func (test *MessageValidationTest) GetReceivedAt(msg []byte, mv *validation.MessageValidator) time.Time {
	receivedAt := test.ReceivedAt
	if receivedAt.Equal(time.Time{}) {
		// If not set, get slot start time for message so that it's valid
		slot := testingutils.GetMessageSlot(msg)
		receivedAt = time.Unix(mv.Beacon.EstimatedTimeAtSlot(slot), 0)
	}
	return receivedAt
}

// SpecTest with list of MessageValidationTest
type MultiMessageValidationTest struct {
	Name  string
	Tests []*MessageValidationTest
}

func (tests *MultiMessageValidationTest) TestName() string {
	return tests.Name
}

func (tests *MultiMessageValidationTest) Run(t *testing.T) {
	for _, test := range tests.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func (tests *MultiMessageValidationTest) GetPostState() (interface{}, error) {
	return nil, nil
}
