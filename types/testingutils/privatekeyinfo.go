package testingutils

import (
	"encoding/hex"
	"fmt"

	"github.com/ssvlabs/ssv-spec/types"
)

type PrivateKeyInfo struct {
	ValidatorSK  string
	Shares       map[types.OperatorID]string
	OperatorKeys map[types.OperatorID]string
}

func BuildPrivateKeyInfo(ks *TestKeySet) *PrivateKeyInfo {
	privateKeyInfo := &PrivateKeyInfo{
		ValidatorSK:  hex.EncodeToString(ks.ValidatorSK.Serialize()),
		Shares:       make(map[types.OperatorID]string),
		OperatorKeys: make(map[types.OperatorID]string),
	}

	for operatorID, shareSK := range ks.Shares {
		privateKeyInfo.Shares[operatorID] = hex.EncodeToString(shareSK.Serialize())
	}

	for operatorID, operatorKey := range ks.OperatorKeys {
		privateKeyInfo.OperatorKeys[operatorID] = fmt.Sprintf("N:%s,E:%d", operatorKey.N.String(), operatorKey.E)
	}

	return privateKeyInfo
}
