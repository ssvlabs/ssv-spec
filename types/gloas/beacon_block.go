package gloas

import (
	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// PayloadAttestation is the aggregated PTC attestation the proposer includes in the block body for the
// previous slot's payload (consensus-specs gloas) — distinct from the single-member
// PayloadAttestationMessage SSV signs in §3. AggregationBits is a Bitvector[PTC_SIZE], PTC_SIZE = 512.
type PayloadAttestation struct {
	AggregationBits bitfield.Bitvector512   `dynssz-size:"PTC_SIZE/8" ssz-index:"0" ssz-size:"64"`
	Data            *PayloadAttestationData `ssz-index:"1"`
	Signature       phase0.BLSSignature     `ssz-index:"2"            ssz-size:"96"`
}

// BeaconBlockBody is the Gloas (ePBS) block body. Versus Electra it drops the inline execution payload,
// execution requests, and blob KZG commitments (the payload and blobs now ship in the §6 envelope) and
// adds SignedExecutionPayloadBid (the payload commitment), PayloadAttestations (the previous slot's PTC
// aggregate), and ParentExecutionRequests. Field order/tags match the pinned spec / go-eth2-client
// PR #280; everything else reuses the existing fork types.
type BeaconBlockBody struct {
	RANDAOReveal              phase0.BLSSignature                   `ssz-index:"0"  ssz-size:"96"`
	ETH1Data                  *phase0.ETH1Data                      `ssz-index:"1"`
	Graffiti                  [32]byte                              `ssz-index:"2"  ssz-size:"32"`
	ProposerSlashings         []*phase0.ProposerSlashing            `ssz-index:"3"  ssz-type:"progressive-list"`
	AttesterSlashings         []*eth2gloas.AttesterSlashing         `ssz-index:"4"  ssz-type:"progressive-list"`
	Attestations              []*eth2gloas.Attestation              `ssz-index:"5"  ssz-type:"progressive-list"`
	Deposits                  []*phase0.Deposit                     `ssz-index:"6"  ssz-type:"progressive-list"`
	VoluntaryExits            []*phase0.SignedVoluntaryExit         `ssz-index:"7"  ssz-type:"progressive-list"`
	SyncAggregate             *altair.SyncAggregate                 `ssz-index:"8"`
	BLSToExecutionChanges     []*capella.SignedBLSToExecutionChange `ssz-index:"9"  ssz-type:"progressive-list"`
	SignedExecutionPayloadBid *SignedExecutionPayloadBid            `ssz-index:"10"`
	PayloadAttestations       []*PayloadAttestation                 `ssz-index:"11" ssz-type:"progressive-list"`
	// Gloas execution requests — the EIP-8282 five-list variant, not electra's three (see execution_requests.go).
	ParentExecutionRequests *ExecutionRequests `ssz-index:"12"`
}

// BeaconBlock is the Gloas (ePBS) beacon block.
type BeaconBlock struct {
	Slot          phase0.Slot
	ProposerIndex phase0.ValidatorIndex
	ParentRoot    phase0.Root `ssz-size:"32"`
	StateRoot     phase0.Root `ssz-size:"32"`
	Body          *BeaconBlockBody
}

// SignedBeaconBlock wraps a Gloas BeaconBlock with the proposer's signature.
type SignedBeaconBlock struct {
	Message   *BeaconBlock
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// Encode/Decode are the convenience wrappers the proposer runner uses to marshal the block into the
// QBFT DataSSZ and to publish the signed block.
func (b *BeaconBlock) Encode() ([]byte, error)  { return b.MarshalSSZ() }
func (b *BeaconBlock) Decode(data []byte) error { return b.UnmarshalSSZ(data) }

func (b *SignedBeaconBlock) Encode() ([]byte, error)  { return b.MarshalSSZ() }
func (b *SignedBeaconBlock) Decode(data []byte) error { return b.UnmarshalSSZ(data) }

// DecodeBeaconBlock unmarshals a Gloas BeaconBlock from QBFT consensus DataSSZ. The Gloas proposer
// path uses it in place of ProposerConsensusData.GetBlockData (which has no Gloas version); the
// returned block doubles as the types.HashRoot the proposer signs.
func DecodeBeaconBlock(dataSSZ []byte) (*BeaconBlock, error) {
	b := &BeaconBlock{}
	if err := b.UnmarshalSSZ(dataSSZ); err != nil {
		return nil, err
	}
	return b, nil
}
