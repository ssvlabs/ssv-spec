package types

import (
	"crypto/sha256"
	"fmt"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (s Signature) Verify(data Root, domain DomainType, sigType SignatureType, pkByts []byte) error {
	computedRoot, err := ComputeSigningRoot(data, ComputeSignatureDomain(domain, sigType))
	if err != nil {
		return errors.Wrap(err, "could not compute signing root")
	}

	sign := &bls.Sign{}
	if err := sign.Deserialize(s); err != nil {
		return errors.Wrap(err, "failed to deserialize signature")
	}

	pk := &bls.PublicKey{}
	if err := pk.Deserialize(pkByts); err != nil {
		return errors.Wrap(err, "failed to deserialize public key")
	}

	if res := sign.VerifyByte(pk, computedRoot[:]); !res {
		return errors.New("failed to verify signature")
	}
	return nil
}

// ComputeSigningRoot returns a singable/ verifiable root calculated from the a provided data and signature domain
func ComputeSigningRoot(data Root, domain SignatureDomain) ([32]byte, error) {
	dataRoot, err := data.GetRoot()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not get root from Root")
	}

	ret := sha256.Sum256(append(dataRoot[:], domain...))
	return ret, nil
}

// ComputeSignatureDomain returns a signature domain based on the domain type and signature type
func ComputeSignatureDomain(domain DomainType, sigType SignatureType) SignatureDomain {
	return SignatureDomain(append(domain[:], sigType[:]...))
}

// ReconstructSignatures receives a map of user indexes and serialized bls.Sign.
// It then reconstructs the original threshold signature using lagrange interpolation
func ReconstructSignatures(signatures map[OperatorID][]byte) (*bls.Sign, error) {
	reconstructedSig := bls.Sign{}

	idVec := make([]bls.ID, 0)
	sigVec := make([]bls.Sign, 0)

	for index, signature := range signatures {
		blsID := bls.ID{}
		err := blsID.SetDecString(fmt.Sprintf("%d", index))
		if err != nil {
			return nil, err
		}

		idVec = append(idVec, blsID)
		blsSig := bls.Sign{}

		err = blsSig.Deserialize(signature)
		if err != nil {
			return nil, err
		}

		sigVec = append(sigVec, blsSig)
	}
	err := reconstructedSig.Recover(sigVec, idVec)
	return &reconstructedSig, err
}

func VerifyReconstructedSignature(sig *bls.Sign, validatorPubKey []byte, root [32]byte) error {
	pk := &bls.PublicKey{}
	if err := pk.Deserialize(validatorPubKey); err != nil {
		return errors.Wrap(err, "could not deserialize validator pk")
	}

	// verify reconstructed sig
	if res := sig.VerifyByte(pk, root[:]); !res {
		return errors.New("could not reconstruct a valid signature")
	}
	return nil
}
