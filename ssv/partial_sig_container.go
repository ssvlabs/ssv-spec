package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type PartialSignatureContainer map[types.OperatorID]*types.SignedPartialSignatureMessage

func NewPartialSignatureContainer() PartialSignatureContainer {
	return make(PartialSignatureContainer)
}

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

func (ps PartialSignatureContainer) All() []*types.SignedPartialSignatureMessage {
	ret := make([]*types.SignedPartialSignatureMessage, 0)
	for _, m := range ps {
		ret = append(ret, m)
	}
	return ret
}
