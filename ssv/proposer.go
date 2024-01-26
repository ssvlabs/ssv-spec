package ssv

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type ProposerRunner struct {
	BaseRunner *BaseRunner
	// ProducesBlindedBlocks is true when the runner will only produce blinded blocks
	ProducesBlindedBlocks bool

	beacon   BeaconNode
	network  Network
	signer   types.KeyManager
	valCheck qbft.ProposedValueCheckF
}

func NewProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {

	runner := &ProposerRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleProposer,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:   beacon,
		network:  network,
		signer:   signer,
		valCheck: valCheck,
	}

	qbftController.WithCommitExtraLoadManagerF(runner.NewCommitExtraLoadProposerManager)

	return runner
}

func (r *ProposerRunner) StartNewDuty(duty *types.Duty) error {
	return r.BaseRunner.baseStartNewDuty(r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ProposerRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ProposerRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
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
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "could not reconstruct randao sig")
	}

	duty := r.GetState().StartingDuty

	var ver spec.DataVersion
	var obj ssz.Marshaler
	if r.ProducesBlindedBlocks {
		// get block data
		obj, ver, err = r.GetBeaconNode().GetBlindedBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
		if err != nil {
			return errors.Wrap(err, "failed to get Beacon block")
		}
	} else {
		// get block data
		obj, ver, err = r.GetBeaconNode().GetBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
		if err != nil {
			return errors.Wrap(err, "failed to get Beacon block")
		}
	}

	byts, err := obj.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal beacon block")
	}

	input := &types.ConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *ProposerRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	// specific duty sig
	var blkToSign ssz.HashRoot
	if r.decidedBlindedBlock() {
		_, blkToSign, err = decidedValue.GetBlindedBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get blinded block data")
		}
	} else {
		_, blkToSign, err = decidedValue.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block data")
		}
	}

	msg, err := r.BaseRunner.signBeaconObject(
		r,
		blkToSign,
		decidedValue.Duty.Slot,
		types.DomainProposer,
	)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     decidedValue.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	postSignedMsg, err := r.BaseRunner.signPostConsensusMsg(r, postConsensusMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign post consensus msg")
	}

	data, err := postSignedMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *ProposerRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus signature")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		if r.decidedBlindedBlock() {
			vBlindedBlk, _, err := r.GetState().DecidedValue.GetBlindedBlockData()
			if err != nil {
				return errors.Wrap(err, "could not get blinded block")
			}

			if err := r.GetBeaconNode().SubmitBlindedBeaconBlock(vBlindedBlk, specSig); err != nil {
				return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed blinded Beacon block")
			}
		} else {
			vBlk, _, err := r.GetState().DecidedValue.GetBlockData()
			if err != nil {
				return errors.Wrap(err, "could not get block")
			}

			if err := r.GetBeaconNode().SubmitBeaconBlock(vBlk, specSig); err != nil {
				return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block")
			}
		}
	}
	r.GetState().Finished = true
	return nil
}

// decidedBlindedBlock returns true if decided value has a blinded block, false if regular block
// WARNING!! should be called after decided only
func (r *ProposerRunner) decidedBlindedBlock() bool {
	_, _, err := r.BaseRunner.State.DecidedValue.GetBlindedBlockData()
	return err == nil
}

func (r *ProposerRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().StartingDuty.Slot)
	return []ssz.HashRoot{types.SSZUint64(epoch)}, types.DomainRandao, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	if r.decidedBlindedBlock() {
		_, data, err := r.GetState().DecidedValue.GetBlindedBlockData()
		if err != nil {
			return nil, phase0.DomainType{}, errors.Wrap(err, "could not get blinded block data")
		}
		return []ssz.HashRoot{data}, types.DomainProposer, nil
	}

	_, data, err := r.GetState().DecidedValue.GetBlockData()
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
func (r *ProposerRunner) executeDuty(duty *types.Duty) error {
	// sign partial randao
	epoch := r.GetBeaconNode().GetBeaconNetwork().EstimatedEpochAtSlot(duty.Slot)
	msg, err := r.BaseRunner.signBeaconObject(r, types.SSZUint64(epoch), duty.Slot, types.DomainRandao)
	if err != nil {
		return errors.Wrap(err, "could not sign randao")
	}
	msgs := types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	// sign msg
	signature, err := r.GetSigner().SignRoot(msgs, types.PartialSignatureType, r.GetShare().SharePubKey)
	if err != nil {
		return errors.Wrap(err, "could not sign randao msg")
	}
	signedPartialMsg := &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: signature,
		Signer:    r.GetShare().OperatorID,
	}

	// broadcast
	data, err := signedPartialMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode randao pre-consensus signature msg")
	}
	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}
	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
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
	return r.BaseRunner.Share
}

