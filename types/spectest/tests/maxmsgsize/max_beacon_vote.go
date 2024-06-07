package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizeBeaconVote = 112
)

func maxBeaconVote() *types.BeaconVote {

	root := [32]byte{1}

	return &types.BeaconVote{
		BlockRoot: root,
		Source: &phase0.Checkpoint{
			Epoch: 1,
			Root:  root,
		},
		Target: &phase0.Checkpoint{
			Epoch: 2,
			Root:  root,
		},
	}
}

func MaxBeaconVote() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "max BeaconVote",
		Object:                maxBeaconVote(),
		ExpectedEncodedLength: maxSizeBeaconVote,
		IsMaxSize:             true,
	}
}
