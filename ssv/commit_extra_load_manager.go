package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// Manager for CommitExtraLoad
type CommitExtraLoadManager struct {
	PartialSigContainer *PartialSigContainer // Stores validated beacon object signatures
	SigningRoot         []phase0.Root        // Stores signing root for comparison

	// Runner-specific
	BaseRunner *BaseRunner
	beaconNode BeaconNode
	signer     types.KeyManager
	Role       types.BeaconRole
	DomainType phase0.DomainType
}

func NewCommitExtraLoadManagerF(
	baseRunner *BaseRunner,
	role types.BeaconRole,
	beaconNode BeaconNode,
	signer types.KeyManager,
	domainType phase0.DomainType) qbft.CommitExtraLoadManagerF {

	return func() qbft.CommitExtraLoadManagerI {
		return &CommitExtraLoadManager{
			PartialSigContainer: NewPartialSigContainer(baseRunner.Share.Quorum),
			SigningRoot:         make([]phase0.Root, 0),

			// Runner-specific
			BaseRunner: baseRunner,
			beaconNode: beaconNode,
			Role:       role,
			signer:     signer,
			DomainType: domainType,
		}
	}
}

// Returns a CommitExtraLoad with the validator's share signatures over the beacon objects derived from the consensus data
func (c *CommitExtraLoadManager) Create(fullData []byte) (qbft.CommitExtraLoad, error) {
	// Get consensus data
	cd, err := c.GetConsensusData(fullData)
	if err != nil {
		return qbft.CommitExtraLoad{}, err
	}

	// Sign
	sigs, roots, err := c.SignBeaconObjectFromConsensusData(cd)
	if err != nil {
		return qbft.CommitExtraLoad{}, errors.Wrap(err, "could not sign beacon object")
	}

	// Store root for later comparison
	c.SigningRoot = roots

	// Return object
	postConsensusSignatures := make([]*qbft.PostConsensusSignature, 0)
	for idx, sig := range sigs {
		postConsensusSignatures = append(postConsensusSignatures, &qbft.PostConsensusSignature{
			Signature:   sig,
			SigningRoot: roots[idx],
		})
	}
	return qbft.CommitExtraLoad{
		PostConsensusSignatures: postConsensusSignatures,
	}, nil
}

// Validates the CommitExtraLoad data inside a SignedMessage:
// - Validate signer field
// - Validate length of Signatures field
// - Compare roots with expected roots from FullData
// - Checks each validator share signature
func (c *CommitExtraLoadManager) Validate(signedMessage *qbft.SignedMessage, fullData []byte) error {
	// Validate Signers length
	if len(signedMessage.Signers) == 0 {
		return errors.New("commit SignedMessage with empty signers")
	}
	if len(signedMessage.Signers) > 1 {
		return errors.New("commit SignedMessage with more than one signer")
	}
	signer := signedMessage.Signers[0]

	// Validate Signatures length
	if len(signedMessage.Message.CommitExtraLoad.PostConsensusSignatures) == 0 {
		return errors.New("CommitExtraLoad with no signatures")
	}
	if c.Role != types.BNRoleSyncCommitteeContribution {
		if len(signedMessage.Message.CommitExtraLoad.PostConsensusSignatures) > 1 {
			return errors.New("CommitExtraLoad with more than one signature")
		}
	}

	// Fill SigningRoots for comparison if empty
	c.FillSigningRootsIfEmpty(fullData)

	// Compare roots
	msgRoots := make([][32]byte, 0)
	for _, postConsensusSignature := range signedMessage.Message.CommitExtraLoad.PostConsensusSignatures {
		msgRoots = append(msgRoots, postConsensusSignature.SigningRoot)
	}
	expectedRoots := make([][32]byte, 0)
	for _, root := range c.SigningRoot {
		expectedRoots = append(expectedRoots, root)
	}
	c.BaseRunner.compareRoots(msgRoots, expectedRoots)

	// Verify each signature
	for _, postConsensusSignature := range signedMessage.Message.CommitExtraLoad.PostConsensusSignatures {
		err := c.BaseRunner.VerifyBeaconObjectPartialSignature(signer, postConsensusSignature.Signature, postConsensusSignature.SigningRoot)
		if err != nil {
			return err
		}
	}
	return nil
}

// Process the CommitExtraLoad from a SignedMessage by storing the signatures
func (c *CommitExtraLoadManager) Process(signedMessage *qbft.SignedMessage) error {

	// for each SigningRoot, add the respective signature
	for _, postConsensusSignature := range signedMessage.Message.CommitExtraLoad.PostConsensusSignatures {
		c.PartialSigContainer.AddSignatureForSignatureRootAndSigner(postConsensusSignature.Signature, postConsensusSignature.SigningRoot, signedMessage.Signers[0])
	}
	return nil
}

// If SigningRoot is empty, fills it with beacon roots derived from the FullData argument
func (c *CommitExtraLoadManager) FillSigningRootsIfEmpty(fullData []byte) error {
	if len(c.SigningRoot) > 0 {
		return nil
	}

	roots, err := c.GetSigningRootFromFullData(fullData)
	if err != nil {
		return err
	}

	c.SigningRoot = roots
	return nil
}

