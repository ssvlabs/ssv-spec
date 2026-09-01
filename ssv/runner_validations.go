package ssv

import (
	"bytes"
	"sort"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/types"
)

func (b *BaseRunner) ValidatePreConsensusMsg(runner Runner, psigMsgs *types.PartialSignatureMessages) error {
	if !b.hasRunningDuty() {
		return types.NewError(types.NoRunningDutyErrorCode, "no running duty")
	}

	if _, ok := runner.(*AggregatorCommitteeRunner); ok {
		// For aggregator committee runner, a special validation is applied
		// since committee views may differ between operators.
		return b.validatePartialSigMsgForSlot(psigMsgs, b.State.StartingDuty.DutySlot())
	} else {
		// For other runner types, the pre-consensus is for a single validator,
		// and more strict validation can be applied.
		if err := b.validatePartialSigMsgForSlot(psigMsgs, b.State.StartingDuty.DutySlot()); err != nil {
			return err
		}

		if err := b.validateValidatorIndexInPartialSigMsg(psigMsgs); err != nil {
			return err
		}

		roots, domain, err := runner.expectedPreConsensusRootsAndDomain()
		if err != nil {
			return err
		}

		return b.verifyExpectedRoot(runner, psigMsgs, roots, domain)
	}
}

// Verify each signature in container removing the invalid ones
func (b *BaseRunner) FallBackAndVerifyEachSignature(container *PartialSigContainer, root [32]byte,
	committee []*types.ShareMember, validatorIndex phase0.ValidatorIndex) {

	signatures := container.GetSignatures(validatorIndex, root)

	for operatorID, signature := range signatures {
		if err := b.verifyBeaconPartialSignature(operatorID, signature, root, committee); err != nil {
			container.Remove(validatorIndex, operatorID, root)
		}
	}
}

func (b *BaseRunner) ValidatePostConsensusMsg(runner Runner, psigMsgs *types.PartialSignatureMessages) error {
	if !b.hasRunningDuty() {
		// A finished Gloas proposer keeps accepting post-consensus packets while an expected §6 envelope
		// root has not yet reached quorum (SIP #94 §4), so an envelope quorum reached only after the block's
		// still publishes the reveal; every other finished duty (and a never-started one) rejects.
		if p, ok := runner.(*ProposerRunner); !ok || !p.awaitingEnvelope() {
			return types.NewError(types.NoRunningDutyErrorCode, "no running duty")
		}
	}

	// TODO https://github.com/ssvlabs/ssv-spec/issues/142 need to fix with this issue solution instead.
	if len(b.State.DecidedValue) == 0 {
		return types.NewError(types.NoDecidedValueErrorCode, "no decided value")
	}

	if b.State.RunningInstance == nil {
		return types.NewError(types.NoRunningConsensusInstanceErrorCode, "no running consensus instance")
	}
	decided, decidedValueBytes := b.State.RunningInstance.IsDecided()
	if !decided {
		return types.NewError(types.ConsensusInstanceNotDecidedErrorCode, "consensus instance not decided")
	}

	switch runner.(type) {
	case *CommitteeRunner:
		// The decided value is a GloasBeaconVote at Gloas slots (SIP #94 §2), a BeaconVote before;
		// decode the shape the duty's fork mandates so a cross-fork value is rejected here too.
		decidedVote := committeeVoteForSlot(runner.GetBeaconNode(), b.State.StartingDuty.DutySlot())
		if err := decidedVote.Decode(decidedValueBytes); err != nil {
			return errors.Wrap(err, "failed to parse decided value to BeaconData")
		}

		return b.validatePartialSigMsgForSlot(psigMsgs, b.State.StartingDuty.DutySlot())
	case *AggregatorCommitteeRunner:
		decidedValue := &types.AggregatorCommitteeConsensusData{}
		if err := decidedValue.Decode(decidedValueBytes); err != nil {
			return errors.Wrap(err, "failed to parse decided value to AggregatorCommitteeConsensusData")
		}

		return b.validatePartialSigMsgForSlot(psigMsgs, b.State.StartingDuty.DutySlot())
	default:
		// Decode still runs to reject a malformed stored value, but the slot to validate partials against
		// is the running duty's — as in the committee/aggregator cases above — not the decided value's own
		// Duty.Slot, so the check never depends on trusting the decided value's contents.
		decidedValue := &types.ProposerConsensusData{}
		if err := decidedValue.Decode(decidedValueBytes); err != nil {
			return errors.Wrap(err, "failed to parse decided value to ProposerConsensusData")
		}

		if err := b.validatePartialSigMsgForSlot(psigMsgs, b.State.StartingDuty.DutySlot()); err != nil {
			return err
		}

		if err := b.validateValidatorIndexInPartialSigMsg(psigMsgs); err != nil {
			return err
		}

		expected, err := runner.expectedPostConsensusRootsAndDomains()
		if err != nil {
			return err
		}

		return b.verifyExpectedPostConsensusRoots(runner, psigMsgs, expected)
	}
}

