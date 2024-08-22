package testingutils

import (
	"bytes"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
)

func validatorPKFromString(pks string) types.ValidatorPK {
	pk := types.ValidatorPK{}
	copy(pk[:], pks)
	return pk
}

// Mock values
var (
	TestingNonAttestingValidator  = validatorPKFromString("non active validator")
	TestingLiquidatedValidator    = validatorPKFromString("liquidated validator")
	TestingNonExistentValidator   = validatorPKFromString("non existent validator")
	TestingNonExistentCommitteeID = validatorPKFromString("non existent CommitteeID")

	TestingValidatorPK                                 = types.ValidatorPK([48]byte{1, 2, 3, 4})
	TestingValidatorPKWithoutProposerDuty              = validatorPKFromString("no proposer duty")
	TestingValidatorPKWithoutSyncCommitteeContribution = validatorPKFromString("no sync committee contribution duty")
	TestingValidatorPKWithSyncCommittee                = validatorPKFromString("with sync committee duty")

	TestingCommitteeID                      = TestingCommitteeMember(Testing4SharesSet()).CommitteeID
	TestingCommitteeIDWithSyncCommitteeDuty = TestingCommitteeMember(Testing13SharesSet()).CommitteeID

	TestingTopic      = "valid topic"
	TestingWrongTopic = "wrong topic"
)

type TestingNetworkDataFetcher struct {
	Committees                map[types.CommitteeID]*validation.CommitteeInfo
	CommitteeIDForValidatorPK map[types.ValidatorPK]types.CommitteeID
}

func (n *TestingNetworkDataFetcher) addCommitteeInfo(committeeInfo *validation.CommitteeInfo, validatorPKs []types.ValidatorPK) {
	n.Committees[committeeInfo.CommitteeID] = committeeInfo

	for _, validatorPK := range validatorPKs {
		n.CommitteeIDForValidatorPK[validatorPK] = committeeInfo.CommitteeID
	}
}

func NewTestingNetworkDataFetcher() *TestingNetworkDataFetcher {

	netDF := &TestingNetworkDataFetcher{
		Committees:                make(map[types.CommitteeID]*validation.CommitteeInfo),
		CommitteeIDForValidatorPK: make(map[types.ValidatorPK]types.CommitteeID),
	}

	// Add a committee with 10 validators
	netDF.addCommitteeInfo(keySetToCommitteeInfo(Testing4SharesSet(),
		[]phase0.ValidatorIndex{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		[]types.ValidatorPK{
			TestingValidatorPK, validatorPKFromString("2"), validatorPKFromString("3"),
			validatorPKFromString("4"), validatorPKFromString("5"), validatorPKFromString("6"),
			validatorPKFromString("7"), validatorPKFromString("8"), validatorPKFromString("9"),
			validatorPKFromString("10"),
		},
	)

	// Add a standard committee with the validator that doesn't have a proposer duty, according to the DutyFetcher
	netDF.addCommitteeInfo(keySetToCommitteeInfo(Testing7SharesSet(), []phase0.ValidatorIndex{ValidatorIndexWithoutProposerDuty}), []types.ValidatorPK{TestingValidatorPKWithoutProposerDuty})

	// Add a committee with the validator that doesn't have a sync committee contribution duty, according to the DutyFetcher
	netDF.addCommitteeInfo(keySetToCommitteeInfo(Testing10SharesSet(), []phase0.ValidatorIndex{ValidatorIndexWithoutSyncCommitteeContributionDuty}), []types.ValidatorPK{TestingValidatorPKWithoutSyncCommitteeContribution})

	// Add a committee with a validator that has a sync committee duty, according to the DutyFetcher
	netDF.addCommitteeInfo(keySetToCommitteeInfo(Testing13SharesSet(), []phase0.ValidatorIndex{ValidatorIndexWithSyncCommitteeDuty}), []types.ValidatorPK{TestingValidatorPKWithSyncCommittee})

	return netDF
}

func (n *TestingNetworkDataFetcher) ValidDomain(domain []byte) bool {
	return bytes.Equal(domain, TestingSSVDomainType[:])
}
func (n *TestingNetworkDataFetcher) CorrectTopic(committee []types.OperatorID, topic string) bool {
	return (topic == TestingTopic)
}
func (n *TestingNetworkDataFetcher) GetCommitteeInfo(msgID types.MessageID) *validation.CommitteeInfo {
	if msgID.GetRoleType() == types.RoleCommittee {
		committeeID := msgID.GetDutyExecutorID()
		return n.Committees[types.CommitteeID(committeeID[16:])]
	} else {
		validatorPK := msgID.GetDutyExecutorID()
		committeeID := n.CommitteeIDForValidatorPK[types.ValidatorPK(validatorPK)]
		return n.Committees[committeeID]
	}
}
func (n *TestingNetworkDataFetcher) ExistingValidator(validatorPK types.ValidatorPK) bool {
	return !bytes.Equal(validatorPK[:], TestingNonExistentValidator[:])
}
func (n *TestingNetworkDataFetcher) ActiveValidator(validatorPK types.ValidatorPK) bool {
	return !bytes.Equal(validatorPK[:], TestingNonAttestingValidator[:])
}
func (n *TestingNetworkDataFetcher) ValidatorLiquidated(validatorPK types.ValidatorPK) bool {
	return bytes.Equal(validatorPK[:], TestingLiquidatedValidator[:])
}
func (n *TestingNetworkDataFetcher) ExistingCommitteeID(committeeID types.CommitteeID) bool {
	return !bytes.Equal(committeeID[:], TestingNonExistentCommitteeID[:])
}

func keySetToCommitteeInfo(ks *TestKeySet, validatorIndex []phase0.ValidatorIndex) *validation.CommitteeInfo {
	cb := TestingCommitteeMember(ks)
	return &validation.CommitteeInfo{
		Validators:  validatorIndex,
		Operators:   cb.Committee,
		CommitteeID: cb.CommitteeID,
	}
}
