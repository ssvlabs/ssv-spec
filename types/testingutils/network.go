package testingutils

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ssvlabs/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
	OperatorID      types.OperatorID
	OperatorSK      *rsa.PrivateKey
}

func NewTestingNetwork(operatorID types.OperatorID, sk *rsa.PrivateKey) *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
		OperatorID:      operatorID,
		OperatorSK:      sk,
	}
}

func (net *TestingNetwork) Broadcast(msgID types.MessageID, message *types.SignedSSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func ConvertBroadcastedMessagesToSSVMessages(signedMessages []*types.SignedSSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range signedMessages {
		ret = append(ret, msg.SSVMessage)
	}
	return ret
}

// RSA private key cannot be unmarshalled directly as it uses big.Int
func (t *TestingNetwork) MarshalJSON() ([]byte, error) {
	aux := struct {
		BroadcastedMsgs []*types.SignedSSVMessage
		OperatorID      types.OperatorID
		OperatorSK      string
	}{}
	if t != nil {
		aux.BroadcastedMsgs = t.BroadcastedMsgs
		aux.OperatorID = t.OperatorID
		aux.OperatorSK = fmt.Sprintf("N:%s,E:%s,D:%s", t.OperatorSK.N.String(), strconv.Itoa(t.OperatorSK.E), t.OperatorSK.D.String())
	}
	return json.Marshal(aux)
}

func (t *TestingNetwork) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nil
	}

	aux := &struct {
		BroadcastedMsgs []*types.SignedSSVMessage
		OperatorID      types.OperatorID
		OperatorSK      string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.OperatorSK != "" {
		// Parse the private key string format "N:x,E:y,D:z"
		var nStr, eStr, dStr string
		parts := strings.Split(aux.OperatorSK, ",")
		if len(parts) != 3 {
			return fmt.Errorf("invalid private key format: expected 3 parts, got %d", len(parts))
		}
		fmt.Sscanf(parts[0], "N:%s", &nStr)
		fmt.Sscanf(parts[1], "E:%s", &eStr)
		fmt.Sscanf(parts[2], "D:%s", &dStr)

		// Convert strings back to big integers
		n := new(big.Int)
		n.SetString(nStr, 10)

		e, err := strconv.Atoi(eStr)
		if err != nil {
			return err
		}

		d := new(big.Int)
		d.SetString(dStr, 10)

		// Reconstruct the RSA private key
		t.OperatorSK = &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: n,
				E: e,
			},
			D: d,
		}
	}

	t.BroadcastedMsgs = aux.BroadcastedMsgs
	t.OperatorID = aux.OperatorID

	return nil
}
