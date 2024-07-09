package qbft

import (
	"bytes"

	"github.com/ssvlabs/ssv-spec/types"
)

type MsgContainer struct {
	Msgs map[Round][]*ProcessingMessage
}

func NewMsgContainer() *MsgContainer {
	return &MsgContainer{
		Msgs: map[Round][]*ProcessingMessage{},
	}
}

// AllMessages returns all messages
func (c *MsgContainer) AllMessages() []*ProcessingMessage {
	ret := make([]*ProcessingMessage, 0)
	for _, roundMsgs := range c.Msgs {
		ret = append(ret, roundMsgs...)
	}
	return ret
}

// MessagesForRound returns all msgs for Height and round, empty slice otherwise
func (c *MsgContainer) MessagesForRound(round Round) []*ProcessingMessage {
	if c.Msgs[round] != nil {
		return c.Msgs[round]
	}
	return make([]*ProcessingMessage, 0)
}

// MessagesForRoundAndRoot returns all msgs for round and value, empty slice otherwise
func (c *MsgContainer) MessagesForRoundAndRoot(round Round, root [32]byte) []*ProcessingMessage {
	if c.Msgs[round] != nil {
		ret := make([]*ProcessingMessage, 0)
		for i := 0; i < len(c.Msgs[round]); i++ {
			msg := c.Msgs[round][i]

			if !bytes.Equal(msg.QBFTMessage.Root[:], root[:]) {
				continue
			}
			ret = append(ret, msg)
		}
		return ret
	}
	return make([]*ProcessingMessage, 0)
}

// LongestUniqueSignersForRoundAndRoot returns the longest set of unique signers and msgs for a specific round and value
func (c *MsgContainer) LongestUniqueSignersForRoundAndRoot(round Round, root [32]byte) ([]types.OperatorID, []*ProcessingMessage) {
	signersRet := make([]types.OperatorID, 0)
	msgsRet := make([]*ProcessingMessage, 0)
	if c.Msgs[round] == nil {
		return signersRet, msgsRet
	}

	for i := 0; i < len(c.Msgs[round]); i++ {
		msg := c.Msgs[round][i]

		if !bytes.Equal(msg.QBFTMessage.Root[:], root[:]) {
			continue
		}

		currentSigners := make([]types.OperatorID, 0)
		currentMsgs := make([]*ProcessingMessage, 0)
		currentMsgs = append(currentMsgs, msg)
		currentSigners = append(currentSigners, msg.SignedMessage.OperatorIDs...)
		for j := i + 1; j < len(c.Msgs[round]); j++ {
			msg2 := c.Msgs[round][j]

			if !bytes.Equal(msg2.QBFTMessage.Root[:], root[:]) {
				continue
			}

			if !msg2.SignedMessage.CommonSigners(currentSigners) {
				currentMsgs = append(currentMsgs, msg2)
				currentSigners = append(currentSigners, msg2.SignedMessage.OperatorIDs...)
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
func (c *MsgContainer) AddFirstMsgForSignerAndRound(msg *ProcessingMessage) (bool, error) {

	if c.Msgs[msg.QBFTMessage.Round] == nil {
		c.Msgs[msg.QBFTMessage.Round] = make([]*ProcessingMessage, 0)
	}

	for _, existingMsg := range c.Msgs[msg.QBFTMessage.Round] {
		if existingMsg.SignedMessage.MatchedSigners(msg.SignedMessage.OperatorIDs) {
			return false, nil
		}
	}

	// add msg
	c.Msgs[msg.QBFTMessage.Round] = append(c.Msgs[msg.QBFTMessage.Round], msg)
	return true, nil
}

// AddMsg will add any message regardless of signers
func (c *MsgContainer) AddMsg(msg *ProcessingMessage) error {

	if c.Msgs[msg.QBFTMessage.Round] == nil {
		c.Msgs[msg.QBFTMessage.Round] = make([]*ProcessingMessage, 0)
	}
	// add msg
	c.Msgs[msg.QBFTMessage.Round] = append(c.Msgs[msg.QBFTMessage.Round], msg)

	return nil
}
