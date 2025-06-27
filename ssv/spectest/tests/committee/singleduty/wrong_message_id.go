package committeesingleduty

import (
	"crypto/rsa"
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongMessageID tests a message that is processed by the committee with wrong message ID
func WrongMessageID() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	// Incorrect ID for SSVMessage
	getPubkeyID := func() types.MessageID {
		ret := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleCommittee)
		return ret
	}

	// Function to get decided message with wrong ID for role
	decidedMessage := func() *types.SignedSSVMessage {
		msgID := getPubkeyID()
		signedMessage := testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{1, 2, 3},
			testingutils.TestingDutySlot,
			msgID[:])

		signedMessage.SSVMessage.MsgID = msgID

		sig1 := testingutils.SignedSSVMessageWithSigner(1, ks.OperatorKeys[1], signedMessage.SSVMessage).Signatures[0]
		sig2 := testingutils.SignedSSVMessageWithSigner(2, ks.OperatorKeys[2], signedMessage.SSVMessage).Signatures[0]
		sig3 := testingutils.SignedSSVMessageWithSigner(3, ks.OperatorKeys[3], signedMessage.SSVMessage).Signatures[0]

		signedMessage.Signatures = [][]byte{sig1, sig2, sig3}

		return signedMessage
	}

	expectedError := "Message invalid: msg ID doesn't match committee ID"

	validatorsIndexList := testingutils.ValidatorIndexList(1)
	ksMap := testingutils.KeySetMapForValidators(1)
	tests := []*committee.CommitteeSpecTest{}

	for _, version := range testingutils.SupportedAttestationVersions {
		tests = append(tests, []*committee.CommitteeSpecTest{
			{
				Name:      fmt.Sprintf("sync committees (%s)", version.String()),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					decidedMessage(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("attestation (%s)", version.String()),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					decidedMessage(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		}...)
	}

	multiSpecTest := committee.NewMultiCommitteeSpecTest(
		"wrong message ID",
		"Tests committee behavior when processing messages with incorrect message IDs that don't match the committee ID",
		tests,
	)

	return multiSpecTest
}
