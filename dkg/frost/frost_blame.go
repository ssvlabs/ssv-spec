package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

func (fr *FROST) processBlame() (*dkg.BlameOutput, error) {

	for operatorID, msg := range fr.state.msgs[Blame] {

		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(msg.Message.Data); err != nil {
			return nil, errors.New("failed to decode blame data")
		}

		var (
			valid bool
			err   error
		)

		switch protocolMessage.BlameMessage.Type {
		case InvalidShare:
			valid, _ = fr.processBlameTypeInvalidShare(operatorID, protocolMessage.BlameMessage)

		case InconsistentMessage:
			valid, _ = fr.processBlameTypeInconsistentMessage(operatorID, protocolMessage.BlameMessage)
		}

		serializedSigneMessage, err := msg.Encode()
		if err != nil {
			return nil, err
		}
		blameOutput := &dkg.BlameOutput{
			Valid:        valid,
			BlameMessage: serializedSigneMessage,
		}
		return blameOutput, nil
	}

	return nil, nil
}

func (fr *FROST) processBlameTypeInvalidShare(operatorID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {

	if len(blameMessage.BlameData) != 1 {
		return false, errors.New("invalid blame data")
	}
	signedMessage := &dkg.SignedMessage{}
	if err := signedMessage.Decode(blameMessage.BlameData[0]); err != nil {
		return false, errors.Wrap(err, "unable to decode BlameData")
	}
	if err := fr.validateSignedMessage(signedMessage); err != nil {
		return false, errors.Wrap(err, "failed to validate signature for blame data")
	}
	if signedMessage.Message.Identifier != fr.state.identifier {
		return false, errors.New("the message doesn't belong to this session")
	}

	protocolMessage := ProtocolMsg{}
	if err := protocolMessage.Decode(signedMessage.Message.Data); err != nil {
		return false, errors.Wrap(err, "unable to decode protocolMessage")
	}
	round1Message := protocolMessage.Round1Message

	blamesPrepMessage := fr.state.msgs[Preparation][operatorID]
	prepProtocolMessage := &ProtocolMsg{}
	err := prepProtocolMessage.Decode(blamesPrepMessage.Message.Data)
	if err != nil || prepProtocolMessage.PreparationMessage == nil {
		return false, errors.New("unable to decode blamer's PreparationMessage")
	}

	blamerSessionSK := ecies.NewPrivateKeyFromBytes(blameMessage.BlamerSessionSk)
	blamerSessionPK := blamerSessionSK.PublicKey.Bytes(true)
	if !bytes.Equal(blamerSessionPK, prepProtocolMessage.PreparationMessage.SessionPk) {
		return false, errors.New("blame's session pubkey is invalid")
	}

	verifiers := new(sharing.FeldmanVerifier)
	for _, commitmentBytes := range round1Message.Commitment {
		commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
		if err != nil {
			return false, err
		}
		verifiers.Commitments = append(verifiers.Commitments, commitment)
	}

	shareBytes, err := ecies.Decrypt(blamerSessionSK, round1Message.Shares[operatorID])
	if err != nil {
		return false, err
	}

	share := &sharing.ShamirShare{
		Id:    operatorID,
		Value: shareBytes,
	}

	if err = verifiers.Verify(share); err == nil {
		return false, nil
	}
	return true, err
}

func (fr *FROST) processBlameTypeInconsistentMessage(operatorID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {

	if len(blameMessage.BlameData) != 2 {
		return false, errors.New("invalid blame data")
	}

	var originalMessage, newMessage dkg.SignedMessage
	if err := originalMessage.Decode(blameMessage.BlameData[0]); err != nil {
		return false, err
	}
	if err := newMessage.Decode(blameMessage.BlameData[1]); err != nil {
		return false, err
	}

	if err := fr.validateSignedMessage(&originalMessage); err != nil {
		return false, errors.Wrap(err, "failed to validate signature for blame data")
	}

	if err := fr.validateSignedMessage(&newMessage); err != nil {
		return false, errors.Wrap(err, "failed to validate signature for blame data")
	}

	if originalMessage.Message.Identifier != fr.state.identifier {
		return false, errors.New("the message doesn't belong to this session")
	}

	if originalMessage.Message.Identifier != newMessage.Message.Identifier {
		return false, errors.New("the two messages don't belong to same session")
	}

	if fr.haveSameRoot(&originalMessage, &newMessage) {
		return false, errors.New("the two messages are consistent")
	}

	protocolMessage1, protocolMessage2 := &ProtocolMsg{}, &ProtocolMsg{}
	if err := protocolMessage1.Decode(originalMessage.Message.Data); err != nil {
		return false, errors.Wrap(err, "failed to decode protocol msg")
	}
	if err := protocolMessage2.Decode(newMessage.Message.Data); err != nil {
		return false, errors.Wrap(err, "failed to decode protocol msg")
	}
	if protocolMessage1.Round != protocolMessage2.Round {
		return false, errors.New("the two messages don't belong the the same round")
	}

	return true, nil
}

func (fr *FROST) createAndBroadcastBlameOfInconsistentMessage(existingMessage, newMessage *dkg.SignedMessage) error {
	fr.state.currentRound = Blame
	existingMessageBytes, err := existingMessage.Encode()
	if err != nil {
		return err
	}
	newMessageBytes, err := newMessage.Encode()
	if err != nil {
		return err
	}
	msg := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InconsistentMessage,
			TargetOperatorID: uint32(newMessage.Signer),
			BlameData:        [][]byte{existingMessageBytes, newMessageBytes},
			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
		},
	}
	return fr.broadcastDKGMessage(msg)
}

func (fr *FROST) createAndBroadcastBlameOfInvalidShare(operatorID uint32) error {
	fr.state.currentRound = Blame
	round1Bytes, err := fr.state.msgs[Round1][operatorID].Encode()
	if err != nil {
		return err
	}
	msg := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidShare,
			TargetOperatorID: operatorID,
			BlameData:        [][]byte{round1Bytes},
			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
		},
	}
	return fr.broadcastDKGMessage(msg)
}

func (fr *FROST) haveSameRoot(existingMessage, newMessage *dkg.SignedMessage) bool {
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

type ErrInvalidShare struct {
	BlameOutput *dkg.BlameOutput
}

func (e ErrInvalidShare) Error() string {
	return "invalid share"
}
