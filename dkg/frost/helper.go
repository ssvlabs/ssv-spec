package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/pkg/errors"
)

// decode message deserializes bytes data into signed msg and protocol msg
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

// haveSameRoot compares root (sha256 sum of data) of given 2 signed message
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
