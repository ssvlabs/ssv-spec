package keygen

import "C"
import (
	"bytes"
	"crypto"
	_ "crypto/sha256"
	"errors"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg/vss"
	"github.com/herumi/bls-eth-go-binary/bls"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Round = uint8

const (
	SECURITY = 256
)

var (
	ErrInvalidRound  = errors.New("invalid round")
	ErrExpectMessage = errors.New("expected message not found")
)

func init() {
	_ = bls.Init(bls.BLS12_381)
	_ = bls.SetETHmode(bls.EthModeDraft07)
}

type Keygen struct {
	Round        Round
	Coefficients vss.Coefficients
	BlindFactor  [32]byte // A random number
	DlogR        *bls.Fr
	PartyI       uint16
	PartyCount   uint16
	skI          *bls.SecretKey
	Round1Msgs   Messages
	Round2Msgs   Messages
	Round3Msgs   Messages
	Round4Msgs   Messages
	Outgoing     Messages
	Output       *LocalKeyShare
	ownShare     *bls.Fr
	inMutex      sync.Mutex
	outMutex     sync.Mutex
}

func NewKeygen(i, t, n uint16) (*Keygen, error) {
	coefficients := vss.CreatePolynomial(int(t + 1))
	bf := MustGetRandomInt(SECURITY)
	kg := &Keygen{
		Round:        0,
		Coefficients: coefficients,
		BlindFactor:  [32]byte{},
		DlogR:        new(bls.Fr),
		PartyI:       i,
		PartyCount:   n,
		skI:          nil,
		Round1Msgs:   make(Messages, n),
		Round2Msgs:   make(Messages, n),
		Round3Msgs:   make(Messages, n),
		Round4Msgs:   make(Messages, n),
		Outgoing:     nil,
		Output:       nil,
		ownShare:     nil,
		inMutex:      sync.Mutex{},
		outMutex:     sync.Mutex{},
	}
	copy(kg.BlindFactor[:], bf.Bytes())
	kg.DlogR.SetByCSPRNG()
	return kg, nil
}

func (k *Keygen) Proceed() error {

	k.inMutex.Lock()
	defer k.inMutex.Unlock()
	var err error
	switch k.Round {
	case 0:
		if err = k.r0CanProceed(); err == nil {
			if err = k.r0Proceed(); err != nil {
				return err
			}
		}
	case 1:
		if err = k.r1CanProceed(); err == nil {
			if err = k.r1Proceed(); err != nil {
				return err
			}
		}
	case 2:
		if err = k.r2CanProceed(); err == nil {
			if err = k.r2Proceed(); err != nil {
				return err
			}
		}
	case 3:
		if err = k.r3CanProceed(); err == nil {
			if err = k.r3Proceed(); err != nil {
				return err
			}
		}
	case 4:
		if err = k.r4CanProceed(); err == nil {
			if err = k.r4Proceed(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid round of state machine: %d", k.Round)
	}
	return nil
}

func (k *Keygen) PushMessage(msg *Message) error {
	if msg == nil || !msg.IsValid() {
		return errors.New("invalid message")
	}
	if msg.Sender <= 0 || msg.Sender > k.PartyCount {
		return errors.New("invalid sender")
	}
	rn, err := msg.GetRoundNumber()
	if err != nil {
		return err
	}
	k.inMutex.Lock()
	defer k.inMutex.Unlock()
	switch rn {
	case 1:
		k.Round1Msgs[msg.Sender-1] = msg
		return nil
	case 2:
		k.Round2Msgs[msg.Sender-1] = msg
		return nil
	case 3:
		k.Round3Msgs[msg.Sender-1] = msg
		return nil
	case 4:
		k.Round4Msgs[msg.Sender-1] = msg
		return nil
	}
	return errors.New("unable to handle message")

}

func (k *Keygen) GetOutgoing() (Messages, error) {
	if success := k.outMutex.TryLock(); success {
		defer k.outMutex.Unlock()
		out := k.Outgoing[:]
		if len(out) > 0 {
			k.trace(fmt.Sprintf("outgoing has %d items", len(out)))
		}
		k.Outgoing = nil
		return out, nil
	} else {
		return nil, errors.New("failed to acquire lock, try again later")
	}
}

func (k *Keygen) pushOutgoing(msg *Message) {
	k.outMutex.Lock()
	defer k.outMutex.Unlock()
	k.Outgoing = append(k.Outgoing, msg)
}

func (k *Keygen) GetDecommitment() [][]byte {
	decomm := make([][]byte, len(k.Coefficients))
	for i, coefficient := range k.Coefficients {
		decomm[i] = bls.CastToSecretKey(&coefficient).GetPublicKey().Serialize()
	}
	return decomm
}

func (k *Keygen) GetCommitment() []byte {

	var data []byte
	decomm := k.GetDecommitment()
	data = append(data, Uint16ToBytes(k.PartyI)...)
	data = append(data, k.BlindFactor[:]...)
	for _, bytes := range decomm {
		data = append(data, bytes...)
	}
	hash := crypto.SHA256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (k *Keygen) VerifyCommitment(r1 Round1Msg, r2 Round2Msg, partyI uint16) bool {
	if len(k.Coefficients) != len(r2.DeCommmitment) {
		return false
	}
	var data []byte
	data = append(data, Uint16ToBytes(partyI)...)
	data = append(data, r2.BlindFactor...)
	for _, bytes := range r2.DeCommmitment {
		data = append(data, bytes...)
	}
	hash := crypto.SHA256.New()
	hash.Write(data)
	comm := hash.Sum(nil)
	return bytes.Compare(r1.Commitment, comm) == 0
}

func (k *Keygen) trace(msg interface{}) {
	log.WithFields(log.Fields{
		"participant": k.PartyI,
	}).Trace(msg)
}
