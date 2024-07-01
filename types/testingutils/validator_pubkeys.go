package testingutils

import (
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

var TestingValidatorPubKeyForValidatorIndex = func(ValidatorIndex phase0.ValidatorIndex) spec.BLSPubKey {
	ks, exists := TestingKeySetMap[ValidatorIndex]
	if !exists {
		panic(fmt.Sprintf("Validator index %v does not exist in TestingKeySetMap", ValidatorIndex))
	}
	pk := ks.ValidatorPK
	pkHexString := pk.SerializeToHexStr()
	pkString, _ := hex.DecodeString(pkHexString)
	blsPK := spec.BLSPubKey{}
	copy(blsPK[:], pkString)
	return blsPK
}

var TestingValidatorPubKeyList = func() []spec.BLSPubKey {
	ret := make([]spec.BLSPubKey, len(TestingKeySetMap))
	listIndex := 0
	for valIdx := range TestingKeySetMap {
		pk := TestingValidatorPubKeyForValidatorIndex(valIdx)
		ret[listIndex] = pk
		listIndex += 1
	}
	return ret
}()
