package keygen

import (
	"errors"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/golang/protobuf/proto"
)

func (x *ParsedMessage) GetRoot() ([]byte, error) {
	baseMsg, err := x.ToBase()
	if err != nil {
		return nil, err
	}
	return baseMsg.GetRoot()
}

func (x *ParsedMessage) FromBase(base *base.Message) error {
	raw, err := proto.Marshal(base)
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, x)
}

func (x *ParsedMessage) ToBase() (*base.Message, error) {
	raw, err := proto.Marshal(x)
	if err != nil {
		return nil, err
	}
	base := &base.Message{}
	err = proto.Unmarshal(raw, base)
	if err != nil {
		return nil, err
	}
	return base, nil
}

func (x *ParsedMessage) IsValid() bool {
	cnt := 0
	if x.Body.Round1 != nil {
		cnt += 1
	}
	if x.Body.Round2 != nil {
		cnt += 1
	}
	if x.Body.Round3 != nil {
		cnt += 1
	}
	if x.Body.Round4 != nil {
		cnt += 1
	}
	return cnt == 1
}

func (x *ParsedMessage) GetRoundNumber() (int, error) {
	if x.Body.Round1 != nil {
		return 1, nil
	}
	if x.Body.Round2 != nil {
		return 2, nil
	}
	if x.Body.Round3 != nil {
		return 3, nil
	}
	if x.Body.Round4 != nil {
		return 4, nil
	}
	return 0, errors.New("invalid round")
}

type LocalKeyShare struct {
	Index           uint64   `json:"i"`
	Threshold       uint64   `json:"threshold"`
	ShareCount      uint64   `json:"share_count"`
	PublicKey       []byte   `json:"vk"`
	SecretShare     []byte   `json:"sk_i"`
	Committee       []uint64 `json:"committee"`
	SharePublicKeys [][]byte `json:"vk_vec"`
}

type ParsedMessages = []*ParsedMessage
