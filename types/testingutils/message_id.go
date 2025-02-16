package testingutils

import "github.com/ssvlabs/ssv-spec/types"

var AggregatorMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleAggregator)
	return ret[:]
}()

var CommitteeMsgID = func(keySet *TestKeySet) []byte {

	// Identifier
	committee := make([]uint64, 0)
	for _, op := range keySet.Committee() {
		committee = append(committee, op.Signer)
	}
	committeeID := types.GetCommitteeID(committee)

	ret := types.NewMsgID(TestingSSVDomainType, committeeID[:], types.RoleCommittee)
	return ret[:]
}

var ProposerMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleProposer)
	return ret[:]
}()

var SyncCommitteeContributionMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleSyncCommitteeContribution)
	return ret[:]
}()

var ValidatorRegistrationMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleValidatorRegistration)
	return ret[:]
}()

var VoluntaryExitMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleVoluntaryExit)
	return ret[:]
}()
