package dkg

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// KeyGenOutput is the bare minimum output from the protocol
type KeyGenOutput struct {
	Share           *bls.SecretKey
	OperatorPubKeys map[types.OperatorID]*bls.PublicKey
	ValidatorPK     types.ValidatorPK
	Threshold       uint64
}

type keygenOutputJSON struct {
	Share           []byte
	OperatorPubKeys map[types.OperatorID][]byte
	ValidatorPK     types.ValidatorPK
	Threshold       uint64
}

func (k *KeyGenOutput) MarshalJSON() ([]byte, error) {
	rep := keygenOutputJSON{
		Share:           k.Share.Serialize(),
		OperatorPubKeys: map[types.OperatorID][]byte{},
		ValidatorPK:     k.ValidatorPK,
		Threshold:       k.Threshold,
	}
	for id, pk := range k.OperatorPubKeys {
		rep.OperatorPubKeys[id] = pk.Serialize()
	}

	return json.Marshal(rep)
}

func (k *KeyGenOutput) UnmarshalJSON(b []byte) error {
	var rep keygenOutputJSON
	if err := json.Unmarshal(b, &rep); err != nil {
		return err
	}
	k.ValidatorPK = rep.ValidatorPK
	k.Threshold = rep.Threshold

	k.OperatorPubKeys = map[types.OperatorID]*bls.PublicKey{}
	k.Share = &bls.SecretKey{}

	for id, bytes := range rep.OperatorPubKeys {
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(bytes); err != nil {
			return err
		}
		k.OperatorPubKeys[id] = pk
	}

	if err := k.Share.Deserialize(rep.Share); err != nil {
		return err
	}

	return nil
}

// KeyGenProtocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type KeyGenProtocol interface {
	Start(init *Init) error
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *SignedMessage) (bool, *KeyGenOutput, error)
}
