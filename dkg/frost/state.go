package frost

import (
	"context"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/coinbase/kryptology/pkg/dkg/frost"
	ecies "github.com/ecies/go/v2"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// State tracks protocol's current round, stores messages in MsgContainer, stores
// session key and operator's secret shares
type State struct {
	// round mutex ensures atomic access to current round
	roundMutex   *sync.Mutex
	currentRound common.ProtocolRound

	// underlying participant from frost lib
	participant *frost.DkgParticipant
	// session keypair for other operators to encrypt messages sent to this operator
	sessionSK *ecies.PrivateKey
	// a container to store messages for each round from each operator
	msgContainer common.IMsgContainer
	// shares generated for each operator using shamir secret sharing in round 1
	operatorShares map[uint32]*bls.SecretKey
	// underlying timer for timeout
	roundTimer *RoundTimer
}

func initState() *State {
	return &State{
		currentRound:   common.Uninitialized,
		msgContainer:   common.NewMsgContainer(),
		operatorShares: make(map[uint32]*bls.SecretKey),
		roundTimer:     NewRoundTimer(context.Background(), nil),
		roundMutex:     new(sync.Mutex),
	}
}

func (state *State) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
	msg, err := GetPreparationMsg(state.msgContainer, operatorID)
	if err != nil {
		return nil, errors.Wrapf(err, "no session pk found for the operator")
	}
	sessionPK, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
	if err != nil {
		return nil, err
	}
	return ecies.Encrypt(sessionPK, data)
}

func (state *State) GetCurrentRound() common.ProtocolRound {
	state.roundMutex.Lock()
	defer state.roundMutex.Unlock()

	return state.currentRound
}

func (state *State) SetCurrentRound(round common.ProtocolRound) {
	state.roundMutex.Lock()
	defer state.roundMutex.Unlock()

	state.currentRound = round
}
