package comparable

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

func SetMessages(instance *qbft.Instance, messages []*types.SSVMessage) {
	InitContainers(instance)

	for _, ssvMsg := range messages {
		if ssvMsg.MsgType != types.SSVConsensusMsgType {
			continue
		}

		msg := &qbft.SignedMessage{}
		if err := msg.Decode(ssvMsg.Data); err != nil {
			panic(err.Error())
		}
		setMessage(instance, msg)
	}
}

func SetSignedMessages(instance *qbft.Instance, messages []*qbft.SignedMessage) {
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

func setMessage(instance *qbft.Instance, msg *qbft.SignedMessage) {
	switch msg.Message.MsgType {
	case qbft.ProposalMsgType:
		if instance.State.ProposeContainer.Msgs[msg.Message.Round] == nil {
			instance.State.ProposeContainer.Msgs[msg.Message.Round] = []*qbft.SignedMessage{}
		}
		instance.State.ProposeContainer.Msgs[msg.Message.Round] = append(instance.State.ProposeContainer.Msgs[msg.Message.Round], msg)
	case qbft.PrepareMsgType:
		if instance.State.PrepareContainer.Msgs[msg.Message.Round] == nil {
			instance.State.PrepareContainer.Msgs[msg.Message.Round] = []*qbft.SignedMessage{}
		}
		instance.State.PrepareContainer.Msgs[msg.Message.Round] = append(instance.State.PrepareContainer.Msgs[msg.Message.Round], msg)
	case qbft.CommitMsgType:
		if instance.State.CommitContainer.Msgs[msg.Message.Round] == nil {
			instance.State.CommitContainer.Msgs[msg.Message.Round] = []*qbft.SignedMessage{}
		}
		instance.State.CommitContainer.Msgs[msg.Message.Round] = append(instance.State.CommitContainer.Msgs[msg.Message.Round], msg)
	case qbft.RoundChangeMsgType:
		if instance.State.RoundChangeContainer.Msgs[msg.Message.Round] == nil {
			instance.State.RoundChangeContainer.Msgs[msg.Message.Round] = []*qbft.SignedMessage{}
		}
		instance.State.RoundChangeContainer.Msgs[msg.Message.Round] = append(instance.State.RoundChangeContainer.Msgs[msg.Message.Round], msg)
	}
}
