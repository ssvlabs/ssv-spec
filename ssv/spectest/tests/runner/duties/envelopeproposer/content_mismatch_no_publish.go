package envelopeproposer

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ContentMismatchNoPublish tests publish-by-content-match (SIP #94 §6): consensus decides another
// operator's envelope — valid (same §4 block root, self-build) but not byte-identical to this
// operator's own production (a different payload root). The duty finishes on post-consensus quorum
// without this operator publishing: only the producer of the decided envelope holds its body.
func ContentMismatchNoPublish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	slot := testingutils.TestingEnvelopeProposerDuty().Slot

	// Another operator's envelope for the same decided block: only the payload root differs.
	decidedEnvelope := testingutils.TestingBlindedExecutionPayloadEnvelope(slot)
	decidedEnvelope.PayloadRoot = phase0.Root{0x0f, 0xf0}
	decidedBytes, err := decidedEnvelope.Encode()
	if err != nil {
		panic(err.Error())
	}
	decidedCd := testingutils.TestingEnvelopeConsensusData(slot)
	decidedCd.DataSSZ = decidedBytes
	fullData, err := decidedCd.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err.Error())
	}
	id := types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleEnvelopeProposer)

	// The post-consensus partials sign the decided envelope's root, not this operator's own.
	postConsensusMsg := func(opID types.OperatorID) *types.PartialSignatureMessages {
		return testingutils.PostConsensusEnvelopeProposerMsgForEnvelope(ks.Shares[opID], opID, decidedEnvelope)
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "envelope proposer content mismatch no publish",
		Documentation: testdoc.EnvelopeProposerContentMismatchNoPublishDoc,
		Runner:        testingutils.EnvelopeProposerRunner(ks),
		Duty:          testingutils.TestingEnvelopeProposerDuty(),
		Messages: append(
			testingutils.SSVDecidingMsgsForHeightWithRoot(root, fullData, id[:], qbft.Height(slot), ks),
			[]*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, postConsensusMsg(1))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, postConsensusMsg(2))),
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgEnvelopeProposer(nil, postConsensusMsg(3))),
			}...,
		),
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PostConsensusEnvelopeProposerMsgForEnvelope(ks.Shares[1], 1, decidedEnvelope), // broadcasts when consensus decides
		},
		BeaconBroadcastedRoots: []string{},
	}
}
