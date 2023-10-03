package forkdata

import (
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Encoding tests encoding of a ssv message
func Encoding() *EncodingTest {
	fork_data := testingutils.TestingForkData

	byts, err := fork_data.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := fork_data.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "ForkData encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
