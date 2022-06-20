package stubdkg

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/bls12_381"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// DKG is a stub dkg protocol simulating a real DKG protocol with 3 stages in it
type DKG struct {
	identifier dkg.RequestID
	operatorID types.OperatorID
	threshold  uint16

	validatorPK    []byte
	operatorShares map[types.OperatorID]*bls.SecretKey
	msgs           map[Round][]*KeygenProtocolMsg
	state          *bls12_381.KeygenWrapper
}

func (s *DKG) Output() (*dkg.KeygenOutput, error) {
	localKeyShare := s.state.Output()
	if localKeyShare == nil {
		return nil, errors.New("unable to find output")
	}
	var sharePubKeys [][]byte
	for _, pk48 := range localKeyShare.SharePublicKeys {
		var pk []byte
		copy(pk, pk48[:])
		sharePubKeys = append(sharePubKeys, pk)
	}
	return &dkg.KeygenOutput{
		Index:           localKeyShare.Index,
		Threshold:       localKeyShare.Threshold,
		ShareCount:      localKeyShare.ShareCount,
		PublicKey:       localKeyShare.PublicKey[:],
		SecretShare:     localKeyShare.SecretShare[:],
		SharePublicKeys: sharePubKeys,
	}, nil
}

func New(operatorID types.OperatorID, identifier dkg.RequestID) dkg.Protocol {
	return &DKG{
		identifier: identifier,
		operatorID: operatorID,
		msgs:       map[Round][]*KeygenProtocolMsg{},
	}
}

func (s *DKG) SetOperators(validatorPK []byte, operatorShares map[types.OperatorID]*bls.SecretKey) {
	s.validatorPK = validatorPK
	s.operatorShares = operatorShares
}

func (s *DKG) Start(init *dkg.Init) ([]dkg.Message, error) {
	var myIndex = -1
	for i, id := range init.OperatorIDs {
		if id == s.operatorID {
			myIndex = i + 1
		}
	}
	s.threshold = init.Threshold
	s.state = bls12_381.New(myIndex, int(init.Threshold), len(init.OperatorIDs))
	outgoing0, err := s.state.Init()
	if err != nil {
		return nil, err
	}
	outgoing, err := s.packMessages(outgoing0)
	if err != nil {
		return nil, err
	}
	return outgoing, nil
}

func (s *DKG) ProcessMsg(msg0 *dkg.Message) (bool, []dkg.Message, error) {
	msg := &KeygenProtocolMsg{}
	err := msg.Decode(msg0.Data)
	if err != nil {
		return false, nil, err
	}
	if s.msgs[msg.RoundNumber] == nil {
		s.msgs[msg.RoundNumber] = []*KeygenProtocolMsg{}
	}
	s.msgs[msg.RoundNumber] = append(s.msgs[msg.RoundNumber], msg)
	if msg.RoundNumber < 1 || msg.RoundNumber > 4 {
		return false, nil, errors.New("wrong round number")
	}

	finished, outgoing0, err := s.state.HandleMessage(msg)
	if err != nil {
		return false, nil, err
	}
	outgoing, err := s.packMessages(outgoing0)
	if err != nil {
		return false, nil, err
	}
	return finished, outgoing, nil

}

func (s *DKG) packMessages(msgs []KeygenProtocolMsg) ([]dkg.Message, error) {
	var outgoing []dkg.Message
	for _, outMsg0 := range msgs {
		data, err := outMsg0.Encode()
		if err != nil {
			return nil, err
		}
		outMsg := dkg.Message{
			MsgType:    dkg.KeygenProtocolMsgType,
			Identifier: s.identifier,
			Data:       data,
		}
		outgoing = append(outgoing, outMsg)
	}
	return outgoing, nil
}

func (s *DKG) signDKGMsg(data []byte) *dkg.SignedMessage {
	return &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: s.identifier,
			Data:       data,
		},
		Signer: s.operatorID,
		// TODO - how do we sign?
	}
}
