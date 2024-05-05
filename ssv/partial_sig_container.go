package ssv

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type PartialSigContainer struct {
	Signatures map[string]map[types.OperatorID][]byte
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64
}

func NewPartialSigContainer(quorum uint64) *PartialSigContainer {
	return &PartialSigContainer{
		Quorum:     quorum,
		Signatures: make(map[string]map[types.OperatorID][]byte),
	}
}

func (ps *PartialSigContainer) AddSignature(sigMsg *types.PartialSignatureMessage) {
	if ps.Signatures[rootHex(sigMsg.SigningRoot)] == nil {
		ps.Signatures[rootHex(sigMsg.SigningRoot)] = make(map[types.OperatorID][]byte)
	}
	m := ps.Signatures[rootHex(sigMsg.SigningRoot)]

	if m[sigMsg.Signer] == nil {
		m[sigMsg.Signer] = make([]byte, 96)
		copy(m[sigMsg.Signer], sigMsg.PartialSignature)
	}
}

// Returns if container has signature for signer and signing root
func (ps *PartialSigContainer) HasSigner(signer types.OperatorID, signingRoot [32]byte) bool {
	if ps.Signatures[rootHex(signingRoot)] == nil {
		return false
	}
	return ps.Signatures[rootHex(signingRoot)][signer] != nil
}

// Return signature for given root and signer
func (ps *PartialSigContainer) GetSignature(signer types.OperatorID, signingRoot [32]byte) (types.Signature, error) {
	if ps.Signatures[rootHex(signingRoot)] == nil {
		return nil, errors.New("Dont have signature for the given signing root")
	}
	if ps.Signatures[rootHex(signingRoot)][signer] == nil {
		return nil, errors.New("Dont have signature on signing root for the given signer")
	}
	return ps.Signatures[rootHex(signingRoot)][signer], nil
}

// Return signature map for given root
func (ps *PartialSigContainer) GetSignatures(signingRoot [32]byte) map[types.OperatorID][]byte {
	return ps.Signatures[rootHex(signingRoot)]
}

// Remove signer from signature map
func (ps *PartialSigContainer) Remove(signer uint64, signingRoot [32]byte) {
	if ps.Signatures[rootHex(signingRoot)] == nil {
		return
	}
	if ps.Signatures[rootHex(signingRoot)][signer] == nil {
		return
	}
	delete(ps.Signatures[rootHex(signingRoot)], signer)
}

func (ps *PartialSigContainer) ReconstructSignature(root [32]byte, validatorPubKey []byte) ([]byte, error) {
	// Reconstruct signatures
	signature, err := types.ReconstructSignatures(ps.Signatures[rootHex(root)])
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}
	if err := types.VerifyReconstructedSignature(signature, validatorPubKey, root); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}
	return signature.Serialize(), nil
}

func (ps *PartialSigContainer) HasQuorum(root [32]byte) bool {
	return uint64(len(ps.Signatures[rootHex(root)])) >= ps.Quorum
}

func rootHex(r [32]byte) string {
	return hex.EncodeToString(r[:])
}
