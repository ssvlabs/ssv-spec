package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/types"
)

func validateShareMap(share map[phase0.ValidatorIndex]*types.Share) error {
	for validatorIndex, validatorShare := range share {
		if err := validatorShare.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid share for validator index %d", validatorIndex))
		}
	}
	return nil
}
