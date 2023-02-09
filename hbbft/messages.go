package hbbft

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

// HasQuorum returns true if a unique set of signers has quorum
func HasQuorum(share *types.Share, msgs []*SignedMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasQuorum(len(uniqueSigners))
}

// HasPartialQuorum returns true if a unique set of signers has partial quorum
func HasPartialQuorum(share *types.Share, msgs []*SignedMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasPartialQuorum(len(uniqueSigners))
}

type MessageType int

const (
	TransactionMsgType MessageType = iota
	ABAInitMsgType
	ABAAuxMsgType
	ABAConfMsgType
	ABAFinishMsgType
	RBCValMsgType
	RBCEchoMsgType
	RBCReadyMsgType
	DECMsgType
)

// =========================
//			Transaction
// =========================

type TransactionData struct {
	Data []byte
}

// Encode returns a msg encoded bytes or error
func (d *TransactionData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *TransactionData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// =========================
//			ABAInit
// =========================

type ABAInitData struct {
	Vote     byte
	Round    Round
	ACSRound ACSRound
}

// Encode returns a msg encoded bytes or error
func (d *ABAInitData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *ABAInitData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *ABAInitData) Validate() error {
	if d.Vote != 0 && d.Vote != 1 {
		return errors.New("ABAInitData: vote not 0 or 1")
	}
	return nil
}

// =========================
//			ABAAux
// =========================

type ABAAuxData struct {
	Vote     byte
	Round    Round
	ACSRound ACSRound
}

// Encode returns a msg encoded bytes or error
func (d *ABAAuxData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *ABAAuxData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *ABAAuxData) Validate() error {
	if d.Vote != 0 && d.Vote != 1 {
		return errors.New("ABAAuxData: vote not 0 or 1")
	}
	return nil
}

// =========================
//			ABAConf
// =========================

type ABAConfData struct {
	Votes    []byte
	Round    Round
	ACSRound ACSRound
}

// Encode returns a msg encoded bytes or error
func (d *ABAConfData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *ABAConfData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *ABAConfData) Validate() error {
	if len(d.Votes) == 0 {
		return errors.New("ABAConfData: empty votes")
	}
	for _, vote := range d.Votes {
		if vote != 0 && vote != 1 {
			return errors.New("ABAConfData: vote not 0 or 1")
		}
	}
	return nil
}

// =========================
//			ABAFinish
// =========================

type ABAFinishData struct {
	Vote     byte
	ACSRound ACSRound
}

// Encode returns a msg encoded bytes or error
func (d *ABAFinishData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *ABAFinishData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *ABAFinishData) Validate() error {
	if d.Vote != 0 && d.Vote != 1 {
		return errors.New("ABAFinishData: vote not 0 or 1")
	}
	return nil
}

// =========================
//			RBCVal
// =========================

type RBCValData struct {
	MarkleTree   []byte
	Branch       []byte
	ErasureShare []byte
}

// Encode returns a msg encoded bytes or error
func (d *RBCValData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *RBCValData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, its variables for validity.
func (d *RBCValData) Validate() error {
	if len(d.MarkleTree) == 0 {
		return errors.New("RBCValData: MarkleTree with length 0")
	}
	if len(d.Branch) == 0 {
		return errors.New("RBCValData: Branch with length 0")
	}
	if len(d.ErasureShare) == 0 {
		return errors.New("RBCValData: ErasureShare with length 0")
	}
	return nil
}

// =========================
//			RBCEcho
// =========================

type RBCEchoData struct {
	MarkleTree   []byte
	Branch       []byte
	ErasureShare []byte
}

// Encode returns a msg encoded bytes or error
func (d *RBCEchoData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *RBCEchoData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, its variables for validity.
func (d *RBCEchoData) Validate() error {
	if len(d.MarkleTree) == 0 {
		return errors.New("RBCValData: MarkleTree with length 0")
	}
	if len(d.Branch) == 0 {
		return errors.New("RBCValData: Branch with length 0")
	}
	if len(d.ErasureShare) == 0 {
		return errors.New("RBCValData: ErasureShare with length 0")
	}
	return nil
}

// =========================
//			RBCReady
// =========================

type RBCReadyData struct {
	MarkleTree []byte
}

// Encode returns a msg encoded bytes or error
func (d *RBCReadyData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *RBCReadyData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, its variables for validity.
func (d *RBCReadyData) Validate() error {
	if len(d.MarkleTree) == 0 {
		return errors.New("RBCValData: MarkleTree with length 0")
	}
	return nil
}

// =========================
//			DEC
// =========================

type DECData struct {
	ACSRound       ACSRound
	OperatorID     types.OperatorID
	Sender         types.OperatorID
	DecryptedShare []byte
}

// Encode returns a msg encoded bytes or error
func (d *DECData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *DECData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, its variables for validity.
func (d *DECData) Validate() error {
	if len(d.DecryptedShare) == 0 {
		return errors.New("RBCValData: DecryptedShare with length 0")
	}
	return nil
}

// =========================
//			Message
// =========================

type Message struct {
	MsgType    MessageType
	Height     Height // QBFT instance Height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte // instance Identifier this msg belongs to
	Data       []byte
}

// GetTransactionData returns transaction specific data
func (msg *Message) GetTransactionData() (*TransactionData, error) {
	ret := &TransactionData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode TransactionData from message")
	}
	return ret, nil
}

// GetABAInitData returns abainit specific data
func (msg *Message) GetABAInitData() (*ABAInitData, error) {
	ret := &ABAInitData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode ABAInit data from message")
	}
	return ret, nil
}

// GetABAAuxData returns abainit specific data
func (msg *Message) GetABAAuxData() (*ABAAuxData, error) {
	ret := &ABAAuxData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode ABAAuxData from message")
	}
	return ret, nil
}

// GetABAConfData returns abainit specific data
func (msg *Message) GetABAConfData() (*ABAConfData, error) {
	ret := &ABAConfData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode ABAConfData from message")
	}
	return ret, nil
}

// GetABAFinishData returns abainit specific data
func (msg *Message) GetABAFinishData() (*ABAFinishData, error) {
	ret := &ABAFinishData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode ABAFinishData from message")
	}
	return ret, nil
}

