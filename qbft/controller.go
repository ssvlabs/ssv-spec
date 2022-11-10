package qbft

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT InstanceContainer
type Controller struct {
	Identifier []byte
	Height     Height // incremental Height for InstanceContainer
	// FutureMsgsContainer holds all msgs from a higher height
	FutureMsgsContainer map[types.OperatorID]Height // maps msg signer to height of higher height received msgs
	Domain              types.DomainType
	Share               *types.Share
	config              IConfig
}

func NewController(
	identifier []byte,
	share *types.Share,
	domain types.DomainType,
	config IConfig,
) *Controller {
	return &Controller{
		Identifier:          identifier,
		Height:              -1, // as we bump the height when starting the first instance
		Domain:              domain,
		Share:               share,
		FutureMsgsContainer: make(map[types.OperatorID]Height),
		config:              config,
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(value []byte) error {
	if err := c.canStartInstance(c.Height+1, value); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	c.bumpHeight()
	newInstance, err := c.addAndStoreNewInstance()
	if err != nil {
		return errors.Wrap(err, "failed to add instance")
	}
	newInstance.Start(value, c.Height)

	return nil
}

// ProcessMsg processes a new msg, returns decided message or error
func (c *Controller) ProcessMsg(msg *SignedMessage) (*SignedMessage, error) {
	if err := c.baseMsgValidation(msg); err != nil {
		return nil, errors.Wrap(err, "invalid msg")
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger highest decided futuremsg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	if isDecidedMsg(c.Share, msg) {
		return c.UponDecided(msg)
	} else if msg.Message.Height > c.Height {
		return c.UponFutureMsg(msg)
	} else {
		return c.UponExistingInstanceMsg(msg)
	}
}

func (c *Controller) UponExistingInstanceMsg(msg *SignedMessage) (*SignedMessage, error) {
	inst, err := c.InstanceForHeight(msg.Message.Height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get instance")
	}
	if inst == nil {
		return nil, errors.New("instance not found")
	}

	prevDecided, _ := inst.IsDecided()

	decided, _, decidedMsg, err := inst.ProcessMsg(msg)
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
	return msg, nil
}

func (c *Controller) baseMsgValidation(msg *SignedMessage) error {
	// verify msg belongs to controller
	if !bytes.Equal(c.Identifier, msg.Message.Identifier) {
		return errors.New("message doesn't belong to Identifier")
	}

	return nil
}

func (c *Controller) InstanceForHeight(height Height) (*Instance, error) {
	state, err := c.GetConfig().GetStorage().GetInstanceState(c.Identifier, height)
	if err != nil {
		return nil, err
	}
	return NewInstanceFromState(c.config, state), nil
}

func (c *Controller) bumpHeight() {
	c.Height++
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.Identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() (*Instance, error) {
	i := NewInstance(c.GetConfig(), c.Share, c.Identifier, c.Height)
	err := c.GetConfig().GetStorage().SaveInstanceState(i.State)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *Controller) canStartInstance(height Height, value []byte) error {
	if height > FirstHeight {
		// check prev instance if prev instance is not the first instance
		state, err := c.GetConfig().GetStorage().GetInstanceState(c.Identifier, height-1)
		if err != nil {
			return errors.Wrap(err, "failed to get instance state")
		}
		inst := NewInstanceFromState(c.GetConfig(), state)
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
	marshaledRoot, err := c.Encode()
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
	return json.Unmarshal(data, &c)
}

func (c *Controller) saveAndBroadcastDecided(aggregatedCommit *SignedMessage) error {
	if err := c.GetConfig().GetStorage().SaveHighestDecided(aggregatedCommit); err != nil {
		return errors.Wrap(err, "could not save decided")
	}

	// Broadcast Decided msg
	byts, err := aggregatedCommit.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided message")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   ControllerIdToMessageID(c.Identifier),
		Data:    byts,
	}
	if err := c.GetConfig().GetNetwork().Broadcast(msgToBroadcast); err != nil {
		// We do not return error here, just Log broadcasting error.
		return errors.Wrap(err, "could not broadcast decided")
	}
	return nil
}

func (c *Controller) GetConfig() IConfig {
	return c.config
}
