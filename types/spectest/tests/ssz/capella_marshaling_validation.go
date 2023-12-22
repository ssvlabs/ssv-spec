package ssz

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CapellaBeaconBlockMarshalling tests a valid capella marshaling
func CapellaBeaconBlockMarshalling() *SSZSpecTest {

	root, err := testingutils.TestingBeaconBlockCapella.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &SSZSpecTest{
		Name:         "capella beacon block marshalling",
		Data:         testingutils.TestProposerConsensusDataBytsV(spec.DataVersionCapella),
		ExpectedRoot: root,
	}
}
