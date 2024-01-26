package ssv

import (
	"bytes"
	"sort"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

func (b *BaseRunner) ValidatePreConsensusMsg(runner Runner, signedMsg *types.SignedPartialSignatureMessage) error {
	if !b.hasRunningDuty() {
		return errors.New("no running duty")
	}

	if err := b.validatePartialSigMsgForSlot(signedMsg, b.State.StartingDuty.Slot); err != nil {
		return err
	}

	roots, domain, err := runner.expectedPreConsensusRootsAndDomain()
	if err != nil {
		return err
	}

	return b.verifyExpectedRoot(runner, signedMsg, roots, domain)
}

func (b *BaseRunner) ValidatePostConsensusMsg(runner Runner, signedMsg *types.SignedPartialSignatureMessage) error {
	if !b.hasRunningDuty() {
		return errors.New("no running duty")
	}

	// TODO https://github.com/bloxapp/ssv-spec/issues/142 need to fix with this issue solution instead.
	if b.State.DecidedValue == nil {
		return errors.New("no decided value")
	}

	if b.State.RunningInstance == nil {
		return errors.New("no running consensus instance")
	}
	decided, decidedValueByts := b.State.RunningInstance.IsDecided()
	if !decided {
		return errors.New("consensus instance not decided")
	}

	decidedValue := &types.ConsensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return errors.Wrap(err, "failed to parse decided value to ConsensusData")
	}

	if err := b.validatePartialSigMsgForSlot(signedMsg, decidedValue.Duty.Slot); err != nil {
		return err
	}

	roots, domain, err := runner.expectedPostConsensusRootsAndDomain()
	if err != nil {
		return err
	}

	return b.verifyExpectedRoot(runner, signedMsg, roots, domain)
}

func (b *BaseRunner) validateDecidedConsensusData(runner Runner, val *types.ConsensusData) error {
	byts, err := val.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	if err := runner.GetValCheckF()(byts); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	return nil
}

func (b *BaseRunner) verifyExpectedRoot(runner Runner, signedMsg *types.SignedPartialSignatureMessage, expectedRootObjs []ssz.HashRoot, domain spec.DomainType) error {
	// Check length
	if len(expectedRootObjs) != len(signedMsg.Message.Messages) {
		return errors.New("wrong expected roots count")
	}

	// Transform expectedRoots ([]ssz.HashRoot) into [][32]byte (or []phase0.Root)
	epoch := b.BeaconNetwork.EstimatedEpochAtSlot(b.State.StartingDuty.Slot)
	d, err := runner.GetBeaconNode().DomainData(epoch, domain)
	if err != nil {
		return errors.Wrap(err, "could not get pre consensus root domain")
	}
	expectedBeaconRoots := make([][32]byte, 0)
	for _, expectedRoot := range expectedRootObjs {

		beaconRoot, err := b.GetBeaconSigningRoot(expectedRoot, d)
		if err != nil {
			return errors.Wrap(err, "could not compute ETH signing root")
		}
		expectedBeaconRoots = append(expectedBeaconRoots, beaconRoot)
	}

	// Get roots from SignedPartialSignatureMessage
	receivedRoots := make([][32]byte, 0)
	for _, msg := range signedMsg.Message.Messages {
		receivedRoots = append(receivedRoots, msg.SigningRoot)
	}

	// Compare roots
	return b.compareRoots(receivedRoots, expectedBeaconRoots)
}

// Compares the sorted version of two lists of roots
func (b *BaseRunner) compareRoots(roots [][32]byte, expectedRoots [][32]byte) error {
	// Check length
	if len(expectedRoots) != len(roots) {
		return errors.New("wrong expected roots count")
	}

	// copy and sort function
	sortCopy := func(r [][32]byte) [][32]byte {
		ret := make([][32]byte, len(r))
		for i, ri := range r {
			copy(ret[i][:], ri[:])
		}
		sort.Slice(ret, func(i, j int) bool {
			return string(ret[i][:]) < string(ret[j][:])
		})
		return ret
	}

	// Sort both lists
	sortedExpectedRoots := sortCopy(expectedRoots)
	sortedRoots := sortCopy(roots)

	// Compare each root
	for i, r := range sortedRoots {
		if !bytes.Equal(sortedExpectedRoots[i][:], r[:]) {
			return errors.New("wrong signing root")
		}
	}
	return nil
}

func (b *BaseRunner) GetBeaconSigningRoot(root ssz.HashRoot, domain spec.Domain) (spec.Root, error) {
	return types.ComputeETHSigningRoot(root, domain)
}
