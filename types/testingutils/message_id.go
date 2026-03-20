package testingutils

import "github.com/ssvlabs/ssv-spec/types"

var CommitteeMsgID = func(keySet *TestKeySet) []byte {
	// Identifier
	committee := make([]uint64, 0)
	for _, op := range keySet.Committee() {
		committee = append(committee, op.Signer)
	}
	committeeID := types.GetCommitteeID(committee)

	ret := types.NewCommitteeMsgID(TestingSSVDomainType, committeeID, types.RoleCommittee)
	return ret[:]
}

var AggregatorCommitteeMsgID = func(keySet *TestKeySet) []byte {
	// Identifier
	committee := make([]uint64, 0)
	for _, op := range keySet.Committee() {
		committee = append(committee, op.Signer)
	}
	committeeID := types.GetCommitteeID(committee)

	ret := types.NewCommitteeMsgID(TestingSSVDomainType, committeeID, types.RoleAggregatorCommittee)
	return ret[:]
}

var ProposerMsgID = func() []byte {
	ret := types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleProposer)
	return ret[:]
}()

var ValidatorRegistrationMsgID = func() []byte {
	ret := types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleValidatorRegistration)
	return ret[:]
}()

var VoluntaryExitMsgID = func() []byte {
	ret := types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleVoluntaryExit)
	return ret[:]
}()
