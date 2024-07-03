package ssv

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type SigningRoot string

type PartialSigContainer struct {
	// Signature map: validator index -> signing root -> operator id (signer) -> signature (from the signer for the validator's signing root)
	Signatures map[phase0.ValidatorIndex]map[SigningRoot]map[types.OperatorID]types.Signature
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64
}

func NewPartialSigContainer(quorum uint64) *PartialSigContainer {
	return &PartialSigContainer{
		Quorum:     quorum,
		Signatures: make(map[phase0.ValidatorIndex]map[SigningRoot]map[types.OperatorID]types.Signature),
	}
}

func (ps *PartialSigContainer) AddSignature(sigMsg *types.PartialSignatureMessage) {
	if ps.Signatures[sigMsg.ValidatorIndex] == nil {
		ps.Signatures[sigMsg.ValidatorIndex] = make(map[SigningRoot]map[types.OperatorID]types.Signature)
	}
	if ps.Signatures[sigMsg.ValidatorIndex][signingRootHex(sigMsg.SigningRoot)] == nil {
		ps.Signatures[sigMsg.ValidatorIndex][signingRootHex(sigMsg.SigningRoot)] = make(map[types.OperatorID]types.Signature)
	}
	m := ps.Signatures[sigMsg.ValidatorIndex][signingRootHex(sigMsg.SigningRoot)]

	if m[sigMsg.Signer] == nil {
		m[sigMsg.Signer] = make([]byte, 96)
		copy(m[sigMsg.Signer], sigMsg.PartialSignature)
	}
}

// HasSigner returns true if container has signature for signer and signing root, else it returns false
func (ps *PartialSigContainer) HasSignature(validatorIndex phase0.ValidatorIndex, signer types.OperatorID, signingRoot [32]byte) bool {
	_, err := ps.GetSignature(validatorIndex, signer, signingRoot)
	return err == nil
}

// GetSignature returns the signature for a given root and signer
func (ps *PartialSigContainer) GetSignature(validatorIndex phase0.ValidatorIndex, signer types.OperatorID, signingRoot [32]byte) (types.Signature, error) {
	if ps.Signatures[validatorIndex] == nil {
		return nil, errors.New("Dont have signature for the given validator index")
	}
	if ps.Signatures[validatorIndex][signingRootHex(signingRoot)] == nil {
		return nil, errors.New("Dont have signature for the given signing root")
	}
	if ps.Signatures[validatorIndex][signingRootHex(signingRoot)][signer] == nil {
		return nil, errors.New("Dont have signature on signing root for the given signer")
	}
	return ps.Signatures[validatorIndex][signingRootHex(signingRoot)][signer], nil
}

// Return signature map for given root
func (ps *PartialSigContainer) GetSignatures(validatorIndex phase0.ValidatorIndex, signingRoot [32]byte) map[types.OperatorID]types.Signature {
	if ps.Signatures[validatorIndex] == nil {
		return nil
	}
	return ps.Signatures[validatorIndex][signingRootHex(signingRoot)]
}

// Remove signer from signature map
func (ps *PartialSigContainer) Remove(validatorIndex phase0.ValidatorIndex, signer uint64, signingRoot [32]byte) {
	if ps.Signatures[validatorIndex] == nil {
		return
	}
	if ps.Signatures[validatorIndex][signingRootHex(signingRoot)] == nil {
		return
	}
	if ps.Signatures[validatorIndex][signingRootHex(signingRoot)][signer] == nil {
		return
	}
	delete(ps.Signatures[validatorIndex][signingRootHex(signingRoot)], signer)
}

func (ps *PartialSigContainer) ReconstructSignature(root [32]byte, validatorPubKey []byte, validatorIndex phase0.ValidatorIndex) ([]byte, error) {
	// Reconstruct signatures
	if ps.Signatures[validatorIndex] == nil {
		return nil, errors.New("no signatures for the given validator index")
	}
	if ps.Signatures[validatorIndex][signingRootHex(root)] == nil {
		return nil, errors.New("no signatures for the given signing root")
	}

	operatorsSignatures := make(map[uint64][]byte)
	for operatorID, sig := range ps.Signatures[validatorIndex][signingRootHex(root)] {
		operatorsSignatures[operatorID] = sig
	}
	signature, err := types.ReconstructSignatures(operatorsSignatures)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}

	// Get validator pub key copy (This avoids cgo Go pointer to Go pointer issue)
	validatorPubKeyCopy := make([]byte, len(validatorPubKey))
	copy(validatorPubKeyCopy, validatorPubKey)

	if err := types.VerifyReconstructedSignature(signature, validatorPubKeyCopy, root); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}
	return signature.Serialize(), nil
}

func (ps *PartialSigContainer) HasQuorum(validatorIndex phase0.ValidatorIndex, root [32]byte) bool {
	return uint64(len(ps.Signatures[validatorIndex][signingRootHex(root)])) >= ps.Quorum
}

func signingRootHex(r [32]byte) SigningRoot {
	return SigningRoot(hex.EncodeToString(r[:]))
}
