package testingutils

import (
	"bytes"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

var TestingConfigAlea = func(keySet *TestKeySet) *alea.Config {
	return &alea.Config{
		Signer:    NewTestingKeyManager(),
		SigningPK: keySet.Shares[1].GetPublicKey().Serialize(),
		Domain:    types.PrimusTestnet,
		ValueCheckF: func(data []byte) error {
			if bytes.Equal(data, TestingInvalidValueCheck) {
				return errors.New("invalid value")
			}

			// as a base validation we do not accept nil values
			if len(data) == 0 {
				return errors.New("invalid value")
			}
			return nil
		},
		ProposerF: func(state *alea.State, round alea.Round) types.OperatorID {
			ans := int(round)%len(state.Share.Committee) + 1
			return types.OperatorID(ans)
		},
		Network: NewTestingNetworkAlea(),
		Timer:   NewTestingTimerAlea(),
	}
}

var TestingInvalidValueCheckAlea = []byte{1, 1, 1, 1}

var TestingShareAlea = func(keysSet *TestKeySet) *types.Share {
	return &types.Share{
		OperatorID:          1,
		ValidatorPubKey:     keysSet.ValidatorPK.Serialize(),
		SharePubKey:         keysSet.Shares[1].GetPublicKey().Serialize(),
		DomainType:          types.PrimusTestnet,
		Quorum:              keysSet.Threshold,
		PartialQuorum:       keysSet.PartialThreshold,
		Committee:           keysSet.Committee(),
		FeeRecipientAddress: TestingFeeRecipient,
	}
}

var BaseInstanceAlea = func() *alea.Instance {
	return baseInstanceAlea(TestingShareAlea(Testing4SharesSet()), Testing4SharesSet(), []byte{1, 2, 3, 4})
}

var SevenOperatorsInstanceAlea = func() *alea.Instance {
	return baseInstanceAlea(TestingShareAlea(Testing7SharesSet()), Testing7SharesSet(), []byte{1, 2, 3, 4})
}

var TenOperatorsInstanceAlea = func() *alea.Instance {
	return baseInstanceAlea(TestingShareAlea(Testing10SharesSet()), Testing10SharesSet(), []byte{1, 2, 3, 4})
}

var ThirteenOperatorsInstanceAlea = func() *alea.Instance {
	return baseInstanceAlea(TestingShareAlea(Testing13SharesSet()), Testing13SharesSet(), []byte{1, 2, 3, 4})
}

var baseInstanceAlea = func(share *types.Share, keySet *TestKeySet, identifier []byte) *alea.Instance {
	ret := alea.NewInstance(TestingConfigAlea(keySet), share, identifier, alea.FirstHeight)
	ret.StartValue = []byte{1, 2, 3, 4}
	return ret
}

func NewTestingAleaController(
	identifier []byte,
	share *types.Share,
	config alea.IConfig,
) *alea.Controller {
	return alea.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		config,
	)
}
