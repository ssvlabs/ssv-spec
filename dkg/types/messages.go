package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
)

type MsgType int

type SessionId = []byte

type Signable interface {
	types.Root
	SetSignature([]byte) error
}

func toRequestID(id SessionId) RequestID {
	// TODO: Check size
	reqID := new(RequestID)
	copy(reqID[:], id)
	return *reqID
}

// Encode returns a msg encoded bytes or error
func (x *MessageHeader) RequestID() RequestID {
	return toRequestID(x.SessionId)
}

// Encode returns a msg encoded bytes or error
func (x *Message) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

// Decode returns error if decoding failed
func (x *Message) Decode(data []byte) error {
	return proto.Unmarshal(data, x)
}

func (x *Message) Validate() error {
	// TODO: Implement
	return nil
}

func (x *Message) GetRoot() ([]byte, error) {
	raw, err := x.Encode()
	if err != nil {
		return nil, err
	}
	newMsg := &Message{}
	err = newMsg.Decode(raw)
	if err != nil {
		return nil, err
	}
	newMsg.Signature = nil
	bytes, err := newMsg.Encode()
	if err != nil {
		return nil, err
	}
	var root []byte
	rootFixed := sha256.Sum256(bytes)
	copy(root, rootFixed[:])

	return root, nil
}

func (x *Message) SetSignature(bytes []byte) error {
	x.Signature = bytes
	return nil
}

func (x *Init) Validate() error {
	// TODO len(operators == 4,7,10,13
	// threshold equal to 2/3 of 4,7,10,13
	// len(WithdrawalCredentials) is valid
	return nil
}

// Encode returns a msg encoded bytes or error
func (x *Init) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

// Decode returns error if decoding failed
func (x *Init) Decode(data []byte) error {
	return json.Unmarshal(data, x)
}

func (x *ParsedInitMessage) FromBase(base *Message) error {
	raw, err := proto.Marshal(base)
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, x)
}

func (x *ParsedInitMessage) ToBase() (*Message, error) {
	raw, err := proto.Marshal(x)
	if err != nil {
		return nil, err
	}
	base := &Message{}
	err = proto.Unmarshal(raw, base)
	if err != nil {
		return nil, err
	}
	return base, nil
}

// SignedMessage Deprecated
type SignedMessage struct {
	Message   *Message
	Signer    types.OperatorID
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, signedMsg)
}

func (signedMsg *SignedMessage) Validate() error {
	// TODO len(sig) == ecdsa sig lenth

	return signedMsg.Message.Validate()
}

func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
	return signedMsg.Message.GetRoot()
}

func (x *LocalKeyShare) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *LocalKeyShare) Decode(data []byte) error {
	return proto.Unmarshal(data, x)
}

// Encode returns a msg encoded bytes or error
func (x *PartialSigMsgBody) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

// Decode returns error if decoding failed
func (x *PartialSigMsgBody) Decode(data []byte) error {
	return proto.Unmarshal(data, x)
}

func (x *ParsedPartialSigMessage) FromBase(base *Message) error {
	raw, err := proto.Marshal(base)
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, x)
}

func (x *ParsedPartialSigMessage) ToBase() (*Message, error) {
	raw, err := proto.Marshal(x)
	if err != nil {
		return nil, err
	}
	base := &Message{}
	err = proto.Unmarshal(raw, base)
	if err != nil {
		return nil, err
	}
	return base, nil
}

func (x *SignedDepositDataMsgBody) ToExtendedDepositData(forkVersion spec.Version, cliVersion string) (*types.ExtendedDepositData, error) {
	_, depData, err := types.GenerateETHDepositData(
		x.ValidatorPublicKey,
		x.WithdrawalCredentials,
		forkVersion,
		types.DomainDeposit,
	)
	if err != nil {
		return nil, err
	}
	copy(depData.DepositData.Signature[:], x.DepositDataSignature)
	depData.CliVersion = cliVersion
	return depData, nil
}

func (x *SignedDepositDataMsgBody) SameDepositData(other *SignedDepositDataMsgBody) bool {
	if bytes.Compare(x.RequestID, other.RequestID) != 0 {
		return false
	}
	if x.Threshold != other.Threshold {
		return false
	}
	if len(x.Committee) != len(other.Committee) {
		return false
	}
	for i, member := range x.Committee {
		if other.Committee[i] != member {
			return false
		}
	}
	if bytes.Compare(x.ValidatorPublicKey, other.ValidatorPublicKey) != 0 {
		return false
	}
	if bytes.Compare(x.WithdrawalCredentials, other.WithdrawalCredentials) != 0 {
		return false
	}
	if bytes.Compare(x.DepositDataSignature, other.DepositDataSignature) != 0 {
		return false
	}
	return true
}

func (x *SignedDepositDataMsgBody) GetRoot() ([]byte, error) {
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
		x.EncryptedShare,
		//o.SharePubKeys, // TODO: Add later
		x.ValidatorPublicKey,
		x.DepositDataSignature,
	)

	return crypto.Keccak256(bytes), nil
}

// Encode returns a msg encoded bytes or error
func (x *SignedDepositDataMsgBody) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

// Decode returns error if decoding failed
func (x *SignedDepositDataMsgBody) Decode(data []byte) error {
	return proto.Unmarshal(data, x)
}

func (x *ParsedSignedDepositDataMessage) FromBase(base *Message) error {
	raw, err := proto.Marshal(base)
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, x)
}

func (x *ParsedSignedDepositDataMessage) ToBase() (*Message, error) {
	raw, err := proto.Marshal(x)
	if err != nil {
		return nil, err
	}
	base := &Message{}
	err = proto.Unmarshal(raw, base)
	if err != nil {
		return nil, err
	}
	return base, nil
}
