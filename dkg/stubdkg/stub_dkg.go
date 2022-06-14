package stubdkg

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// s is a stub dkg protocol simulating a real DKG protocol with 3 stages in it
type s struct {
	identifier types.MessageID
	network    dkg.Network
	operatorID types.OperatorID
	threshold  uint16

	msgs map[Round][]*KeygenProtocolMsg
}

func New(network dkg.Network, operatorID types.OperatorID, identifier types.MessageID) dkg.Protocol {
	return &s{
		identifier: identifier,
		network:    network,
		operatorID: operatorID,
		msgs:       map[Round][]*KeygenProtocolMsg{},
	}
}

func (s *s) Start(init *dkg.Init) error {
	s.threshold = init.Threshold
	// TODO send stage 1 msg
	return nil
}

func (s *s) ProcessMsg(msg *dkg.SignedMessage) (bool, []dkg.Message, error) {
	// TODO validate msg signature is valid and i'm in the audience

	dataMsg := &KeygenProtocolMsg{}
	if err := dataMsg.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "could not decode protocol msg")
	}

	if s.msgs[dataMsg.RoundNumber] == nil {
		s.msgs[dataMsg.RoundNumber] = []*KeygenProtocolMsg{}
	}
	s.msgs[dataMsg.RoundNumber] = append(s.msgs[dataMsg.RoundNumber], dataMsg)

	switch dataMsg.RoundNumber {
	case KG_R1:
		data := dataMsg.GetRound1Data()
		// assert it's for me
		if len(s.msgs[stubStage1]) >= int(s.threshold) {
			// TODO send stage 2 msg
		}
	case KG_R2:
		data := dataMsg.GetRound2Data()
		if len(s.msgs[stubStage2]) >= int(s.threshold) {
			// TODO send stage 3 msg
		}
	case KG_R3:
		data := dataMsg.GetRound3Data()
		// pass
	case KG_R4:
		data := dataMsg.GetRound4Data()
		if len(s.msgs[stubStage3]) >= int(s.threshold) {
			ret := &dkg.ProtocolOutput{
				Share:       testingutils.Testing4SharesSet().Shares[s.operatorID],
				ValidatorPK: testingutils.Testing4SharesSet().PK.Serialize(),
				OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
					1: testingutils.Testing4SharesSet().Shares[1].GetPublicKey(),
					2: testingutils.Testing4SharesSet().Shares[2].GetPublicKey(),
					3: testingutils.Testing4SharesSet().Shares[3].GetPublicKey(),
					4: testingutils.Testing4SharesSet().Shares[4].GetPublicKey(),
				},
			}
			return true, ret, nil
		}
	}
	return false, nil, nil
}

func (s *s) signDKGMsg(data []byte) *dkg.SignedMessage {
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
