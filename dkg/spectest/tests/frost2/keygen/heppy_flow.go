package keygen

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	SessionPKs map[types.OperatorID]string = map[types.OperatorID]string{
		1: "036ff75a45bb43f1190f89838326ed4f2e090293184e56ff4a01a1a6db548fbae6",
		2: "038680ce08d663c436ddb98265dd26a0c775bf4728ab5ae385671eeb5b87ab08e7",
		3: "0204470b016f243d34ff27d8c869c3b8012612232390d8d3259bc40bf4dc3c4551",
		4: "0328893f709ce7ad1ee70f393cf5ba152fc11043043f0a0acb1591923ebea52dbd",
	}

	PreparationMessageBytes = func(id types.OperatorID) []byte {
		pk, _ := hex.DecodeString(SessionPKs[id])
		msg := &frost.ProtocolMsg{
			Round: frost.Preparation,
			PreparationMessage: &frost.PreparationMessage{
				SessionPk: pk,
			},
		}
		byts, _ := msg.Encode()
		return byts
	}
)

func HappyFlow() *frost2.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	init := &dkg.Init{
		OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
		Threshold:             3,
		WithdrawalCredentials: testingutils.TestingWithdrawalCredentials,
		Fork:                  testingutils.TestingForkVersion,
	}
	initBytes, _ := init.Encode()
	root := testingutils.DespositDataSigningRoot(ks, init)

	return &frost2.MsgProcessingSpecTest{
		Name: "happy flow",
		Operator: &dkg.Operator{
			OperatorID:       1,
			ETHAddress:       ks.DKGOperators[1].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[1].EncryptionKey.PublicKey,
		},
		NodeConfig: &dkg.Config{
			KeygenProtocol:  frost.New,
			ReshareProtocol: frost.NewResharing,
			Network:         network,
			Signer:          keyManager,
			Storage:         storage,
		},
		InputMessages: []*dkg.SignedMessage{
			SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: identifier,
				Data:       initBytes,
			}),
			SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(2),
			}),
			SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(3),
			}),
			SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(4),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(1, root, ks.Shares[1]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 1, root),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			1: ks.SignedOutputObject(identifier, 1, root),
			2: ks.SignedOutputObject(identifier, 2, root),
			3: ks.SignedOutputObject(identifier, 3, root),
			4: ks.SignedOutputObject(identifier, 4, root),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}

func SignDKGMsg(sk *ecdsa.PrivateKey, opID types.OperatorID, msg *dkg.Message) *dkg.SignedMessage {
	signedMessage := &dkg.SignedMessage{
		Message: msg,
		Signer:  opID,
	}

	root, err := signedMessage.GetRoot()
	if err != nil {
		panic(err)
	}

	sig, err := crypto.Sign(root, sk)
	if err != nil {
		panic(err)
	}

	signedMessage.Signature = sig
	return signedMessage
}
