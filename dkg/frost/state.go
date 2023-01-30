package frost

import (
	"context"

	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// ProtocolRound is enum for all the rounds in the protocol
type ProtocolRound int

const (
	Uninitialized ProtocolRound = iota
	Preparation
	Round1
	Round2
	KeygenOutput
	Blame
	Timeout
)

var rounds = []ProtocolRound{
	Uninitialized,
	Preparation,
	Round1,
	Round2,
	KeygenOutput,
	Blame,
	Timeout,
}

func (round ProtocolRound) String() string {
	m := map[ProtocolRound]string{
		Uninitialized: "Uninitialized",
		Preparation:   "Preparation",
		Round1:        "Round1",
		Round2:        "Round2",
		KeygenOutput:  "KeygenOutput",
		Blame:         "Blame",
		Timeout:       "Timeout",
	}
	return m[round]
}

// State tracks protocol's current round, stores messages in MsgContainer, stores
// session key and operator's secret shares
type State struct {
	currentRound ProtocolRound
	// underlying participant from frost lib
	participant *frost.DkgParticipant
	// session keypair for other operators to encrypt messages sent to this operator
	sessionSK *ecies.PrivateKey
	// a container to store messages for each round from each operator
	msgContainer IMsgContainer
	// shares generated for each operator using shamir secret sharing in round 1
	operatorShares map[uint32]*bls.SecretKey
	// underlying timer for timeout
	roundTImer *RoundTimer
}

func initState() *State {
	return &State{
		currentRound:   Uninitialized,
		msgContainer:   newMsgContainer(),
		operatorShares: make(map[uint32]*bls.SecretKey),
		roundTImer:     NewRoundTimer(context.Background(), nil),
	}
}

func (state *State) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
	msg, err := state.msgContainer.GetPreparationMsg(operatorID)
	if err != nil {
		return nil, errors.Wrapf(err, "no session pk found for the operator")
	}
	sessionPK, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
	if err != nil {
		return nil, err
	}
	return ecies.Encrypt(sessionPK, data)
}