// Returns the signing root of a decoded beacon object from FullData
func (c *CommitExtraLoadManager) GetSigningRootFromFullData(fullData []byte) ([]phase0.Root, error) {
	cd, err := c.GetConsensusData(fullData)
	if err != nil {
		return []phase0.Root{}, err
	}
	return c.GetSigningRootsFromConsensusData(cd)
}

// Returns a ConsensusData decoded from FullData
func (c *CommitExtraLoadManager) GetConsensusData(fullData []byte) (*types.ConsensusData, error) {
	cd := &types.ConsensusData{}
	err := cd.Decode(fullData)
	if err != nil {
		return nil, errors.Wrap(err, "could not get consensus data")
	}
	return cd, nil
}

// Computes the signing root of the beacon object taken from a consensus data
func (c *CommitExtraLoadManager) GetSigningRootsFromConsensusData(cd *types.ConsensusData) ([]phase0.Root, error) {

	// Get beacon objects
	objs, err := c.GetBeaconObjects(cd)
	if err != nil {
		return []phase0.Root{}, errors.Wrap(err, "could not get beacon object")
	}

	// Get domain to create roots
	epoch := c.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(cd.Duty.Slot)
	domain, err := c.beaconNode.DomainData(epoch, types.DomainProposer)
	if err != nil {
		return []phase0.Root{}, errors.Wrap(err, "could not get beacon domain")
	}

	// Get root from each object
	ret := make([]phase0.Root, 0)
	for _, obj := range objs {
		signingRoot, err := c.BaseRunner.GetBeaconSigningRoot(obj, domain)
		if err != nil {
			return []phase0.Root{}, errors.Wrap(err, "could not get beacon signing root")
		}
		ret = append(ret, signingRoot)
	}
	return ret, nil

}

// Computes the signing root and signature of the beacon object taken from a consensus data
func (c *CommitExtraLoadManager) SignBeaconObjectFromConsensusData(cd *types.ConsensusData) ([]types.Signature, []phase0.Root, error) {

	// Get beacon objects
	objs, err := c.GetBeaconObjects(cd)
	if err != nil {
		return nil, []phase0.Root{}, errors.Wrap(err, "could not get beacon objects")
	}

	// Get domain to create signatures
	epoch := c.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(cd.Duty.Slot)
	domain, err := c.beaconNode.DomainData(epoch, c.DomainType)
	if err != nil {
		return nil, []phase0.Root{}, errors.Wrap(err, "could not get beacon domain")
	}

	// Compute the root and signature for each beacon object
	sigs := make([]types.Signature, 0)
	roots := make([]phase0.Root, 0)
	for _, obj := range objs {
		sig, root, err := c.signer.SignBeaconObject(obj, domain, c.BaseRunner.Share.SharePubKey, c.DomainType)
		if err != nil {
			return nil, []phase0.Root{}, errors.Wrap(err, "could not sign beacon object")
		}
		sigs = append(sigs, sig)
		roots = append(roots, root)
	}
	return sigs, roots, err
}

// Returns the beacon object taken from the ConsensusData
func (c *CommitExtraLoadManager) GetBeaconObjects(cd *types.ConsensusData) ([]ssz.HashRoot, error) {

	switch c.Role {

	case types.BNRoleProposer:

		var blkToSign ssz.HashRoot
		var err error

		// Try to get blinded block
		_, blkToSign, err = cd.GetBlindedBlockData()
		// If can't get blinded block, try to get block data
		if err != nil {
			_, blkToSign, err = cd.GetBlockData()
			if err != nil {
				return []ssz.HashRoot{}, errors.Wrap(err, "could not get block data or blinded block data")
			}
		}
		return []ssz.HashRoot{blkToSign}, nil

	case types.BNRoleAttester:

		attestationData, err := cd.GetAttestationData()
		if err != nil {
			return []ssz.HashRoot{}, errors.Wrap(err, "could not get attestation data")
		}
		return []ssz.HashRoot{attestationData}, nil

	case types.BNRoleAggregator:

		aggregateAndProof, err := cd.GetAggregateAndProof()
		if err != nil {
			return []ssz.HashRoot{}, errors.Wrap(err, "could not get aggregate and proof")
		}
		return []ssz.HashRoot{aggregateAndProof}, nil

	case types.BNRoleSyncCommittee:

		root, err := cd.GetSyncCommitteeBlockRoot()
		if err != nil {
			return []ssz.HashRoot{}, errors.Wrap(err, "could not get sync committee block root")
		}
		return []ssz.HashRoot{types.SSZBytes(root[:])}, nil

	case types.BNRoleSyncCommitteeContribution:

		contributions, err := cd.GetSyncCommitteeContributions()
		if err != nil {
			return []ssz.HashRoot{}, errors.Wrap(err, "could not get contributions")
		}

		// specific duty sig
		ret := make([]ssz.HashRoot, 0)
		for _, contrib := range contributions {

			contribAndProof := &altair.ContributionAndProof{
				AggregatorIndex: cd.Duty.ValidatorIndex,
				Contribution:    &contrib.Contribution,
				SelectionProof:  contrib.SelectionProofSig,
			}
			if err != nil {
				return []ssz.HashRoot{}, errors.Wrap(err, "could not generate contribution and proof")
			}
			ret = append(ret, contribAndProof)
		}
		return ret, nil
	}
	return []ssz.HashRoot{}, errors.New("unexpected type")
}
