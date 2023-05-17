package ssv

import (
	"bytes"
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type PartialSignatureContainer map[types.OperatorID]*types.SignedPartialSignatureMessage

func (ps PartialSignatureContainer) ReconstructSignature(root [32]byte, validatorPubKey []byte) ([]byte, error) {
	// collect signatures
	sigs := ps.SignatureForRoot(root)

	// reconstruct
	signature, err := types.ReconstructSignatures(sigs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}
	if err := types.VerifyReconstructedSignature(signature, validatorPubKey, root); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}
	return signature.Serialize(), nil
}

func (ps PartialSignatureContainer) SignatureForRoot(root [32]byte) map[types.OperatorID][]byte {
	sigs := make(map[types.OperatorID][]byte, 0)
	for operatorID, sigMsg := range ps {
		for _, msg := range sigMsg.Message.Messages {
			if bytes.Equal(msg.SigningRoot[:], root[:]) {
				sigs[operatorID] = msg.PartialSignature
			}
		}
	}
	return sigs
}

func (ps PartialSignatureContainer) Roots() [][32]byte {
	if len(ps) > 0 {
		ret := make([][32]byte, 0)
		for _, sigMsg := range ps {
			for _, msg := range sigMsg.Message.Messages {
				ret = append(ret, msg.SigningRoot)
			}
			break // only need to iterate first msg
		}
		return ret
	}

	return [][32]byte{}
}

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
