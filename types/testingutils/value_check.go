package testingutils

import (
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

func UnknownDutyValueCheck() qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}

// GetSlashableRootForBeaconVote returns the AttestationData's root created from the message's BeaconVote as done in ssv/value_check.go
func GetSlashableRootForBeaconVote(msg *types.SignedSSVMessage) []byte {
	// Get BeaconVote from FullData
	beaconVote := &types.BeaconVote{}
	err := beaconVote.Decode(msg.FullData)
	if err != nil {
		panic(err)
	}
	// Get slot
	consensusMessage, err := qbft.DecodeMessage(msg.SSVMessage.Data)
	if err != nil {
		panic(err)
	}
	slot := phase0.Slot(consensusMessage.Height)
	// Construct AttestationData
	attestationData := &phase0.AttestationData{
		Slot:            slot,
		Index:           math.MaxUint64,
		BeaconBlockRoot: beaconVote.BlockRoot,
		Target:          beaconVote.Target,
		Source:          beaconVote.Source,
	}
	// Get root
	dataRoot, err := attestationData.HashTreeRoot()
	if err != nil {
		panic(err)
	}

	return dataRoot[:]
}
