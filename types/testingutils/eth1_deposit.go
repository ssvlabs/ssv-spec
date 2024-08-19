package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// GenerateETHDepositData returns un-signed deposit data and deposit data root for signature
func GenerateETHDepositData(
	validatorPK, withdrawalCredentials []byte,
	fork phase0.Version,
	domain phase0.DomainType) ([]byte, *phase0.DepositData, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], validatorPK)

	ret := &phase0.DepositMessage{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(types.MaxEffectiveBalanceInGwei),
	}

	domainR, err := types.ComputeETHDomain(domain, fork, types.GenesisValidatorsRoot)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit domain")
	}
	signingRoot, err := types.ComputeETHSigningRoot(ret, domainR)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit signing root")
	}
	return signingRoot[:], &phase0.DepositData{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(types.MaxEffectiveBalanceInGwei),
	}, nil
}
