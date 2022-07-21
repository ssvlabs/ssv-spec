package types

import (
	"encoding/json"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type ExtendedDepositData struct {
	phase0.DepositData
	DepositMessageRoot phase0.Root
	DepositDataRoot    phase0.Root
	ForkVersion        phase0.Version
	CliVersion         string
}

type extendedDepositDataJSON struct {
	PublicKey             string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                uint64 `json:"amount"`
	Signature             string `json:"signature"`
	DepositMessageRoot    string `json:"deposit_message_root"`
	DepositDataRoot       string `json:"deposit_data_root"`
	ForkVersion           string `json:"fork_version"`
	CliVersion            string `json:"deposit_cli_version"`
}

// MarshalJSON implements json.Marshaler.
func (d *ExtendedDepositData) MarshalJSON() ([]byte, error) {
	depDataRoot, err := d.DepositData.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&extendedDepositDataJSON{
		PublicKey:             fmt.Sprintf("%x", d.PublicKey),
		WithdrawalCredentials: fmt.Sprintf("%x", d.WithdrawalCredentials),
		Amount:                uint64(d.Amount),
		Signature:             fmt.Sprintf("%x", d.Signature),
		DepositMessageRoot:    fmt.Sprintf("%x", d.DepositMessageRoot),
		DepositDataRoot:       fmt.Sprintf("%x", depDataRoot),
		ForkVersion:           fmt.Sprintf("%x", d.ForkVersion),
		CliVersion:            d.CliVersion,
	})
}

// GenerateETHDepositData returns un-signed deposit data and deposit data root for signature
func GenerateETHDepositData(
	validatorPK, withdrawalCredentials []byte,
	fork phase0.Version,
	domain phase0.DomainType) ([]byte, *ExtendedDepositData, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], validatorPK)

	ret := phase0.DepositMessage{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
	}
	messageRoot, err := ret.HashTreeRoot()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get DepositData root")
	}

	domainR, err := ComputeETHDomain(domain, fork)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit domain")
	}
	signingRoot, err := ComputeETHSigningRoot(messageRoot, domainR)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit signing root")
	}
	return signingRoot[:], &ExtendedDepositData{
		DepositData: phase0.DepositData{
			PublicKey:             pk,
			WithdrawalCredentials: withdrawalCredentials,
			Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
		},
		DepositMessageRoot: messageRoot,
		ForkVersion:        fork,
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
