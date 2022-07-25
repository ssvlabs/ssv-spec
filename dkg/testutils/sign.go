package testutils

import (
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
)

type DepositSignDataSet struct {
	DkgPartyDataSet
	PublicKey         []byte
	Root              []byte
	PartialSignatures map[types.OperatorID][]byte
	FinalSignature    []byte
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
			OperatorIDs:           s.IndicesVec(),
			Threshold:             s.Threshold(),
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
		RequestID:             reqId[:],
		OperatorID:            uint64(operatorId),
		EncryptedShare:        FakeEncryption(s.SecretShares[operatorId]),
		Committee:             s.IndicesVec(),
		Threshold:             uint64(len(s.PartyData[operatorId].Coefficients) - 1),
		ValidatorPublicKey:    s.PublicKey,
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
	ds := TestSuiteFourOperators()
	return DepositSignDataSet{
		DkgPartyDataSet: ds,
		PublicKey:       h2b("8adbbb94ab3b4741e651e20255ad33e73483d0c83181b3aedad5fec9d648e952bfd4baeef8236781ce00300d17ae31ad"),
		Root:            h2b("fd42dbe9b7d9a89973682d8a0e4f106bacc6c0c9a0a274b984179460f8a15494"),
		PartialSignatures: map[types.OperatorID][]byte{
			1: h2b("a718e6bde435b2c451d2563e0e4191365ebf343062d89cc178b93a6d5bc3cd00627b03bdf2270cb955d85c1dadafb2050364f04a813a12b71de5a8ea9bd59f2b95a1402b00bd6d8628e6d3617c59511b7e543b46e341d9c88815972bd22be009"),
			2: h2b("aefeffccb022b9e262142a6b9bc52a1b8cdfa2235b27754c73baf53efa90c4163d35eb3d316177bbc93dbadc4976c6a4107a8631d615ff07d47c0d39636fe9478d1b5cb3b1b7bfe06b59d04fc76b504a79cc7a9b70d61a17989e59575d37eb04"),
			3: h2b("a547ebdbe91074474230ec6c05ef289c7eeb2156b2ae4aabdb20b676c65d3666ebf8e3e63a0c38842616986413d40ce90deb431608230756dab597c7673e26c4b1b8ef6dbd63213c8cbe73ea2edf3f514022f54f21566821d5068f04ad635260"),
			4: h2b("b13373f8cb291e083e437fa805950731e8ead7629afdd29ed38e63155b32fb98cabf731644465068d94bac0d7f118e940c27ae8e296710f32e578225cbef9dc657808a5f0c00b70da00db8d47b525ae076e58285fd7ac17ec1aee4d1d5b46ac6"),
		},
		FinalSignature: h2b("8367b439510a7eca0921c81cb05fb39e8a4d7460ca638aa77db4dc5f48085d2c23eaf9ef99eb919ba86ad7910145ac1709422a5604a2181e81324136e11a7810dc193d2f8bf480e7899cb1edb9e9b032a26d3260ff531972a4757cf042d61f7c"),
	}
}

func (s *DepositSignDataSet) MakeOutput(operatorId types.OperatorID) *dkgtypes.SignedDepositDataMsgBody {
	reqID := TestingRequestID
	out := &dkgtypes.SignedDepositDataMsgBody{
		RequestID:             make([]byte, len(reqID)),
		OperatorID:            uint64(operatorId),
		EncryptedShare:        FakeEncryption(s.SecretShares[operatorId]),
		Committee:             s.DkgPartyDataSet.IndicesVec(),
		Threshold:             uint64(len(s.PartyData[operatorId].Coefficients) - 1),
		ValidatorPublicKey:    s.PublicKey,
		WithdrawalCredentials: TestingWithdrawalCredentials,
		DepositDataSignature:  s.FinalSignature,
	}
	copy(out.RequestID, reqID[:])
	return out
}
