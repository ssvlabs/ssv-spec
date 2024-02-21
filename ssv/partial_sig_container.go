package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"sort"
)

type PartialSignatureContainer map[types.OperatorID]*types.SignedPartialSignatureMessage

func NewPartialSignatureContainer() PartialSignatureContainer {
	return make(PartialSignatureContainer)
}

// ReconstructSignature return reconstructed signature for a root
func (ps PartialSignatureContainer) ReconstructSignature(root [32]byte, validatorPubKey []byte) ([]byte, error) {
	// collect signatures
	sigs := ps.SignaturesForRoot(root)

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

// SignaturesForRoot returns a map of signer and signature for a specific root
func (ps PartialSignatureContainer) SignaturesForRoot(root [32]byte) map[types.OperatorID][]byte {
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

// Roots returns roots for the partial sigs
func (ps PartialSignatureContainer) Roots() [][32]byte {
	ret := make([][32]byte, 0)

	if len(ps) <= 0 {
		return ret
	}

	for _, sigMsg := range ps {
		for _, msg := range sigMsg.Message.Messages {
			ret = append(ret, msg.SigningRoot)
		}
		break // only need to iterate first msg
	}
	return ret
}

// AllSorted returns ordered by signer array of signed messages
func (ps PartialSignatureContainer) AllSorted() []*types.SignedPartialSignatureMessage {
	ret := make([]*types.SignedPartialSignatureMessage, len(ps))
	for _, m := range ps {
		ret = append(ret, m)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Signer < ret[j].Signer
	})
	return ret
}
