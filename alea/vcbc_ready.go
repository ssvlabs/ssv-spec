package alea

import (
	"bytes"
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponVCBCReady(signedMessage *SignedMessage) error {
	if i.verbose {
		fmt.Println("uponVCBCReady")
	}
	// get Data
	vcbcReadyData, err := signedMessage.Message.GetVCBCReadyData()
	if err != nil {
		return errors.Wrap(err, "uponVCBCReady: could not get vcbcReadyData data from signedMessage")
	}

	// get sender ID
	senderID := signedMessage.GetSigners()[0]
	if i.verbose {
		fmt.Println("\tgod senderID:", senderID)
	}

	// check if it's the first time. If not, return. If yes, update map and continue
	if i.State.VCBCState.HasReceivedReady(vcbcReadyData.Author, vcbcReadyData.Priority, senderID) {
		return nil
	} else {
		i.State.VCBCState.SetReceivedReady(vcbcReadyData.Author, vcbcReadyData.Priority, senderID, true)
	}

	// If this is the author of the VCBC proposals -> aggregate signature
	if vcbcReadyData.Author == i.State.Share.OperatorID {
		if i.verbose {
			fmt.Println("\tgoing to update W and r")
		}

		// update W, the list of signedMessages to be aggregated later
		i.State.VCBCState.AppendToW(vcbcReadyData.Author, vcbcReadyData.Priority, signedMessage)
		W := i.State.VCBCState.GetW(vcbcReadyData.Author, vcbcReadyData.Priority)

		// update counter associated with author and priority
		i.State.VCBCState.IncrementR(vcbcReadyData.Author, vcbcReadyData.Priority)
		r := i.State.VCBCState.GetR(vcbcReadyData.Author, vcbcReadyData.Priority)

		if i.verbose {
			fmt.Println("\tW:", W)
		}
		if i.verbose {
			fmt.Println("\tr:", r)
		}

		// if reached quorum, aggregate signatures and broadcast FINAL message
		if r >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println("\treached quorum")
			}

			aggregatedMessage, err := AggregateMsgs(W)
			if err != nil {
				return errors.Wrap(err, "uponVCBCReady: unable to aggregate messages to produce VCBCFinal")
			}
			if i.verbose {
				fmt.Println("\tgot aggregatedMessage")
			}

			aggregatedMsgEncoded, err := aggregatedMessage.Encode()
			if err != nil {
				return errors.Wrap(err, "uponVCBCReady: could not encode aggregated msg")
			}

			i.State.VCBCState.SetU(vcbcReadyData.Author, vcbcReadyData.Priority, aggregatedMsgEncoded)

			vcbcFinalMsg, err := CreateVCBCFinal(i.State, i.config, vcbcReadyData.Hash, vcbcReadyData.Priority, aggregatedMsgEncoded, vcbcReadyData.Author)
			if err != nil {
				return errors.Wrap(err, "uponVCBCReady: failed to create VCBCReady message with proof")
			}
			if i.verbose {
				fmt.Println("\tBroadcasting VCBCFinal")
			}
			i.Broadcast(vcbcFinalMsg)

			i.uponVCBCFinal(vcbcFinalMsg)

		}
	}

	return nil
}

func AggregateMsgs(msgs []*SignedMessage) (*SignedMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("AggregateMsgs: can't aggregate zero msgs")
	}

	var ret *SignedMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.DeepCopy()
		} else {
			if err := ret.Aggregate(m); err != nil {
				return nil, errors.Wrap(err, "AggregateMsgs: could not aggregate msg")
			}
		}
	}
	return ret, nil
}

func isValidVCBCReady(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != VCBCReadyMsgType {
		return errors.New("msg type is not VCBCReadyMsgType")
	}
	if signedMsg.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	VCBCReadyData, err := signedMsg.Message.GetVCBCReadyData()
	if err != nil {
		return errors.Wrap(err, "could not get VCBCReadyData data")
	}
	if err := VCBCReadyData.Validate(); err != nil {
		return errors.Wrap(err, "VCBCReadyData invalid")
	}

	// author
	author := VCBCReadyData.Author
	authorInCommittee := false
	for _, opID := range operators {
		if opID.OperatorID == author {
			authorInCommittee = true
		}
	}
	if !authorInCommittee {
		return errors.New("author (OperatorID) doesn't exist in Committee")
	}

	// priority & hash
	priority := VCBCReadyData.Priority
	if state.VCBCState.HasM(author, priority) {
		localHash, err := GetProposalsHash(state.VCBCState.GetM(author, priority))
		if err != nil {
			return errors.Wrap(err, "could not get local hash")
		}
		if !bytes.Equal(localHash, VCBCReadyData.Hash) {
			return errors.New("existing (priority,author) proposals have different hash")
		}
	}

	return nil
}

func CreateVCBCReady(state *State, config IConfig, hash []byte, priority Priority, author types.OperatorID) (*SignedMessage, error) {
	vcbcReadyData := &VCBCReadyData{
		Hash:     hash,
		Priority: priority,
		// Proof:			proof,
		Author: author,
	}
	dataByts, err := vcbcReadyData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCReady: could not encode vcbcReadyData")
	}
	msg := &Message{
		MsgType:    VCBCReadyMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateVCBCReady: failed signing filler msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
