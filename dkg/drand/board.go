package drand

import (
	"encoding/json"
	ssvdkg "github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/drand/kyber/share/dkg"
	"github.com/pkg/errors"
)

type Board struct {
	DealsC         chan dkg.DealBundle
	ResponseC      chan dkg.ResponseBundle
	JustificationC chan dkg.JustificationBundle

	operatorID types.OperatorID
	identifier ssvdkg.RequestID
	network    ssvdkg.Network
	signer     types.DKGSigner
	storage    ssvdkg.Storage
}

func (b *Board) PushDeals(bundle *dkg.DealBundle) {
	b.signAndBroadcast(&Message{
		MsgType:    DealBundleMsg,
		DealBundle: bundle,
	})
}

func (b *Board) IncomingDeal() <-chan dkg.DealBundle {
	if b.DealsC == nil {
		b.DealsC = make(chan dkg.DealBundle)
	}
	return b.DealsC
}

func (b *Board) PushResponses(bundle *dkg.ResponseBundle) {
	b.signAndBroadcast(&Message{
		MsgType:        ResponseBundleMsg,
		ResponseBundle: bundle,
	})
}

func (b *Board) IncomingResponse() <-chan dkg.ResponseBundle {
	if b.ResponseC == nil {
		b.ResponseC = make(chan dkg.ResponseBundle)
	}
	return b.ResponseC
}

func (b *Board) PushJustifications(bundle *dkg.JustificationBundle) {
	b.signAndBroadcast(&Message{
		MsgType:             JustificationBundleMsg,
		JustificationBundle: bundle,
	})
}

func (b *Board) IncomingJustification() <-chan dkg.JustificationBundle {
	if b.JustificationC == nil {
		b.JustificationC = make(chan dkg.JustificationBundle)
	}
	return b.JustificationC
}

func (b *Board) signAndBroadcast(message *Message) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "could not marshal message")
	}

	bcastMessage := &ssvdkg.SignedMessage{
		Message: &ssvdkg.Message{
			MsgType:    ssvdkg.ProtocolMsgType,
			Identifier: b.identifier,
			Data:       msgBytes,
		},
		Signer: b.operatorID,
	}

	exist, operator, err := b.storage.GetDKGOperator(b.operatorID)
	if err != nil {
		return errors.Wrap(err, "could not get operator")
	}
	if !exist {
		return errors.Errorf("operator not found")
	}

	sig, err := b.signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
	if err != nil {
		return errors.Wrap(err, "could not sign message")
	}
	bcastMessage.Signature = sig

	return b.network.BroadcastDKGMessage(bcastMessage)
}
