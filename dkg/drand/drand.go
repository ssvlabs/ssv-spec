package drand

import (
	"github.com/bloxapp/ssv-spec/dkg"
	dranddkg "github.com/drand/kyber/share/dkg"
	"github.com/pkg/errors"
)

type DRand struct {
	board    *Board
	protocol dranddkg.Protocol
}

func (d *DRand) Start() error {
	d.protocol.Start()
	go func() {}()
	return nil
}

// ProcessMsg returns true and a bls share if finished
func (d *DRand) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
	if err := d.validateSignedMessage(msg); err != nil {
		return false, nil, errors.Wrap(err, "invalid signed message")
	}

	protocolMsg := &Message{}
	if err := protocolMsg.Decode(msg.Message.Data); err != nil {
		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
	}
	if valid := protocolMsg.validate(); !valid {
		return false, nil, errors.New("failed to validate protocol message")
	}

	switch protocolMsg.MsgType {
	case DealBundleMsg:
		if err := d.processDealBundle(*protocolMsg.DealBundle); err != nil {
			return false, nil, errors.Wrap(err, "failed processing deal bundle")
		}
	case ResponseBundleMsg:
	case JustificationBundleMsg:
	default:
		return false, nil, errors.New("unknown protocol message type")
	}

	how to wait and return for result?
}

func (d *DRand) validateSignedMessage(msg *dkg.SignedMessage) error {

}
