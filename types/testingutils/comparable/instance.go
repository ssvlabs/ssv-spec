package comparable

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func SetMessages(instance *qbft.Instance, messages []*types.SignedSSVMessage) {
	InitContainers(instance)

	for _, signedSSVMessage := range messages {
		if signedSSVMessage.SSVMessage.MsgType != types.SSVConsensusMsgType {
			continue
		}

		setMessage(instance, signedSSVMessage)
	}
}

func SetSignedMessages(instance *qbft.Instance, messages []*types.SignedSSVMessage) {
	InitContainers(instance)

	for _, msg := range messages {
		setMessage(instance, msg)
	}
}

// InitContainers initializes empty containers for Propose, Prepare, Commit and RoundChange messages
func InitContainers(instance *qbft.Instance) {
	instance.State.ProposeContainer = qbft.NewMsgContainer()
	instance.State.PrepareContainer = qbft.NewMsgContainer()
	instance.State.CommitContainer = qbft.NewMsgContainer()
	instance.State.RoundChangeContainer = qbft.NewMsgContainer()
}

func setMessage(instance *qbft.Instance, signedMsg *types.SignedSSVMessage) {

	msg, err := qbft.DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		panic(err)
	}

	switch msg.MsgType {
	case qbft.ProposalMsgType:
		if instance.State.ProposeContainer.Msgs[msg.Round] == nil {
			instance.State.ProposeContainer.Msgs[msg.Round] = []*qbft.ProcessingMessage{}
		}
		instance.State.ProposeContainer.Msgs[msg.Round] = append(instance.State.ProposeContainer.Msgs[msg.Round], testingutils.ToProcessingMessage(signedMsg))
	case qbft.PrepareMsgType:
		if instance.State.PrepareContainer.Msgs[msg.Round] == nil {
			instance.State.PrepareContainer.Msgs[msg.Round] = []*qbft.ProcessingMessage{}
		}
		instance.State.PrepareContainer.Msgs[msg.Round] = append(instance.State.PrepareContainer.Msgs[msg.Round], testingutils.ToProcessingMessage(signedMsg))
	case qbft.CommitMsgType:
		if instance.State.CommitContainer.Msgs[msg.Round] == nil {
			instance.State.CommitContainer.Msgs[msg.Round] = []*qbft.ProcessingMessage{}
		}
		instance.State.CommitContainer.Msgs[msg.Round] = append(instance.State.CommitContainer.Msgs[msg.Round], testingutils.ToProcessingMessage(signedMsg))
	case qbft.RoundChangeMsgType:
		if instance.State.RoundChangeContainer.Msgs[msg.Round] == nil {
			instance.State.RoundChangeContainer.Msgs[msg.Round] = []*qbft.ProcessingMessage{}
		}
		instance.State.RoundChangeContainer.Msgs[msg.Round] = append(instance.State.RoundChangeContainer.Msgs[msg.Round], testingutils.ToProcessingMessage(signedMsg))
	}
}
