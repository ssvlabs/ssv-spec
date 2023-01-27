package alea

import (
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
		errors.Wrap(err, "uponVCBCReady: could not get vcbcReadyData data from signedMessage")
	}

	// get sender ID
	senderID := signedMessage.GetSigners()[0]
	if i.verbose {
		fmt.Println("\tgod senderID:", senderID)
	}

	// check if it's the first time. If not, return. If yes, update map and continue
	if _, exists := i.State.VCBCState.ReceivedReady[vcbcReadyData.Author]; !exists {
		i.State.VCBCState.ReceivedReady[vcbcReadyData.Author] = make(map[Priority]map[types.OperatorID]bool)
	}
	if _, exists := i.State.VCBCState.ReceivedReady[vcbcReadyData.Author][vcbcReadyData.Priority]; !exists {
		i.State.VCBCState.ReceivedReady[vcbcReadyData.Author][vcbcReadyData.Priority] = make(map[types.OperatorID]bool)
	}
	if i.State.VCBCState.ReceivedReady[vcbcReadyData.Author][vcbcReadyData.Priority][senderID] {
		return nil
	} else {
		i.State.VCBCState.ReceivedReady[vcbcReadyData.Author][vcbcReadyData.Priority][senderID] = true
	}

	// If this is the author of the VCBC proposals -> aggregate signature
	if vcbcReadyData.Author == i.State.Share.OperatorID {
		if i.verbose {
			fmt.Println("\tgoing to update W and r")
		}

		// update W, the list of signedMessages to be aggregated later
		i.State.VCBCState.appendToW(vcbcReadyData.Author, vcbcReadyData.Priority, signedMessage)
		W := i.State.VCBCState.getW(vcbcReadyData.Author, vcbcReadyData.Priority)

		// update counter associated with author and priority
		i.State.VCBCState.incrementR(vcbcReadyData.Author, vcbcReadyData.Priority)
		r := i.State.VCBCState.getR(vcbcReadyData.Author, vcbcReadyData.Priority)

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
			i.State.VCBCState.setU(vcbcReadyData.Author, vcbcReadyData.Priority, aggregatedMessage.Signature)

			vcbcFinalMsg, err := CreateVCBCFinal(i.State, i.config, vcbcReadyData.Hash, vcbcReadyData.Priority, aggregatedMessage.Signature, vcbcReadyData.Author)
			if err != nil {
				return errors.Wrap(err, "uponVCBCReady: failed to create VCBCReady message with proof")
			}
			if i.verbose {
				fmt.Println("\tBroadcasting VCBCFinal")
			}
			i.Broadcast(vcbcFinalMsg)

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
