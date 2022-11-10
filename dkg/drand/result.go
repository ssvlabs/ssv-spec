package drand

import (
	dkg2 "github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/drand/kyber/share/dkg"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (d *DRand) getResult() *dkg.OptionResult {
	// TODO thread safe
	return d.result
}

func (d *DRand) setResult(result *dkg.OptionResult) {
	// TODO thread safe
	d.result = result
}

func (d *DRand) validateResult(result *dkg.OptionResult) error {
	if result == nil {
		return errors.New("no result")
	}
	return errors.Wrap(result.Error, "error on option result")
}

func (d *DRand) getProtocolOutcome(result *dkg.OptionResult) (*dkg2.ProtocolOutcome, error) {
	blame, err := d.getBlameOutput(result)
	if err != nil {
		return nil, errors.Wrap(err, "could not get blame output")
	}
	keygen, err := d.getKeygenOutput(result.Result.Key)
	if err != nil {
		return nil, errors.Wrap(err, "could not get keygen output")
	}
	return &dkg2.ProtocolOutcome{
		BlameOutput:    blame,
		ProtocolOutput: keygen,
	}, nil
}

func (d *DRand) getBlameOutput(result *dkg.OptionResult) (*dkg2.BlameOutput, error) {
	//if len(result.Result.QUAL) != len(d.operators) {
	//	return errors.New("not all operators qualified")
	//}
	return nil, nil
}

func (d *DRand) getKeygenOutput(result *dkg.DistKeyShare) (*dkg2.KeyGenOutput, error) {
	// generate validator pubkey and shares pub keys
	pks, err := pubKeysFromPolyEvaluation(d.config.Suite, result.Commitments(), append([]uint32{0}, d.operators...))
	if err != nil {
		return nil, errors.Wrap(err, "could not get pub keys from commitments")
	}
	validatorPK := pks[0]
	pks = pks[1:]
	operatorPks := make(map[types.OperatorID]*bls.PublicKey)
	for i, id := range d.operators {
		operatorPks[types.OperatorID(id)] = pks[i]
	}

	// get share sk
	byts, err := result.Share.V.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not get share secret key bytes")
	}
	sk := &bls.SecretKey{}
	if err := sk.Deserialize(byts); err != nil {
		return nil, errors.Wrap(err, "could not get share secret key")
	}

	return &dkg2.KeyGenOutput{
		Share:           sk,
		OperatorPubKeys: operatorPks,
		ValidatorPK:     validatorPK.Serialize(),
		Threshold:       d.threshold,
	}, nil
}

func (d *DRand) getKeyGenOutput() {

}
