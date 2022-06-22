package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// GenerateETHDepositData returns un-signed deposit data and deposit data root for signature
func GenerateETHDepositData(
	validatorPK, withdrawalCredentials []byte,
	fork phase0.Version,
	domain phase0.DomainType) ([]byte, *phase0.DepositData, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], validatorPK)

	ret := phase0.DepositMessage{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
	}
	r, err := ret.HashTreeRoot()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get DepositData root")
	}

	domainR, err := ComputeETHDomain(domain, fork)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit domain")
	}
	signingRoot, err := ComputeETHSigningRoot(r, domainR)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit signing root")
	}
	return signingRoot[:], &phase0.DepositData{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
	}, nil
}

// ComputeETHDomain returns computed domain
func ComputeETHDomain(domain phase0.DomainType, fork phase0.Version) (phase0.Domain, error) {
	ret := phase0.Domain{}
	copy(ret[0:4], domain[:])

	forkData := phase0.ForkData{
		CurrentVersion:        fork,
		GenesisValidatorsRoot: phase0.Root{},
	}
	forkDataRoot, err := forkData.HashTreeRoot()
	if err != nil {

	}
	copy(ret[4:32], forkDataRoot[0:28])
	return ret, nil
}

func ComputeETHSigningRoot(root phase0.Root, domain phase0.Domain) (phase0.Root, error) {
	signingContainer := phase0.SigningData{
		ObjectRoot: root,
		Domain:     domain,
	}
	return signingContainer.HashTreeRoot()
}
