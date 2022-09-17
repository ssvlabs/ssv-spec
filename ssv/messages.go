package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type PartialSignatures struct {
	Messages []*PartialSignature `ssz-max:"13"`
}

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
	if len(pss.Messages) == 0 {
		return errors.New("no PartialSignatures messages")
	}
	for _, m := range pss.Messages {
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
	Slot             phase0.Slot     // Slot represents the slot for which the partial BN signature is for
	PartialSignature types.Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot      []byte          `ssz-size:"32"` // the root signed in PartialSignature
	Signer           types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (ps *PartialSignature) Encode() ([]byte, error) {
	return json.Marshal(ps)
}

// Decode returns error if decoding failed
func (ps *PartialSignature) Decode(data []byte) error {
	return json.Unmarshal(data, ps)
}

func (ps *PartialSignature) GetRoot() ([]byte, error) {
	marshaledRoot, err := ps.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PartialSignature")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

func (ps *PartialSignature) Validate() error {
	if len(ps.PartialSignature) != 96 {
		return errors.New("PartialSignature sig invalid")
	}
	if len(ps.SigningRoot) != 32 {
		return errors.New("SigningRoot invalid")
	}
	return nil
}

// SignedPartialSignature is an operator's signature over PartialSignature
type SignedPartialSignature struct {
	Message   PartialSignatures
	Signature types.Signature `ssz-size:"96"`
	Signer    types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (sps *SignedPartialSignature) Encode() ([]byte, error) {
	return json.Marshal(sps)
}

// Decode returns error if decoding failed
func (sps *SignedPartialSignature) Decode(data []byte) error {
	return json.Unmarshal(data, &sps)
}

func (sps *SignedPartialSignature) GetSignature() types.Signature {
	return sps.Signature
}

func (sps *SignedPartialSignature) GetSigners() []types.OperatorID {
	return []types.OperatorID{sps.Signer}
}

func (sps *SignedPartialSignature) GetRoot() ([]byte, error) {
	return sps.Message.GetRoot()
}

func (sps *SignedPartialSignature) Aggregate(signedMsg types.MessageSignature) error {
	//if !bytes.Equal(sps.GetRoot(), signedMsg.GetRoot()) {
	//	return errors.New("can't aggregate msgs with different roots")
	//}
	//
	//// verify no matching Signer
	//for _, signerID := range sps.Signer {
	//	for _, toMatchID := range signedMsg.GetSigners() {
	//		if signerID == toMatchID {
	//			return errors.New("Signer IDs partially/ fully match")
	//		}
	//	}
	//}
	//
	//allSigners := append(sps.Signer, signedMsg.GetSigners()...)
	//
	//// verify and aggregate
	//sig1, err := blsSig(sps.Signature)
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
	//sps.Signature = sig1.Serialize()
	//sps.Signer = allSigners
	//return nil
	panic("implement")
}

// MatchedSigners returns true if the provided Signer ids are equal to GetSignerIds() without order significance
func (sps *SignedPartialSignature) MatchedSigners(ids []types.OperatorID) bool {
	toMatchCnt := make(map[types.OperatorID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[types.OperatorID]int)
	for _, id := range sps.GetSigners() {
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

func (sps *SignedPartialSignature) Validate() error {
	if len(sps.Signature) != 96 {
		return errors.New("SignedPartialSignature sig invalid")
	}
	return sps.Message.Validate()
}
