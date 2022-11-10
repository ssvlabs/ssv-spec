package drand

import (
	"github.com/bloxapp/ssv-spec/dkg"
	dranddkg "github.com/drand/kyber/share/dkg"
	"github.com/pkg/errors"
)

type DRand struct {
	board    *Board
	config   *dranddkg.Config
	protocol *dranddkg.Protocol
	result   *dranddkg.OptionResult

	operators []uint32
	threshold uint64
}

func (d *DRand) Start() error {
	d.protocol.Start()
	go func() {
		d.protocol.WaitEnd()
	}()
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
		d.board.DealsC <- *protocolMsg.DealBundle
	case ResponseBundleMsg:
		d.board.ResponseC <- *protocolMsg.ResponseBundle
	case JustificationBundleMsg:
		d.board.JustificationC <- *protocolMsg.JustificationBundle
	default:
		return false, nil, errors.New("unknown protocol message type")
	}

	if result := d.getResult(); result != nil {
		if err := d.validateResult(result); err != nil {
			return true, nil, errors.Wrap(err, "invalid result")
		}
		outcome, err := d.getProtocolOutcome(result)
		return true, outcome, err
	}
	return false, nil, nil
}

func (d *DRand) validateSignedMessage(msg *dkg.SignedMessage) error {
	// TODO
	return nil
}
