package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type ProposerRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}

	return &ProposerRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RoleProposer,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}, nil
}

func (r *ProposerRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ProposerRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ProposerRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing randao message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	// randao is relevant only for block proposals, no need to check type
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
			r.GetShare().ValidatorIndex)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}

	duty := r.GetState().StartingDuty.(*types.ValidatorDuty)

	// get block data
	vBlk, err := r.GetBeaconNode().GetBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to get Beacon block")
	}

	// Proposer consensus always agrees on the blinded block form. If the beacon
	// node returned a full block, derive the blinded form locally before QBFT.
	blindedVBlk, blindedObj, err := ensureBlindedProposal(vBlk)
	if err != nil {
		return errors.Wrap(err, "could not blind beacon block")
	}

	byts, err := blindedObj.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal blinded beacon block")
	}

	input := &types.ValidatorConsensusData{
		Duty:    *duty,
		Version: blindedVBlk.Version,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input.Duty.DutySlot(), input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *ProposerRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg, &types.ValidatorConsensusData{})
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	// specific duty sig
	var blkToSign ssz.HashRoot

	cd := decidedValue.(*types.ValidatorConsensusData)
	_, blkToSign, err = cd.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "could not get block data")
	}

	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), blkToSign,
		cd.Duty.Slot,
		types.DomainProposer)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     cd.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := postConsensusMsg.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *ProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root,
					r.GetShare().Committee, r.GetShare().ValidatorIndex)
			}
			return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		validatorConsensusData := &types.ValidatorConsensusData{}
		err = validatorConsensusData.Decode(r.GetState().DecidedValue)
		if err != nil {
			return errors.Wrap(err, "could not create consensus data")
		}
		vBlk, _, err := validatorConsensusData.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block")
		}

		if err := r.GetBeaconNode().SubmitBeaconBlock(vBlk, specSig); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *ProposerRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().StartingDuty.DutySlot())
	return []ssz.HashRoot{types.SSZUint64(epoch)}, types.DomainRandao, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	validatorConsensusData := &types.ValidatorConsensusData{}
	err := validatorConsensusData.Decode(r.GetState().DecidedValue)
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not create consensus data")
	}

	_, data, err := validatorConsensusData.GetBlockData()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get block data")
	}
	return []ssz.HashRoot{data}, types.DomainProposer, nil
}

// executeDuty steps:
// 1) sign a partial randao sig and wait for 2f+1 partial sigs from peers
// 2) reconstruct randao and send GetBeaconBlock to BN
// 3) start consensus on duty + block data
// 4) Once consensus decides, sign partial block and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid block sig to the BN
func (r *ProposerRunner) executeDuty(duty types.Duty) error {
	// sign partial randao
	epoch := r.GetBeaconNode().GetBeaconNetwork().EstimatedEpochAtSlot(duty.DutySlot())
	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), types.SSZUint64(epoch), duty.DutySlot(),
		types.DomainRandao)
	if err != nil {
		return errors.Wrap(err, "could not sign randao")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := msgs.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
}

