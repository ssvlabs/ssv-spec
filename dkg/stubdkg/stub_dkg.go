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
	network    dkg.Network
	operatorID types.OperatorID
	threshold  uint16

	validatorPK[] byte
	operatorShares map[types.OperatorID]*bls.SecretKey
	msgs  map[Round][]*KeygenProtocolMsg
	state *bls12_381.KeygenWrapper
}

func (s DKG) Output() *LocalKeyShare {
	panic("implement me")
}

func New(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID) dkg.Protocol {
	return &DKG{
		identifier: identifier,
		network:    network,
		operatorID: operatorID,
		msgs:       map[Round][]*KeygenProtocolMsg{},
	}
}

func (s *DKG) SetOperators(validatorPK []byte, operatorShares map[types.OperatorID]*bls.SecretKey) {
	s.validatorPK = validatorPK
	s.operatorShares = operatorShares
}

func (s *DKG) Start(init *dkg.Init) error {
	var myIndex = -1
	for i, id := range init.OperatorIDs {
		if id == s.operatorID {
			myIndex = i + 1
		}
	}
	s.threshold = init.Threshold
	s.state = bls12_381.New(myIndex, int(init.Threshold), len(init.OperatorIDs))
	outgoing, err := s.state.Init()
	if err != nil {
		return err
	}
	for _, msg := range outgoing {
		s.network.Broadcast(&msg)
	}
	return nil
}

func (s *DKG) ProcessMsg(msg *KeygenProtocolMsg) (bool, []KeygenProtocolMsg, error) {

	if s.msgs[msg.RoundNumber] == nil {
		s.msgs[msg.RoundNumber] = []*KeygenProtocolMsg{}
	}
	s.msgs[msg.RoundNumber] = append(s.msgs[msg.RoundNumber], msg)
	if msg.RoundNumber < 1 || msg.RoundNumber > 4 {
		return false, nil, errors.New("wrong round number")
	}

	finished, outgoing, err := s.state.HandleMessage(msg)
	if err != nil {
		return false, nil, err
	}
	for _, outMsg := range outgoing {
		s.network.Broadcast(&outMsg)
	}
	return finished, outgoing, nil

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