// GetData returns RBCVal specific data
func (msg *Message) GetRBCValData() (*RBCValData, error) {
	ret := &RBCValData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode RBCValData from message")
	}
	return ret, nil
}

// GetData returns RBCEcho specific data
func (msg *Message) GetRBCEchoData() (*RBCEchoData, error) {
	ret := &RBCEchoData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode RBCEchoData from message")
	}
	return ret, nil
}

// GetData returns RBCReady specific data
func (msg *Message) GetRBCReadyData() (*RBCReadyData, error) {
	ret := &RBCReadyData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode RBCReadyData from message")
	}
	return ret, nil
}

// GetData returns DEC specific data
func (msg *Message) GetDECData() (*DECData, error) {
	ret := &DECData{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode DECData from message")
	}
	return ret, nil
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, &msg)
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() ([]byte, error) {
	marshaledRoot, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode message")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (msg *Message) Validate() error {
	if len(msg.Identifier) == 0 {
		return errors.New("message identifier is invalid")
	}
	if len(msg.Data) == 0 {
		return errors.New("message data is invalid")
	}
	if msg.MsgType > 5 {
		return errors.New("message type is invalid")
	}
	return nil
}

type SignedMessage struct {
	Signature types.Signature
	Signers   []types.OperatorID
	Message   *Message // message for which this signature is for
}

func (signedMsg *SignedMessage) GetSignature() types.Signature {
	return signedMsg.Signature
}
func (signedMsg *SignedMessage) GetSigners() []types.OperatorID {
	return signedMsg.Signers
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (signedMsg *SignedMessage) MatchedSigners(ids []types.OperatorID) bool {
	if len(signedMsg.Signers) != len(ids) {
		return false
	}

	for _, id := range signedMsg.Signers {
		found := false
		for _, id2 := range ids {
			if id == id2 {
				found = true
			}
		}

		if !found {
			return false
		}
	}
	return true
}

// CommonSigners returns true if there is at least 1 common signer
func (signedMsg *SignedMessage) CommonSigners(ids []types.OperatorID) bool {
	for _, id := range signedMsg.Signers {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (signedMsg *SignedMessage) Aggregate(sig types.MessageSignature) error {
	if signedMsg.CommonSigners(sig.GetSigners()) {
		return errors.New("duplicate signers")
	}

	r1, err := signedMsg.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	r2, err := sig.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	if !bytes.Equal(r1, r2) {
		return errors.New("can't aggregate, roots not equal")
	}

	aggregated, err := signedMsg.Signature.Aggregate(sig.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not aggregate signatures")
	}
	signedMsg.Signature = aggregated
	signedMsg.Signers = append(signedMsg.Signers, sig.GetSigners()...)

	return nil
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &signedMsg)
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
	return signedMsg.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedMessage) DeepCopy() *SignedMessage {
	ret := &SignedMessage{
		Signers:   make([]types.OperatorID, len(signedMsg.Signers)),
		Signature: make([]byte, len(signedMsg.Signature)),
	}
	copy(ret.Signers, signedMsg.Signers)
	copy(ret.Signature, signedMsg.Signature)

	ret.Message = &Message{
		MsgType:    signedMsg.Message.MsgType,
		Height:     signedMsg.Message.Height,
		Round:      signedMsg.Message.Round,
		Identifier: make([]byte, len(signedMsg.Message.Identifier)),
		Data:       make([]byte, len(signedMsg.Message.Data)),
	}
	copy(ret.Message.Identifier, signedMsg.Message.Identifier)
	copy(ret.Message.Data, signedMsg.Message.Data)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (signedMsg *SignedMessage) Validate() error {
	if len(signedMsg.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(signedMsg.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range signedMsg.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		if signer == 0 {
			return errors.New("signer ID 0 not allowed")
		}
		signed[signer] = true
	}

	return signedMsg.Message.Validate()
}
