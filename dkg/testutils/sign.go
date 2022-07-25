package testutils

import (
	"encoding/hex"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/gogo/protobuf/sortkeys"
)

func h2b(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

type DepositSignDataSet struct {
	testingutils.TestKeySet
	Root              []byte
	PartialSignatures map[types.OperatorID][]byte
	FinalSignature    []byte
}

func (s DepositSignDataSet) GetVkVec() [][]byte {
	var vkVec [][]byte
	for _, opId := range s.GetOperatorIds() {
		vkVec = append(vkVec, s.Shares[types.OperatorID(opId)].GetPublicKey().Serialize())
	}
	return vkVec
}

func (s DepositSignDataSet) MakeLocalKeyShare(index uint64) *dkgtypes.LocalKeyShare {
	return &dkgtypes.LocalKeyShare{
		Index:           index,
		Threshold:       s.GetDkgThreshold(),
		PublicKey:       s.ValidatorPK.Serialize(),
		SecretShare:     s.Shares[types.OperatorID(index)].Serialize(),
		Committee:       s.GetOperatorIds(),
		SharePublicKeys: s.GetVkVec(),
	}
}

func (s DepositSignDataSet) GetDkgThreshold() uint64 {
	return s.TestKeySet.Threshold - 1 // Different threshold standard
}

func (s DepositSignDataSet) GetOperatorIds() []uint64 {
	ids := make([]uint64, len(s.Shares))
	count := 0
	for id, _ := range s.Shares {
		ids[count] = uint64(id)
		count++
	}
	sortkeys.Uint64s(ids)
	return ids
}

func (s DepositSignDataSet) Operator(operatorId types.OperatorID) *dkgtypes.Operator {
	return &dkgtypes.Operator{
		OperatorID:       operatorId,
		ETHAddress:       s.DKGOperators[operatorId].ETHAddress,
		EncryptionPubKey: &s.DKGOperators[operatorId].EncryptionKey.PublicKey,
	}
}

func (s DepositSignDataSet) ParsedInitMessage(operatorId types.OperatorID) *dkgtypes.Message {
	reqId := TestingRequestID
	msg := dkgtypes.ParsedInitMessage{
		Header: &dkgtypes.MessageHeader{
			SessionId: reqId[:],
			MsgType:   int32(dkgtypes.InitMsgType),
			Sender:    uint64(operatorId),
			Receiver:  0,
		},
		Body: &dkgtypes.Init{
			OperatorIds:           s.GetOperatorIds(),
			Threshold:             s.GetDkgThreshold(), // Different threshold standard
			WithdrawalCredentials: TestingWithdrawalCredentials,
			Fork:                  TestingForkVersion[:],
		},
		Signature: nil,
	}
	base, _ := msg.ToBase()
	return base
}

func (s DepositSignDataSet) ParsedPartialSigMessage(operatorId types.OperatorID) *dkgtypes.Message {
	msg := &dkgtypes.ParsedPartialSigMessage{
		Header: &dkgtypes.MessageHeader{
			SessionId: TestingRequestID[:],
			MsgType:   int32(dkgtypes.PartialSingatureMsgType),
			Sender:    uint64(operatorId),
		},
		Body: &dkgtypes.PartialSigMsgBody{
			Signer:    uint64(operatorId),
			Root:      s.Root,
			Signature: s.PartialSignatures[operatorId],
		},
		Signature: nil,
	}
	base, _ := msg.ToBase()
	return base
}

func (s DepositSignDataSet) ParsedSignedDepositDataMessage(operatorId types.OperatorID) *dkgtypes.Message {
	reqId := TestingRequestID
	body := &dkgtypes.SignedDepositDataMsgBody{
		RequestId:             reqId[:],
		OperatorId:            uint64(operatorId),
		EncryptedShare:        FakeEncryption(s.Shares[operatorId].Serialize()),
		Committee:             s.GetOperatorIds(),
		Threshold:             s.GetDkgThreshold(),
		ValidatorPublicKey:    s.ValidatorPK.Serialize(),
		WithdrawalCredentials: TestingWithdrawalCredentials,
		DepositDataSignature:  s.FinalSignature,
		OperatorSignature:     nil,
	}
	root, _ := body.GetRoot()
	sig := FakeEcdsaSign(root, s.DKGOperators[operatorId].ETHAddress.Bytes())
	body.OperatorSignature = sig
	msg := &dkgtypes.ParsedSignedDepositDataMessage{
		Header: &dkgtypes.MessageHeader{
			SessionId: reqId[:],
			MsgType:   int32(dkgtypes.SignedDepositDataMsgType),
			Sender:    uint64(operatorId),
		},
		Body:      body,
		Signature: nil,
	}
	base, _ := msg.ToBase()
	return base
}

func TestDepositSignDataSetFourOperators() DepositSignDataSet {
	ks := testingutils.Testing4SharesSet()
	return DepositSignDataSet{
		TestKeySet: *ks,
		Root:       h2b("1fa0068233c6c0ffedd8fb1c6dea0fd13d67d5a558b803acecf15601e36dd8d9"),
		PartialSignatures: map[types.OperatorID][]byte{
			1: h2b("a2eba343f78bb70fca2f0fe45880f1770d0797b96f26f393243a3941aebf47a55205a1a0cda6558120a067cb11a653f418fa27946c315d67a3d86c48db27d2977e58c59acd8fe60219d89e68c4a8151307a8372be015c3d25a01d0e5a0263d24"),
			2: h2b("b2fce5a84f7c8c641ec0e59a12840a8ffc4e7b95078cb8919f960dac2f5aba157b71cc6d562d2aa7126d0e19e142b1a616c8c03a6274df51a50a856c0ee732605ccbea74091b1e6f7f819cf7d3c57e074a9d227aad2be0bca5448008782e433a"),
			3: h2b("865c9e5e9041a8ecfa3a04d940de408cd6a99be48c81ca69b4f22ed57a9ec3b406bc60e65e7118c6b770e41d456c54b40aaad0c6362176136bb0d1cd845bc4bce13ec113ff2b291e556a44f13a2bed864a21eb3034c1111767ec3b6c2a909e1e"),
			4: h2b("94a7257ac6be9bd2110059ff2990700ed7390315779758920f2af1d1a86dbdeca7af022624e6a624840a786e88b6538506be9fd6f1e25b633707dd0cc5c36bcceb3b4cb9512f4f36a7b8f53a8e1d7efa16b6b860b35301347ee85a7c4f50462f"),
		},
		FinalSignature: h2b("8367b439510a7eca0921c81cb05fb39e8a4d7460ca638aa77db4dc5f48085d2c23eaf9ef99eb919ba86ad7910145ac1709422a5604a2181e81324136e11a7810dc193d2f8bf480e7899cb1edb9e9b032a26d3260ff531972a4757cf042d61f7c"),
	}
}

func (s *DepositSignDataSet) MakeOutput(operatorId types.OperatorID) *dkgtypes.SignedDepositDataMsgBody {
	reqID := TestingRequestID
	out := &dkgtypes.SignedDepositDataMsgBody{
		RequestId:             make([]byte, len(reqID)),
		OperatorId:            uint64(operatorId),
		EncryptedShare:        FakeEncryption(s.Shares[operatorId].Serialize()),
		Committee:             s.GetOperatorIds(),
		Threshold:             s.GetDkgThreshold(),
		ValidatorPublicKey:    s.ValidatorPK.Serialize(),
		WithdrawalCredentials: TestingWithdrawalCredentials,
		DepositDataSignature:  s.FinalSignature,
	}
	copy(out.RequestId, reqID[:])
	return out
}