// ensureBlindedProposal returns a blinded proposal and the concrete SSZ-marshaler
// that must be encoded into proposer consensus data. If the input is already
// blinded, it is returned unchanged. For full Deneb/Electra/Fulu proposals this
// intentionally keeps only the blinded block contents needed for consensus; blob
// sidecars and proofs from the outer BlockContents wrapper are not propagated.
func ensureBlindedProposal(p *api.VersionedProposal) (*api.VersionedProposal, ssz.Marshaler, error) {
	if p == nil {
		return nil, nil, fmt.Errorf("nil proposal")
	}

	if p.Blinded {
		switch p.Version {
		case spec.DataVersionCapella:
			if p.CapellaBlinded == nil {
				return nil, nil, fmt.Errorf("capella blinded block is nil")
			}
			return p, p.CapellaBlinded, nil
		case spec.DataVersionDeneb:
			if p.DenebBlinded == nil {
				return nil, nil, fmt.Errorf("deneb blinded block is nil")
			}
			return p, p.DenebBlinded, nil
		case spec.DataVersionElectra:
			if p.ElectraBlinded == nil {
				return nil, nil, fmt.Errorf("electra blinded block is nil")
			}
			return p, p.ElectraBlinded, nil
		case spec.DataVersionFulu:
			if p.FuluBlinded == nil {
				return nil, nil, fmt.Errorf("fulu blinded block is nil")
			}
			return p, p.FuluBlinded, nil
		default:
			return nil, nil, fmt.Errorf("unsupported proposal version %d", p.Version)
		}
	}

	switch p.Version {
	case spec.DataVersionCapella:
		if p.Capella == nil || p.Capella.Body == nil || p.Capella.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("capella block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Capella.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute capella transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Capella.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute capella withdrawals root: %w", err)
		}

		blinded := &apiv1capella.BlindedBeaconBlock{
			Slot:          p.Capella.Slot,
			ProposerIndex: p.Capella.ProposerIndex,
			ParentRoot:    p.Capella.ParentRoot,
			StateRoot:     p.Capella.StateRoot,
			Body: &apiv1capella.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Capella.Body.RANDAOReveal,
				ETH1Data:          p.Capella.Body.ETH1Data,
				Graffiti:          p.Capella.Body.Graffiti,
				ProposerSlashings: p.Capella.Body.ProposerSlashings,
				AttesterSlashings: p.Capella.Body.AttesterSlashings,
				Attestations:      p.Capella.Body.Attestations,
				Deposits:          p.Capella.Body.Deposits,
				VoluntaryExits:    p.Capella.Body.VoluntaryExits,
				SyncAggregate:     p.Capella.Body.SyncAggregate,
				ExecutionPayloadHeader: &capella.ExecutionPayloadHeader{
					ParentHash:       p.Capella.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Capella.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Capella.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Capella.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Capella.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Capella.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Capella.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Capella.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Capella.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Capella.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Capella.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Capella.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Capella.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
				},
				BLSToExecutionChanges: p.Capella.Body.BLSToExecutionChanges,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, CapellaBlinded: blinded}, blinded, nil
	case spec.DataVersionDeneb:
		if p.Deneb == nil || p.Deneb.Block == nil || p.Deneb.Block.Body == nil || p.Deneb.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("deneb block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Deneb.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute deneb transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Deneb.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute deneb withdrawals root: %w", err)
		}

		blinded := &apiv1deneb.BlindedBeaconBlock{
			Slot:          p.Deneb.Block.Slot,
			ProposerIndex: p.Deneb.Block.ProposerIndex,
			ParentRoot:    p.Deneb.Block.ParentRoot,
			StateRoot:     p.Deneb.Block.StateRoot,
			Body: &apiv1deneb.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Deneb.Block.Body.RANDAOReveal,
				ETH1Data:          p.Deneb.Block.Body.ETH1Data,
				Graffiti:          p.Deneb.Block.Body.Graffiti,
				ProposerSlashings: p.Deneb.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Deneb.Block.Body.AttesterSlashings,
				Attestations:      p.Deneb.Block.Body.Attestations,
				Deposits:          p.Deneb.Block.Body.Deposits,
				VoluntaryExits:    p.Deneb.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Deneb.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Deneb.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Deneb.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Deneb.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Deneb.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Deneb.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Deneb.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Deneb.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Deneb.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Deneb.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Deneb.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Deneb.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Deneb.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Deneb.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Deneb.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Deneb.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Deneb.Block.Body.BLSToExecutionChanges,
				// Blob commitments live on the beacon block body, not on the outer BlockContents wrapper.
				BlobKZGCommitments: p.Deneb.Block.Body.BlobKZGCommitments,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, DenebBlinded: blinded}, blinded, nil
	case spec.DataVersionElectra:
		if p.Electra == nil || p.Electra.Block == nil || p.Electra.Block.Body == nil || p.Electra.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("electra block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Electra.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute electra transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Electra.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute electra withdrawals root: %w", err)
		}

		blinded := &apiv1electra.BlindedBeaconBlock{
			Slot:          p.Electra.Block.Slot,
			ProposerIndex: p.Electra.Block.ProposerIndex,
			ParentRoot:    p.Electra.Block.ParentRoot,
			StateRoot:     p.Electra.Block.StateRoot,
			Body: &apiv1electra.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Electra.Block.Body.RANDAOReveal,
				ETH1Data:          p.Electra.Block.Body.ETH1Data,
				Graffiti:          p.Electra.Block.Body.Graffiti,
				ProposerSlashings: p.Electra.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Electra.Block.Body.AttesterSlashings,
				Attestations:      p.Electra.Block.Body.Attestations,
				Deposits:          p.Electra.Block.Body.Deposits,
				VoluntaryExits:    p.Electra.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Electra.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Electra.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Electra.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Electra.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Electra.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Electra.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Electra.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Electra.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Electra.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Electra.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Electra.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Electra.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Electra.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Electra.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Electra.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Electra.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Electra.Block.Body.BLSToExecutionChanges,
				BlobKZGCommitments:    p.Electra.Block.Body.BlobKZGCommitments,
				ExecutionRequests:     p.Electra.Block.Body.ExecutionRequests,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, ElectraBlinded: blinded}, blinded, nil
	case spec.DataVersionFulu:
		if p.Fulu == nil || p.Fulu.Block == nil || p.Fulu.Block.Body == nil || p.Fulu.Block.Body.ExecutionPayload == nil {
			return nil, nil, fmt.Errorf("fulu block or payload is nil")
		}

		txRoot, err := types.SSZTransactions(p.Fulu.Block.Body.ExecutionPayload.Transactions).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute fulu transactions root: %w", err)
		}
		withdrawalsRoot, err := types.SSZWithdrawals(p.Fulu.Block.Body.ExecutionPayload.Withdrawals).HashTreeRoot()
		if err != nil {
			return nil, nil, fmt.Errorf("could not compute fulu withdrawals root: %w", err)
		}

		// Fulu currently reuses Electra's blinded block structure.
		blinded := &apiv1electra.BlindedBeaconBlock{
			Slot:          p.Fulu.Block.Slot,
			ProposerIndex: p.Fulu.Block.ProposerIndex,
			ParentRoot:    p.Fulu.Block.ParentRoot,
			StateRoot:     p.Fulu.Block.StateRoot,
			Body: &apiv1electra.BlindedBeaconBlockBody{
				RANDAOReveal:      p.Fulu.Block.Body.RANDAOReveal,
				ETH1Data:          p.Fulu.Block.Body.ETH1Data,
				Graffiti:          p.Fulu.Block.Body.Graffiti,
				ProposerSlashings: p.Fulu.Block.Body.ProposerSlashings,
				AttesterSlashings: p.Fulu.Block.Body.AttesterSlashings,
				Attestations:      p.Fulu.Block.Body.Attestations,
				Deposits:          p.Fulu.Block.Body.Deposits,
				VoluntaryExits:    p.Fulu.Block.Body.VoluntaryExits,
				SyncAggregate:     p.Fulu.Block.Body.SyncAggregate,
				ExecutionPayloadHeader: &deneb.ExecutionPayloadHeader{
					ParentHash:       p.Fulu.Block.Body.ExecutionPayload.ParentHash,
					FeeRecipient:     p.Fulu.Block.Body.ExecutionPayload.FeeRecipient,
					StateRoot:        p.Fulu.Block.Body.ExecutionPayload.StateRoot,
					ReceiptsRoot:     p.Fulu.Block.Body.ExecutionPayload.ReceiptsRoot,
					LogsBloom:        p.Fulu.Block.Body.ExecutionPayload.LogsBloom,
					PrevRandao:       p.Fulu.Block.Body.ExecutionPayload.PrevRandao,
					BlockNumber:      p.Fulu.Block.Body.ExecutionPayload.BlockNumber,
					GasLimit:         p.Fulu.Block.Body.ExecutionPayload.GasLimit,
					GasUsed:          p.Fulu.Block.Body.ExecutionPayload.GasUsed,
					Timestamp:        p.Fulu.Block.Body.ExecutionPayload.Timestamp,
					ExtraData:        p.Fulu.Block.Body.ExecutionPayload.ExtraData,
					BaseFeePerGas:    p.Fulu.Block.Body.ExecutionPayload.BaseFeePerGas,
					BlockHash:        p.Fulu.Block.Body.ExecutionPayload.BlockHash,
					TransactionsRoot: txRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      p.Fulu.Block.Body.ExecutionPayload.BlobGasUsed,
					ExcessBlobGas:    p.Fulu.Block.Body.ExecutionPayload.ExcessBlobGas,
				},
				BLSToExecutionChanges: p.Fulu.Block.Body.BLSToExecutionChanges,
				BlobKZGCommitments:    p.Fulu.Block.Body.BlobKZGCommitments,
				ExecutionRequests:     p.Fulu.Block.Body.ExecutionRequests,
			},
		}

		return &api.VersionedProposal{Version: p.Version, Blinded: true, FuluBlinded: blinded}, blinded, nil
	default:
		return nil, nil, fmt.Errorf("unsupported proposal version %d", p.Version)
	}
}

func (r *ProposerRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ProposerRunner) GetNetwork() Network {
	return r.network
}

func (r *ProposerRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *ProposerRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *ProposerRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ProposerRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *ProposerRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *ProposerRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
