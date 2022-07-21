package types

import (
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
)

// Init is the first message in a DKG which initiates a DKG
type Init struct {
	// Nonce is used to differentiate DKG tasks of the same OperatorIDs and WithdrawalCredentials
	Nonce int64
	// OperatorIDs are the operators selected for the DKG
	OperatorIDs []types.OperatorID
	// Threshold DKG threshold for signature reconstruction
	Threshold uint64
	// WithdrawalCredentials used when signing the deposit data
	WithdrawalCredentials []byte
	// Fork is eth2 fork version
	Fork spec.Version
}

func (msg *Init) Validate() error {
	// TODO len(operators == 4,7,10,13
	// threshold equal to 2/3 of 4,7,10,13
	// len(WithdrawalCredentials) is valid
	return nil
}

// Encode returns a msg encoded bytes or error
func (msg *Init) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Init) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