func (r *ProposerRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ProposerRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *ProposerRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *ProposerRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *ProposerRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *ProposerRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

func (r *ProposerRunner) SignBeaconObjectFromConsensusData(cd types.ConsensusData) (types.Signature, phase0.Root, error) {

	// specific duty sig
	var blkToSign ssz.HashRoot
	var err error

	// Try to get blinded block
	_, blkToSign, err = cd.GetBlindedBlockData()
	// If can't get blinded block, try to get block data
	if err != nil {
		_, blkToSign, err = cd.GetBlockData()
		if err != nil {
			return nil, phase0.Root{}, errors.Wrap(err, "could not get block data or blinded block data")
		}
	}

	sig, root, err := r.BaseRunner.GetBeaconObjectSignature(r, blkToSign, cd.Duty.Slot, types.DomainProposer)
	if err != nil {
		return nil, phase0.Root{}, errors.Wrap(err, "failed signing proposer beacon object")
	}

	return sig, root, nil
}

// Proposer manager for CommitExtraLoad
type CommitExtraLoadProposerManager struct {
	Signatures     map[types.OperatorID]types.Signature // Stores validated beacon object signatures
	ProposerRunner *ProposerRunner
	SigningRoot    phase0.Root // Stores signing root for comparison
}

func (r *ProposerRunner) NewCommitExtraLoadProposerManager() qbft.CommitExtraLoadManagerI {
	return &CommitExtraLoadProposerManager{
		ProposerRunner: r,
		Signatures:     make(map[uint64]types.Signature),
		SigningRoot:    phase0.Root{},
	}
}

// Returns a CommitExtraLoad with the validator's share signature over the beacon object
func (c *CommitExtraLoadProposerManager) Create(fullData []byte) (qbft.CommitExtraLoad, error) {
	// Get consensus data
	cd, err := c.GetConsensusData(fullData)
	if err != nil {
		return qbft.CommitExtraLoad{}, err
	}

	// Sign
	sig, root, err := c.ProposerRunner.SignBeaconObjectFromConsensusData(*cd)
	if err != nil {
		return qbft.CommitExtraLoad{}, errors.Wrap(err, "could not sign becon object")
	}

	// Store root for later comparison
	c.SigningRoot = root

	// Returns object
	return qbft.CommitExtraLoad{
		Signatures: []types.Signature{sig},
	}, nil
}

// Validates the CommitExtraLoad data inside a SignedMessage.
// - Validate fields
// - Checks the sender's validator share signature
func (c *CommitExtraLoadProposerManager) Validate(signedMessage *qbft.SignedMessage, fullData []byte) error {
	// Validate Signers length
	if len(signedMessage.Signers) == 0 {
		return errors.New("commit SignedMessage with empty signers")
	}
	if len(signedMessage.Signers) > 1 {
		return errors.New("commit SignedMessage with more than one signer")
	}
	// Validate Signatures length
	if len(signedMessage.Message.CommitExtraLoad.Signatures) == 0 {
		return errors.New("CommitExtraLoad with no signatures")
	}
	if len(signedMessage.Message.CommitExtraLoad.Signatures) > 1 {
		return errors.New("CommitExtraLoad with more than one signature")
	}

	signer := signedMessage.Signers[0]
	signature := signedMessage.Message.CommitExtraLoad.Signatures[0]
	root, err := c.GetSigningRootFromFullData(fullData)
	if err != nil {
		return errors.Wrap(err, "could not get ETH signing root from full data")
	}

	// Compare signing root if already instantiated
	if c.SigningRoot != [32]byte{} {
		if !bytes.Equal(root[:], c.SigningRoot[:]) {
			return errors.New("wrong signing root")
		}
	} else {
		c.SigningRoot = root
	}

	// Verify signature
	return c.ProposerRunner.BaseRunner.VerifyBeaconObjectPartialSignature(signer, signature, root)
}

// Process the CommitExtraLoad from a SignedMessage by storing the signature
func (c *CommitExtraLoadProposerManager) Process(signedMessage *qbft.SignedMessage) error {
	c.Signatures[signedMessage.Signers[0]] = signedMessage.Message.CommitExtraLoad.Signatures[0]
	return nil
}

// Returns the signing root of a decoded beacon object from FullData
func (c *CommitExtraLoadProposerManager) GetSigningRootFromFullData(fullData []byte) (phase0.Root, error) {
	cd, err := c.GetConsensusData(fullData)
	if err != nil {
		return phase0.Root{}, err
	}
	return c.GetSigningRoot(cd)
}

// Returns a ConsensusData decoded from FullData
func (c *CommitExtraLoadProposerManager) GetConsensusData(fullData []byte) (*types.ConsensusData, error) {
	cd := &types.ConsensusData{}
	err := cd.Decode(fullData)
	if err != nil {
		return nil, errors.Wrap(err, "could not get consensus data")
	}
	return cd, nil
}

// Returns the proposer beacon object taken from the ConsensusData
func (c *CommitExtraLoadProposerManager) GetBeaconObject(cd *types.ConsensusData) (ssz.HashRoot, error) {
	// specific duty sig
	var blkToSign ssz.HashRoot
	var err error

	// Try to get blinded block
	_, blkToSign, err = cd.GetBlindedBlockData()
	// If can't get blinded block, try to get block data
	if err != nil {
		_, blkToSign, err = cd.GetBlockData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get block data or blinded block data")
		}
	}
	return blkToSign, nil
}

// Computes the signing root of the beacon object taken from a consensus data
func (c *CommitExtraLoadProposerManager) GetSigningRoot(cd *types.ConsensusData) (phase0.Root, error) {

	obj, err := c.GetBeaconObject(cd)
	if err != nil {
		return phase0.Root{}, errors.Wrap(err, "could not get beacon object")
	}

	epoch := c.ProposerRunner.GetBaseRunner().BeaconNetwork.EstimatedEpochAtSlot(cd.Duty.Slot)
	domain, err := c.ProposerRunner.GetBeaconNode().DomainData(epoch, types.DomainProposer)
	if err != nil {
		return phase0.Root{}, errors.Wrap(err, "could not get beacon domain")
	}

	return c.ProposerRunner.BaseRunner.GetBeaconSigningRoot(obj, domain)
}
