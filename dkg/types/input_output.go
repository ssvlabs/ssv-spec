package types

import (
	"crypto/ecdsa"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Output is the last message in every DKG which marks a specific node's end of process
type Output struct {
	// RequestID for the DKG instance (not used for signing)
	RequestID RequestID
	// ShareIndex the 1-based index of the share
	ShareIndex uint64
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// SharePubKeys the public keys corresponding to the shares
	SharePubKeys [][]byte
	// DKGSetSize number of participants in the DKG
	DKGSetSize uint64
	// Threshold DKG threshold for signature reconstruction
	Threshold uint64
	// ValidatorPubKey the resulting public key corresponding to the shared private key
	ValidatorPubKey types.ValidatorPK
	// WithdrawalCredentials same as in Init
	WithdrawalCredentials []byte
	// DepositDataSignature reconstructed signature of DepositMessage according to eth2 spec
	DepositDataSignature types.Signature
}

func (o *Output) ToExtendedDepositData(forkVersion spec.Version, cliVersion string) (*types.ExtendedDepositData, error) {
	_, depData, err := types.GenerateETHDepositData(
		o.ValidatorPubKey,
		o.WithdrawalCredentials,
		forkVersion,
		types.DomainDeposit,
	)
	if err != nil {
		return nil, err
	}
	copy(depData.DepositData.Signature[:], o.DepositDataSignature)
	depData.CliVersion = cliVersion
	return depData, nil
}

func (o *Output) GetRoot() ([]byte, error) {
	bytesSolidity, _ := abi.NewType("bytes", "", nil)

	// TODO: Include RequestID, SharePubKeys and ShareIndex
	arguments := abi.Arguments{
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
	}

	bytes, _ := arguments.Pack(
		o.EncryptedShare,
		//o.SharePubKeys, // TODO: Add later
		o.ValidatorPubKey,
		o.DepositDataSignature,
	)

	return crypto.Keccak256(bytes), nil
}

// Encode returns a msg encoded bytes or error
func (o *Output) Encode() ([]byte, error) {
	return json.Marshal(o)
}

// Decode returns error if decoding failed
func (o *Output) Decode(data []byte) error {
	return json.Unmarshal(data, o)
}

type SignedOutput struct {
	// Data signed
	Data *Output
	// Signer Operator ID which signed
	Signer types.OperatorID
	// Signature over Data.GetRoot()
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (msg *SignedOutput) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *SignedOutput) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

func SignOutput(output *Output, privKey *ecdsa.PrivateKey) (types.Signature, error) {
	root, err := output.GetRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not get root from output message")
	}

	return crypto.Sign(root, privKey)
}

// TODO: What's the difference / intention of this vs Output.
type KeygenOutput struct {
	Index           uint16
	Threshold       uint16
	ShareCount      uint16
	PublicKey       []byte
	SecretShare     []byte // TODO: Maybe only keep the encrypted version?
	SharePublicKeys [][]byte
}

// Encode returns a msg encoded bytes or error
func (d *KeygenOutput) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenOutput) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Encode returns a msg encoded bytes or error
func (msg *PartialSigMsgBody) Encode() ([]byte, error) {
	return proto.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *PartialSigMsgBody) Decode(data []byte) error {
	return proto.Unmarshal(data, msg)
}