func (b *BaseRunner) validateDecidedConsensusData(runner Runner, val types.Encoder) error {
	byts, err := val.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	if err := runner.GetValCheckF()(byts); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	return nil
}

// verifyExpectedPostConsensusRoots validates a post-consensus packet against per-root domains (SIP #94
// §4): every message root must equal an expected root computed under its own domain, each expected root is
// covered at most once, and every required root must be present. An unexpected root, a duplicate, or a
// missing required root rejects the whole packet. Optional roots — the Gloas proposer's §6 envelope entry —
// may be absent, so a block-only packet stays valid. Matching the covered set of expected roots states
// §4's "block plus optional envelope" directly, where a bare entry count would let a duplicate block pass.
func (b *BaseRunner) verifyExpectedPostConsensusRoots(runner Runner, psigMsgs *types.PartialSignatureMessages, expected []PostConsensusRoot) error {
	epoch := b.BeaconNetwork.EstimatedEpochAtSlot(b.State.StartingDuty.DutySlot())

	type expectedSigningRoot struct {
		root     [32]byte
		optional bool
	}
	signingRoots := make([]expectedSigningRoot, 0, len(expected))
	for _, e := range expected {
		d, err := runner.GetBeaconNode().DomainData(epoch, e.Domain)
		if err != nil {
			return errors.Wrap(err, "could not get post-consensus root domain")
		}
		r, err := types.ComputeETHSigningRoot(e.Root, d)
		if err != nil {
			return errors.Wrap(err, "could not compute ETH signing root")
		}
		signingRoots = append(signingRoots, expectedSigningRoot{root: r, optional: e.Optional})
	}

	// Match each message to an expected root, covering each at most once: an unexpected root, or a second
	// entry for an already-covered one, rejects the packet.
	covered := make(map[[32]byte]struct{})
	for _, msg := range psigMsgs.Messages {
		matched := false
		for _, sr := range signingRoots {
			if sr.root != msg.SigningRoot {
				continue
			}
			if _, dup := covered[sr.root]; dup {
				return types.NewError(types.WrongRootsCountErrorCode, "duplicate expected signing root")
			}
			covered[sr.root] = struct{}{}
			matched = true
			break
		}
		if !matched {
			return types.NewError(types.WrongSigningRootErrorCode, "unexpected signing root")
		}
	}

	// Every required root must be covered; optional roots (the §6 envelope) may be absent.
	for _, sr := range signingRoots {
		if _, ok := covered[sr.root]; !sr.optional && !ok {
			return types.NewError(types.WrongRootsCountErrorCode, "missing required signing root")
		}
	}
	return nil
}

func (b *BaseRunner) verifyExpectedRoot(runner Runner, psigMsgs *types.PartialSignatureMessages, expectedRootObjs []types.HashRoot, domain phase0.DomainType) error {
	if len(expectedRootObjs) != len(psigMsgs.Messages) {
		return types.NewError(types.WrongRootsCountErrorCode, "wrong expected roots count")
	}

	// convert expected roots to map and mark unique roots when verified
	sortedExpectedRoots, err := func(expectedRootObjs []types.HashRoot) ([][32]byte, error) {
		epoch := b.BeaconNetwork.EstimatedEpochAtSlot(b.State.StartingDuty.DutySlot())
		d, err := runner.GetBeaconNode().DomainData(epoch, domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not get pre consensus root domain")
		}

		ret := make([][32]byte, 0)
		for _, rootI := range expectedRootObjs {
			r, err := types.ComputeETHSigningRoot(rootI, d)
			if err != nil {
				return nil, errors.Wrap(err, "could not compute ETH signing root")
			}
			ret = append(ret, r)
		}

		sort.Slice(ret, func(i, j int) bool {
			return string(ret[i][:]) < string(ret[j][:])
		})
		return ret, nil
	}(expectedRootObjs)
	if err != nil {
		return err
	}

	sortedRoots := func(msgs types.PartialSignatureMessages) [][32]byte {
		ret := make([][32]byte, 0)
		for _, msg := range msgs.Messages {
			ret = append(ret, msg.SigningRoot)
		}

		sort.Slice(ret, func(i, j int) bool {
			return string(ret[i][:]) < string(ret[j][:])
		})
		return ret
	}(*psigMsgs)

	// verify roots
	for i, r := range sortedRoots {
		if !bytes.Equal(sortedExpectedRoots[i][:], r[:]) {
			return types.NewError(types.WrongSigningRootErrorCode, "wrong signing root")
		}
	}
	return nil
}
