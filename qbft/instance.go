package qbft

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"sync"
)

type ProposedValueCheckF func(data []byte) error
type ProposerF func(state *State, round Round) types.OperatorID

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State  *State
	config IConfig

	processMsgF *types.ThreadSafeF
	startOnce   sync.Once
	StartValue  []byte
}

func NewInstance(
	config IConfig,
	share *types.Share,
	identifier types.MessageID,
	height Height,
) *Instance {
	return &Instance{
		State: &State{
			Share:                share,
			ID:                   identifier,
			Round:                FirstRound,
			Height:               height,
			LastPreparedRound:    NoRound,
			ProposeContainer:     NewMsgContainer(),
			PrepareContainer:     NewMsgContainer(),
			CommitContainer:      NewMsgContainer(),
			RoundChangeContainer: NewMsgContainer(),
		},
		config:      config,
		processMsgF: types.NewThreadSafeF(),
	}
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		// propose if this node is the proposer
		if proposer(i.State, i.GetConfig(), FirstRound) == i.State.Share.OperatorID {
			proposal, err := CreateProposal(i.State, i.config, i.StartValue, nil, nil)
			if err != nil {
				// TODO log
			}

			proposalEncoded, err := proposal.Encode()
			if err != nil {
				return
			}

			msgID := types.PopulateMsgType(i.State.ID, types.ConsensusProposeMsgType)

			broadcastMsg := &types.Message{
				ID:   msgID,
				Data: proposalEncoded,
			}

			if err := i.Broadcast(broadcastMsg); err != nil {
				// TODO - log
			}
		}
	})
}

func (i *Instance) Broadcast(msg *types.Message) error {
	return i.config.GetNetwork().Broadcast(msg)
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msgID types.MessageID, msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
	if err := msg.Validate(msgID.GetMsgType()); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msgID.GetMsgType() {
		case types.ConsensusProposeMsgType:
			return i.uponProposal(msg, i.State.ProposeContainer)
		case types.ConsensusPrepareMsgType:
			return i.uponPrepare(msg, i.State.PrepareContainer, i.State.CommitContainer)
		case types.ConsensusCommitMsgType:
			decided, aggregatedCommit, err = i.UponCommit(msg, i.State.CommitContainer)
			i.State.Decided = decided
			if decided {
				i.State.DecidedValue = msg.Message.Input
			}
			return err
		case types.ConsensusRoundChangeMsgType:
			return i.uponRoundChange(i.StartValue, msg, i.State.RoundChangeContainer, i.config.GetValueCheckF())
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	return i.State.Decided, i.State.DecidedValue
}

// GetConfig returns the instance config
func (i *Instance) GetConfig() IConfig {
	return i.config
}

// GetHeight interface implementation
func (i *Instance) GetHeight() Height {
	return i.State.Height
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}
