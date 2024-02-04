package proposer

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InstanceNonProposer tests starting an instance for the proposer operator
func InstanceProposer() tests.SpecTest {
	pre := testingutils.BaseInstance()

	ks := testingutils.Testing4SharesSet()
	proposalMsg := testingutils.TestingProposalMessage(ks.Shares[1], 1)
	proposalMsg.Message.PrepareJustification = make([][]byte, 0)
	proposalMsg.Message.RoundChangeJustification = make([][]byte, 0)

	return &tests.MsgProcessingSpecTest{
		Name: "instance proposer",
		Pre:  pre,
		OutputMessages: []*qbft.SignedMessage{
			proposalMsg,
		},
		StartInstance: true,
		StartValue:    testingutils.TestingQBFTFullData,
	}
}

// InstanceNonProposer tests starting an instance for the non proposer operator
func InstanceNonProposer() tests.SpecTest {

	pre := testingutils.BaseInstanceWithOperatorID(2)
	return &tests.MsgProcessingSpecTest{
		Name:           "instance non proposer",
		Pre:            pre,
		OutputMessages: []*qbft.SignedMessage{},
		StartInstance:  true,
		StartValue:     testingutils.TestingQBFTFullData,
	}
}
