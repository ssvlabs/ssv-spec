package qbft

import (
	"bytes"
	"encoding/json"

	"github.com/ssvlabs/ssv-spec/types"
)

type MsgContainer struct {
	Msgs map[Round][]*types.SignedSSVMessage
}

func NewMsgContainer() *MsgContainer {
	return &MsgContainer{
		Msgs: map[Round][]*types.SignedSSVMessage{},
	}
}

// AllMessages returns all messages
func (c *MsgContainer) AllMessages() []*types.SignedSSVMessage {
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
			signedMsg := c.Msgs[round][i]

			m, err := DecodeMessage(signedMsg.SSVMessage.Data)
			if err != nil {
				continue
			}

			if !bytes.Equal(m.Root[:], root[:]) {
				continue
			}
			ret = append(ret, signedMsg)
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
		signedMsg := c.Msgs[round][i]

		m, err := DecodeMessage(signedMsg.SSVMessage.Data)
		if err != nil {
			continue
		}

		if !bytes.Equal(m.Root[:], root[:]) {
			continue
		}

		currentSigners := make([]types.OperatorID, 0)
		currentMsgs := make([]*types.SignedSSVMessage, 0)
		currentMsgs = append(currentMsgs, signedMsg)
		currentSigners = append(currentSigners, signedMsg.GetOperatorIDs()...)
		for j := i + 1; j < len(c.Msgs[round]); j++ {
			signedMsg2 := c.Msgs[round][j]

			m2, err := DecodeMessage(signedMsg2.SSVMessage.Data)
			if err != nil {
				continue
			}

			if !bytes.Equal(m2.Root[:], root[:]) {
				continue
			}

			if !signedMsg2.CommonSigners(currentSigners) {
				currentMsgs = append(currentMsgs, signedMsg2)
				currentSigners = append(currentSigners, signedMsg2.GetOperatorIDs()...)
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
func (c *MsgContainer) AddFirstMsgForSignerAndRound(signedMsg *types.SignedSSVMessage) (bool, error) {

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return false, err
	}

	if c.Msgs[msg.Round] == nil {
		c.Msgs[msg.Round] = make([]*types.SignedSSVMessage, 0)
	}

	for _, existingMsg := range c.Msgs[msg.Round] {
		if existingMsg.MatchedSigners(signedMsg.GetOperatorIDs()) {
			return false, nil
		}
	}

	// add msg
	c.Msgs[msg.Round] = append(c.Msgs[msg.Round], signedMsg)
	return true, nil
}

// AddMsg will add any message regardless of signers
func (c *MsgContainer) AddMsg(signedMsg *types.SignedSSVMessage) error {

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return err
	}

	if c.Msgs[msg.Round] == nil {
		c.Msgs[msg.Round] = make([]*types.SignedSSVMessage, 0)
	}
	// add msg
	c.Msgs[msg.Round] = append(c.Msgs[msg.Round], signedMsg)

	return nil
}

// Encode returns the encoded struct in bytes or error
func (c *MsgContainer) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode returns error if decoding failed
func (c *MsgContainer) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}
