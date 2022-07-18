package keygen

import (
	"bytes"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type KGProtocol struct {
	Identifier dkg.RequestID
	Operator   types.OperatorID
	Init       dkg.Init
	State      *Keygen
}

func (k *KGProtocol) Output() ([]byte, error) {
	if k.State == nil || k.State.Output == nil {
		return nil, nil
	}
	jStr, err := json.Marshal(k.State.Output)
	if err != nil {
		return nil, err
	}
	return jStr, nil
}

func New(init *dkg.Init, identifier dkg.RequestID, config dkg.ProtocolConfig) (dkg.Protocol, error) {
	var myIndex uint64 = 0
	for i, id := range init.OperatorIDs {
		if id == config.Operator.OperatorID {
			myIndex = uint64(i) + 1
		}
	}
	var ids []uint64
	for _, id := range init.OperatorIDs {
		ids = append(ids, uint64(id))
	}
	state, err := NewKeygen(identifier[:], myIndex, uint64(init.Threshold), ids)
	if err != nil {
		return nil, err
	}
	return &KGProtocol{
		Init:       *init,
		Identifier: identifier,
		Operator:   config.Operator.OperatorID,
		State:      state,
	}, nil
}

func (k *KGProtocol) Start() ([]base.Message, error) {
	if err := k.State.Proceed(); err != nil {
		return nil, err
	}
	return k.getAndEncodeOutgoing()
}

func (k *KGProtocol) ProcessMsg(msg *base.Message) ([]base.Message, error) {
	if msg == nil {
		return nil, errors.New("nil message")
	}
	if msg.Header.MsgType != k.State.HandleMessageType {
		return nil, errors.New("not valid message type")
	}
	if bytes.Compare(msg.Header.SessionId, k.State.SessionID) != 0 {
		return nil, errors.New("invalid Identifier")
	}
	pMsg := &ParsedMessage{}
	if err := pMsg.FromBase(msg); err != nil {
		return nil, err
	}

	if err := k.State.PushMessage(pMsg); err != nil {
		return nil, err
	}

	if err := k.State.Proceed(); err != nil {
		return nil, err
	}
	return k.getAndEncodeOutgoing()
}

func (k *KGProtocol) getAndEncodeOutgoing() ([]base.Message, error) {
	outgoingInner, err := k.State.GetOutgoing()
	if err != nil {
		return nil, err
	}
	var outgoing []base.Message
	for _, out := range outgoingInner {
		if msg, err := out.ToBase(); err == nil {
			outgoing = append(outgoing, *msg)
		} else {
			// TODO: Log error
			log.Errorf("error: %v", err)
		}
	}
	return outgoing, nil
}
