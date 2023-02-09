package hbbft

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type ProposedValueCheckF func(data []byte) error
type ProposerF func(state *State, round Round) types.OperatorID
type CoinF func(round Round) byte

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State  *State
	config IConfig

	processMsgF *types.ThreadSafeF
	startOnce   sync.Once
	StartValue  []byte
	verbose     bool
}

func NewInstance(
	config IConfig,
	share *types.Share,
	identifier []byte,
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
			HBBFTState:           NewHBBFTState(),
		},
		config:      config,
		processMsgF: types.NewThreadSafeF(),
		verbose:     false,
	}
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		go i.HbbftLoop()

	})
}

func (i *Instance) HbbftLoop() {

	// take int(b/n) transactions from buffer
	B := i.State.HBBFTState.GetBatchSize()
	N := len(i.State.Share.Committee)
	num_transactions := int(B / N)
	for {
		proposed := i.State.HBBFTState.GetRandomTransacations(num_transactions)

		// proposed_encrypted, err := i.State.HBBFTState.EncryptProposed(proposed, i.State.Share.ValidatorPubKey)
		proposed_encrypted, err := json.Marshal(proposed)
		if err != nil {
			fmt.Println("HbbftLoop: could not encrypt proposed value.", err)
			os.Exit(1)
		}

		v := i.State.HBBFTState.RunACS(proposed_encrypted)

		y := make(map[types.OperatorID][]*TransactionData, len(v))
		for j, vj := range v {

			// ej := i.State.HBBFTState.DecryptShare(vj, i.config, i.State.Share.DomainType, i.config.GetSignatureDomainType(), i.State.Share.SharePubKey)
			ej := vj

			decMsg, err := CreateDEC(i.State, i.config, i.State.HBBFTState.Round, j, i.State.Share.OperatorID, ej)
			if err != nil {
				fmt.Println("HbbftLoop: Error: could not create DEC msg.", err)
				os.Exit(1)
			}
			i.Broadcast(decMsg)
			for {
				if i.State.HBBFTState.GetLenDECmsgs(i.State.HBBFTState.Round, j) >= i.State.Share.PartialQuorum {
					break
				}
			}
			// y[j] = i.State.HBBFTState.DecriptDECSet(i.State.HBBFTState.GetDECs(i.State.HBBFTState.Round, j), i.State.Share.ValidatorPubKey)
			// FIX ME ! adapted version, do not run DEC on the set, just get any element
			DECmap := i.State.HBBFTState.GetDECMsgsMap(i.State.HBBFTState.Round, j)
			for _, value := range DECmap {
				err = json.Unmarshal(value, y[j])
				break
			}
		}
		// FIX ME ! sort block?
		// block_r := i.State.HBBFTState.SortY(y)
		i.State.HBBFTState.StoreBlockAndUpdateBuffer(y)

		i.State.HBBFTState.IncrementRound()
	}
}

func (i *Instance) Broadcast(msg *SignedMessage) error {
	byts, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode message")
	}

	msgID := types.MessageID{}
	copy(msgID[:], msg.Message.Identifier)

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    byts,
	}
	return i.config.GetNetwork().Broadcast(msgToBroadcast)
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
	if err := i.BaseMsgValidation(msg); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msg.Message.MsgType {
		case TransactionMsgType:
			return i.uponTransaction(msg)
		case ABAInitMsgType:
			return i.uponABAInit(msg)
		case ABAAuxMsgType:
			return i.uponABAAux(msg)
		case ABAConfMsgType:
			return i.uponABAConf(msg)
		case ABAFinishMsgType:
			return i.uponABAFinish(msg)
		case RBCValMsgType:
			return i.uponRBCVal(msg)
		case RBCEchoMsgType:
			return i.uponRBCEcho(msg)
		case RBCReadyMsgType:
			return i.uponRBCReady(msg)
		case DECMsgType:
			return i.uponDEC(msg)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

func (i *Instance) BaseMsgValidation(msg *SignedMessage) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid signed message")
	}

	if msg.Message.Round < i.State.Round {
		return errors.New("past round")
	}

	switch msg.Message.MsgType {
	case TransactionMsgType:
		return isValidTransaction(i.State, i.config, msg, i.State.Share.Committee)
	case ABAInitMsgType:
		return isValidABAInit(i.State, i.config, msg, i.State.Share.Committee)
	case ABAAuxMsgType:
		return isValidABAAux(i.State, i.config, msg, i.State.Share.Committee)
	case ABAConfMsgType:
		return isValidABAConf(i.State, i.config, msg, i.State.Share.Committee)
	case ABAFinishMsgType:
		return isValidABAFinish(i.State, i.config, msg, i.State.Share.Committee)
	case RBCValMsgType:
		return isValidRBCVal(i.State, i.config, msg, i.State.Share.Committee)
	case RBCEchoMsgType:
		return isValidRBCEcho(i.State, i.config, msg, i.State.Share.Committee)
	case RBCReadyMsgType:
		return isValidRBCReady(i.State, i.config, msg, i.State.Share.Committee)
	case DECMsgType:
		return isValidDEC(i.State, i.config, msg, i.State.Share.Committee)
	default:
		return errors.New("signed message type not supported")
	}
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	if state := i.State; state != nil {
		return state.Decided, state.DecidedValue
	}
	return false, nil
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
