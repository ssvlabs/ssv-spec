package frost

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

// createAndBroadcastBlameOfInconsistentMessage function creates a Protocol
// Message to blame an inconsistent message, including the existing message and
// the new message as blame data.
func (fr *Instance) createAndBroadcastBlameOfInconsistentMessage(existingMessage, newMessage *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.SetCurrentRound(common.Blame)

	existingMessageBytes, err := existingMessage.Encode()
	if err != nil {
		return false, nil, err
	}
	newMessageBytes, err := newMessage.Encode()
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: common.Blame,
		BlameMessage: &BlameMessage{
			Type:             InconsistentMessage,
			TargetOperatorID: uint32(newMessage.Signer),
			BlameData:        [][]byte{existingMessageBytes, newMessageBytes},
			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	if err := fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg); err != nil {
		return false, nil, err
	}

	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: bcastMsg,
		},
	}, nil
}

// createAndBroadcastBlameOfInvalidShare function creates a Protocol
// Message to blame an invalid share, including the round 1 message from culprit
// operator
func (fr *Instance) createAndBroadcastBlameOfInvalidShare(culpritOID uint32) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.SetCurrentRound(common.Blame)

	round1Msg, err := fr.state.msgContainer.GetSignedMsg(common.Round1, culpritOID)
	if err != nil {
		return false, nil, err
	}
	round1Bytes, err := round1Msg.Encode()
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: common.Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidShare,
			TargetOperatorID: culpritOID,
			BlameData:        [][]byte{round1Bytes},
			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	if err := fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg); err != nil {
		return false, nil, err
	}

	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: bcastMsg,
		},
	}, nil
}

// createAndBroadcastBlameOfInvalidMessage function creates a Protocol Message to
// blame an invalid message, including the operatorID of the culprit and the
// received signed message.
func (fr *Instance) createAndBroadcastBlameOfInvalidMessage(culpritOID uint32, message *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	fr.state.SetCurrentRound(common.Blame)

	bytes, err := message.Encode()
	if err != nil {
		return false, nil, err
	}

	msg := &ProtocolMsg{
		Round: common.Blame,
		BlameMessage: &BlameMessage{
			Type:             InvalidMessage,
			TargetOperatorID: culpritOID,
			BlameData:        [][]byte{bytes},
			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
		},
	}

	bcastMsg, err := fr.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	if err := fr.config.GetNetwork().BroadcastDKGMessage(bcastMsg); err != nil {
		return false, nil, err
	}

	return true, &dkg.ProtocolOutcome{
		BlameOutput: &dkg.BlameOutput{
			Valid:        true,
			BlameMessage: bcastMsg,
		},
	}, nil
}

// checkBlame checks validity of the blame message as per its blame type
func (fr *Instance) checkBlame(blamerOID uint32, protocolMessage *ProtocolMsg, signedMessage *dkg.SignedMessage) (finished bool, protocolOutcome *dkg.ProtocolOutcome, err error) {
	fr.state.SetCurrentRound(common.Blame)

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

// processBlameTypeInvalidShare checks if blame message for invalid share is
// valid by verifying commitments in blame message
func (fr *Instance) processBlameTypeInvalidShare(blamerOID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}
	if len(blameMessage.BlameData) != 1 {
		return false, errors.New("invalid blame data")
	}
	signedMessage, protocolMessage, err := decodeMessage(blameMessage.BlameData[0])
	if err != nil {
		return false, errors.Wrap(err, "failed to decode signed message")
	}

	if err := fr.validateSignedMessage(signedMessage); err != nil {
		return false, errors.Wrap(err, "failed to Validate signature for blame data")
	}

	blamerPrepMsg, err := GetPreparationMsg(fr.state.msgContainer, blamerOID)
	if err != nil {
		return false, errors.New("unable to retrieve blamer's PreparationMessage")
	}

	blamerSessionSK := ecies.NewPrivateKeyFromBytes(blameMessage.BlamerSessionSk)
	blamerSessionPK := blamerSessionSK.PublicKey.Bytes(true)
	if !bytes.Equal(blamerSessionPK, blamerPrepMsg.SessionPk) {
		return false, errors.New("blame's session pubkey is invalid")
	}

	round1Message := protocolMessage.Round1Message

	shareBytes, err := ecies.Decrypt(blamerSessionSK, round1Message.Shares[blamerOID])
	if err != nil {
		return true, nil
	}
	share := &sharing.ShamirShare{
		Id:    blamerOID,
		Value: shareBytes,
	}

	verifiers := new(sharing.FeldmanVerifier)
	for _, commitmentBytes := range round1Message.Commitment {
		commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
		if err != nil {
			return false, err
		}
		verifiers.Commitments = append(verifiers.Commitments, commitment)
	}

	// verification fails -> blame is valid
	return verifiers.Verify(share) != nil, nil
}

// processBlameTypeInconsistentMessage verifies blame of inconsisstent message
// type by comparing roots of both messages
func (fr *Instance) processBlameTypeInconsistentMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}

	if len(blameMessage.BlameData) != 2 {
		return false, errors.New("invalid blame data")
	}

	signedMsg1, protocolMessage1, err := decodeMessage(blameMessage.BlameData[0])

	if err != nil {
		return false, err
	} else if err := fr.validateSignedMessage(signedMsg1); err != nil {
		return false, errors.Wrap(err, "failed to validate signed message in blame data")
	} else if err := protocolMessage1.Validate(); err != nil {
		return false, errors.New("invalid protocol message")
	}

	signedMsg2, protocolMessage2, err := decodeMessage(blameMessage.BlameData[1])

	if err != nil {
		return false, err
	} else if err := fr.validateSignedMessage(signedMsg2); err != nil {
		return false, errors.Wrap(err, "failed to validate signed message in blame data")
	} else if err := protocolMessage2.Validate(); err != nil {
		return false, errors.New("invalid protocol message")
	}

	if haveSameRoot(signedMsg1, signedMsg2) {
		return false, errors.New("the two messages are consistent")
	}

	if protocolMessage1.Round != protocolMessage2.Round {
		return false, errors.New("the two messages don't belong the the same round")
	}

	return true, nil
}

// processBlameTypeInvalidMessage verifies blame of invalid message type by
// validating signed message and protocol message
func (fr *Instance) processBlameTypeInvalidMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
	if err := blameMessage.Validate(); err != nil {
		return false, errors.Wrap(err, "invalid blame message")
	}
	if len(blameMessage.BlameData) != 1 {
		return false, errors.New("invalid blame data")
	}
	signedMsg, pMsg, err := decodeMessage(blameMessage.BlameData[0])
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
