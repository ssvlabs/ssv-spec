package stubdkg

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// DKG is a stub dkg protocol simulating a real DKG protocol with 3 stages in it
type DKG struct {
	identifier dkg.RequestID
	network    dkg.Network
	operatorID types.OperatorID
	operators  []types.OperatorID

	validatorPK    []byte
	operatorShares map[types.OperatorID]*bls.SecretKey

	msgs map[Stage][]*ProtocolMsg
}

func New(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID) dkg.KeyGenProtocol {
	return &DKG{
		identifier: identifier,
		network:    network,
		operatorID: operatorID,
		msgs:       map[Stage][]*ProtocolMsg{},
	}
}

func (s *DKG) SetOperators(validatorPK []byte, operatorShares map[types.OperatorID]*bls.SecretKey) {
	s.validatorPK = validatorPK
	s.operatorShares = operatorShares
}

func (s *DKG) Start(init *dkg.Init) error {
	s.operators = init.OperatorIDs
	// TODO send Stage 1 msg
	return nil
}

func (s *DKG) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error) {
	// TODO validate msg

	dataMsg := &ProtocolMsg{}
	if err := dataMsg.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "could not decode protocol msg")
	}

	if s.msgs[dataMsg.Stage] == nil {
		s.msgs[dataMsg.Stage] = []*ProtocolMsg{}
	}
	s.msgs[dataMsg.Stage] = append(s.msgs[dataMsg.Stage], dataMsg)

	switch dataMsg.Stage {
	case StubStage1:
		if len(s.msgs[StubStage1]) == len(s.operators) {
			fmt.Printf("stage 1 done\n")
			// TODO send Stage 2 msg
		}
	case StubStage2:
		if len(s.msgs[StubStage2]) == len(s.operators) {
			fmt.Printf("stage 2 done\n")
			// TODO send Stage 3 msg
		}
	case StubStage3:
		if len(s.msgs[StubStage3]) == len(s.operators) {
			ret := &dkg.KeyGenOutput{
				Share:       s.operatorShares[s.operatorID],
				ValidatorPK: s.validatorPK,
				OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
					1: s.operatorShares[1].GetPublicKey(),
					2: s.operatorShares[2].GetPublicKey(),
					3: s.operatorShares[3].GetPublicKey(),
					4: s.operatorShares[4].GetPublicKey(),
				},
			}
			return true, ret, nil
		}
	}
	return false, nil, nil
}

//func (s *DKG) signDKGMsg(data []byte) *dkg.SignedMessage {
//	return &dkg.SignedMessage{
//		Message: &dkg.Message{
//			MsgType:    dkg.ProtocolMsgType,
//			Identifier: s.identifier,
//			Data:       data,
//		},
//		Signer: s.operatorID,
//		// TODO - how do we sign?
//	}
//}
