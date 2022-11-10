package qbft

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/types"
)

// HistoricalInstanceCapacity represents the upper bound of InstanceContainer a processmsg can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the processmsg will process messages for
const HistoricalInstanceCapacity int = 5

type InstanceContainer [HistoricalInstanceCapacity]*Instance

func (i InstanceContainer) FindInstance(height Height) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// addNewInstance will add the new instance at index 0, pushing all other stored InstanceContainer one index up (ejecting last one if existing)
func (i *InstanceContainer) addNewInstance(instance *Instance) {
	for idx := HistoricalInstanceCapacity - 1; idx > 0; idx-- {
		i[idx] = i[idx-1]
	}
	i[0] = instance
}

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT InstanceContainer
type Controller struct {
	Identifier types.MessageID
	Height     Height // incremental Height for InstanceContainer
	// StoredInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	StoredInstances InstanceContainer
	// FutureMsgsContainer holds all msgs from a higher height
	FutureMsgsContainer map[types.OperatorID]Height // maps msg signer to height of higher height received msgs
	Domain              types.DomainType
	Share               *types.Share
	config              IConfig
}

func NewController(
	identifier types.MessageID,
	share *types.Share,
	domain types.DomainType,
	config IConfig,
) *Controller {
	return &Controller{
		Identifier: identifier,
		// TODO<olegshmuelov>: fastssz does not support int but only uint.
		// this why the Height type was changed from int64 to uint64
		// The initial value of the height was changed from -1 to math.MaxUint64
		// when we bump the height at the first instance we get 0 (math.MaxUint64 + 1 = 0)
		// if we want to keep the Height to be int64 there are 2 possible options:
		// 1. create a messageSSZ with uint64 instead of Height and transform to message with int64
		// 2. open a pr for fastssz with the implementation that supports int
		Height:              math.MaxUint64, // as we bump the height when starting the first instance
		Domain:              domain,
		Share:               share,
		StoredInstances:     InstanceContainer{},
		FutureMsgsContainer: make(map[types.OperatorID]Height),
		config:              config,
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(inputData *Data) error {
	if err := c.canStartInstance(c.Height+1, inputData); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	c.bumpHeight()
	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(inputData, c.Height)

	return nil
}

// ProcessMsg processes a new msg, returns decided message or error
func (c *Controller) ProcessMsg(msg *types.Message) (*SignedMessage, error) {
	msgID := msg.GetID()
	msgType := msgID.GetMsgType()
	if err := c.baseMsgValidation(msgID); err != nil {
		return nil, errors.Wrap(err, "invalid msg")
	}

	signedMsg := &SignedMessage{}
	if err := signedMsg.Decode(msg.GetData()); err != nil {
		return nil, errors.Wrap(err, "could not decode consensus msg from network msg")
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger highest decided futuremsg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	if isDecidedMsg(c.Share, signedMsg, msgType) {
		return c.UponDecided(signedMsg)
	} else if signedMsg.Message.Height > c.Height {
		return c.UponFutureMsg(msgType, signedMsg)
	} else {
		return c.UponExistingInstanceMsg(msgType, signedMsg)
	}
}

func (c *Controller) UponExistingInstanceMsg(msgType types.MsgType, signedMsg *SignedMessage) (*SignedMessage, error) {
	inst := c.InstanceForHeight(signedMsg.Message.Height)
	if inst == nil {
		return nil, errors.New("instance not found")
	}

	prevDecided, _ := inst.IsDecided()

	decided, _, decidedMsg, err := inst.ProcessMsg(msgType, signedMsg)
	if err != nil {
		return nil, errors.Wrap(err, "could not process msg")
	}

	// if previously Decided we do not return Decided true again
	if prevDecided {
		return nil, err
	}

	// save the highest Decided
	if !decided {
		return nil, nil
	}

	if err := c.saveAndBroadcastDecided(decidedMsg); err != nil {
		// no need to fail processing instance deciding if failed to save/ broadcast
		fmt.Printf("%s\n", err.Error())
	}
	return decidedMsg, nil
}

func (c *Controller) baseMsgValidation(msgID types.MessageID) error {
	// verify msg belongs to controller
	if !msgID.Compare(c.Identifier) {
		return errors.New("message doesn't belong to Identifier")
	}

	return nil
}

func (c *Controller) InstanceForHeight(height Height) *Instance {
	return c.StoredInstances.FindInstance(height)
}

func (c *Controller) bumpHeight() {
	c.Height++
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() types.MessageID {
	return c.Identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() *Instance {
	i := NewInstance(c.GetConfig(), c.Share, c.Identifier, c.Height)
	c.StoredInstances.addNewInstance(i)
	return i
}

func (c *Controller) canStartInstance(height Height, value *Data) error {
	if height > FirstHeight {
		// check prev instance if prev instance is not the first instance
		inst := c.StoredInstances.FindInstance(height - 1)
		if inst == nil {
			return errors.New("could not find previous instance")
		}
		if decided, _ := inst.IsDecided(); !decided {
			return errors.New("previous instance hasn't Decided")
		}
	}

	// check value
	if err := c.GetConfig().GetValueCheckF()(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}

	return nil
}

// GetRoot returns the state's deterministic root
func (c *Controller) GetRoot() ([]byte, error) {
	rootStruct := struct {
		Identifier             types.MessageID
		Height                 Height
		InstanceRoots          [][]byte
		HigherReceivedMessages map[types.OperatorID]Height
		Domain                 types.DomainType
		Share                  *types.Share
	}{
		Identifier:             c.Identifier,
		Height:                 c.Height,
		InstanceRoots:          make([][]byte, len(c.StoredInstances)),
		HigherReceivedMessages: c.FutureMsgsContainer,
		Domain:                 c.Domain,
		Share:                  c.Share,
	}

	for i, inst := range c.StoredInstances {
		if inst != nil {
			r, err := inst.GetRoot()
			if err != nil {
				return nil, errors.Wrap(err, "failed getting instance root")
			}
			rootStruct.InstanceRoots[i] = r
		}
	}

	marshaledRoot, err := json.Marshal(rootStruct)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode implementation
func (c *Controller) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode implementation
func (c *Controller) Decode(data []byte) error {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return errors.Wrap(err, "could not decode controller")
	}

	config := c.GetConfig()
	for _, i := range c.StoredInstances {
		if i != nil {
			i.config = config
		}
	}
	return nil
}

func (c *Controller) saveAndBroadcastDecided(aggregatedCommit *SignedMessage) error {
	if err := c.GetConfig().GetStorage().SaveHighestDecided(c.GetIdentifier(), aggregatedCommit); err != nil {
		return errors.Wrap(err, "could not save decided")
	}

	// Broadcast Decided msg
	decidedEncoded, err := aggregatedCommit.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided message")
	}

	msgID := types.PopulateMsgType(c.Identifier, types.DecidedMsgType)

	broadcastMsg := &types.Message{
		ID:   msgID,
		Data: decidedEncoded,
	}

	if err := c.GetConfig().GetNetwork().Broadcast(broadcastMsg); err != nil {
		// We do not return error here, just Log broadcasting error.
		return errors.Wrap(err, "could not broadcast decided")
	}
	return nil
}

func (c *Controller) GetConfig() IConfig {
	return c.config
}
