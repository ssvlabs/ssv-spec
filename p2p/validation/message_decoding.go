package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Decodes a pubsub.Message.Data into a SignedSSVMessage
func decodeSignedSSVMessage(pMsg *pubsub.Message) (*types.SignedSSVMessage, error) {
	msg := &types.SignedSSVMessage{}
	err := msg.Decode(pMsg.Data)
	if err != nil {
		return nil, ErrMalformedPubSubMessage
	}
	return msg, nil
}

// Get a SignedSSVMessage's slot
func getMessageSlot(msg *types.SignedSSVMessage) (phase0.Slot, error) {
	if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
		qbftMessage := &qbft.Message{}
		if err := qbftMessage.Decode(msg.SSVMessage.Data); err != nil {
			return 0, ErrUndecodableData
		}
		return phase0.Slot(qbftMessage.Height), nil
	} else if msg.SSVMessage.MsgType == types.SSVPartialSignatureMsgType {
		pSigMsgs := &types.PartialSignatureMessages{}
		if err := pSigMsgs.Decode(msg.SSVMessage.Data); err != nil {
			return 0, ErrUndecodableData
		}
		return pSigMsgs.Slot, nil
	} else {
		return 0, ErrUnknownSSVMessageType
	}
}
