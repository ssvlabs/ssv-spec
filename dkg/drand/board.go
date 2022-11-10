package drand

import (
	"encoding/json"
	ssvdkg "github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/drand/kyber/share/dkg"
)

type Board struct {
	DealsC     chan dkg.DealBundle
	operatorID types.OperatorID
	identifier ssvdkg.RequestID
	network    ssvdkg.Network
	signer     types.DKGSigner
	storage    ssvdkg.Storage
}

func (b *Board) PushDeals(bundle *dkg.DealBundle) {
	msgBytes, _ := json.Marshal(bundle)
	//if err != nil {
	//	return nil, err
	//}

	bcastMessage := &ssvdkg.SignedMessage{
		Message: &ssvdkg.Message{
			MsgType:    ssvdkg.ProtocolMsgType,
			Identifier: b.identifier,
			Data:       msgBytes,
		},
		Signer: b.operatorID,
	}

	exist, operator, _ := b.storage.GetDKGOperator(b.operatorID)
	//if err != nil {
	//	return nil, err
	//}
	if !exist {
		return
		//return nil, errors.Errorf("operator with id %d not found", fr.state.operatorID)
	}

	sig, _ := b.signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
	//if err != nil {
	//	return nil, err
	//}
	bcastMessage.Signature = sig

	b.network.BroadcastDKGMessage(bcastMessage)
}

func (b *Board) IncomingDeal() <-chan dkg.DealBundle {
	if b.DealsC == nil {
		b.DealsC = make(chan dkg.DealBundle)
	}
	return b.DealsC
}

func (b *Board) PushResponses(*dkg.ResponseBundle) {

}

func (b *Board) IncomingResponse() <-chan dkg.ResponseBundle {

}

func (b *Board) PushJustifications(*dkg.JustificationBundle) {

}

func (b *Board) IncomingJustification() <-chan dkg.JustificationBundle {

}
