package types

// MaxCommitteeSize is the maximum number of operators in a committee.
//
// Note: SSZ struct tags and generated encoders still embed the numeric limit,
// but validation code should reference this constant to avoid drift.
const MaxCommitteeSize = 13
