package hbbft

import (
	"bytes"
	"encoding/json"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type MsgContainer struct {
	Msgs map[Round][]*SignedMessage
}

func NewMsgContainer() *MsgContainer {
	return &MsgContainer{
		Msgs: map[Round][]*SignedMessage{},
	}
}

// AllMessaged returns all messages
func (c *MsgContainer) AllMessaged() []*SignedMessage {
	ret := make([]*SignedMessage, 0)
	for _, roundMsgs := range c.Msgs {
		ret = append(ret, roundMsgs...)
	}
	return ret
}

// MessagesForRound returns all msgs for Height and round, empty slice otherwise
func (c *MsgContainer) MessagesForRound(round Round) []*SignedMessage {
	if c.Msgs[round] != nil {
		return c.Msgs[round]
	}
	return make([]*SignedMessage, 0)
}

// MessagesForRoundAndValue returns all msgs for round and value, empty slice otherwise
func (c *MsgContainer) MessagesForRoundAndValue(round Round, value []byte) []*SignedMessage {
	if c.Msgs[round] != nil {
		ret := make([]*SignedMessage, 0)
		for i := 0; i < len(c.Msgs[round]); i++ {
			m := c.Msgs[round][i]

			if !bytes.Equal(m.Message.Data, value) {
				continue
			}
			ret = append(ret, m)
		}
		return ret
	}
	return make([]*SignedMessage, 0)
}

// LongestUniqueSignersForRoundAndValue returns the longest set of unique signers and msgs for a specific round and value
func (c *MsgContainer) LongestUniqueSignersForRoundAndValue(round Round, value []byte) ([]types.OperatorID, []*SignedMessage) {
	signersRet := make([]types.OperatorID, 0)
	msgsRet := make([]*SignedMessage, 0)
	if c.Msgs[round] == nil {
		return signersRet, msgsRet
	}

	for i := 0; i < len(c.Msgs[round]); i++ {
		m := c.Msgs[round][i]

		if !bytes.Equal(m.Message.Data, value) {
			continue
		}

		currentSigners := make([]types.OperatorID, 0)
		currentMsgs := make([]*SignedMessage, 0)
		currentMsgs = append(currentMsgs, m)
		currentSigners = append(currentSigners, m.GetSigners()...)
		for j := i + 1; j < len(c.Msgs[round]); j++ {
			m2 := c.Msgs[round][j]

			if !bytes.Equal(m2.Message.Data, value) {
				continue
			}

			if !m2.CommonSigners(currentSigners) {
				currentMsgs = append(currentMsgs, m2)
				currentSigners = append(currentSigners, m2.GetSigners()...)
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
func (c *MsgContainer) AddFirstMsgForSignerAndRound(msg *SignedMessage) (bool, error) {
	if c.Msgs[msg.Message.Round] == nil {
		c.Msgs[msg.Message.Round] = make([]*SignedMessage, 0)
	}

	for _, existingMsg := range c.Msgs[msg.Message.Round] {
		if existingMsg.MatchedSigners(msg.Signers) {
			return false, nil
		}
	}

	// add msg
	c.Msgs[msg.Message.Round] = append(c.Msgs[msg.Message.Round], msg)
	return true, nil
}

// AddMsg will add any message regardless of signers
func (c *MsgContainer) AddMsg(msg *SignedMessage) {
	if c.Msgs[msg.Message.Round] == nil {
		c.Msgs[msg.Message.Round] = make([]*SignedMessage, 0)
	}
	// add msg
	c.Msgs[msg.Message.Round] = append(c.Msgs[msg.Message.Round], msg)
}

// Encode returns the encoded struct in bytes or error
func (c *MsgContainer) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode returns error if decoding failed
func (c *MsgContainer) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}
