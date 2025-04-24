package testingutils

import (
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

var TestingValidatorPubKeyForValidatorIndex = func(ValidatorIndex phase0.ValidatorIndex) phase0.BLSPubKey {
	ks, exists := TestingKeySetMap[ValidatorIndex]
	if !exists {
		panic(fmt.Sprintf("Validator index %v does not exist in TestingKeySetMap", ValidatorIndex))
	}
	pk := ks.ValidatorPK
	pkHexString := pk.SerializeToHexStr()
	pkString, _ := hex.DecodeString(pkHexString)
	blsPK := phase0.BLSPubKey{}
	copy(blsPK[:], pkString)
	return blsPK
}

var TestingValidatorPubKeyList = func() []phase0.BLSPubKey {
	ret := make([]phase0.BLSPubKey, len(TestingKeySetMap))
	listIndex := 0
	for valIdx := range TestingKeySetMap {
		pk := TestingValidatorPubKeyForValidatorIndex(valIdx)
		ret[listIndex] = pk
		listIndex += 1
	}
	return ret
}()
