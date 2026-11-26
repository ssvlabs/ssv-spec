package runner

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasProposerJustifiedValueReproposed tests that a justified value is re-proposed unchanged, payload_root
// included, whatever the re-proposing leader's own bid or payload was (SIP #94 §4). This operator's own produce
// is an external bid, so the value it proposes in the first round has a zero payload_root. A quorum of peers
// then changes round, each having prepared the self-build value. Leading the next round, this operator proposes
// that value, with its payload_root, and not its own.
func GloasProposerJustifiedValueReproposed() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas
	height := qbft.Height(testingutils.TestingDutySlotV(version))

	ownValue, err := testingutils.TestProposerExternalBidConsensusDataV(version).Encode()
	if err != nil {
		panic(err.Error())
	}
	justifiedValue := testingutils.TestProposerConsensusDataBytsV(version)
	justifiedRoot, err := qbft.HashDataRoot(justifiedValue)
	if err != nil {
		panic(err.Error())
	}

	// The peers' prepares for the self-build value in the first round: the justification of their round changes.
	peers := []types.OperatorID{2, 3, 4}
	var prepares []*types.SignedSSVMessage
	for _, opID := range peers {
		prepares = append(prepares, testingutils.SignQBFTMsg(ks.OperatorKeys[opID], opID, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Root:       justifiedRoot,
		}))
	}

	messages := []*types.SignedSSVMessage{
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
	}
	for _, opID := range peers {
		roundChange := testingutils.SignQBFTMsg(ks.OperatorKeys[opID], opID, &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   height,
			Round:                    2,
			Identifier:               testingutils.ProposerMsgID,
			Root:                     justifiedRoot,
			DataRound:                qbft.FirstRound,
			RoundChangeJustification: testingutils.MarshalJustifications(prepares),
		})
		roundChange.FullData = justifiedValue
		messages = append(messages, roundChange)
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "gloas proposer justified value reproposed",
		Documentation: testdoc.GloasProposerJustifiedValueReproposedDoc,
		Runner:        testingutils.ProposerRunner(ks),
		Duty:          testingutils.TestingProposerDutyV(version),
		Messages:      messages,
		OutputMessages: []*types.PartialSignatureMessages{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
		},
		// The first round's proposal is this operator's own value; the second round's is the justified one.
		QBFTProposals: [][]byte{ownValue, justifiedValue},
		BeaconNode:    &tests.BeaconNodeBehaviour{ProducesExternalBid: true},
	}
}
