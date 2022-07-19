package base

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// Output is the last message in every DKG which marks a specific node's end of process
type Output struct {
	// RequestID for the DKG instance (not used for signing)
	RequestID RequestID
	// ShareIndex the 1-based index of the share
	ShareIndex uint16
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// SharePubKeys the public keys corresponding to the shares
	SharePubKeys [][]byte
	// DKGSetSize number of participants in the DKG
	DKGSetSize uint16
	// Threshold DKG threshold for signature reconstruction
	Threshold uint16
	// ValidatorPubKey the resulting public key corresponding to the shared private key
	ValidatorPubKey types.ValidatorPK
	// WithdrawalCredentials same as in Init
	WithdrawalCredentials []byte
	// DepositDataSignature reconstructed signature of DepositMessage according to eth2 spec
	DepositDataSignature types.Signature
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


// PartialDepositData contains a partial deposit data signature
type PartialDepositData struct {
	Signer    types.OperatorID
	Root      []byte
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (msg *PartialDepositData) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *PartialDepositData) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}