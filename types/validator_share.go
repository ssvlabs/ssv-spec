package types

type ValidatorShare struct {
	OperatorID     OperatorID
	OperatorPubKey []byte           `ssz-size:"294"`
	SharePubKey    ShareValidatorPK `ssz-size:"48"`
}
