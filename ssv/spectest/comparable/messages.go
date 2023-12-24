package ssvcomparable

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

func SetMessagesInContainer(container ssv.PartialSignatureContainer, messages []*types.SSVMessage) ssv.PartialSignatureContainer {
	for _, ssvMsg := range messages {
		if ssvMsg.MsgType != types.SSVPartialSignatureMsgType {
			continue
		}

		msg := &types.SignedPartialSignatureMessage{}
		if err := msg.Decode(ssvMsg.Data); err != nil {
			panic(err.Error())
		}
		container[msg.Signer] = msg
	}
	return container
}
