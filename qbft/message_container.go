package qbft

import (
	"bytes"
	"encoding/json"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type MsgContainer struct {
	Msgs map[Round][]*types.SignedSSVMessage
}

func NewMsgContainer() *MsgContainer {
	return &MsgContainer{
		Msgs: map[Round][]*types.SignedSSVMessage{},
	}
}

// AllMessaged returns all messages
func (c *MsgContainer) AllMessaged() []*types.SignedSSVMessage {
	ret := make([]*types.SignedSSVMessage, 0)
	for _, roundMsgs := range c.Msgs {
		ret = append(ret, roundMsgs...)
	}
	return ret
}

// MessagesForRound returns all msgs for Height and round, empty slice otherwise
func (c *MsgContainer) MessagesForRound(round Round) []*types.SignedSSVMessage {
	if c.Msgs[round] != nil {
		return c.Msgs[round]
	}
	return make([]*types.SignedSSVMessage, 0)
}

// MessagesForRoundAndRoot returns all msgs for round and value, empty slice otherwise
func (c *MsgContainer) MessagesForRoundAndRoot(round Round, root [32]byte) []*types.SignedSSVMessage {
	if c.Msgs[round] != nil {
		ret := make([]*types.SignedSSVMessage, 0)
		for i := 0; i < len(c.Msgs[round]); i++ {
			m := c.Msgs[round][i]

			// Decode
			message := &Message{}
			if err := message.Decode(m.SSVMessage.Data); err != nil {
				continue
			}

			if !bytes.Equal(message.Root[:], root[:]) {
				continue
			}
			ret = append(ret, m)
		}
		return ret
	}
	return make([]*types.SignedSSVMessage, 0)
}

// LongestUniqueSignersForRoundAndRoot returns the longest set of unique signers and msgs for a specific round and value
func (c *MsgContainer) LongestUniqueSignersForRoundAndRoot(round Round, root [32]byte) ([]types.OperatorID, []*types.SignedSSVMessage) {
	signersRet := make([]types.OperatorID, 0)
	msgsRet := make([]*types.SignedSSVMessage, 0)
	if c.Msgs[round] == nil {
		return signersRet, msgsRet
	}

	for i := 0; i < len(c.Msgs[round]); i++ {
		m := c.Msgs[round][i]

		// Decode
		message := &Message{}
		if err := message.Decode(m.SSVMessage.Data); err != nil {
			continue
		}

		if !bytes.Equal(message.Root[:], root[:]) {
			continue
		}

		currentSigners := make([]types.OperatorID, 0)
		currentMsgs := make([]*types.SignedSSVMessage, 0)
		currentMsgs = append(currentMsgs, m)
		currentSigners = append(currentSigners, m.GetOperatorIDs()...)
		for j := i + 1; j < len(c.Msgs[round]); j++ {
			m2 := c.Msgs[round][j]

			if !bytes.Equal(message.Root[:], root[:]) {
				continue
			}

			if !CommonSigners(m2, currentSigners) {
				currentMsgs = append(currentMsgs, m2)
				currentSigners = append(currentSigners, m2.GetOperatorIDs()...)
			}
		}

		if len(signersRet) < len(currentSigners) {
			signersRet = currentSigners
			msgsRet = currentMsgs
		}
	}

	return signersRet, msgsRet
}

// AddFirstMsgForSignerAndRound will add the first msg for each signer for a specific round, consequent msgs will not be added
func (c *MsgContainer) AddFirstMsgForSignerAndRound(msg *types.SignedSSVMessage) (bool, error) {

	// Decode
	message := &Message{}
	if err := message.Decode(msg.SSVMessage.Data); err != nil {
		return false, errors.Wrap(err, "Could not decode Message")
	}

	if c.Msgs[message.Round] == nil {
		c.Msgs[message.Round] = make([]*types.SignedSSVMessage, 0)
	}

	for _, existingMsg := range c.Msgs[message.Round] {
		if MatchedSigners(existingMsg, msg.GetOperatorIDs()) {
			return false, nil
		}
	}

	// add msg
	c.Msgs[message.Round] = append(c.Msgs[message.Round], msg)
	return true, nil
}

// AddMsg will add any message regardless of signers
func (c *MsgContainer) AddMsg(msg *types.SignedSSVMessage) {

	// Decode
	message := &Message{}
	if err := message.Decode(msg.SSVMessage.Data); err != nil {
		return
	}

	if c.Msgs[message.Round] == nil {
		c.Msgs[message.Round] = make([]*types.SignedSSVMessage, 0)
	}
	// add msg
	c.Msgs[message.Round] = append(c.Msgs[message.Round], msg)
}

// Encode returns the encoded struct in bytes or error
func (c *MsgContainer) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode returns error if decoding failed
func (c *MsgContainer) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}

// CommonSigners returns true if there is at least one common signer betwen msg and signers
func CommonSigners(msg *types.SignedSSVMessage, signers []types.OperatorID) bool {
	for _, id1 := range msg.GetOperatorIDs() {
		for _, id2 := range signers {
			if id1 == id2 {
				return true
			}
		}
	}
	return false
}

// MatchedSigners returns true if the messages have equal signers without order significance
func MatchedSigners(msg1 *types.SignedSSVMessage, signers []types.OperatorID) bool {
	if len(msg1.GetOperatorIDs()) != len(signers) {
		return false
	}

	for _, id := range msg1.GetOperatorIDs() {
		found := false
		for _, id2 := range signers {
			if id == id2 {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}
