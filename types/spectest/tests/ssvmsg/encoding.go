package ssvmsg

import (
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	signedMsg := testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1))

	msg := signedMsg.SSVMessage

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := msg.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "SSVMessage encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
