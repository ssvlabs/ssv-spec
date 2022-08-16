package _basemsgtests

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type IBaseTest interface {
	BuildTest(role types.BeaconRole) *BasePartialSigTest
}

var AllBaseTests = []IBaseTest{
	//ValidMsg(),
	//WrongSigner(),
	//DuplicateMsg(),
	//InvalidMsg(),
	//NoSignedMsgs(),
	//UnknownSigner(),
	UnknownBeaconSigner(),
}

type baseTest struct {
	Name string
	Msgs []struct {
		MsgSignerSKs        []byte
		MsgSignerIDs        types.OperatorID
		BeaconRoots         [][]byte
		BeaconSignerSKs     [][]byte
		BeaconRootSignerIDs []types.OperatorID
		Slots               []phase0.Slot
	}
	msgAccepted       bool // indicates if a msg has been accepted to change the post root
	ExpectedBaseError string
}

type BasePartialSigTest struct {
	Name                    string
	Msgs                    []*ssv.SignedPartialSignatureMessage
	PostDutyRunnerStateRoot string
	ExpectedError           string
}

func (b *baseTest) BuildTest(role types.BeaconRole) *BasePartialSigTest {
	msgs := make([]*ssv.SignedPartialSignatureMessage, 0)
	for _, msgData := range b.Msgs {
		msgSigner := msgData.MsgSignerIDs
		msgSignerSK := &bls.SecretKey{}
		msgSignerSK.Deserialize(msgData.MsgSignerSKs)

		msgType := b.getMsgType(role)

		beaconMsgs := make([]*ssv.PartialSignatureMessage, 0)
		for j := range msgData.BeaconRoots {
			beaconSigner := msgData.BeaconRootSignerIDs[j]
			beaconSignerSK := &bls.SecretKey{}
			beaconSignerSK.Deserialize(msgData.BeaconSignerSKs[j])
			beaconRoot := msgData.BeaconRoots[j]
			slot := msgData.Slots[j]
			sig := beaconSignerSK.SignByte(beaconRoot)

			beaconMsg := &ssv.PartialSignatureMessage{
				Slot:             slot,
				PartialSignature: sig.Serialize(),
				SigningRoot:      beaconRoot,
				Signer:           beaconSigner,
			}

			if msgType == ssv.ContributionProofs {
				// TOOD - add meta data
			}

			beaconMsgs = append(beaconMsgs, beaconMsg)
		}

		msg := ssv.PartialSignatureMessages{
			Type:     msgType,
			Messages: beaconMsgs,
		}
		r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PartialSignatureType))
		sig := msgSignerSK.SignByte(r)
		msgs = append(msgs, &ssv.SignedPartialSignatureMessage{
			Message:   msg,
			Signature: sig.Serialize(),
			Signer:    msgSigner,
		})
	}

	return &BasePartialSigTest{
		Name:                    b.Name,
		Msgs:                    msgs,
		PostDutyRunnerStateRoot: b.getPostRoot(role),
		ExpectedError:           b.getExpectedError(role, b.ExpectedBaseError),
	}
}

func (b *baseTest) getExpectedError(role types.BeaconRole, baseErr string) string {
	if len(baseErr) == 0 {
		return ""
	}
	switch role {
	case types.BNRoleProposer:
		return "failed processing randao message: invalid randao message: " + baseErr
	default:
		panic("beacon role not supporter")
	}
}

func (b *baseTest) getPostRoot(role types.BeaconRole) string {
	switch role {
	case types.BNRoleProposer:
		if b.msgAccepted {
			return "bf0fc94c6f5e64c39cd7f83e62e28db3fc15a4c8fa3e85f775062a14679bde87"
		}
		return "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4"
	default:
		panic("beacon role not supporter")
	}
}

func (b *baseTest) getMsgType(role types.BeaconRole) ssv.PartialSigMsgType {
	switch role {
	case types.BNRoleProposer:
		return ssv.RandaoPartialSig
	default:
		panic("beacon role not supporter")
	}
}
