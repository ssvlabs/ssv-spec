package qbft

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"sync"
)

type ProposedValueCheckF func(data *Data) error
type ProposerF func(state *State, round Round) types.OperatorID

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State  *State
	config IConfig

	processMsgF *types.ThreadSafeF
	startOnce   sync.Once
	StartValue  *Data
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
func (i *Instance) Start(value *Data, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		// propose if this node is the proposer
		if proposer(i.State, i.GetConfig(), FirstRound) == i.State.Share.OperatorID {
			proposeMsg, err := CreateProposal(i.State, i.config, i.StartValue, nil, nil)
			// nolint
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}

			proposalEncoded, err := proposeMsg.Encode()
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}

			// nolint
			if err := i.Broadcast(proposalEncoded, types.ConsensusProposeMsgType); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}

		if err := i.config.GetNetwork().SyncHighestRoundChange(i.State.ID, i.State.Height); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	})
}

func (i *Instance) Broadcast(data []byte, msgType types.MsgType) error {
	broadcastMsg := &types.Message{
		ID:   types.PopulateMsgType(i.State.ID, msgType),
		Data: data,
	}

	return i.config.GetNetwork().Broadcast(broadcastMsg)
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(
	msgType types.MsgType,
	msg *SignedMessage,
) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
	if err := msg.Validate(msgType); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msgType {
		case types.ConsensusProposeMsgType:
			return i.uponProposal(msg)
		case types.ConsensusPrepareMsgType:
			return i.uponPrepare(msg)
		case types.ConsensusCommitMsgType:
			decided, decidedValue, aggregatedCommit, err = i.UponCommit(msg)
			if decided {
				i.State.Decided = decided
				i.State.DecidedValue = decidedValue
			}
			return err
		case types.ConsensusRoundChangeMsgType:
			return i.uponRoundChange(msg)
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

// GetRoot returns the state's deterministic root
func (i *Instance) GetRoot() ([]byte, error) {
	return i.State.GetRoot()
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}
