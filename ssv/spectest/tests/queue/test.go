package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/aquasecurity/table"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/queue"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

type MessagePrioritySliceTest struct {
	Name  string
	State *queue.State

	// Messages is the correctly-ordered slice of messages.
	Messages []mockMessage

	// Permutations is a slice of permutations, each of which is
	// used to deterministically shuffle Messages.
	Permutations [][]int
}

func (test *MessagePrioritySliceTest) TestName() string {
	return test.Name
}

func (test *MessagePrioritySliceTest) Run(t *testing.T) {
	messages := make(messageSlice, len(test.Messages))
	for i, m := range test.Messages {
		var err error
		messages[i], err = queue.DecodeSSVMessage(m.toSSVMessage(test.State))
		require.NoError(t, err)
	}

	shuffles := make([]messageSlice, len(test.Permutations))
	for i, perm := range test.Permutations {
		shuffles[i] = messages.shuffle(perm)
	}

	prioritizer := queue.NewMessagePrioritizer(test.State)

	for _, shuffle := range shuffles {
		require.False(t, messages.equal(shuffle), "shuffle shouldn't be in correct order!")
		shuffle.sort(prioritizer)
		require.True(t, messages.equal(shuffle), "incorrect order:\n"+shuffle.dump(test.State))
	}
}

type mockMessage struct {
	Consensus    *mockConsensusMessage
	NonConsensus *mockNonConsensusMessage
}

func (m mockMessage) toSSVMessage(state *queue.State) *types.SSVMessage {
	if m.Consensus != nil {
		return m.Consensus.toSSVMessage(state)
	}
	return m.NonConsensus.toSSVMessage(state)
}

type mockConsensusMessage struct {
	Type    qbft.MessageType
	Decided bool
	Height  qbft.Height
}

func (m mockConsensusMessage) toSSVMessage(state *queue.State) *types.SSVMessage {
	var (
		typ         = m.Type
		signerCount = 1
	)
	if m.Decided {
		typ = qbft.CommitMsgType
		signerCount = int(state.Quorum) + 1
	}

	var signers []types.OperatorID
	for i := 0; i < signerCount; i++ {
		signers = append(signers, types.OperatorID(i))
	}

	return testingutils.SSVMsgAttester(&qbft.SignedMessage{
		Message: &qbft.Message{
			MsgType:    typ,
			Height:     m.Height,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       []byte{1, 2, 3, 4},
		},
		Signature: []byte{1, 2, 3, 4},
		Signers:   signers,
	}, nil)
}

type mockNonConsensusMessage struct {
	Type ssv.PartialSigMsgType
	Slot phase0.Slot
}

func (m mockNonConsensusMessage) toSSVMessage(state *queue.State) *types.SSVMessage {
	return testingutils.SSVMsgAttester(nil, &ssv.SignedPartialSignatureMessage{
		Message: ssv.PartialSignatureMessages{
			Type: m.Type,
			Messages: []*ssv.PartialSignatureMessage{
				{
					Slot:             m.Slot,
					PartialSignature: []byte{},
					SigningRoot:      []byte{},
					Signer:           0,
				},
			},
		},
		Signature: []byte{1, 2, 3, 4},
		Signer:    types.OperatorID(1),
	})
}

type messageSlice []*queue.DecodedSSVMessage

func (m messageSlice) shuffle(permutation []int) messageSlice {
	shuffled := make([]*queue.DecodedSSVMessage, len(m))
	for i, j := range permutation {
		shuffled[i] = m[j]
	}
	return shuffled
}

func (m messageSlice) sort(prioritizer queue.MessagePrioritizer) {
	sort.Slice(m, func(i, j int) bool {
		return prioritizer.Prior(m[i], m[j])
	})
}

func (m messageSlice) equal(m2 messageSlice) bool {
	if len(m) != len(m2) {
		return false
	}
	for i := range m {
		a, err := json.Marshal(m[i])
		if err != nil {
			panic(err)
		}
		b, err := json.Marshal(m2[i])
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(a, b) {
			return false
		}
	}
	return true
}

func (m messageSlice) dump(s *queue.State) string {
	b := &strings.Builder{}
	tbl := table.New(b)
	tbl.SetHeaders("#", "Kind", "Height/Slot", "Type", "Decided")
	for i, msg := range m {
		var (
			kind         string
			typ          interface{}
			heightOrSlot interface{}
			relation     string
		)

		n := 0
		switch m := msg.Body.(type) {
		case *qbft.SignedMessage:
			n = int(m.Message.Height) - int(s.Height)
		case *ssv.SignedPartialSignatureMessage:
			n = int(m.Message.Messages[0].Slot) - int(s.Slot)
		}
		switch {
		case n < 0:
			relation = "lower"
		case n == 0:
			relation = "current"
		case n > 0:
			relation = "higher"
		}

		switch mm := msg.Body.(type) {
		case *ssv.SignedPartialSignatureMessage:
			heightOrSlot = mm.Message.Messages[0].Slot
			typ = mm.Message.Type
			if typ == ssv.PostConsensusPartialSig {
				kind = "post-consensus"
			} else {
				kind = "pre-consensus"
			}
		case *qbft.SignedMessage:
			kind = "consensus"
			heightOrSlot = mm.Message.Height
			typ = mm.Message.MsgType
		}

		tbl.AddRow(
			fmt.Sprint(i),
			kind,
			fmt.Sprintf("%d (%s)", heightOrSlot, relation),
			fmt.Sprint(typ),
			// fmt.Sprintf("%t", isDecidedMesssage(s, msg)),
			"",
		)
	}
	tbl.Render()
	return b.String()
}
