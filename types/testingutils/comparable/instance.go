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

		msg := &types.SignedSSVMessage{}
		if err := msg.Decode(ssvMsg.Data); err != nil {
			panic(err.Error())
		}
		setMessage(instance, msg)
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

func setMessage(instance *qbft.Instance, msg *types.SignedSSVMessage) {

	qbftMessage := &qbft.Message{}
	if err := qbftMessage.Decode(msg.SSVMessage.Data); err != nil {
		panic(err)
	}

	switch qbftMessage.MsgType {
	case qbft.ProposalMsgType:
		if instance.State.ProposeContainer.Msgs[qbftMessage.Round] == nil {
			instance.State.ProposeContainer.Msgs[qbftMessage.Round] = []*types.SignedSSVMessage{}
		}
		instance.State.ProposeContainer.Msgs[qbftMessage.Round] = append(instance.State.ProposeContainer.Msgs[qbftMessage.Round], msg)
	case qbft.PrepareMsgType:
		if instance.State.PrepareContainer.Msgs[qbftMessage.Round] == nil {
			instance.State.PrepareContainer.Msgs[qbftMessage.Round] = []*types.SignedSSVMessage{}
		}
		instance.State.PrepareContainer.Msgs[qbftMessage.Round] = append(instance.State.PrepareContainer.Msgs[qbftMessage.Round], msg)
	case qbft.CommitMsgType:
		if instance.State.CommitContainer.Msgs[qbftMessage.Round] == nil {
			instance.State.CommitContainer.Msgs[qbftMessage.Round] = []*types.SignedSSVMessage{}
		}
		instance.State.CommitContainer.Msgs[qbftMessage.Round] = append(instance.State.CommitContainer.Msgs[qbftMessage.Round], msg)
	case qbft.RoundChangeMsgType:
		if instance.State.RoundChangeContainer.Msgs[qbftMessage.Round] == nil {
			instance.State.RoundChangeContainer.Msgs[qbftMessage.Round] = []*types.SignedSSVMessage{}
		}
		instance.State.RoundChangeContainer.Msgs[qbftMessage.Round] = append(instance.State.RoundChangeContainer.Msgs[qbftMessage.Round], msg)
	}
}
