package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

func (state *ProtocolState) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
	msg, err := state.msgContainer.GetPreparationMsg(operatorID)
	if err != nil {
		return nil, errors.Wrapf(err, "no session pk found for the operator")
	}
	sessionPK, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
	if err != nil {
		return nil, err
	}
	return ecies.Encrypt(sessionPK, data)
}

func (fr *FROST) toSignedMessage(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
	msgBytes, err := msg.Encode()
	if err != nil {
		return nil, err
	}

	bcastMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: fr.config.identifier,
			Data:       msgBytes,
		},
		Signer: fr.config.operatorID,
	}

	exist, operator, err := fr.storage.GetDKGOperator(fr.config.operatorID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.Errorf("operator with id %d not found", fr.config.operatorID)
	}

	sig, err := fr.signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
	if err != nil {
		return nil, err
	}
	bcastMessage.Signature = sig
	return bcastMessage, nil
}

func (fr *FROST) broadcastDKGMessage(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
	bcastMessage, err := fr.toSignedMessage(msg)
	if err != nil {
		return nil, err
	}
	if _, err = fr.state.msgContainer.SaveMsg(fr.state.currentRound, bcastMessage); err != nil {
		return nil, err
	}
	return bcastMessage, fr.network.BroadcastDKGMessage(bcastMessage)
}

func decodeMessage(data []byte) (*dkg.SignedMessage, *ProtocolMsg, error) {
	signedMsg := &dkg.SignedMessage{}
	if err := signedMsg.Decode(data); err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode signed message")
	}
	pMsg := &ProtocolMsg{}
	if err := pMsg.Decode(signedMsg.Message.Data); err != nil {
		return signedMsg, nil, errors.Wrap(err, "failed to decode protocol msg")
	}
	return signedMsg, pMsg, nil
}

func haveSameRoot(existingMessage, newMessage *dkg.SignedMessage) bool {
	r1, err := existingMessage.GetRoot()
	if err != nil {
		return false
	}
	r2, err := newMessage.GetRoot()
	if err != nil {
		return false
	}
	return bytes.Equal(r1, r2)
}
