package types

// ShareMember holds ShareValidatorPK and ValidatorIndex
type ValidatorShare struct {
	OperatorID     OperatorID
	OperatorPubKey []byte           `ssz-size:"294"`
	SharePubKey    ShareValidatorPK `ssz-size:"48"`
}
