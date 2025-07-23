package testingutils

import (
	"crypto/rsa"
	"encoding/hex"

	"github.com/ssvlabs/ssv-spec/types"
)

type PrivateKeyInfo struct {
	ValidatorSK                             string
	ValidatorPK                             string
	ShareCount, Threshold, PartialThreshold uint64
	Shares                                  map[types.OperatorID]string
	OperatorKeys                            map[types.OperatorID]string
	DKGOperators                            map[types.OperatorID]struct {
		SK            string
		ETHAddress    string
		EncryptionKey string
	}
}

func BuildPrivateKeyInfo(ks *TestKeySet) *PrivateKeyInfo {
	privateKeyInfo := &PrivateKeyInfo{
		ValidatorSK:      hex.EncodeToString(ks.ValidatorSK.Serialize()),
		ValidatorPK:      hex.EncodeToString(ks.ValidatorPK.Serialize()),
		ShareCount:       ks.ShareCount,
		Threshold:        ks.Threshold,
		PartialThreshold: ks.PartialThreshold,
		Shares:           make(map[types.OperatorID]string),
		OperatorKeys:     make(map[types.OperatorID]string),
		DKGOperators: make(map[types.OperatorID]struct {
			SK            string
			ETHAddress    string
			EncryptionKey string
		}),
	}

	for operatorID, shareSK := range ks.Shares {
		privateKeyInfo.Shares[operatorID] = hex.EncodeToString(shareSK.Serialize())
	}

	for operatorID, operatorKey := range ks.OperatorKeys {
		privateKeyInfo.OperatorKeys[operatorID] = exportRsaPrivateKey(operatorKey)
	}

	for operatorID, dkgOperator := range ks.DKGOperators {
		privateKeyInfo.DKGOperators[operatorID] = struct {
			SK            string
			ETHAddress    string
			EncryptionKey string
		}{
			SK:            dkgOperator.SK.D.String(),
			ETHAddress:    dkgOperator.ETHAddress.Hex(),
			EncryptionKey: exportRsaPrivateKey(dkgOperator.EncryptionKey),
		}
	}

	return privateKeyInfo
}

func exportRsaPrivateKey(sk *rsa.PrivateKey) string {
	return hex.EncodeToString(types.PrivateKeyToPem(sk))
}
