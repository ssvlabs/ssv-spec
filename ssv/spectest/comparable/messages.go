package ssvcomparable

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

func SetMessagesInContainer(container *ssv.PartialSigContainer, messages []*types.SignedSSVMessage) *ssv.PartialSigContainer {
	for _, signedSSVMsg := range messages {
		ssvMsg := &types.SSVMessage{}
		if err := ssvMsg.Decode(signedSSVMsg.Data); err != nil {
			panic(err)
		}
		if ssvMsg.MsgType != types.SSVPartialSignatureMsgType {
			continue
		}

		msg := &types.SignedPartialSignatureMessage{}
		if err := msg.Decode(ssvMsg.Data); err != nil {
			panic(err.Error())
		}

		for _, partialSigMsg := range msg.Message.Messages {
			root := hex.EncodeToString(partialSigMsg.SigningRoot[:])
			if container.Signatures[root] == nil {
				container.Signatures[root] = map[types.OperatorID][]byte{}
			}
			container.Signatures[root][partialSigMsg.Signer] = partialSigMsg.PartialSignature
		}
	}
	return container
}
