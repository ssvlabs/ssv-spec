package qbft

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/types"
)

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT InstanceContainer
type Controller struct {
	Identifier []byte
	Height     Height // incremental Height for InstanceContainer
	// StoredInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	StoredInstances InstanceContainer
	CommitteeMember *types.CommitteeMember
	config          IConfig
}

func NewController(identifier []byte, committeeMember *types.CommitteeMember, config IConfig) *Controller {
	return &Controller{
		Identifier:      identifier,
		Height:          FirstHeight,
		CommitteeMember: committeeMember,
		StoredInstances: InstanceContainer{},
		config:          config,
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(height Height, value []byte) error {
	if err := c.GetConfig().GetValueCheckF()(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}

	// can't use <= because of height == 0 case
	if height < c.Height {
		return errors.New("attempting to start an instance with a past height")
	}

	// covers height == 0 case
	if c.StoredInstances.FindInstance(height) != nil {
		return errors.New("instance already running")
	}

	c.Height = height
	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(value, height)

	c.forceStopAllInstanceExceptCurrent()

	return nil
}

// ProcessMsg processes a new msg, returns decided message or error
func (c *Controller) ProcessMsg(signedMessage *types.SignedSSVMessage) (*types.SignedSSVMessage, error) {

	msg, err := NewProcessingMessage(signedMessage)
	if err != nil {
		return nil, errors.New("could not create ProcessingMessage from signed message")
	}

	if err := c.BaseMsgValidation(msg); err != nil {
		return nil, errors.Wrap(err, "invalid msg")
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger highest decided futuremsg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	isDecided, err := IsDecidedMsg(c.CommitteeMember, msg)
	if err != nil {
		return nil, err
	}
	if isDecided {
		return c.UponDecided(msg)
	}

	isFuture, err := c.isFutureMessage(msg)
	if err != nil {
		return nil, err
	}
	if isFuture {
		return nil, fmt.Errorf("future msg from height, could not process")
	}

	return c.UponExistingInstanceMsg(msg)

}

func (c *Controller) UponExistingInstanceMsg(msg *ProcessingMessage) (*types.SignedSSVMessage, error) {

	inst := c.InstanceForHeight(msg.QBFTMessage.Height)
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

	if err := c.broadcastDecided(decidedMsg); err != nil {
		// no need to fail processing instance deciding if failed to save/ broadcast
		fmt.Printf("%s\n", err.Error())
	}
	return decidedMsg, nil
}

// BaseMsgValidation returns error if msg is invalid (base validation)
func (c *Controller) BaseMsgValidation(msg *ProcessingMessage) error {
	// verify msg belongs to controller
	if !bytes.Equal(c.Identifier, msg.QBFTMessage.Identifier) {
		return errors.New("message doesn't belong to Identifier")
	}
	return nil
}

func (c *Controller) InstanceForHeight(height Height) *Instance {
	return c.StoredInstances.FindInstance(height)
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.Identifier
}

// isFutureMessage returns true if message height is from a future instance.
// It takes into consideration a special case where FirstHeight didn't start but  c.Height == FirstHeight (since we bump height on start instance)
func (c *Controller) isFutureMessage(msg *ProcessingMessage) (bool, error) {
	if c.Height == FirstHeight && c.StoredInstances.FindInstance(c.Height) == nil {
		return true, nil
	}
	return msg.QBFTMessage.Height > c.Height, nil
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() *Instance {
	i := NewInstance(c.GetConfig(), c.CommitteeMember, c.Identifier, c.Height)
	c.StoredInstances.addNewInstance(i)
	return i
}

func (c *Controller) forceStopAllInstanceExceptCurrent() {
	for _, i := range c.StoredInstances {
		if i.State.Height != c.Height {
			i.ForceStop()
		}
	}
}

func (c *Controller) broadcastDecided(aggregatedCommit *types.SignedSSVMessage) error {

	if err := c.GetConfig().GetNetwork().Broadcast(aggregatedCommit.SSVMessage.GetID(), aggregatedCommit); err != nil {
		// We do not return error here, just Log broadcasting error.
		return errors.Wrap(err, "could not broadcast decided")
	}
	return nil
}

func (c *Controller) GetConfig() IConfig {
	return c.config
}
