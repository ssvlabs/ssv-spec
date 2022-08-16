package qbft

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// HistoricalInstanceCapacity represents the upper bound of InstanceContainer a controller can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the controller will process messages for
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
	Identifier []byte
	Height     Height // incremental Height for InstanceContainer
	// StoredInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	StoredInstances InstanceContainer
	Domain          types.DomainType
	Share           *types.Share
	signer          types.SSVSigner
	valueCheck      ProposedValueCheckF
	storage         Storage
	network         Network
	proposerF       ProposerF
}

func NewController(
	identifier []byte,
	share *types.Share,
	domain types.DomainType,
	signer types.SSVSigner,
	valueCheck ProposedValueCheckF,
	storage Storage,
	network Network,
	proposerF ProposerF,
) *Controller {
	return &Controller{
		Identifier:      identifier,
		Height:          -1, // as we bump the height when starting the first instance
		Domain:          domain,
		Share:           share,
		StoredInstances: InstanceContainer{},
		signer:          signer,
		valueCheck:      valueCheck,
		storage:         storage,
		network:         network,
		proposerF:       proposerF,
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(value []byte) error {
	if err := c.canStartInstance(c.Height+1, value); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	c.bumpHeight()
	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(value, c.Height)

	return nil
}

// ProcessMsg processes a new msg, returns true if Decided, non nil byte slice if Decided (Decided value) and error
// Decided returns just once per instance as true, following messages (for example additional commit msgs) will not return Decided true
func (c *Controller) ProcessMsg(msg *SignedMessage) (bool, []byte, error) {
	if !bytes.Equal(c.Identifier, msg.Message.Identifier) {
		return false, nil, errors.New(fmt.Sprintf("message doesn't belong to Identifier"))
	}

	inst := c.InstanceForHeight(msg.Message.Height)
	if inst == nil {
		return false, nil, errors.New(fmt.Sprintf("instance not found"))
	}

	prevDecided, _ := inst.IsDecided()
	decided, decidedValue, aggregatedCommit, err := inst.ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "could not process msg")
	}

	// if previously Decided we do not return Decided true again
	if prevDecided {
		return false, nil, err
	}

	// save the highest Decided
	if !decided {
		return false, nil, nil
	}

	if err := c.saveAndBroadcastDecided(aggregatedCommit); err != nil {
		// TODO - we do not return error, should log?
	}
	return decided, decidedValue, nil
}

func (c *Controller) InstanceForHeight(height Height) *Instance {
	return c.StoredInstances.FindInstance(height)
}

func (c *Controller) bumpHeight() {
	c.Height++
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.Identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() *Instance {
	i := NewInstance(c.GenerateConfig(), c.Share, c.Identifier, c.Height)
	c.StoredInstances.addNewInstance(i)
	return i
}

func (c *Controller) canStartInstance(height Height, value []byte) error {
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
	if err := c.valueCheck(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}

	return nil
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

	config := c.GenerateConfig()
	for _, i := range c.StoredInstances {
		if i != nil {
			i.config = config
		}
	}
	return nil
}

func (c *Controller) saveAndBroadcastDecided(aggregatedCommit *SignedMessage) error {
	if err := c.storage.SaveHighestDecided(aggregatedCommit); err != nil {
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
	if err := c.network.BroadcastDecided(msgToBroadcast); err != nil {
		// We do not return error here, just Log broadcasting error.
		return errors.Wrap(err, "could not broadcast decided")
	}
	return nil
}

func (c *Controller) GenerateConfig() IConfig {
	return &Config{
		Signer:      c.signer,
		SigningPK:   c.Share.ValidatorPubKey,
		Domain:      c.Domain,
		ValueCheckF: c.valueCheck,
		Storage:     c.storage,
		Network:     c.network,
		ProposerF:   c.proposerF,
	}
}
