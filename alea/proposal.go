package alea

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(signedProposal *SignedMessage, proposeMsgContainer *MsgContainer) error {

	if i.verbose {
		fmt.Println("uponProposal")
	}

	// get Data
	proposalDataReceived, err := signedProposal.Message.GetProposalData()
	if err != nil {
		return errors.Wrap(err, "uponProposal: could not get proposal data from signedProposal")
	}
	if i.verbose {
		fmt.Println("\tData:", proposalDataReceived.Data)
	}

	// check if message has been already delivered
	if i.State.Delivered.HasProposal(proposalDataReceived) {
		if i.verbose {
			fmt.Println("\tDelivered Queue:", i.State.Delivered, ". Has the proposal:", proposalDataReceived, i.State.Delivered.HasProposal(proposalDataReceived))
		}
		return errors.New("proposal already delivered")
	}

	// Add message to container
	proposeMsgContainer.AddMsg(signedProposal)

	// add to vcbc state
	i.State.VCBCState.AppendToM(i.State.Share.OperatorID, i.State.VCBCState.Priority, proposalDataReceived)

	if i.verbose {
		fmt.Println("\tadded to own VCBC messages")
	}

	// Check if container has less maximum size. If so, returns
	if len(i.State.VCBCState.GetM(i.State.Share.OperatorID, i.State.VCBCState.Priority)) < i.State.BatchSize {
		if i.verbose {
			fmt.Println("\tdidn't reach batch size")
		}
		return nil
	}

	// else, starts VCBC
	if i.verbose {
		fmt.Println("\tReached batch size. Starting vcbc with priority:", i.State.VCBCState.Priority)
	}
	i.StartVCBC(i.State.VCBCState.Priority)

	// Increment priority
	i.State.VCBCState.Priority += 1
	return nil
}

// Encode returns the list encoded bytes or error
func EncodeProposals(proposals []*ProposalData) ([]byte, error) {
	return json.Marshal(proposals)
}

// Decode returns error if decoding failed
func DecodeProposals(data []byte) ([]*ProposalData, error) {
	proposals := make([]*ProposalData, 0)
	err := json.Unmarshal(data, &proposals)
	if err != nil {
		return nil, errors.Wrap(err, "DecodeProposals: could not unmarshal proposals")
	}
	return proposals, nil
}

// GetHash returns the SHA-256 hash
func GetProposalsHash(proposals []*ProposalData) ([]byte, error) {
	encoded, err := EncodeProposals(proposals)
	if err != nil {
		return nil, errors.Wrap(err, "GetProposalsHash: could not encode proposals")
	}
	ret := sha256.Sum256(encoded)
	return ret[:], nil
}

func isValidProposal(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedProposal.Signature.VerifyByOperators(signedProposal, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	proposalData, err := signedProposal.Message.GetProposalData()
	if err != nil {
		return errors.Wrap(err, "could not get proposal data")
	}
	if err := proposalData.Validate(); err != nil {
		return errors.Wrap(err, "proposalData invalid")
	}

	return nil
}

// CreateProposal
func CreateProposal(state *State, config IConfig, value []byte /*, roundChanges, prepares []*SignedMessage*/) (*SignedMessage, error) {
	proposalData := &ProposalData{
		Data: value,
		// RoundChangeJustification: roundChanges,
		// PrepareJustification:     prepares,
	}
	dataByts, err := proposalData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposal data")
	}
	msg := &Message{
		MsgType:    ProposalMsgType,
		Height:     state.Height,
		Round:      state.AleaDefaultRound,
		Identifier: state.ID,
		Data:       dataByts,
	}

	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
