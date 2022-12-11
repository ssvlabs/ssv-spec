package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

func (fr *FROST) createAndBroadcastBlameOfInconsistentMessage(existingMessage, newMessage *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.currentRound = Blame

	existingMessageBytes, err := existingMessage.Encode()
	if err != nil {
		return false, nil, err
	}
	newMessageBytes, err := newMessage.Encode()
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InconsistentMessage,
			TargetOperatorID: uint32(newMessage.Signer),
			BlameData:        [][]byte{existingMessageBytes, newMessageBytes},
			BlamerSessionSk:  fr.config.sessionSK.Bytes(),
		},
	}

	signedMessage, err := fr.broadcastDKGMessage(msg)
	if err != nil {
		return false, nil, err
	}

	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: signedMessage,
		},
	}, nil
}

func (fr *FROST) createAndBroadcastBlameOfInvalidShare(culpritOID uint32) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.currentRound = Blame

	round1Bytes, err := fr.state.msgs[Round1][culpritOID].Encode()
	if err != nil {
		return false, nil, err
	}
	msg := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidShare,
			TargetOperatorID: culpritOID,
			BlameData:        [][]byte{round1Bytes},
			BlamerSessionSk:  fr.config.sessionSK.Bytes(),
		},
	}
	signedMessage, err := fr.broadcastDKGMessage(msg)
	if err != nil {
		return false, nil, err
	}
	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: signedMessage,
		},
	}, nil
}

func (fr *FROST) createAndBroadcastBlameOfInvalidMessage(culpritOID uint32, message *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.currentRound = Blame

	bytes, err := message.Encode()
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidMessage,
			TargetOperatorID: culpritOID,
			BlameData:        [][]byte{bytes},
			BlamerSessionSk:  fr.config.sessionSK.Bytes(),
		},
	}

	signedMsg, err := fr.broadcastDKGMessage(msg)
	if err != nil {
		return false, nil, err
	}

	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: signedMsg,
		},
	}, nil
}

func (fr *FROST) checkBlame(blamerOID uint32, protocolMessage *ProtocolMsg, signedMessage *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {
	fr.state.currentRound = Blame

	var valid bool
	switch protocolMessage.BlameMessage.Type {
	case InvalidShare:
		valid, err = fr.processBlameTypeInvalidShare(blamerOID, protocolMessage.BlameMessage)
	case InconsistentMessage:
		valid, err = fr.processBlameTypeInconsistentMessage(protocolMessage.BlameMessage)
	case InvalidMessage:
		valid, err = fr.processBlameTypeInvalidMessage(protocolMessage.BlameMessage)
	default:
		valid, err = false, errors.New("unrecognized blame type")
	}
	if err != nil {
		return false, nil, err
	}
	return true, &dkg.ProtocolOutcome{BlameOutput: &dkg.BlameOutput{Valid: valid, BlameMessage: signedMessage}}, nil
}

func (fr *FROST) processBlameTypeInvalidShare(blamerOID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}
	if len(blameMessage.BlameData) != 1 {
		return false, errors.New("invalid blame data")
	}
	signedMessage, protocolMessage, err := fr.decodeMessage(blameMessage.BlameData[0])
	if err != nil {
		return false, errors.Wrap(err, "failed to decode signed message")
	}

	if err := fr.validateSignedMessage(signedMessage); err != nil {
		return false, errors.Wrap(err, "failed to Validate signature for blame data")
	}

	round1Message := protocolMessage.Round1Message

	blamerPrepSignedMessage := fr.state.msgs[Preparation][blamerOID]
	blamerPrepProtocolMessage := &ProtocolMsg{}
	err = blamerPrepProtocolMessage.Decode(blamerPrepSignedMessage.Message.Data)
	if err != nil || blamerPrepProtocolMessage.PreparationMessage == nil {
		return false, errors.New("unable to decode blamer's PreparationMessage")
	}

	blamerSessionSK := ecies.NewPrivateKeyFromBytes(blameMessage.BlamerSessionSk)
	blamerSessionPK := blamerSessionSK.PublicKey.Bytes(true)
	if !bytes.Equal(blamerSessionPK, blamerPrepProtocolMessage.PreparationMessage.SessionPk) {
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

	shareBytes, err := ecies.Decrypt(blamerSessionSK, round1Message.Shares[blamerOID])
	if err != nil {
		return true, nil
	}

	share := &sharing.ShamirShare{
		Id:    blamerOID,
		Value: shareBytes,
	}

	if err = verifiers.Verify(share); err != nil {
		return true, nil
	}
	return false, err
}

func (fr *FROST) processBlameTypeInconsistentMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}

	if len(blameMessage.BlameData) != 2 {
		return false, errors.New("invalid blame data")
	}

	signedMsg1, protocolMessage1, err := fr.decodeMessage(blameMessage.BlameData[0])

	if err != nil {
		return false, err
	} else if err := fr.validateSignedMessage(signedMsg1); err != nil {
		return false, errors.Wrap(err, "failed to validate signed message in blame data")
	} else if err := protocolMessage1.Validate(); err != nil {
		return false, errors.New("invalid protocol message")
	}

	signedMsg2, protocolMessage2, err := fr.decodeMessage(blameMessage.BlameData[1])

	if err != nil {
		return false, err
	} else if err := fr.validateSignedMessage(signedMsg2); err != nil {
		return false, errors.Wrap(err, "failed to validate signed message in blame data")
	} else if err := protocolMessage2.Validate(); err != nil {
		return false, errors.New("invalid protocol message")
	}

	if fr.haveSameRoot(signedMsg1, signedMsg2) {
		return false, errors.New("the two messages are consistent")
	}

	if protocolMessage1.Round != protocolMessage2.Round {
		return false, errors.New("the two messages don't belong the the same round")
	}

	return true, nil
}

func (fr *FROST) processBlameTypeInvalidMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}
	if len(blameMessage.BlameData) != 1 {
		return false, errors.New("invalid blame data")
	}
	signedMsg, pMsg, err := fr.decodeMessage(blameMessage.BlameData[0])
	if err != nil {
		return false, err
	} else if err := fr.validateSignedMessage(signedMsg); err != nil {
		return false, errors.Wrap(err, "failed to validate signed message in blame data")
	}

	err = pMsg.Validate()
	if err != nil {
		return true, nil
	}
	return false, errors.New("message is valid")
}

type ErrBlame struct {
	BlameOutput *dkg.BlameOutput
}

func (e ErrBlame) Error() string {
	return "detected and processed blame"
}
