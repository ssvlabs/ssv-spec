package keygen

import (
	"errors"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/base"
	log "github.com/sirupsen/logrus"
	"time"
)

type Runner struct {
	Keygen   *Keygen
	incoming <-chan base.Message
	outgoing chan<- base.Message
}

func NewRunner(identifier dkg.RequestID, i, t, n uint64, incoming <-chan base.Message, outgoing chan<- base.Message) (*Runner, error) {
	kg, err := NewKeygen(identifier[:], i, t, n)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Keygen:   kg,
		incoming: incoming,
		outgoing: outgoing,
	}, nil
}

func (r *Runner) Initialize() error {
	if r.Keygen.Round == 0 {
		return r.Keygen.Proceed()
	}
	return errors.New("state machine is not initializable")
}

func (r *Runner) ProcessLoop() {
	finished := r.Keygen != nil && r.Keygen.Output != nil
	for !finished {
		select {
		case msg, ok := <-r.incoming:
			if ok {
				parsed := &ParsedMessage{}
				if err := parsed.FromBase(&msg); err == nil {
					r.Keygen.PushMessage(parsed)
				} else {
					// TODO: Log error
				}

			}
		case <-time.After(1 * time.Second):
			finished = r.Keygen.Output != nil
			_ = r.Keygen.Proceed()
			if outgoing, _ := r.Keygen.GetOutgoing(); outgoing != nil {
				for _, out := range outgoing {
					if msg, err := out.ToBase(); err == nil {
						r.outgoing <- *msg
					} else {
						// TODO: Standardize log error
						log.Errorf("err: %v", err)
					}

				}
			}
			if finished {
				break
			}
		}
	}
}

func (r *Runner) trace(funcName string, result interface{}) {
	log.WithFields(log.Fields{
		"participant": r.Keygen.PartyI,
		"funcName":    funcName,
		"result":      result,
	}).Trace("statusCheck")
}
