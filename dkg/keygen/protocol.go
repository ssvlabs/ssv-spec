package keygen

import (
	"github.com/bloxapp/ssv-spec/dkg"
	types2 "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type KGProtocol struct {
	identifier dkg.RequestID
	operatorID types.OperatorID
	init       dkg.Init
	state      *Keygen
}

func New(init *dkg.Init, identifier dkg.RequestID, config types2.ProtocolConfig) (dkg.Protocol, error) {
	var myIndex uint16 = 0
	for i, id := range init.OperatorIDs {
		if id == config.Operator.OperatorID {
			myIndex = uint16(i) + 1
		}
	}
	state, err := NewKeygen(myIndex, init.Threshold, uint16(len(init.OperatorIDs)))
	if err != nil {
		return nil, err
	}
	return &KGProtocol{
		init:       *init,
		identifier: identifier,
		operatorID: config.Operator.OperatorID,
		state:      state,
	}, nil
}

func (k *KGProtocol) Start() ([]dkg.Message, error) {
	if err := k.state.Proceed(); err != nil {
		return nil, err
	}
	return k.getAndEncodeOutgoing()
}

func (k *KGProtocol) ProcessMsg(msg *dkg.Message) ([]dkg.Message, error) {
	if msg == nil {
		return nil, errors.New("nil message")
	}
	if msg.MsgType != dkg.ProtocolMsgType {
		return nil, errors.New("not valid message type")
	}
	if msg.Identifier != k.identifier {
		return nil, errors.New("invalid identifier")
	}
	pMsg := new(Message)
	pMsg.Decode(msg.Data)
	k.state.PushMessage(pMsg)
	return k.getAndEncodeOutgoing()
}

func (k *KGProtocol) getAndEncodeOutgoing() ([]dkg.Message, error) {
	outgoingInner, err := k.state.GetOutgoing()
	if err != nil {
		return nil, err
	}
	var outgoing []dkg.Message
	for _, msg := range outgoingInner {
		data, err := msg.Encode()
		if err != nil {
			return nil, err
		}
		outMsg := dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: k.identifier,
			Data:       data,
		}
		outgoing = append(outgoing, outMsg)
	}
	return outgoing, nil
}
