package types

import (
	"github.com/pkg/errors"
)

type PartialSigMsgType uint64

const (
	// PostConsensusPartialSig is a partial signature over a decided duty (attestation data, block, etc)
	PostConsensusPartialSig PartialSigMsgType = iota
	// RandaoPartialSig is a partial signature over randao reveal
	RandaoPartialSig
	// SelectionProofPartialSig is a partial signature for aggregator selection proof
	SelectionProofPartialSig
	// ContributionProofs is the partial selection proofs for sync committee contributions (it's an array of sigs)
	ContributionProofs
	// ValidatorRegistrationPartialSig is a partial signature over a ValidatorRegistration object
	ValidatorRegistrationPartialSig
)

type PartialSignatureMessages struct {
	Type     PartialSigMsgType
	Messages []*PartialSignatureMessage `ssz-max:"13"`
}

// Encode returns a msg encoded bytes or error
func (msgs *PartialSignatureMessages) Encode() ([]byte, error) {
	return msgs.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msgs *PartialSignatureMessages) Decode(data []byte) error {
	return msgs.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msgs PartialSignatureMessages) GetRoot() ([32]byte, error) {
	return msgs.HashTreeRoot()
}

func (msgs PartialSignatureMessages) Validate() error {
	if len(msgs.Messages) == 0 {
		return errors.New("no PartialSignatureMessages messages")
	}
	for _, m := range msgs.Messages {
		if err := m.Validate(); err != nil {
			return errors.Wrap(err, "message invalid")
		}
	}
	return nil
}

// PartialSignatureMessage is a msg for partial Beacon chain related signatures (like partial attestation, block, randao sigs)
type PartialSignatureMessage struct {
	PartialSignature Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot      []byte    `ssz-size:"32"` // the root signed in PartialSignature
	Signer           OperatorID
}

// Encode returns a msg encoded bytes or error
func (pcsm *PartialSignatureMessage) Encode() ([]byte, error) {
	return pcsm.MarshalSSZ()
}

// Decode returns error if decoding failed
func (pcsm *PartialSignatureMessage) Decode(data []byte) error {
	return pcsm.UnmarshalSSZ(data)
}

func (pcsm *PartialSignatureMessage) GetRoot() ([32]byte, error) {
	return pcsm.HashTreeRoot()
}

func (pcsm *PartialSignatureMessage) Validate() error {
	if pcsm.Signer == 0 {
		return errors.New("signer ID 0 not allowed")
	}
	return nil
}

// SignedPartialSignatureMessage is an operator's signature over PartialSignatureMessage
type SignedPartialSignatureMessage struct {
	Message   PartialSignatureMessages
	Signature Signature `ssz-size:"96"`
	Signer    OperatorID
}

// Encode returns a msg encoded bytes or error
func (spcsm *SignedPartialSignatureMessage) Encode() ([]byte, error) {
	return spcsm.MarshalSSZ()
}

// Decode returns error if decoding failed
func (spcsm *SignedPartialSignatureMessage) Decode(data []byte) error {
	return spcsm.UnmarshalSSZ(data)
}

func (spcsm *SignedPartialSignatureMessage) GetSignature() Signature {
	return spcsm.Signature
}

func (spcsm *SignedPartialSignatureMessage) GetSigners() []OperatorID {
	return []OperatorID{spcsm.Signer}
}

func (spcsm *SignedPartialSignatureMessage) GetRoot() ([32]byte, error) {
	return spcsm.Message.GetRoot()
}

func (spcsm *SignedPartialSignatureMessage) Aggregate(signedMsg MessageSignature) error {
	//if !bytes.Equal(spcsm.GetRoot(), signedMsg.GetRoot()) {
	//	return errors.New("can't aggregate msgs with different roots")
	//}
	//
	//// verify no matching Signer
	//for _, signerID := range spcsm.Signer {
	//	for _, toMatchID := range signedMsg.GetSigners() {
	//		if signerID == toMatchID {
	//			return errors.New("Signer IDs partially/ fully match")
	//		}
	//	}
	//}
	//
	//allSigners := append(spcsm.Signer, signedMsg.GetSigners()...)
	//
	//// verify and aggregate
	//sig1, err := blsSig(spcsm.Signature)
	//if err != nil {
	//	return errors.Wrap(err, "could not parse PartialSignature")
	//}
	//
	//sig2, err := blsSig(signedMsg.GetSignature())
	//if err != nil {
	//	return errors.Wrap(err, "could not parse PartialSignature")
	//}
	//
	//sig1.Add(sig2)
	//spcsm.Signature = sig1.Serialize()
	//spcsm.Signer = allSigners
	//return nil
	panic("implement")
}

// MatchedSigners returns true if the provided Signer ids are equal to GetSignerIds() without order significance
func (spcsm *SignedPartialSignatureMessage) MatchedSigners(ids []OperatorID) bool {
	toMatchCnt := make(map[OperatorID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[OperatorID]int)
	for _, id := range spcsm.GetSigners() {
		foundCnt[id]++
	}

	for id, cnt := range toMatchCnt {
		if cnt != foundCnt[id] {
			return false
		}
	}
	return true
}

//
//func blsSig(sig []byte) (*bls.Sign, error) {
//	ret := &bls.Sign{}
//	if err := ret.Deserialize(sig); err != nil {
//		return nil, errors.Wrap(err, "could not covert PartialSignature byts to bls.sign")
//	}
//	return ret, nil
//}

func (spcsm *SignedPartialSignatureMessage) Validate() error {
	if len(spcsm.Signature) != 96 {
		return errors.New("SignedPartialSignatureMessage sig invalid")
	}
	if spcsm.Signer == 0 {
		return errors.New("signer ID 0 not allowed")
	}

	for _, msg := range spcsm.Message.Messages {
		if spcsm.Signer != msg.Signer {
			return errors.New("inconsistent signers")
		}
	}

	return spcsm.Message.Validate()
}
