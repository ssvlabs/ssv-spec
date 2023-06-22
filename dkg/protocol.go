package dkg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type ProtocolOutcome struct {
	ProtocolOutput *KeyGenOutput
	KeySignOutput  *KeySignOutput
	BlameOutput    *BlameOutput
}

func (o *ProtocolOutcome) IsFailedWithBlame() (bool, error) {
	if o.ProtocolOutput == nil && o.BlameOutput == nil {
		return false, errors.New("invalid outcome - missing KeyGenOutput and BlameOutput")
	}
	if o.ProtocolOutput != nil && o.BlameOutput != nil {
		return false, errors.New("invalid outcome - has both KeyGenOutput and BlameOutput")
	}
	return o.BlameOutput != nil, nil
}

// KeyGenOutput is the bare minimum output from the protocol
type KeyGenOutput struct {
	Share           *bls.SecretKey
	OperatorPubKeys map[types.OperatorID]*bls.PublicKey
	ValidatorPK     types.ValidatorPK
	Threshold       uint64
}

// KeySignOutput is the output of signature protocol
type KeySignOutput struct {
	RequestID   RequestID
	Signature   []byte
	ValidatorPK types.ValidatorPK
}

func (o *KeySignOutput) GetRoot() ([]byte, error) {
	bytesSolidity, _ := abi.NewType("bytes", "", nil)
	arguments := abi.Arguments{
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
	}
	bytes, err := arguments.Pack(
		[]byte(o.ValidatorPK),
		[]byte(o.Signature),
	)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(bytes), nil
}

// BlameOutput is the output of blame round
type BlameOutput struct {
	Valid        bool
	BlameMessage *SignedMessage
}

// Protocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type Protocol interface {
	Start() error
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *SignedMessage) (bool, *ProtocolOutcome, error)
}
