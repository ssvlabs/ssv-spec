package ssz

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DenebBeaconBlockMarshalling tests a valid deneb beacon block marshaling
func DenebBeaconBlockMarshalling() *SSZSpecTest {

	root, err := testingutils.TestingBeaconBlockDeneb.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &SSZSpecTest{
		Name:         "deneb beacon block marshalling",
		Data:         testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
		ExpectedRoot: root,
	}
}
