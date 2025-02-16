package testingutils

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

const (
	TestingDutySlot            = 12
	TestingDutySlot2           = 50
	TestingDutyEpoch           = 0
	TestingValidatorIndex      = 1
	TestingWrongValidatorIndex = 100

	UnknownDutyType = 100
)

// Util function for signing a beacon object with the appropriate domain
var signBeaconObject = func(obj ssz.HashRoot, domainType phase0.DomainType, ks *TestKeySet) phase0.BLSSignature {
	domain, _ := NewTestingBeaconNode().DomainData(1, domainType)
	ret, _, _ := NewTestingKeyManager().SignBeaconObject(obj, domain, ks.ValidatorPK.Serialize(), domainType)

	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], ret)

	return blsSig
}

// Util inline function for getting the SSZ root of an object
func GetSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

// Unknown duty type
var TestingUnknownDutyType = types.ValidatorDuty{
	Type:                    UnknownDutyType,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

// Wrong validator pub key
var TestingWrongDutyPK = types.ValidatorDuty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingWrongValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}
