package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))

	byts, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := msg.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return NewEncodingTest(
		"SSVMessage encoding",
		"Test SSVMessage encoding",
		byts,
		root,
	)
}
