package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/aquasecurity/table"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

var (
	mockStateRunning = &State{
		HasRunningInstance: true,
		Height:             100,
		Slot:               64,
		Quorum:             4,
	}
)

func TestMessagePrioritizerSlice(t *testing.T) {
	tests := []struct {
		name     string
		state    *State
		messages []mockMessage
	}{
		{
			name:  "All-in-one",
			state: mockStateRunning,
			messages: []mockMessage{
				// 1. Current height/slot:
				// 1.1. Consensus
				mockConsensusMessage{Height: 100},
				// 1.2. Pre-consensus
				mockNonConsensusMessage{Slot: 64, Type: ssv.SelectionProofPartialSig},
				// 1.3. Post-consensus
				mockNonConsensusMessage{Slot: 64, Type: ssv.PostConsensusPartialSig},

				// 2. Higher height/slot:
				// 2.1 Decided
				mockConsensusMessage{Height: 101, Decided: true},
				// 2.2. Pre-consensus
				mockNonConsensusMessage{Slot: 65, Type: ssv.SelectionProofPartialSig},
				// 2.3. Consensus
				mockConsensusMessage{Height: 101},
				// 2.4. Post-consensus
				mockNonConsensusMessage{Slot: 65, Type: ssv.PostConsensusPartialSig},

				// 3. Lower height/slot:
				// 3.1 Decided
				mockConsensusMessage{Height: 99, Decided: true},
				// 3.2. Commit
				mockConsensusMessage{Height: 99, Type: qbft.CommitMsgType},
				// 3.3. Pre-consensus
				mockNonConsensusMessage{Slot: 63, Type: ssv.SelectionProofPartialSig},
			},
		},
	}

	for testIndex, test := range tests {
		t.Run(fmt.Sprintf("#%d: %s", testIndex, test.name), func(t *testing.T) {
			messages := make(messageSlice, len(test.messages))
			for i, m := range test.messages {
				var err error
				messages[i], err = DecodeSSVMessage(m.ssvMessage(test.state))
				require.NoError(t, err)
			}

			var scrambles []messageSlice
			for {
				scrambledMessages := messages.scramble()
				if scrambledMessages.equal(messages) {
					continue
				}
				scrambles = append(scrambles, scrambledMessages)
				if len(scrambles) == 10 {
					break
				}
			}

			prioritizer := NewSSVMessagePrioritizer(test.state)
			for _, scramble := range scrambles {
				scramble.sort(prioritizer)
				correctOrder := messages.equal(scramble)
				if !correctOrder {
					require.Fail(t, "incorrect order:\n"+scramble.dump(test.state))
				}
			}
		})
	}
}

type mockMessage interface {
	ssvMessage(*State) *types.SSVMessage
}

type mockConsensusMessage struct {
	Type    qbft.MessageType
	Decided bool
	Height  qbft.Height
}

func (m mockConsensusMessage) ssvMessage(state *State) *types.SSVMessage {
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

func (m mockNonConsensusMessage) ssvMessage(state *State) *types.SSVMessage {
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

type messageSlice []*DecodedSSVMessage

func (m messageSlice) scramble() messageSlice {
	scrambled := make([]*DecodedSSVMessage, len(m))
	for i, j := range rand.Perm(len(m)) {
		scrambled[i] = m[j]
	}
	return scrambled
}

func (m messageSlice) sort(prioritizer *SSVMessagePrioritizer) {
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

func (m messageSlice) dump(s *State) string {
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

		switch compareHeightOrSlot(s, msg) {
		case -1:
			relation = "lower"
		case 0:
			relation = "current"
		case 1:
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
			fmt.Sprintf("%t", isDecidedMesssage(s, msg)),
		)
	}
	tbl.Render()
	return b.String()
}
