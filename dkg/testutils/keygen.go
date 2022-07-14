package testutils

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/vss"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

type KeygenPartyData struct {
	Index         uint32
	Coefficients  [][]byte
	BlindFactor   []byte
	DlogR         []byte
	R1Message     keygen.ParsedMessage
	R2Message     keygen.ParsedMessage
	R3Messages    []keygen.ParsedMessage
	R4Message     keygen.ParsedMessage
	OwnShare      []byte
	LocalKeyShare keygen.LocalKeyShare
}

func h2b(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

func TestSuiteThreeOfFour() []KeygenPartyData {
	suite := make([]KeygenPartyData, 4)
	suite[0] = KeygenPartyData{
		Index: 1,
		Coefficients: [][]byte{
			h2b("0000000000000000000000000000000000000000000000000000000000000056"),
			h2b("0000000000000000000000000000000000000000000000000000000000000084"),
			h2b("000000000000000000000000000000000000000000000000000000000000007a"),
		},
		BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000f2"),
		DlogR:       h2b("00000000000000000000000000000000000000000000000000000000000000fe"),
		R1Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  1,
			},
			Body: &keygen.KeygenMsgBody{
				Round1: &keygen.Round1Msg{Commitment: h2b("0f005be0af9f824763987d264b839b1b77f8124c54799f9fc1c4ebae043d8e59")},
			},
		},
		R2Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  1,
			},
			Body: &keygen.KeygenMsgBody{
				Round2: &keygen.Round2Msg{
					DeCommmitment: [][]byte{
						h2b("997b2de22feea1fb11d265cedac9b02020c54ebf7cbc76ffdfe2dbfda93696e5f83af8d2c4ff54ce8ee987edbab19252"),
						h2b("8c62ca6abda1a9af02d5c477d2bbf4c00900328f3f03c45f5e1e6bc69a5be2b7acc2532a923f19cb4d4ab43d0d2f42ec"),
						h2b("91d2fe0eded16c39a891ba065319dabfe2c0c300f5e5f5c84f31f6c52344084f0bb60d79650fc1dfe8d2a26fe34bd1fa"),
					},
					BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000f2"),
				},
			},
		},
		R3Messages: []keygen.ParsedMessage{
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   1,
					Receiver: 2,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000346")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   1,
					Receiver: 3,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("000000000000000000000000000000000000000000000000000000000000062c")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   1,
					Receiver: 4,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000a06")},
				},
			},
		},
		R4Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  1,
			},
			Body: &keygen.KeygenMsgBody{
				Round4: &keygen.Round4Msg{
					Commitment:        h2b("b422f8004e8e7c47cf4bc69c3a551b3491916e415b824c2d064204d55c465fb6839834a3f37d8a9271c75e5e2d1f3718"),
					PubKey:            h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
					ChallengeResponse: h2b("68e9c65b695f8adcaea2817d6cb341fa413b2a81f44b93857f6f4d587a5d011e"),
				},
			},
		},
		OwnShare: h2b("0000000000000000000000000000000000000000000000000000000000000154"),
		LocalKeyShare: keygen.LocalKeyShare{
			Index:           1,
			Threshold:       2,
			ShareCount:      4,
			PublicKey:       h2b("b5f7fb0f225f4efee8fafb9117f21e4c0a82f1b5e31e9a4aad46ca618a1dbec125f76480792bc7665f5ec2ad265642c5"),
			SecretShare:     h2b("00000000000000000000000000000000000000000000000000000000000004fe"),
			SharePublicKeys: [][]byte {
				h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
				h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
				h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
				h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
			},
		},
	}
	suite[1] = KeygenPartyData{
		Index: 2,
		Coefficients: [][]byte{
			h2b("0000000000000000000000000000000000000000000000000000000000000097"),
			h2b("0000000000000000000000000000000000000000000000000000000000000099"),
			h2b("0000000000000000000000000000000000000000000000000000000000000073"),
		},
		BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000057"),
		DlogR:       h2b("000000000000000000000000000000000000000000000000000000000000009b"),
		R1Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  2,
			},
			Body: &keygen.KeygenMsgBody{
				Round1: &keygen.Round1Msg{Commitment: h2b("904f8b42da19b72760cdcf1833b0abfea0e0c3acea8853aa242303155d23a2f9")},
			},
		},
		R2Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  2,
			},
			Body: &keygen.KeygenMsgBody{
				Round2: &keygen.Round2Msg{
					DeCommmitment: [][]byte{
						h2b("8623144b531c2852fb755a4d8b4c9b303a026de6f99b1e88a1e91fa82bc10d6c7a9d8dad7926b6b7afd21ca4edb92408"),
						h2b("a82f4819a86b89c9cbd6d164e959fe0061e6a9b705862be2952d3cf642b515bd5edae4e6338e4eeb975a9082ff205bb7"),
						h2b("8a5898f52fe9b20f089d2aa31e9e0a3fe26c272ce087ffdfd3490d3f4fa1cacbec4879f5f7cd7708e241a658be5e4a2f"),
					},
					BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000057"),
				},
			},
		},
		R3Messages: []keygen.ParsedMessage{
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   2,
					Receiver: 1,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("00000000000000000000000000000000000000000000000000000000000001a3")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   2,
					Receiver: 3,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("000000000000000000000000000000000000000000000000000000000000066d")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   2,
					Receiver: 4,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000a2b")},
				},
			},
		},
		R4Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  2,
			},
			Body: &keygen.KeygenMsgBody{
				Round4: &keygen.Round4Msg{
					Commitment:        h2b("a69ec7c89252e2531c057ebeb86098e3b59ca01558afd5f6de4ec40370cb40de07856334770ecacbf23e123201266f67"),
					PubKey:            h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
					ChallengeResponse: h2b("5faad3dcac6fb4f59d43a115f7d626d5120ce01821566b76db0dda178b8e0ebb"),
				},
			},
		},
		OwnShare: h2b("0000000000000000000000000000000000000000000000000000000000000395"),
		LocalKeyShare: keygen.LocalKeyShare{
			Index:           2,
			Threshold:       2,
			ShareCount:      4,
			PublicKey:       h2b("b5f7fb0f225f4efee8fafb9117f21e4c0a82f1b5e31e9a4aad46ca618a1dbec125f76480792bc7665f5ec2ad265642c5"),
			SecretShare:     h2b("0000000000000000000000000000000000000000000000000000000000000c20"),
			SharePublicKeys: [][]byte {
				h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
				h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
				h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
				h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
			},
		},
	}
	suite[2] = KeygenPartyData{
		Index: 3,
		Coefficients: [][]byte{
			h2b("000000000000000000000000000000000000000000000000000000000000003d"),
			h2b("0000000000000000000000000000000000000000000000000000000000000096"),
			h2b("0000000000000000000000000000000000000000000000000000000000000077"),
		},
		BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000d8"),
		DlogR:       h2b("0000000000000000000000000000000000000000000000000000000000000097"),
		R1Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  3,
			},
			Body: &keygen.KeygenMsgBody{
				Round1: &keygen.Round1Msg{Commitment: h2b("cfc682867677f33047d6ab7d067c6bd3910e439cd7bb3d428285064073103735")},
			},
		},
		R2Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  3,
			},
			Body: &keygen.KeygenMsgBody{
				Round2: &keygen.Round2Msg{
					DeCommmitment: [][]byte{
						h2b("912b440c4d3c8177a012cea1cc58115cbc6795afc389363c7769bf419b9451bcde764586cf26c15e9906ea54837d031a"),
						h2b("b194e855fa3d9ab53cbfbc97e7e0ce463723428bb1ad25952713eac04d086bf2407bdb78f8b8173f07aa795bd5e491dc"),
						h2b("a8e1bc8a6493fc7ed293f44c99b28d31561c4818984891e5817c92d270c9408241ceaca44ab079409d13cc0df9e2e187"),
					},
					BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000d8"),
				},
			},
		},
		R3Messages: []keygen.ParsedMessage{
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   3,
					Receiver: 1,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("000000000000000000000000000000000000000000000000000000000000014a")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   3,
					Receiver: 2,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000345")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   3,
					Receiver: 4,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000a05")},
				},
			},
		},
		R4Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  3,
			},
			Body: &keygen.KeygenMsgBody{
				Round4: &keygen.Round4Msg{
					Commitment:        h2b("8623144b531c2852fb755a4d8b4c9b303a026de6f99b1e88a1e91fa82bc10d6c7a9d8dad7926b6b7afd21ca4edb92408"),
					PubKey:            h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
					ChallengeResponse: h2b("696405ff3bfbf720e771b7b5ab4646e62addafcc5cbce479fa1e700a9211bb8e"),
				},
			},
		},
		OwnShare: h2b("000000000000000000000000000000000000000000000000000000000000062e"),
		LocalKeyShare: keygen.LocalKeyShare{
			Index:           3,
			Threshold:       2,
			ShareCount:      4,
			PublicKey:       h2b("b5f7fb0f225f4efee8fafb9117f21e4c0a82f1b5e31e9a4aad46ca618a1dbec125f76480792bc7665f5ec2ad265642c5"),
			SecretShare:     h2b("00000000000000000000000000000000000000000000000000000000000016ba"),
			SharePublicKeys: [][]byte {
				h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
				h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
				h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
				h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
			},
		},
	}
	suite[3] = KeygenPartyData{
		Index: 4,
		Coefficients: [][]byte{
			h2b("000000000000000000000000000000000000000000000000000000000000002a"),
			h2b("000000000000000000000000000000000000000000000000000000000000003b"),
			h2b("0000000000000000000000000000000000000000000000000000000000000058"),
		},
		BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000095"),
		DlogR:       h2b("00000000000000000000000000000000000000000000000000000000000000ef"),
		R1Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  4,
			},
			Body: &keygen.KeygenMsgBody{
				Round1: &keygen.Round1Msg{Commitment: h2b("733eb4134c09f1b43b85af773fa70a572dd71b0fa8db0fb2104daa5cc252866b")},
			},
		},
		R2Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  4,
			},
			Body: &keygen.KeygenMsgBody{
				Round2: &keygen.Round2Msg{
					DeCommmitment: [][]byte{
						h2b("8ce3b57b791798433fd323753489cac9bca43b98deaafaed91f4cb010730ae1e38b186ccd37a09b8aed62ce23b699c48"),
						h2b("98536b398e5b7f1276f7cb426fba0ec2b8b0b64fba7785ea528bebed6ae56c0dee59f5d295fa4c97a1c621ecacfc4ec3"),
						h2b("a19dd710fbf120dbd2ce410c1abeb52c639d2c3be0ec285dc444d6edea01cee272988e051d5c9c37f06fea79b96ba57b"),
					},
					BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000095"),
				},
			},
		},
		R3Messages: []keygen.ParsedMessage{
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   4,
					Receiver: 1,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("00000000000000000000000000000000000000000000000000000000000000bd")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   4,
					Receiver: 2,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("0000000000000000000000000000000000000000000000000000000000000200")},
				},
			},
			{
				Header: &base.MessageHeader{
					MsgType:  int32(dkg.ProtocolMsgType),
					Sender:   4,
					Receiver: 3,
				},
				Body: &keygen.KeygenMsgBody{
					Round3: &keygen.Round3Msg{Share: h2b("00000000000000000000000000000000000000000000000000000000000003f3")},
				},
			},
		},
		R4Message: keygen.ParsedMessage{
			Header: &base.MessageHeader{
				MsgType: int32(dkg.ProtocolMsgType),
				Sender:  4,
			},
			Body: &keygen.KeygenMsgBody{
				Round4: &keygen.Round4Msg{
					Commitment:        h2b("a7e53203bbed6adaa99c54f786622592dcaa4cd702e9aaaa355b8dcf302301f8b8dfec87625a9560079d3f8daf076c5d"),
					PubKey:            h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
					ChallengeResponse: h2b("71cba0135f3b81949434f0ef92aecb74f9feaacde42ae3ad2bb5e27cb83e121d"),
				},
			},
		},
		OwnShare: h2b("0000000000000000000000000000000000000000000000000000000000000696"),
		LocalKeyShare: keygen.LocalKeyShare{
			Index:           3,
			Threshold:       2,
			ShareCount:      4,
			PublicKey:       h2b("b5f7fb0f225f4efee8fafb9117f21e4c0a82f1b5e31e9a4aad46ca618a1dbec125f76480792bc7665f5ec2ad265642c5"),
			SecretShare:     h2b("00000000000000000000000000000000000000000000000000000000000024cc"),
			SharePublicKeys: [][]byte {
				h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
				h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
				h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
				h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
			},
		},
	}
	return suite
}

var BaseInstance = func() dkg.Protocol {
	state := keygen.EmptyKeygen(2, 4)
	partyData := TestSuiteThreeOfFour()[0]
	state.PartyI = 1
	state.Coefficients = make(vss.Coefficients, len(partyData.Coefficients))
	for i, bytes := range partyData.Coefficients {
		state.Coefficients[i].Deserialize(bytes)
	}
	copy(state.BlindFactor[:], partyData.BlindFactor)
	state.DlogR.Deserialize(partyData.DlogR)
	return &keygen.KGProtocol{
		Identifier: dkg.RequestID{},
		Operator:   1,
		Init: dkg.Init{
			Nonce:                 0,
			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
			Threshold:             2,
			WithdrawalCredentials: testingutils.TestingWithdrawalCredentials,
			Fork:                  spec.Version{},
		},
		State: &state,
	}
}
