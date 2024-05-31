package qbft

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/types"
)

type signing interface {
	// GetSigner returns an operator signer instance
	GetOperatorSigner() types.OperatorSigner
}

type IConfig interface {
	signing
	// GetValueCheckF returns value check function
	GetValueCheckF() ProposedValueCheckF
	// GetProposerF returns func used to calculate proposer
	GetProposerF() ProposerF
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetTimer returns round timer
	GetTimer() Timer
	// GetSignatureVerifier returns the signature verifier for operator signatures
	GetSignatureVerifier() types.SignatureVerifier
	// GetCutOffRound returns the round that stops the instance
	GetCutOffRound() Round
}

type Config struct {
	OperatorSigner    types.OperatorSigner
	SigningPK         []byte
	Domain            types.DomainType
	ValueCheckF       ProposedValueCheckF
	ProposerF         ProposerF
	Network           Network
	Timer             Timer
	SignatureVerifier types.SignatureVerifier
	CutOffRound       Round
}

// GetSigner returns a Signer instance
func (c *Config) GetOperatorSigner() types.OperatorSigner {
	return c.OperatorSigner
}

// GetSigningPubKey returns the public key used to sign all QBFT messages
func (c *Config) GetSigningPubKey() []byte {
	return c.SigningPK
}

// GetSignatureDomainType returns the Domain type used for signatures
func (c *Config) GetSignatureDomainType() types.DomainType {
	return c.Domain
}

// GetValueCheckF returns value check instance
func (c *Config) GetValueCheckF() ProposedValueCheckF {
	return c.ValueCheckF
}

// GetProposerF returns func used to calculate proposer
func (c *Config) GetProposerF() ProposerF {
	return c.ProposerF
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() Network {
	return c.Network
}

// GetTimer returns round timer
func (c *Config) GetTimer() Timer {
	return c.Timer
}

func (c *Config) GetCutOffRound() Round {
	return c.CutOffRound
}

// GetSignatureVerifier returns the verifier for operator's signatures
func (c *Config) GetSignatureVerifier() types.SignatureVerifier {
	return c.SignatureVerifier
}

type State struct {
	Share                           *types.Share
	ID                              []byte // instance Identifier
	Round                           Round
	Height                          Height
	LastPreparedRound               Round
	LastPreparedValue               []byte
	ProposalAcceptedForCurrentRound *types.SignedSSVMessage
	Decided                         bool
	DecidedValue                    []byte

	ProposeContainer     *MsgContainer
	PrepareContainer     *MsgContainer
	CommitContainer      *MsgContainer
	RoundChangeContainer *MsgContainer
}

// GetRoot returns the state's deterministic root
func (s *State) GetRoot() ([32]byte, error) {
	marshaledRoot, err := s.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

// Encode returns a msg encoded bytes or error
func (s *State) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode returns error if decoding failed
func (s *State) Decode(data []byte) error {
	return json.Unmarshal(data, &s)
}
