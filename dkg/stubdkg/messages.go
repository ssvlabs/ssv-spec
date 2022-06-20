package stubdkg

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Sha256Hash [32]byte
type BlsSignature [96]byte
type BlsPublicKey [48]byte
type BlsScalar [32]byte
type BigInt []byte
type Round uint16

const (
	KG_R1 = 1
	KG_R2 = 2
	KG_R3 = 3
	KG_R4 = 4
)

type KeygenProtocolMsg struct {
	RoundNumber Round
	Sender      uint16
	Receiver    uint16
	Data        []byte
}

// Encode returns a msg encoded bytes or error
func (d *KeygenProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (msg *KeygenProtocolMsg) GetRound1Data() (*KeygenRound1Data, error) {
	if msg.RoundNumber != 1 {
		return nil, errors.New("the message doesn't contain a round 1 data")
	}
	ret := &KeygenRound1Data{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode round 1 data from message")
	}
	return ret, nil
}

func (msg *KeygenProtocolMsg) GetRound2Data() (*KeygenRound2Data, error) {
	if msg.RoundNumber != 2 {
		return nil, errors.New("the message doesn't contain a round 2 data")
	}
	ret := &KeygenRound2Data{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode round 2 data from message")
	}
	return ret, nil
}

func (msg *KeygenProtocolMsg) GetRound3Data() (*KeygenRound3Data, error) {
	if msg.RoundNumber != 3 {
		return nil, errors.New("the message doesn't contain a round 3 data")
	}
	ret := &KeygenRound3Data{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode round 3 data from message")
	}
	return ret, nil
}

func (msg *KeygenProtocolMsg) GetRound4Data() (*KeygenRound4Data, error) {
	if msg.RoundNumber != 4 {
		return nil, errors.New("the message doesn't contain a round 3 data")
	}
	ret := &KeygenRound4Data{}
	if err := ret.Decode(msg.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode round4 data from message")
	}
	return ret, nil
}

func (msg *KeygenProtocolMsg) SetRound1Data(data *KeygenRound1Data) error {
	bytes, err := data.Encode()
	if err != nil {
		return errors.New("unable to encode data")
	}
	msg.RoundNumber = 1
	msg.Data = bytes
	return nil
}

func (msg *KeygenProtocolMsg) SetRound2Data(data *KeygenRound2Data) error {
	bytes, err := data.Encode()
	if err != nil {
		return errors.New("unable to encode data")
	}
	msg.RoundNumber = 2
	msg.Data = bytes
	return nil
}

func (msg *KeygenProtocolMsg) SetRound3Data(data *KeygenRound3Data) error {
	bytes, err := data.Encode()
	if err != nil {
		return errors.New("unable to encode data")
	}
	msg.RoundNumber = 3
	msg.Data = bytes
	return nil
}

func (msg *KeygenProtocolMsg) SetRound4Data(data *KeygenRound4Data) error {
	bytes, err := data.Encode()
	if err != nil {
		return errors.New("unable to encode data")
	}
	msg.RoundNumber = 4
	msg.Data = bytes
	return nil
}

// KeygenRound1Data contains the commitment data
type KeygenRound1Data struct {
	Commitment Sha256Hash `json:"com"`
}

// Encode returns a msg encoded bytes or error
func (d *KeygenRound1Data) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenRound1Data) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// KeygenRound2Data contains the decommitment data
type KeygenRound2Data struct {
	BlindFactor BigInt       `json:"blind_factor"`
	YI          BlsPublicKey `json:"y_i"`
}

// Encode returns a msg encoded bytes or error
func (d *KeygenRound2Data) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenRound2Data) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// KeygenRound3Data distributes shares for u_i
type KeygenRound3Data struct {
	Parameters struct {
		Threshold  int `json:"threshold"`
		ShareCount int `json:"share_count"`
	} `json:"parameters"` // The vss parameters
	Commitments []BlsPublicKey `json:"commitments"`
	ShareIJ     BlsScalar      `json:"share_i_j"`
}

// Encode returns a msg encoded bytes or error
func (d *KeygenRound3Data) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenRound3Data) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// KeygenRound4Data proves knowledge of x_i
type KeygenRound4Data struct {
	Pk                BlsPublicKey `json:"pk"`
	PkTRandCommitment BlsPublicKey `json:"pk_t_rand_commitment"`
	ChallengeResponse BlsScalar    `json:"challenge_response"`
}

// Encode returns a msg encoded bytes or error
func (d *KeygenRound4Data) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenRound4Data) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

type LocalKeyShare struct {
	Index           uint16         `json:"i"`
	Threshold       uint16         `json:"threshold"`
	ShareCount      uint16         `json:"share_count"`
	PublicKey       BlsPublicKey   `json:"vk"`
	SecretShare     BlsScalar      `json:"sk_i"`
	SharePublicKeys []BlsPublicKey `json:"vk_vec"`
}

type PartialSignature struct {
	I      uint16       `json:"i"`
	SigmaI BlsSignature `json:"sigma_i"`
}

// Encode returns a msg encoded bytes or error
func (d *PartialSignature) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *PartialSignature) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}
