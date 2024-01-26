package ssv

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
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
	ps.AddSignatureForSignatureRootAndSigner(sigMsg.PartialSignature, sigMsg.SigningRoot, sigMsg.Signer)
}

func (ps *PartialSigContainer) AddSignatureForSignatureRootAndSigner(signature types.Signature, signingRoot [32]byte, signer types.OperatorID) {
	if ps.Signatures[rootHex(signingRoot)] == nil {
		ps.Signatures[rootHex(signingRoot)] = make(map[types.OperatorID][]byte)
	}
	m := ps.Signatures[rootHex(signingRoot)]

	if m[signer] == nil {
		m[signer] = make([]byte, 96)
		copy(m[signer], signature)
	}
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
