package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type PartialSignatures []*PartialSignature

// Encode returns a msg encoded bytes or error
func (pss *PartialSignatures) Encode() ([]byte, error) {
	return json.Marshal(pss)
}

// Decode returns error if decoding failed
func (pss *PartialSignatures) Decode(data []byte) error {
	return json.Unmarshal(data, pss)
}

// GetRoot returns the root used for signing and verification
func (pss PartialSignatures) GetRoot() ([]byte, error) {
	marshaledRoot, err := pss.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PartialSignatures")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

func (pss PartialSignatures) Validate() error {
	if len(pss) == 0 {
		return errors.New("no PartialSignatures messages")
	}
	for _, m := range pss {
		if err := m.Validate(); err != nil {
			return errors.Wrap(err, "message invalid")
		}
	}
	return nil
}

type PartialSignatureMetaData struct {
	ContributionSubCommitteeIndex uint64
}

// PartialSignature is a msg for partial Beacon chain related signatures (like partial attestation, block, randao sigs)
type PartialSignature struct {
	Slot        phase0.Slot     // Slot represents the slot for which the partial BN signature is for
	Signature   types.Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot []byte          `ssz-size:"32"` // the root signed in PartialSignature
}

// Encode returns a msg encoded bytes or error
func (p *PartialSignature) Encode() ([]byte, error) {
	return json.Marshal(p)
}

// Decode returns error if decoding failed
func (p *PartialSignature) Decode(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *PartialSignature) GetRoot() ([]byte, error) {
	marshaledRoot, err := p.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PartialSignature")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

func (p *PartialSignature) Validate() error {
	if len(p.Signature) != 96 {
		return errors.New("PartialSignature sig invalid")
	}
	if len(p.SigningRoot) != 32 {
		return errors.New("SigningRoot invalid")
	}
	return nil
}

// SignedPartialSignatures is an operator's signature over PartialSignature
type SignedPartialSignatures struct {
	PartialSignatures PartialSignatures `ssz-max:"13"`
	Signature         types.Signature   `ssz-size:"96"`
	Signer            types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (s *SignedPartialSignatures) Encode() ([]byte, error) {
	return s.MarshalSSZ()
}

// Decode returns error if decoding failed
func (s *SignedPartialSignatures) Decode(data []byte) error {
	return s.UnmarshalSSZ(data)
}

func (s *SignedPartialSignatures) GetSignature() types.Signature {
	return s.Signature
}

func (s *SignedPartialSignatures) GetSigners() []types.OperatorID {
	return []types.OperatorID{s.Signer}
}

func (s *SignedPartialSignatures) GetRoot() ([]byte, error) {
	return s.PartialSignatures.GetRoot()
}

func (s *SignedPartialSignatures) Aggregate(signedMsg types.MessageSignature) error {
	//if !bytes.Equal(s.GetRoot(), signedMsg.GetRoot()) {
	//	return errors.New("can't aggregate msgs with different roots")
	//}
	//
	//// verify no matching Signer
	//for _, signerID := range s.Signer {
	//	for _, toMatchID := range signedMsg.GetSigners() {
	//		if signerID == toMatchID {
	//			return errors.New("Signer IDs partially/ fully match")
	//		}
	//	}
	//}
	//
	//allSigners := append(s.Signer, signedMsg.GetSigners()...)
	//
	//// verify and aggregate
	//sig1, err := blsSig(s.Signature)
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
	//s.Signature = sig1.Serialize()
	//s.Signer = allSigners
	//return nil
	panic("implement")
}

// MatchedSigners returns true if the provided Signer ids are equal to GetSignerIds() without order significance
func (s *SignedPartialSignatures) MatchedSigners(ids []types.OperatorID) bool {
	toMatchCnt := make(map[types.OperatorID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[types.OperatorID]int)
	for _, id := range s.GetSigners() {
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

func (s *SignedPartialSignatures) Validate() error {
	if len(s.Signature) != 96 {
		return errors.New("SignedPartialSignatures sig invalid")
	}
	if s.Signer == 0 {
		return errors.New("signer ID 0 not allowed")
	}
	return s.PartialSignatures.Validate()
}
