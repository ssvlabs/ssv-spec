package keygen

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type Runner struct {
	Keygen   *Keygen
	incoming <-chan Message
	outgoing chan<- Message
}

func NewRunner(i, t, n uint16, incoming <-chan Message, outgoing chan<- Message) (*Runner, error) {
	kg, err := NewKeygen(i, t, n)
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
				r.Keygen.PushMessage(&msg)
			}
		case <-time.After(1 * time.Second):
			finished = r.Keygen.Output != nil
			_ = r.Keygen.Proceed()
			if outgoing, _ := r.Keygen.GetOutgoing(); outgoing != nil {
				for _, out := range outgoing {
					r.outgoing <- *out
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
