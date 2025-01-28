package types

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/go-bitfield"
)

// VersionedAttestationResponse contains a versioned repsonse to the attestation duty.
type VersionedAttestationResponse struct {
	Version   spec.DataVersion
	Phase0    *phase0.Attestation
	Altair    *phase0.Attestation
	Bellatrix *phase0.Attestation
	Capella   *phase0.Attestation
	Deneb     *phase0.Attestation
	Electra   *electra.SingleAttestation
}

// Signature returns the signature of the attestation response.
func (v *VersionedAttestationResponse) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 attestation response")
		}

		return v.Phase0.Signature, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no altair attestation response")
		}

		return v.Altair.Signature, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix attestation response")
		}

		return v.Bellatrix.Signature, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella attestation response")
		}

		return v.Capella.Signature, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no deneb attestation response")
		}

		return v.Deneb.Signature, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no electra attestation response")
		}

		return v.Electra.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// Sets the signature field of the attestation response.
func (v *VersionedAttestationResponse) WithSignature(signature phase0.BLSSignature) error {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return errors.New("no phase0 attestation response")
		}
		v.Phase0.Signature = signature
		return nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return errors.New("no altair attestation response")
		}
		v.Altair.Signature = signature
		return nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return errors.New("no bellatrix attestation response")
		}
		v.Bellatrix.Signature = signature
		return nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return errors.New("no capella attestation response")
		}
		v.Capella.Signature = signature
		return nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return errors.New("no deneb attestation response")
		}
		v.Deneb.Signature = signature
		return nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return errors.New("no electra attestation response")
		}
		v.Electra.Signature = signature
		return nil
	default:
		return errors.New("unknown version")
	}
}

// AttestationData returns the AttestationData of the attestation response.
func (v *VersionedAttestationResponse) AttestationData() (*phase0.AttestationData, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no phase0 attestation resoponse")
		}

		return v.Phase0.Data, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no altair attestation resoponse")
		}

		return v.Altair.Data, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix attestation resoponse")
		}

		return v.Bellatrix.Data, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella attestation resoponse")
		}

		return v.Capella.Data, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb attestation resoponse")
		}

		return v.Deneb.Data, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra attestation resoponse")
		}

		return v.Electra.Data, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedAttestationResponse) String() string {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unknown version"
	}
}

// Functions for the ssz.HashRoot interface
func (v *VersionedAttestationResponse) GetTree() (*ssz.Node, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no phase0 attestation response")
		}

		return v.Phase0.GetTree()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no altair attestation resoponse")
		}

		return v.Altair.GetTree()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix attestation resoponse")
		}

		return v.Bellatrix.GetTree()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella attestation resoponse")
		}

		return v.Capella.GetTree()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb attestation resoponse")
		}

		return v.Deneb.GetTree()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra attestation resoponse")
		}

		return v.Electra.GetTree()
	default:
		return nil, errors.New("unknown version")
	}
}
func (v *VersionedAttestationResponse) HashTreeRoot() ([32]byte, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return [32]byte{}, errors.New("no phase0 attestation response")
		}

		return v.Phase0.HashTreeRoot()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return [32]byte{}, errors.New("no altair attestation resoponse")
		}

		return v.Altair.HashTreeRoot()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, errors.New("no bellatrix attestation resoponse")
		}

		return v.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella attestation resoponse")
		}

		return v.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, errors.New("no deneb attestation resoponse")
		}

		return v.Deneb.HashTreeRoot()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return [32]byte{}, errors.New("no electra attestation resoponse")
		}

		return v.Electra.HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}
func (v *VersionedAttestationResponse) HashTreeRootWith(hh ssz.HashWalker) error {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return errors.New("no phase0 attestation response")
		}

		return v.Phase0.HashTreeRootWith(hh)
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return errors.New("no altair attestation resoponse")
		}

		return v.Altair.HashTreeRootWith(hh)
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return errors.New("no bellatrix attestation resoponse")
		}

		return v.Bellatrix.HashTreeRootWith(hh)
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return errors.New("no capella attestation resoponse")
		}

		return v.Capella.HashTreeRootWith(hh)
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return errors.New("no deneb attestation resoponse")
		}

		return v.Deneb.HashTreeRootWith(hh)
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return errors.New("no electra attestation resoponse")
		}

		return v.Electra.HashTreeRootWith(hh)
	default:
		return errors.New("unknown version")
	}
}

// Returns a versioned attestation response without a signature
func ConstrusctVersionedAttestationResponseWithoutSignature(attestationData *phase0.AttestationData, version spec.DataVersion, validatorDuty *ValidatorDuty) *VersionedAttestationResponse {

	aggregationBitfield := bitfield.NewBitlist(validatorDuty.CommitteeLength)
	aggregationBitfield.SetBitAt(validatorDuty.ValidatorCommitteeIndex, true)

	switch version {
	case spec.DataVersionPhase0:
		return &VersionedAttestationResponse{
			Version: version,
			Phase0: &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			},
		}
	case spec.DataVersionAltair:
		return &VersionedAttestationResponse{
			Version: version,
			Altair: &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			},
		}
	case spec.DataVersionBellatrix:
		return &VersionedAttestationResponse{
			Version: version,
			Bellatrix: &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			},
		}
	case spec.DataVersionCapella:
		return &VersionedAttestationResponse{
			Version: version,
			Capella: &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			},
		}
	case spec.DataVersionDeneb:
		return &VersionedAttestationResponse{
			Version: version,
			Deneb: &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			},
		}
	case spec.DataVersionElectra:
		return &VersionedAttestationResponse{
			Version: version,
			Electra: &electra.SingleAttestation{
				CommitteeIndex: validatorDuty.CommitteeIndex,
				AttesterIndex:  validatorDuty.ValidatorIndex,
				Data:           attestationData,
			},
		}
	default:
		return nil
	}
}
