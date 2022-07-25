package testutils

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/gg20"
	"github.com/bloxapp/ssv-spec/gg20/algorithms/vss"
	types2 "github.com/bloxapp/ssv-spec/gg20/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/gogo/protobuf/sortkeys"
)

type DkgPartyDataSet struct {
	testingutils.TestKeySet
	PublicKey       []byte
	SecretShares    map[types.OperatorID][]byte
	SharePublicKeys map[types.OperatorID][]byte
	PartyData       map[types.OperatorID]*DkgPartyData
}

func (s DkgPartyDataSet) Threshold() uint64 {
	out := uint64(0)
	for _, data := range s.PartyData {
		out = uint64(len(data.Coefficients) - 1)
		break
	}
	return out
}

func (s DkgPartyDataSet) OperatorIDs() []types.OperatorID {
	var opIds []types.OperatorID
	for _, id := range s.IndicesVec() {
		opIds = append(opIds, types.OperatorID(id))
	}
	return opIds
}

func (s DkgPartyDataSet) IndicesVec() []uint64 {
	ids := make([]uint64, len(s.SharePublicKeys))
	count := 0
	for id, _ := range s.SharePublicKeys {
		ids[count] = uint64(id)
		count++
	}
	sortkeys.Uint64s(ids)
	return ids
}

func (s DkgPartyDataSet) VkVec() [][]byte {
	ids := s.IndicesVec()
	var out [][]byte

	for _, id := range ids {
		out = append(out, s.SharePublicKeys[types.OperatorID(id)])
	}
	return out
}

func (s DkgPartyDataSet) MakeLocalKeyShare(index uint64) *dkgtypes.LocalKeyShare {
	threshold := len(s.PartyData[types.OperatorID(index)].Coefficients) - 1
	return &dkgtypes.LocalKeyShare{
		Index:           index,
		Threshold:       uint64(threshold),
		PublicKey:       s.PublicKey,
		SecretShare:     s.SecretShares[types.OperatorID(index)],
		Committee:       s.IndicesVec(),
		SharePublicKeys: s.VkVec(),
	}
}

func (s DkgPartyDataSet) R1(operatorId types.OperatorID) *types2.ParsedMessage {
	return &types2.ParsedMessage{
		Header: &dkgtypes.MessageHeader{
			MsgType: int32(dkgtypes.ProtocolMsgType),
			Sender:  uint64(operatorId),
		},
		Body: &types2.KeygenMsgBody{
			Round1: &types2.Round1Msg{
				Commitment: s.PartyData[operatorId].Commitment,
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R2(operatorId types.OperatorID) *types2.ParsedMessage {
	return &types2.ParsedMessage{
		Header: &dkgtypes.MessageHeader{
			MsgType: int32(dkgtypes.ProtocolMsgType),
			Sender:  uint64(operatorId),
		},
		Body: &types2.KeygenMsgBody{
			Round2: &types2.Round2Msg{
				Decommitment: s.PartyData[operatorId].DeCommitment,
				BlindFactor:  s.PartyData[operatorId].BlindFactor,
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R3(operatorId types.OperatorID, receiver types.OperatorID) *types2.ParsedMessage {
	return &types2.ParsedMessage{
		Header: &dkgtypes.MessageHeader{
			MsgType:  int32(dkgtypes.ProtocolMsgType),
			Sender:   uint64(operatorId),
			Receiver: uint64(receiver),
		},
		Body: &types2.KeygenMsgBody{
			Round3: &types2.Round3Msg{
				Share: s.PartyData[operatorId].Shares[receiver],
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R4(operatorId types.OperatorID) *types2.ParsedMessage {
	return &types2.ParsedMessage{
		Header: &dkgtypes.MessageHeader{
			MsgType: int32(dkgtypes.ProtocolMsgType),
			Sender:  uint64(operatorId),
		},
		Body: &types2.KeygenMsgBody{
			Round4: &types2.Round4Msg{
				Commitment:        s.PartyData[operatorId].ProofCommitment,
				PubKey:            s.SharePublicKeys[operatorId],
				ChallengeResponse: s.PartyData[operatorId].ProofResponse,
			},
		},
		Signature: nil,
	}
}

type DkgPartyData struct {
	Coefficients    [][]byte
	BlindFactor     []byte
	DlogR           []byte
	Commitment      []byte
	DeCommitment    [][]byte
	Shares          map[types.OperatorID][]byte
	ProofCommitment []byte
	ProofResponse   []byte
}

type KeygenPartyData struct {
	Index         uint32
	Coefficients  [][]byte
	BlindFactor   []byte
	DlogR         []byte
	R1Message     types2.ParsedMessage
	R2Message     types2.ParsedMessage
	R3Messages    []types2.ParsedMessage
	R4Message     types2.ParsedMessage
	OwnShare      []byte
	LocalKeyShare dkgtypes.LocalKeyShare
}

func h2b(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

var baseInstance = func(dataset DkgPartyDataSet) dkgtypes.Protocol {
	partyData := dataset.PartyData[1]
	threshold := uint64(len(partyData.Coefficients) - 1)
	state := gg20.EmptyKeygen(threshold, uint64(len(dataset.PartyData)))
	committee := make([]uint64, len(dataset.PartyData))
	count := 0
	for id, _ := range dataset.PartyData {
		committee[count] = uint64(id)
		count++
	}
	sortkeys.Uint64s(committee)
	state.Committee = committee
	state.PartyI = 1
	state.Coefficients = make(vss.Coefficients, len(partyData.Coefficients))
	for i, bytes := range partyData.Coefficients {
		state.Coefficients[i].Deserialize(bytes)
	}
	copy(state.BlindFactor[:], partyData.BlindFactor)
	state.DlogR.Deserialize(partyData.DlogR)
	ids := make([]types.OperatorID, len(dataset.PartyData))
	count = 0
	for id, _ := range committee {
		ids[count] = types.OperatorID(id)
		count++
	}

	return &gg20.KGProtocol{
		Identifier: dkgtypes.RequestID{},
		Operator:   1,
		Init: dkgtypes.Init{
			OperatorIds:           dataset.IndicesVec(),
			Threshold:             threshold,
			WithdrawalCredentials: testutils.TestingWithdrawalCredentials,
			Fork:                  testutils.TestingForkVersion[:],
		},
		State: &state,
	}
}

var BaseInstance = baseInstance(TestSuiteFourOperators())
var BaseInstanceNonContinuous = baseInstance(TestSuiteFourOperatorsNonContinuous())
var SevenOperatorsInstance = baseInstance(TestSuiteSevenOperators())
var TenOperatorsInstance = baseInstance(TestSuiteTenOperators())
var ThirteenOperatorsInstance = baseInstance(TestSuiteThirteenOperators())

func TestSuiteThreeOfFourSmallValues() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("b5f7fb0f225f4efee8fafb9117f21e4c0a82f1b5e31e9a4aad46ca618a1dbec125f76480792bc7665f5ec2ad265642c5"),
		SecretShares: map[types.OperatorID][]byte{
			1: h2b("00000000000000000000000000000000000000000000000000000000000004fe"),
			2: h2b("0000000000000000000000000000000000000000000000000000000000000c20"),
			3: h2b("00000000000000000000000000000000000000000000000000000000000016ba"),
			4: h2b("00000000000000000000000000000000000000000000000000000000000024cc"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1: h2b("b8164b884c361ba68b639002bc3538e2390df1f9dffe4ece9c762bb1a1fee3c3c9e8eb5023b7a4ad2709f3dc5af04546"),
			2: h2b("9131d754d6a8835180d705d3abec2ed38a9f6f4251f5fe8780bac24313c73dc7b958aeff554d3c78a6a2a2e460de3d50"),
			3: h2b("849d872d0de74efc91250ed55d5548a050f9fff365c9ba28a370472f43bd791ac2f120f74e5c223e8495fb0a0c5196fc"),
			4: h2b("b67aaeea82cb25c782cf96aff800f8421032134c6fac788eb1e2f5080554c7f207aa8e8be0bc547421465ab60fee3914"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("0000000000000000000000000000000000000000000000000000000000000056"),
					h2b("0000000000000000000000000000000000000000000000000000000000000084"),
					h2b("000000000000000000000000000000000000000000000000000000000000007a"),
				},
				BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000f2"),
				DlogR:       h2b("00000000000000000000000000000000000000000000000000000000000000fe"),
				Commitment:  h2b("0f005be0af9f824763987d264b839b1b77f8124c54799f9fc1c4ebae043d8e59"),
				DeCommitment: [][]byte{
					h2b("997b2de22feea1fb11d265cedac9b02020c54ebf7cbc76ffdfe2dbfda93696e5f83af8d2c4ff54ce8ee987edbab19252"),
					h2b("8c62ca6abda1a9af02d5c477d2bbf4c00900328f3f03c45f5e1e6bc69a5be2b7acc2532a923f19cb4d4ab43d0d2f42ec"),
					h2b("91d2fe0eded16c39a891ba065319dabfe2c0c300f5e5f5c84f31f6c52344084f0bb60d79650fc1dfe8d2a26fe34bd1fa"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("0000000000000000000000000000000000000000000000000000000000000154"),
					2: h2b("0000000000000000000000000000000000000000000000000000000000000346"),
					3: h2b("000000000000000000000000000000000000000000000000000000000000062c"),
					4: h2b("0000000000000000000000000000000000000000000000000000000000000a06"),
				},
				ProofCommitment: h2b("b422f8004e8e7c47cf4bc69c3a551b3491916e415b824c2d064204d55c465fb6839834a3f37d8a9271c75e5e2d1f3718"),
				ProofResponse:   h2b("68e9c65b695f8adcaea2817d6cb341fa413b2a81f44b93857f6f4d587a5d011e"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("0000000000000000000000000000000000000000000000000000000000000097"),
					h2b("0000000000000000000000000000000000000000000000000000000000000099"),
					h2b("0000000000000000000000000000000000000000000000000000000000000073"),
				},
				BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000057"),
				DlogR:       h2b("000000000000000000000000000000000000000000000000000000000000009b"),
				Commitment:  h2b("904f8b42da19b72760cdcf1833b0abfea0e0c3acea8853aa242303155d23a2f9"),
				DeCommitment: [][]byte{
					h2b("8623144b531c2852fb755a4d8b4c9b303a026de6f99b1e88a1e91fa82bc10d6c7a9d8dad7926b6b7afd21ca4edb92408"),
					h2b("a82f4819a86b89c9cbd6d164e959fe0061e6a9b705862be2952d3cf642b515bd5edae4e6338e4eeb975a9082ff205bb7"),
					h2b("8a5898f52fe9b20f089d2aa31e9e0a3fe26c272ce087ffdfd3490d3f4fa1cacbec4879f5f7cd7708e241a658be5e4a2f"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("00000000000000000000000000000000000000000000000000000000000001a3"),
					2: h2b("0000000000000000000000000000000000000000000000000000000000000395"),
					3: h2b("000000000000000000000000000000000000000000000000000000000000066d"),
					4: h2b("0000000000000000000000000000000000000000000000000000000000000a2b"),
				},
				ProofCommitment: h2b("a69ec7c89252e2531c057ebeb86098e3b59ca01558afd5f6de4ec40370cb40de07856334770ecacbf23e123201266f67"),
				ProofResponse:   h2b("5faad3dcac6fb4f59d43a115f7d626d5120ce01821566b76db0dda178b8e0ebb"),
			},
			3: {
				Coefficients: [][]byte{
					h2b("000000000000000000000000000000000000000000000000000000000000003d"),
					h2b("0000000000000000000000000000000000000000000000000000000000000096"),
					h2b("0000000000000000000000000000000000000000000000000000000000000077"),
				},
				BlindFactor: h2b("00000000000000000000000000000000000000000000000000000000000000d8"),
				DlogR:       h2b("0000000000000000000000000000000000000000000000000000000000000097"),
				Commitment:  h2b("cfc682867677f33047d6ab7d067c6bd3910e439cd7bb3d428285064073103735"),
				DeCommitment: [][]byte{
					h2b("912b440c4d3c8177a012cea1cc58115cbc6795afc389363c7769bf419b9451bcde764586cf26c15e9906ea54837d031a"),
					h2b("b194e855fa3d9ab53cbfbc97e7e0ce463723428bb1ad25952713eac04d086bf2407bdb78f8b8173f07aa795bd5e491dc"),
					h2b("a8e1bc8a6493fc7ed293f44c99b28d31561c4818984891e5817c92d270c9408241ceaca44ab079409d13cc0df9e2e187"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("000000000000000000000000000000000000000000000000000000000000014a"),
					2: h2b("0000000000000000000000000000000000000000000000000000000000000345"),
					3: h2b("000000000000000000000000000000000000000000000000000000000000062e"),
					4: h2b("0000000000000000000000000000000000000000000000000000000000000a05"),
				},
				ProofCommitment: h2b("8623144b531c2852fb755a4d8b4c9b303a026de6f99b1e88a1e91fa82bc10d6c7a9d8dad7926b6b7afd21ca4edb92408"),
				ProofResponse:   h2b("696405ff3bfbf720e771b7b5ab4646e62addafcc5cbce479fa1e700a9211bb8e"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("000000000000000000000000000000000000000000000000000000000000002a"),
					h2b("000000000000000000000000000000000000000000000000000000000000003b"),
					h2b("0000000000000000000000000000000000000000000000000000000000000058"),
				},
				BlindFactor: h2b("0000000000000000000000000000000000000000000000000000000000000095"),
				DlogR:       h2b("00000000000000000000000000000000000000000000000000000000000000ef"),
				Commitment:  h2b("733eb4134c09f1b43b85af773fa70a572dd71b0fa8db0fb2104daa5cc252866b"),
				DeCommitment: [][]byte{
					h2b("8ce3b57b791798433fd323753489cac9bca43b98deaafaed91f4cb010730ae1e38b186ccd37a09b8aed62ce23b699c48"),
					h2b("98536b398e5b7f1276f7cb426fba0ec2b8b0b64fba7785ea528bebed6ae56c0dee59f5d295fa4c97a1c621ecacfc4ec3"),
					h2b("a19dd710fbf120dbd2ce410c1abeb52c639d2c3be0ec285dc444d6edea01cee272988e051d5c9c37f06fea79b96ba57b"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("00000000000000000000000000000000000000000000000000000000000000bd"),
					2: h2b("0000000000000000000000000000000000000000000000000000000000000200"),
					3: h2b("00000000000000000000000000000000000000000000000000000000000003f3"),
					4: h2b("0000000000000000000000000000000000000000000000000000000000000696"),
				},
				ProofCommitment: h2b("a7e53203bbed6adaa99c54f786622592dcaa4cd702e9aaaa355b8dcf302301f8b8dfec87625a9560079d3f8daf076c5d"),
				ProofResponse:   h2b("71cba0135f3b81949434f0ef92aecb74f9feaacde42ae3ad2bb5e27cb83e121d"),
			},
		},
	}
}

func TestSuiteFourOperators() DkgPartyDataSet {
	return DkgPartyDataSet{
		TestKeySet: *testingutils.Testing4SharesSet(),
		PublicKey:  h2b("8adbbb94ab3b4741e651e20255ad33e73483d0c83181b3aedad5fec9d648e952bfd4baeef8236781ce00300d17ae31ad"),
		SecretShares: map[types.OperatorID][]byte{
			1: h2b("159208425a2e719697647e32b745f3108505694422f8780304b4c7a0a588c83a"),
			2: h2b("477d0d26c8bda3bc5e77efb961f4105bee19fbb5fbd35c94785a48574f7b12fc"),
			3: h2b("526882bfa229b2ca61166e4e36680ba85e817b9fc0670cb273e4e143735c4d49"),
			4: h2b("3654690ce6729ec09f3ff9f134a1e4f5d63be90170b3885cf7549265112c7721"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1: h2b("8ef0df412670cc9b838e61d29052a2b730c816cd516e84346fa575022296bf62046f8e28433ed2cce18537c64456b3bb"),
			2: h2b("8172a74048867bce1860dfc077a0034aaabc1b4a151efc27806f8e92161efd90bc9546041a6b5142ccfc03550bcf0568"),
			3: h2b("b8ec7b107017fce6cf6c3fa1ec3ce02cdf0b271a10d0559f880006fff0154b3e4e39a0bce54315bbb0f2cb28701d1923"),
			4: h2b("8b72e8714bb9b80416eb11dd2a62144367932eebf22511c1c12150b3d123aee8ef351222ae1c8fccee68a7b8a880b755"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("34a77f6219fe5e3fd77444167eb842e7ca4505b988e049709935c98a7f7423b9"),
					h2b("0c23d2ca0c6cabfe335989e0dec2896ecda50c55a1360e4fe285dcc11d074a5c"),
					h2b("39fc21c50f44805588427be8048c3413b4d929c15bf154383229bd2509ac950f"),
				},
				BlindFactor: h2b("a7eacdba5e614ab343adfdc2e2273ff30e9b31abc27ec84bec67b4862e841386"),
				DlogR:       h2b("14c03aadeea18b95888cef14ce7d577e78f1611c6dfa5c580b271cb3bac3468b"),
				Commitment:  h2b("4c4d00683f5bc74e8cbc40d493dc10e571b5b7791593efa684e511661dcdf5ff"),
				DeCommitment: [][]byte{
					h2b("b736d5228eb746e7338f46682c62ba72afc02c2dc27cb053000abc3ee3ddbb98947cd34d2062bf9073cdef07cd76110f"),
					h2b("8216fc394fc2452dad2dfa99c859ca92240647ba551892e4cd77f14dd026984ca81f44ef208551a13b1d2c6bd6bedd73"),
					h2b("aedde85b96669b526186a576cded3e5d01ed5ea488d058496375eef7255fc5ea44d8c7dfcc58dbdec0db57b8ab2c03b9"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("06d9cc9e0c120d4b5fd671d758652864f90597cd86094ff9ade56371a6280323"),
					2: h2b("4d045d641caebd01f8bd97683b2a760991787d643b14fef326e877a2e0350cab"),
					3: h2b("1f4be30df89972d33bb604b913c47bcaec226e77a8069e5f043f06202d9b404f"),
					4: h2b("658bac41f30d294f8f3369d9f576e9b3b07eb30dccdae63b45e90ee78e5a9e11"),
				},
				ProofCommitment: h2b("b3d1dbe5efcd55f5bfa0d055ab4376270ae2789a695e3f4f0759708079126f0e978a7a23bbe192635c2e69687a81bcf2"),
				ProofResponse:   h2b("3267dd6ffb1f6f08f24be34a9092151d79b6df9521b9305163c17789b57ad4e6"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("1b38be7959a9d6bb404e6348ea87e00feafea1fa5c952ad57b484d5b93646df3"),
					h2b("106def05f6a35564b6d441f3412640efe6b179dd89b98d7f0d6fdaf54c4bc633"),
					h2b("4d2ce8c5db08b12d0fcf77849dd17315cb6b6127c9c1ec887d0c203044e649c3"),
				},
				BlindFactor: h2b("60488a81b35fba615578f46cc1aee060bf3bf660e2ca0d54c4fb93dd85246556"),
				DlogR:       h2b("1eff751418c65956dab0b57573334a752661b1b683fd7abe35cdf771345da1a5"),
				Commitment:  h2b("545f29e48b5a0bc2f5a7189c0ef985232bda7db4f809d41e6b38450ea889533d"),
				DeCommitment: [][]byte{
					h2b("ab9b709d18eb1335e2fba7ca953c43e6a83f3e97e9872c4f89379c661398189d5ab8b63b64b3d44408ac943618c33e89"),
					h2b("8c6527ab4f466abb16ec5682a5b332875e0bb3d6de3d0a913b541a6a2ca300e0cc766b3ccda775c599141562290e72f4"),
					h2b("af64fba748ebb6a3ced81b7a58400d3249b7c4712bc3f672f5f568f5dd4d4740672d69624c05739244c9105d021e1a58"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("04e5eef201b86004d3b844b8bfddbc10495dd8fcb01248de05c4488224967de8"),
					2: h2b("14ff49a3363ace6053873d29c734a636ead62e4b9714e3f88a58840a3f952162"),
					3: h2b("4b84ce8cf73121cdbfbb4c9c008c9e83cf67a1e7119cfc250904fff3e4605861"),
					4: h2b("3488d65c1afddd04e51a9b076243ccf1a3548fcc1fac356481c9bc4012f822e4"),
				},
				ProofCommitment: h2b("8308a8376b2a49ea6dcd239763ae59bdc1f2951fe5be7c0797256abf02f2d08f3857970954de573663ed6d74e0efa08a"),
				ProofResponse:   h2b("13cffd5645b07c00e7cf015ce57758b0dab6836d010cd7b394d9a3cf046f7e7a"),
			},
			3: {
				Coefficients: [][]byte{
					h2b("1b429e2c309895e87e381af5ffb17bb3a910ce2b691cf9a6c1174b5339d8340a"),
					h2b("48314e157ee9fd3ed7c89d55a53b95dd863ff60faa2455ea1bde5360cf706d87"),
					h2b("5a343ac5d3a42c9c78f4430f0cd91f26c1be35c826e26227a9927c8d883f81cb"),
				},
				BlindFactor: h2b("f950232571ad809e466ddad4f2122dbef5e2b45f351e845d4ce5d5702e2c0420"),
				DlogR:       h2b("2f038de1506ff2f2c26016b5cfecdea578e9c271fd049726cc7ada0b1bac1d82"),
				Commitment:  h2b("e2c8cc98479332d5ab88156e508b2b7648f36c7b5140ad4331f71b4ca3ede6f1"),
				DeCommitment: [][]byte{
					h2b("8aa955f38e9d87ea6b8f5ed96c9080455db8a7dbd1fab67d070e4f6f40f84f0f9f43dfc1e0f4a75f4b922a8594981e6d"),
					h2b("925bda96685f2076914495ccc6f14c6d4b14933e5bad1b1d7dfb4d8bd5c706e56f6004cef1f555574af3aaf5edd6e31c"),
					h2b("99a5b3c6674846bdacfb235ea4a06259e3c09d4ddf8cb15d1aa033e7190cd04abd129f282056f897077e198557ec0cfe"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("49ba7fb45989427b9bbb2352a82458b29d5156003a2555b986881b429188235b"),
					2: h2b("44bf8821d6874db744b301bd5705c3f46d93015f58f5be1d9f1de44ef9b71640"),
					3: h2b("0c51b774a792b79b791fb6360c55bd7919d5d048c58e32d30ad8a67872650cb9"),
					4: h2b("145eb4fff648fd706c3b18c4d1b61d45f5d766bf7fed0fd8c9b861bdfb9206c7"),
				},
				ProofCommitment: h2b("868a51a2664ccbd0ac5d1c5aeaf766ecf692dc5d498077b0e555e133f5eeebdb18733122eaf7e1a26384c538940c2150"),
				ProofResponse:   h2b("2c59ab12f76bb4a81884c880acd4afd6bad8d79a7d0d4cb4c6569b024fa9187b"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("395fe6b105764c05dc550774e0afc5256c6a9670e740a90f435efce428d4a74f"),
					h2b("419e1f99e118a6cbc56960cbabb1dc104d6083fca230ef851bf9d14eba48655c"),
					h2b("2cf515595f86c9898ccfec1f7dc0c4bd9300d0122942a8db6b2a32356625172b"),
				},
				BlindFactor: h2b("367dd5b9259f4169fcf372f277d67d645704edd9bd7597a04b75f1d2d2a6e4de"),
				DlogR:       h2b("433963c95f02f4887211840f8c4a37c3f2511a2d48f726d4e44edbb2e5c8bec3"),
				Commitment:  h2b("a8187fdf9eace26e1c9f06d0483655b8ab50ae746cf2276f729c857337bdb3ba"),
				DeCommitment: [][]byte{
					h2b("adf941c332d05482288513313cb0c072502c6a7ac7bc6d770eebe4dd3b900c9a2b5fd163ef03f60de0946ad182c5ae56"),
					h2b("b0221d757dce658b6c219bfecc4fc9ec2814bef1bc480fe84c0c15cd08220eb6c63b22a5d69b727ad944d1c2750875f2"),
					h2b("b0fd7e0c8de1a2f178d10bf205a6feaa56e629788f1008c0a37f7ee2b3216e90f722de4cece1e4a311ff584a38f54159"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("340574511c783f12fb547c5800808dedf90e467cb2b5e570ca830069494223d5"),
					2: h2b("14a78550c8ea47eb00b9f1721231082c57f5f2a9d0b2178a27fb685a35f9ceb0"),
					3: h2b("4f33c1033469e3d61fbf3ecb1f630be5dcdf3efb41339b5a5bc834b5eefba7e1"),
					4: h2b("6fbc80153559958c252a8c5b1e74c115340c876e043c14e265e9657d7447af67"),
				},
				ProofCommitment: h2b("9592a0296fb0ba33474371e80f92739ec454e01c4fc343cf94b4632a115d03fd8484bee8100bd71ecfc2fe48158767b6"),
				ProofResponse:   h2b("25e1299915bcae87b1e89241ecb25535545ccc8fa93d8ad225e6acf620a322ab"),
			},
		},
	}
}

func TestSuiteFourOperatorsNonContinuous() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("8c5a360dd3e42ab46e02f7ba1006d9b8fcaaf0b5a84ee7662f0e8aadf90ba7e4935bb754d44e664661d12303186184ef"),
		SecretShares: map[types.OperatorID][]byte{
			1: h2b("0c9df69c2facd10db55aa107592fc66aebe745ee23d22c4098ea97da90379deb"),
			2: h2b("69045855c46e7b6c03681e5220526742326f1d30418300a6f31ba3538a702104"),
			4: h2b("55282fc57774d31f7a54072714e61812cda710987dde0fe78fc069f0c8e9ba53"),
			8: h2b("264bc1dd26dbfb5b3a047a06dab41e6186e7870dfa6e4c316a14b5d16dff396c"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1: h2b("aee333fb4b2bdce4ea36e96249f9ce4a3ff23017d224a38a5aa34a8ae9f2d347a67d6f7c9cf944a3187f0998946bc5cc"),
			2: h2b("8153a31d24c5dd4615f2bd8d9b63e2cb060ecf8d39d479411226331709406dd07ded47c34023c0dcd8c3b8587f3a9b77"),
			4: h2b("b47f2a14d678d96a23036480136f4afdb5825fb36b3724fb8676a1491c68aad36f6f732378763652d69e5faf09ce7939"),
			8: h2b("a5ad90c6ef31fa3baee65abdad69b0f95ad4225e6ebe95ef1b26fee539a05110a1e3ecb4ce0036b50b3bf9efa748f46a"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("567a7b8011e0aeea0f2f79eea46a4d48ec7d189136582237200a4d400d21e5b9"),
					h2b("3de7bbdfe1c9bfdd9faf9fbc1f4ea3383d28c9c52286cb6269e50b07fa2e8247"),
					h2b("39f5f9c9b41d564fc5683d78acfaba53ff4cb09089780f9a8aae64bf865733ba"),
				},
				BlindFactor: h2b("2cfe4c1c994de75d3fb8007cbc613f75d64d06d923235c60d6df1dd3605024f4"),
				DlogR:       h2b("1da9b4cc252f03a0d60de93fbf9de1a5ea250b5902a90fb046b1e582e38b6801"),
				Commitment:  h2b("fbdde05c7a73dd0a3c3af8642705cd3f5cef3d173ecf26f77d67d95832442074"),
				DeCommitment: [][]byte{
					h2b("b5d488c050f6f8a44c2b013237b52550bbefd74cf33051676226b768c97303d2148d7f52dc0620cd2134457b8d9e18a6"),
					h2b("8c3300a80b89a4f328fbf1f652e93528a09ad582c7ca03f7706e5f251fd5e1208e2aca179cfeb134302f2a1116cd2daf"),
					h2b("85743c7782e54d9de3496e409269dcca1004b484dc409bc884c05cb1a219a21817de7f2ee04fd3f48238e9d808401fc5"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("5a6a89d67e2a47cf410d7f1b6711d2cfd534eee3e258a135149dbd088da79bb9"),
					2: h2b("5e58e46d2911100bca8227317a0cf4f968c88254a14ae3691e8df6511adbb92c"),
					4: h2b("66307e5b3ab62e8ae42f601990fe0f349082e0905804b1737284c5625b4f2a6b"),
					8: h2b("01dd9de723c12658ff5f9cd17929c3458e86d66da8cf18121acbd5867462e64c"),
				},
				ProofCommitment: h2b("91790504c9da3e4e3b2e1c92a916bd70d36ab6bd8d031bfdee5a3f1df1997e5c16f93060b81ef676114bc41f1ff2bc84"),
				ProofResponse:   h2b("331af3ee2a7ee063758be64699cb688dd351fb450fdbb8920661221116a59883"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("09e646f0fc3c9b6636bda28bd1b27cf98ecf2033f1184459ef44fc72d28a75f7"),
					h2b("3728365d4899379d7948bfebc48006c57c8a2b1765f5e0467fc67fe958c51ee2"),
					h2b("6f515b9f0df0d55f7871b4171484c9c5dc1890606a37f08117990a0145b2a68e"),
				},
				BlindFactor: h2b("9add5575b35a0bb1f2952260f270e19873b8664803bce0d224ca5793335117fb"),
				DlogR:       h2b("2176050eb6a420135354e13025249a4ce75871755dd6525e8dd08c8da7129611"),
				Commitment:  h2b("375de77b9c0dd5c7cb7d8527bce7f5136770e3dd9761fa3d9fc9dbbd8323723e"),
				DeCommitment: [][]byte{
					h2b("92c27b010ff62d4d24de8c17cd9a252d382d8b099f597415bd1986bb90447544b693165635a6292078deda699a2a82fc"),
					h2b("aaacfc0274df2befe3b43b465b82f5513e1c9c1f32f99523e876c7b422db33750780ac2b4456f342e1dc023f6a270480"),
					h2b("ab34c204d40be0ca504d476650ffce47c673d90508ba89b1048c1b50708a1bd2cbdba30bdc63670fff3d6a0d3443f948"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("3c72319a29292b1af53e3e86a115757f93b437a8c147b92286a4865e71023b66"),
					2: h2b("65c584db1ebc6afe3e2e929f863e5186a94f27d865ea56ef4d36244e9adf4def"),
					4: h2b("28d4bdd13a397e083c248b23883fdc12b2e8ee652c8ab19667ef9c3d90c95a4e"),
					8: h2b("282765812c3ee0db1bc0cda4a4cd4b42361dec46af2e0b14d3bb7c50055d1046"),
				},
				ProofCommitment: h2b("8a57042b268fccf34a7de02c45f0bbd5ee87b8835d1b211418feb207395d1b8b35883f07a1544e658ed37ea2ada5fd12"),
				ProofResponse:   h2b("6fa7e35dfca1ade4898c8e4b567a0d2fe19b155f48337f5e2d4983d52081d082"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("391a7a8dc036bd52ddbaf7730ec6848cb680931a77ef87b21300fef1cd03be2c"),
					h2b("558d1f3ee7cb62ac90dfd2ed976a2b45bca0548c5a879267154bc43f7ad1df0e"),
					h2b("6938e387a3b1fdf79c18c375192685f9bb8c95316e0c2cf5be8aa721f0d5c7a6"),
				},
				BlindFactor: h2b("9d9dc60286c6c13857194076f1c8b4c309e78ac6b2cfd37ff121ca26c74857bf"),
				DlogR:       h2b("5bc6e03b589331cc1213cce5b59504d4ed934a2c8dcabe447c07511f39ad285f"),
				Commitment:  h2b("8031a39292ea03b94fd06920fcf5c041de86cb991100cd28f6314b9c4d0a0cf9"),
				DeCommitment: [][]byte{
					h2b("849fdd10ee25347b5b687a07390073324df9725b63fbc82dad583469564964d11a34f055f4bbd527f0a16c40ae78760e"),
					h2b("907e354c1a39d0b788b2ceaaa0a9dc12c165206b21caf7ce9476039d98b5900d31f7d4b732cd4b7807f2fe6f6ccd8fa8"),
					h2b("8fba594889e5d9a73aa557a014aea74850ed629da3810bb3f96a94a3bd23b0099a6b130a33ea1a090b44fdc5fc7bb2c5"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("10052eadf8792366a43fddc5ac1385c1873234d240868f10e6d76a5538ab64de"),
					2: h2b("4574028a4e8208216fbc72fa720bbae47b3f5ce9e537945c37c323fd85fe9adb"),
					4: h2b("7015137dd70ed5b37bef21f25b1838e4d233542fc2ec84bb50da821fc5a7b4b3"),
					8: h2b("385281a383b1fe919c746a23ab435bd48f3f8318cfa0589a6008e9a9d9049fdc"),
				},
				ProofCommitment: h2b("949a3a6b8cec3c5ecf6cce0b607e829f9e583c98d4ff2ceaf07373007f9d98212cfbb71d5ad64bf4fdd9aac9f3624321"),
				ProofResponse:   h2b("5d96153830748b86fc4fa1b7c2087e8b89a37bce8e79fdcca3473619dbb497e9"),
			},
			8: {
				Coefficients: [][]byte{
					h2b("07a87215bebc18c9bfda49f99e48e688acef19cb11bad665bf2a28f35751dcad"),
					h2b("321499e8d2a08d5f5069bf999d93800572ff3d68ba2f111be7186472076eb954"),
					h2b("13da4f2551be8f2430feac1c7c5c41d68358db6173be1354708e5cb6fa21cbef"),
				},
				BlindFactor: h2b("ee3c943f538f3ca1d492ecb467b82cb4c2bd38888606a12f3891d5b4728c6c51"),
				DlogR:       h2b("2f2303a7a2e44e4d3cb5ec3ea358a8fe07315691e26613c47bf6bec1568d3947"),
				Commitment:  h2b("33def7fa025524e2ff30048eada71b43c089e5efd38cdbcbb8c6cdfb4867bdf8"),
				DeCommitment: [][]byte{
					h2b("b88c60c7d61a1b1b66cd83da0f7455c6213bfe7c5c39723be00218eb606aa36305d6ec138cb111b6cad2cb3afd350ca5"),
					h2b("84837a022448db65f0ad295b5a45168df4d27c0b0af913e72757215ef2dab360e04fd4ba03fa395c73ecddd0093e376a"),
					h2b("b3d9fef734071652e8d3c3a3021ab2151c9640ca4a112ad29e903c3fc5caa88dd3a69ac6b6579b4eb8ba029335a7c7ce"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("4d975b23e31b354d4142b5afb838a864a34732953fa7fad616d0ea1c58e261f0"),
					2: h2b("474d3b298159f2d0f16ea196c13f15e84c935e1f5512e9f04f9464b44eb67f10"),
					4: h2b("3de92ec17eb14b694484aa07b3d3a3f15f833579365ee0206471862f172980e9"),
					8: h2b("37e1e4247cc772ddb5a97d751b1b8c0a86c0e543d2cf2c6f1b847a501b3aa2ff"),
				},
				ProofCommitment: h2b("83e31877269ecb9f0d8c297f952c7df9a3aade80e14b7edffaabfd183d9e2a5e4dd476e33e5eda0ffaec940899c02f50"),
				ProofResponse:   h2b("6a7575321d862d76d5cd18a436ff4e5d331494c8a9f1271555cd52fb921d4d84"),
			},
		},
	}
}

func TestSuiteSevenOperators() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("8c590972477e53d3a958862db859e2cd2e17c745965db25d92ebfa27247823a6690ce48a1938dfb67eeded83b4f7aa26"),
		SecretShares: map[types.OperatorID][]byte{
			1: h2b("25c05d6154ca4c1535f59a399e56ae3221c3cca27d48ca78e5ad40910268ea32"),
			2: h2b("04557397b2f3330813ac703534893829dd1abaa6490df0c1abee27cf36d20efe"),
			3: h2b("5c6ffcab83d4bb4faaec14a0d845542dc0c229b3393510736355267f79984376"),
			4: h2b("0ff0ebabc8244a62e8144495d9c25d87a41fcace4ce40fcf5a5a97e7960c11ed"),
			5: h2b("477e6a07526e26e4ea0771792e4fedcb2498328db42cd00e30b9c9d37b67150e"),
			6: h2b("5f28f5a0508c111b352228b319a544a5b3e414fdd01bf168d770fc9e3ccbf7d3"),
			7: h2b("2a0fa56c1bb8b556b35c4bd75e0a7ff93ac479b032b09f1390c1632d1146758a"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1: h2b("a9e55560d709e747aaa205857bb1e56d380baa6d3c48cc2d89fab4ee9e0425149365f42098fb1c234b5d5a00ddbc0c22"),
			2: h2b("b4107e0f4ef33a4058ab060062a067f6a8bcf468d7afe4d657d5eceb47f4b089559f020a5e9adafed37d4a63accb3e18"),
			3: h2b("a1703a70b963d584cc4f458a2c59250e77d17a648bf97c4b3f0a3394d4d0dcc3959ffc0e547be073811c277ebf0aa457"),
			4: h2b("8ae459c387a771ae68bb048e7b88168741061989f7e15927e2a0dc6eb50098ce7410e349131077757cdfc07d579cf3de"),
			5: h2b("aaf899348cadc53147b6b12fa0e5d53bcc6563f0a9c65c329b6e082f2544daa39df0763b808bde4465f208bd93a7cc77"),
			6: h2b("a5921b8b113b05b1f3c9ec58661841b6291440463ae97b7872973060384c6cf70f08aed2766d475614a89c7e47944755"),
			7: h2b("862857cc969044b436125cd7faa4701a21104e4e8a1a625d14a8c1ff8d708dc8291ea1ca4529ddfc7b0e7e95907cffdf"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("2793866ce3d5871551049c9ae983d27b947ab5462b28335480f59d876505f252"),
					h2b("0e23cbdf09f1925b6aa8b7ea3c0ee0773b22f0db17a32a9bb4e5875fcbdabc48"),
					h2b("5903a5a02be1e36d5dd62a11ba33fb98b8cde9cf940220009c725409990565e7"),
					h2b("31ccc83fdf0369d69055753c639b3f747cae2b92a382107075508b74696cfd74"),
					h2b("53a37c095f720d9b3e6b296073257583f2b368c391f654e80c6f0fb2bd959899"),
				},
				BlindFactor: h2b("8d96657177334a3cec149210d706fd9d9bf96f667d16edb07d94c8b97ec9d757"),
				DlogR:       h2b("62394aa67ce0c8580100bd23b660d508f102533b8e9bb927558164dc40a948db"),
				Commitment:  h2b("28bb83e28d34d90d3f8aa7b573f259062fbd962dc2ee95ca3f585cad97959e73"),
				DeCommitment: [][]byte{
					h2b("af058c84fd3294ab891b55b0bc6478dc0019dda6ecf0b24796170f2776891327ac4476f39034eba5ac4a65fae39d7ed6"),
					h2b("b2a6b6914b2cc480fa89fa500f069a0645b37703bba53197ddd7bf4f5896e2bbf855ac6cfe8cb0156cd34a4ff8b7bf52"),
					h2b("982a8c4afa4a41a2d053879c90063980fa67bf441e78b799ea4b0d1af7ad9d99743dc45dff9b35d1a0be5812e38c08cf"),
					h2b("8f1833a71143d562e92b40cba75a1c71aee3615cf8f2ec423379ab38783afe4dd22777b1860b82972b079f17a3a2064c"),
					h2b("906ba8f531b62db8c1b0eaf8ef66ab8e8958f1494b0522dab92ef38e0388d8645ef03a29b5277cd3236bb5dc79045881"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("2c4fed8f04e379bf81d06d23a343b3795051dc410c492b4b540d1419f0e8aa8c"),
					2: h2b("49d1f367a96992d66cfbc40fec41a5501b4a9ed2e60a62a4cdff534e8592779c"),
					3: h2b("4bf7c65bc0af06086e0bb971341032ae86392995144bc0c26e4dd51f449abfa1"),
					4: h2b("232b192b6936fbc653b9eddb12059247532b9743a22106d7dee38c59159f3701"),
					5: h2b("584eca451a5e0e8b932e626b4edd7ed6c92794c049d02be9f414e49aa643e06a"),
					6: h2b("3d084678fde455bd7b1818a9be31d70d6533d10674d7ecd2ac9fc2586a330cce"),
					7: h2b("2f63cd98c000d81433ce4ab06ea3d5aa7d0c313a3be7983c31ab82d79b1d5b68"),
				},
				ProofCommitment: h2b("a44834be85ac96618e09f0220357ef821c93650b5ff3d04b6fa40f65668684685e67ac0309ec87261fa3b494c4e9ff36"),
				ProofResponse:   h2b("64935ca34cc4e35181c0d7b0318405ea566323739c0a71a84146ac92da9d09fd"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("6518542ae17c96bdad03ba92e1d9c3624e2cc11ef7f72e093398ff0cb6069a7f"),
					h2b("4629d9d467d481d7179bd82e8572b629bbe756cf7a487b068271c600c3827a5f"),
					h2b("3a151287dc4508c237df1fbfc145b39bb7b8eae202345d6d273e848daab14bba"),
					h2b("2e30171948263a01320c8378d356252b165128eb6882b12bc147a62f0d254384"),
					h2b("213352f8591a1bc1f9fb9413f7258f75c590d7e3476eb6bbfab0d45945f77f3d"),
				},
				BlindFactor: h2b("6d7751c3ef030c184d9403b1638165073da70b8ddc813ca3face9c8cd4aeaccf"),
				DlogR:       h2b("29c23390ca68c9408ed334639f1925db92829f6c87eff4d45b8913e58f1f369d"),
				Commitment:  h2b("5f6ecc574893b1ed0f2bbffa784b1a249bd3a68cf0b2c481db5badea1db2c050"),
				DeCommitment: [][]byte{
					h2b("a8b9bd16cef31ba17dff590a999003fbee17e891e42e3ef5c26b4a12fba25d19934905363aaff3ba935f9c9da02104a6"),
					h2b("ae7478e5966d2b13247220eec277cab34e38c2a0e1aeaaf6a69a59fe4e6857034590f1d546403dce05df3e7aaefd9754"),
					h2b("ac8c209c551df568210a0a3fb566f2cb493bec92ca7c3224c65a10341f85c4cf79f0e6f244448c77933567d3f50708a6"),
					h2b("9930f6d7381202e04cf414b3d5d9a68379130006ccc8a63a7b8e6d768b1f01398efa97c388a7c7e38997def9bfe77cfe"),
					h2b("8b62c2aa59fa126c7ced863529739967042661d406e8baf3cdcedf77b307b878e992728bf090eb802e679ea5070727ec"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("4cdf5bf2739b7c89c21319fddfca31bdf633bb992468b6669941c42577572357"),
					2: h2b("634009b02b48e783b858fe9c94ebd69f174fd3b4b06c9af38ac1145cb072ce0a"),
					3: h2b("11936c5bfa79410716577fa22fb1ed4d6b6762c21abbad33d0a2b16887071436"),
					4: h2b("0e4bbce9871d977e053f3600e38dc405e2fe6f549476704eb4064561aff55d2d"),
					5: h2b("195019a85b102bf850ce0441b8c583e002827933012d4d13fe9f62c86f52fdf0"),
					6: h2b("57c57c3d60ecea0400fc9e9ca959b8dddd836560f6dba5f8fab5848198693932"),
					7: h2b("6c336afc90d66d1cf38adae88dc1a1e16f383618bdeb1f1a73240dd88db53f56"),
				},
				ProofCommitment: h2b("8d0e0741426738911c5c6e6e59b70fb2a6cb423f748fda27a2758342b3311c8e1b54b7e7d61aef25166c533511fbbb7a"),
				ProofResponse:   h2b("619b7111812f4a1e47afb9e58fce1715d5506149233be77f2523f9b1e216e639"),
			},
			3: {
				Coefficients: [][]byte{
					h2b("205af2c52a36beddbca2531cbe310c5385692b14c4465e97ce6bb2732c31db9c"),
					h2b("4a7ff61be71942359d1e0b9f9e3c3a8a2aa39057389e890f69d2fe3962657bb4"),
					h2b("25d108856c1f0afc77aff660a357149e241ef2420a0feb9b40cd7249ce4da5c0"),
					h2b("04e95d16ccd3ef7411c1ee4e9bef18ead74f10234bb052504497b353aaa5f299"),
					h2b("47947ea09acd2860beca9f70fcad080bcb78eb3f9f9449f84f6b7972eac84d85"),
				},
				BlindFactor: h2b("210e1ee349fbf40dfa32ed9b50e4dfc31b70e9adc88220f54603fc04a09bf988"),
				DlogR:       h2b("0c938e1fbc3fa91b5aeeecb10bb29c6ec93e43d3d5d5c6656b7aa979c58c1064"),
				Commitment:  h2b("882700ce8066d60a86324cac3319dd66b7638228e4126a37875b4ace99f3bd51"),
				DeCommitment: [][]byte{
					h2b("9698939427566257c5fe94fa2f3019a66f670fb070394589ac462723f7d700316201e25b99c96302b268134b5622df19"),
					h2b("95d16adb8cce8b49c289ff26042f1a9fe08a4847257f121414d1f53fdde6bb31cd191e00e601a9153b7b71eb49d17ac2"),
					h2b("844ae6a9a41e2855b762f86c0b9420a3c894b4a6296d549b47246e2f9cd08685840be02fa44a09d32f93bd2a109d01ee"),
					h2b("b3254cad4339cd7fd71a533e336bd0b31ca8a3333cfd52da4da64eb103636f2477d37bf6be0618a033ac0a37dfbf6cb8"),
					h2b("9143fdf96a006222a64319466064ed8866e279c2ef0bd6fe3d5b59ed0d0e2cd003a8ff712ba4530f048523868612cd4b"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("693c25cabb72a69c6ec30ad48ebea46d2336050df23b138c0d0f4fbdf2533d2d"),
					2: h2b("0a2054999f58143cb667b4fab51723ae9c91f5b8b49fa537c0bca9e62be7d70f"),
					3: h2b("22d069d1276a886dd21572a9a08560d70f87a29f4291387bb01f091edcfe3f78"),
					4: h2b("624b6b463b567ddd493f2702c32bc360ca22370cc931699213f518730c6c5104"),
					5: h2b("628f6ffcd5c9e10f85f2ff46aee9cb9145a650166ba1600b9710e5c2c3cf2ab8"),
					6: h2b("1c87477b2e118da808144add1efdea7e7a4fa09e43002ace5657e1c20f8b3003"),
					7: h2b("5bf8211ce1b6cbb19895535d01a75a14330e95565e691c14e0c2dff7fecc08bf"),
				},
				ProofCommitment: h2b("a9a9960b028820b0a4f7d09e727ca79b6c6a4cfbd3fc9b2d395d365c06dd5b420c9b9b9b292990980ab9d2ceced55e0d"),
				ProofResponse:   h2b("5193a23a1cd4b7ef7f0ca618330ef68bdb19dacd1c9ef54b78f849eda9c2e138"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("5a31d322df9192bad461bf1d0e56db9a372ba8a98b824643cca638bb199bd435"),
					h2b("5fa2ebb0501f6284dede71e38b436cb13cb0a582d9af0341b6814288a3b39b9a"),
					h2b("4c29dff44b00790c54315c1fc8faecbc4c5f53d14d1f734a64a2c71cf361ca17"),
					h2b("07e3a16ddc09a5d70b266116ddc973fea37980e64d9bcccde8314c84f121db6e"),
					h2b("0869ad809637f2c5940081a9c415aa8f3ef7c502d3286825959b7dd407d8cfd4"),
				},
				BlindFactor: h2b("b4b06d0eac14b567d406138eeaa730d9038407d4b505dbfe6455b7108e666068"),
				DlogR:       h2b("3c6e168b7280ec67f373c4f692cbde5fdf927fe3cba3848f5b15b9c999790f77"),
				Commitment:  h2b("c454404a237594edeb2c1e3dc296777b127a8a81393c30c0dc59a846abb3fbd2"),
				DeCommitment: [][]byte{
					h2b("8d71c20e8b23e1ce507f3f20beeefebcfff06879b829bd68b136cbdfee111169e5ba426f15ae4e075f79de9a2f1a08b5"),
					h2b("ada5b0c6a9b79b3fabdfaa15902193cb476289dab71813ec3faae07f4ab5b1b45c923285abf843ebd8844c797ed3fb08"),
					h2b("b7ef3112a31048644c3b5fa284af18c293a1eff040d6cdd6a0be95364af1659ad573eb25f46210076644e9c2964a9334"),
					h2b("b012d710658ab90feaedfd8b11601b00843b27595eee159fd1efe687d45a41d80c276a0281d5abd3e3005d08b47f7438"),
					h2b("a73c4c188cba8f5b9ab73ba700430a2ab0b9b40e5c3f387a5cec6a37d796c4d6241f62f98b94ec212df8b9d3aaff5fb3"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("2e709f0f99b80c584024bfd0f130a38afb319fe0d31839c565970cbba9abe526"),
					2: h2b("584521d9f5eda75648c426863ea4a0b6f6e0c24212ccdabf67761bae35260c6f"),
					3: h2b("662aecdd95d475673451158595f027258af3398ccdfc2b524d48e2857d52a870"),
					4: h2b("3c982e3106afcd7af5aefef9eeb6b046aba003c153cd952598a8aa15ffcd9368"),
					5: h2b("15fdb0a5c1644bbf2c93a4f09901daf2e09af320bf33e967d1c28715f6862375"),
					6: h2b("0ed9347500da5456f24d4151339f9365f1570a29f6ef04cf87575a1e57bf2a93"),
					7: h2b("25b56f3fc1fe135ad9c285ddac22133ee5067b99ad8bcf9b4ebbcfaad60ef69c"),
				},
				ProofCommitment: h2b("afb2ebde63f2e5ce1e937e858120a645ac5b19b2906d17480349a6133d7faa7a5e9f907a6d7505099ddcf13b6ac05e5a"),
				ProofResponse:   h2b("23ffd4cfed81f222ca5cd69d12a4a84fbb8e53d0e7fba64b70d267d7882a6149"),
			},
			5: {
				Coefficients: [][]byte{
					h2b("4cb0f14b8992cc5d6ff432ad68ec1b94c972f198807189711b38f3d006f67a66"),
					h2b("637f573ece8ceab69afee6f95d131f6cf97332289b3bb66ac9cbdbdc80a78168"),
					h2b("5e82e91494e9a40543addf717fcab1729d6375fa7059cf8cf057150c268c31d4"),
					h2b("70f7850b194fddc6669c1f9ef275f3fcfb02826c2fd6ac9362d98399e12938f9"),
					h2b("542bdb115795a0fcccc01f3b09a49217357f10289014115eb96815ba2e09590d"),
				},
				BlindFactor: h2b("cfb13c78d1a0a813a22815147ce6955dc43cbe618b079b65174ed2322f32faf1"),
				DlogR:       h2b("06ffaa5fb8a48e55e456e2c175e137cdf59d0be900f45f0f0ce50b5a46ee8105"),
				Commitment:  h2b("1a47addc0f4954286c79c32ed9b5a3152a471d47a9f930b39bf2887f0e624983"),
				DeCommitment: [][]byte{
					h2b("a9572e9f4372192c9334a36d9290c434c2a2cf9861c44351feb292ef8b2be806c12be556a4c4ca1aba05449e8751d070"),
					h2b("85798936d068d71b60618b9c57514b3a9cdbce58e82079f7c9dfbd9431f39fca7876bf8835f6827e0b3602bae261bd21"),
					h2b("8b10e01f912ad4d9d942d16eb59412e7bf6a7d0ae85347a8b41b2a7d1452dd8f42827faad9f8c3a961ca2421dc04ea75"),
					h2b("85bddf026112ad68d8ebe5baac65a0e4e52b0982715f6440eec2d1145a2af6da33d8020de1851cd1159f7e76930d0fcf"),
					h2b("80eef357df20f8e410d64815e5db4f5761c214a0d8157e2e040b49c2c89d998c10216b656cbccb3acff9dc0f7ab1a9bf"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("041ff46eb778e4bbb515d7d21b5d127341d49c444bf85d5ef19d7e10bd5cbfa4"),
					2: h2b("05ffc66aadcaf3d2b4e656455f67c90f34673f72f86fb31a1d7a78448c559d05"),
					3: h2b("5093698d21c326cf3ad7cf0d0ab981ef9163af01e7e68789028c06553428ee13"),
					4: h2b("2b96fe97e7f276357714481f236a296315d91662fe6ca486665055b3c5fee77d"),
					5: h2b("6be40960c4541849fe9ee48a10258d9aee612caba20a6cdc70079d53a3e0171c"),
					6: h2b("47f1dd8b46ff141864befc3c54be76856cf0b87fbacc584ba8b41db180b563ec"),
					7: h2b("0f678a7bc5d731246e5b2c34cbf9b5edff496823b2b91b84fb1a20cd60480e12"),
				},
				ProofCommitment: h2b("b88d4ccf44ae97ba4fcc97ea60d8ad2b230b2be48dea14ac25c6de94e6db6ff90f2f9989bec760d4cd99129ffbd52831"),
				ProofResponse:   h2b("37ce228158e299214f9e483326734e862eda46ba2d5ca56cf92e6dd1402104bd"),
			},
			6: {
				Coefficients: [][]byte{
					h2b("5c971677c99df9841c4df302f4952adde9d37f3fc130c3b9ace8a9fdae293346"),
					h2b("34f1b7cd210f5a0b39d35f6b1501aeebbaeef49c915e6c5cdd91a77b6417f92e"),
					h2b("4b36ebee10fe8d77da31ea11f8629ee2f6e4870b9005530f634814e97baf1138"),
					h2b("71bf3e4d7aa431651b8a1bb1c6ed7d084c41fb5a0ba08d1eeae50890697f73ea"),
					h2b("15e5ee5ef51e9d2b63b900def164cebdc16255179679bb506b5cee82bdb1fcca"),
				},
				BlindFactor: h2b("c9fc22ba534500cfbfadb01eb71aed3cf20df1e2d91d96e3a40f3385c03d18a0"),
				DlogR:       h2b("2d6078a1527f3daf508b03a0298a6741c80a7c23a6184c8c4818f4f51bf05d49"),
				Commitment:  h2b("9f97ee92b5789d10a8db7a89f026adc6023c566d2890f6531779950cd51472e1"),
				DeCommitment: [][]byte{
					h2b("ab674bf91f128c0927df20c41ebeec644267d82762e110e1ab11558631d31f91c41d6bc0c647f3d7370f988b81c47e3a"),
					h2b("a269b052554191a631dffc61719ab57d8fdb899eebc829d880ac74183b31f378e79feca53ad68c07574562f1f90ebae9"),
					h2b("93a656602c699cc90e6ad68f1a515db1675ac2866ba201f1bf805eb0b78de6bcea8adf22f550fbd67cf34a61cdee3fbd"),
					h2b("9172732b72a0b1bc6165ade1b7e8ce68751acbbae76caa5385c54bcb9cdf42f94f05027dea39da2d57df185d73fe1794"),
					h2b("a697848cf38114337532cc9024078c4080dc132ba00e7d146e3f1603cb6d05b53347254d50dfe0d3d0008788d63a1441"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("089bf0e5ee9637bf15e8d0f89d663c62ae125f5084b3b79844045d78b521ae5d"),
					2: h2b("14c240460687ea1e10399d25bd60300fcb5c14c3e8bba2be022379588c30d663"),
					3: h2b("5cc9808a4c9af66f1c9124a673090e80fc6f26e05c3f7f4081b5bb6f595aef2f"),
					4: h2b("126d9a98fd250ee0723d394d448f31fc26d03f1069a9aec56de13dd70b53ef82"),
					5: h2b("010b28ec25bd24696ea9e9cb4ff27e8bfd7e5cf1b6cfc67d828276ec597f810c"),
					6: h2b("49fb32f1d527f461954049881adb29e55e3ea24604fac3298c2c3953c3f10067"),
					7: h2b("64932d0c1b5808808f6670a19299bb67519c52f330e7091a6827b3f9936b7d17"),
				},
				ProofCommitment: h2b("98681d7220eb7f402d720dd43609359d1f53385f788f4539f29fa6b1a032847571c62e48fa08ec5575b48140b3137063"),
				ProofResponse:   h2b("4f56a328c68640dfc38d69e1c45a914e6deb105c00a7dc01428ce7be5868ef0f"),
			},
			7: {
				Coefficients: [][]byte{
					h2b("31701076481cda763c7ac2a044fe9ced881fa329532992effb9ae27a1efb70c4"),
					h2b("35eea84c5dec38e31255440b689787626b1e76301cc13350040fe096f9408b0a"),
					h2b("6d068bda79aa8befeda5185f70c2aa3a830137820ca786dc6e86c0de3345d0cf"),
					h2b("58e18dcfafff44b04223052ad9ed1c2ac0c1241e4aa71ab0fa186c7c8b3ee9b8"),
					h2b("1e85dbe3e17214cbf9b4b39a1a797f9238a35359ef55eaabe7cc3fddb4ead5ab"),
				},
				BlindFactor: h2b("d6f97ce469c4330cdab27a5025c19f5b2d240d7cfd22f9dc94578792f293c010"),
				DlogR:       h2b("398882c82b76172069a6c064d20c836583a96adc7b0523bc605346eba7ecd863"),
				Commitment:  h2b("8a534869330e99f526d22ff57d5bf424090817d4c06febd98e2ab8f44386e9d0"),
				DeCommitment: [][]byte{
					h2b("83add110064e1c19b2b43625711080f179e1570b9cc72f207ee5e479c837ab5d9c8c7dc8cd0c9291e173ded4156819a3"),
					h2b("b3acc1e6633dbf7952ec4e93875e4a671b7f43cd87e8ea237573703a9ef7c40b5748d087f1f0c72d5959385aa0939a10"),
					h2b("9778c71959eb8509bae8f66ada6f338882040067400c1a28433492157b23e46d6d3f307bf7a4706cb6909b50d19cc103"),
					h2b("a72a91036c6ad9631833f61ea0b712460a3bb930da0adbd5de58e75c67f571ca64e2ec0d83c3640ad078e585a486fdbd"),
					h2b("9155db29443b5caf65db71081dac1704eabdfea43dfcb03aa768822e56148d836699574afdc6141049b57c71aa833de5"),
				},
				Shares: map[types.OperatorID][]byte{
					1: h2b("63f15faa5de9fe3511d927bfff7bba3cc828804db6929a7b5016304b8bab8bfe"),
					2: h2b("35e4ef5511809702c3b966b4bfbd86c6128327f609fa31570b5d08ef87387275"),
					3: h2b("24557f2719e7d1056266e7c2dd30a3d4a20cd756b5754be2a2baf28bc621aa78"),
					4: h2b("5d5ad8e82d8a5e8900890168e738c043b6c31ffd91325aa246a17015f2eac257"),
					5: h2b("4c2c2923d898fab77fe920579b8ec079420643cee57aee3fe2a80154ae1b5062"),
					6: h2b("68d69c7623805eb95e5a27920cc81e7b358fc5116aad27871d8c22db8e3ef2ed"),
					7: h2b("6886c23eec734c94e8d1104d7bcf89da347bd6624942476959314e0f1fe5504c"),
				},
				ProofCommitment: h2b("8dea97f080048fed7128c1cbd76f0faa601e7945544cb1f58c3749f4205043fa9a101f015cf6297ad37cbbc5462ff248"),
				ProofResponse:   h2b("27fdd0769e37fa4b483a2ca81b2cd932ec6628ba06f56e7ca4a0a412fd074fd1"),
			},
		},
	}
}

func TestSuiteTenOperators() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("81879c900f5ed4534c01c05e7d97f0367ee2d73904090e4eaa46c3372028432dacc3d671bc443fdf7ed48ce6ec3bd068"),
		SecretShares: map[types.OperatorID][]byte{
			1:  h2b("0c55e129d97f6b0b1899ef1bf546a97c31b33b94da5afdae9998bd4a11a6809f"),
			2:  h2b("6c0d358b8490d4fd458b0dd40392f36e369a8f331d6b431a099f361084adafef"),
			3:  h2b("5e5f919dcca1bbe91ece703442188eafcce7fba6951b7415e12651094ea4dc96"),
			4:  h2b("0dfa5ba131421c0223e915810915f86ec26d0685a1376e0f680f1832478d205d"),
			5:  h2b("438fffcefac6f69f220e2aca735a88aaed61d3d780e77187ef27cce40069fc84"),
			6:  h2b("0692cf4ed76bd36fd151b1f45a95d01a3dfd73ce4f2a24bd0e86c5c86be65b82"),
			7:  h2b("4db57ba3cc93af45a8aaab84267bd03e18f20a2ade1c1c500d6ca6a66cc3985f"),
			8:  h2b("295792d19bb4a9f9e7fd72434dc6bd0a2a9396f0713e220a73adf2204a128b8a"),
			9:  h2b("59808df2f83b99bf72a321bfc6f320cc7e1b84cb5684159ec4a3f54409369d4c"),
			10: h2b("45a22f2e84555822466cdeb13e1d39b1b51220cd5e5f438463a50d0cadb2ddb9"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1:  h2b("afab25077fc502f4b0d0bceaebcd869926597a11c36c4d08d2b7afd7b10d84cb92cf987146289466f5422923e6560b94"),
			2:  h2b("85a7413f98bdb049a280596e1e2f40c989787631bfd40380305851e6c6ea23bebc4bf929c65cddc08658e2100bc7c006"),
			3:  h2b("97a5c7140a00fe934f6aafb3acc05f009f288c23d10bea6bc4bc4535488d730da354eea6680335842873ba39d1632714"),
			4:  h2b("b28f05a941683821edb74af472fe2c983022b1ebd2369cbd7c5a508a5fad7463dd1eaa3726a21d25cdbd9a2fda12b239"),
			5:  h2b("8d92ff2ed50202e1d5f85fd9f0a8189b30852477ce3c839022eb417bc5e30e4af85efda18c68b6933bc222a646c80ffc"),
			6:  h2b("900ca609ff9d41fc3803c36e0f7f9d1ecb3f12424f4bebccb68c47de6a042196b4c7c4b42562fcf4835ee1d58a17b41a"),
			7:  h2b("a50f74a1f0264e6da2bffac9641b6464c9a8892b6d739a23718b9cdbfbbfbecca580431f376155a96bc75da30a8c7287"),
			8:  h2b("8b776e9ce4042b033407f1349e1c524aedd44b8909d366f12d9c2eb3db3e14b08411ea7134231eab83bcbdf8372ca4f6"),
			9:  h2b("a8c3572b6378772a315e380f2ec6df0c57df8b391bc12147d700c6edb4ec5289203fb85630058683b40c1f5b7f3d3f82"),
			10: h2b("96e2ddb456df82a1804173da0e64cc2d2f9c8d1bce1dfbf424c44e54c7c18974e9bacad63df62aaba70752ef72113584"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("64c87ae397cc6b321fa473f71088c4f3881fe2593badebcaea50642d41240316"),
					h2b("32694835ea6d6f6ce643bb71e4281e878c770db2a849e5ab0cddd03673bc4fee"),
					h2b("2e7704548f6c30183b75f97d89561c2a432d22f65df5e2646dd38225612e5e22"),
					h2b("655cea94d4eea06b57792ce744aaa907e9c630e3d22947af226f315965795f8c"),
					h2b("3f3cec95d48936c9f2281671336f8f5d88d2279f3afa532170f6a033fbeff484"),
					h2b("39bfcb9a8cbacfd129065bf9c7709610fb1bd96a1cbcf718d8c640e34da795b4"),
					h2b("3ca761187ef8cf7bb424e49a133d5f92c0573408b56ba95bf11bae051fbdf448"),
				},
				BlindFactor: h2b("6abaf09d1d9300d35eaee712e322f297cdce4190be3d43c6a923c1f1de38c619"),
				DlogR:       h2b("390ca389ba25854355384e2926d25c1af24d441bb026d930658f4699a786b685"),
				Commitment:  h2b("c2e15d5a6eec1889c8957d769197ab4ad523639cbc8d2519ffbbe26aebf88e1a"),
				DeCommitment: [][]byte{
					h2b("b298615ac307eab463769ed6ed8dfabdc4c50e44b85f068faa34b6d3fe9f054d8ffcb9c50225b6770d6aa8e3c31650da"),
					h2b("80be95c9f162b7cb8ae6f7f087c0180575781a9f7411ec1beb87b3bfdbde432666a30bfbb829e4b3ca2287396d42c330"),
					h2b("a7548664aadc077780eaff7f1c17b2b2856b32b767fc5a7ebdc624743a9b28eab431a88ac7acb687a7077f7ad8195aee"),
					h2b("8ce3f9b1071609e846ae2cbd6416effe5265c2e9f4dde8fa37478070c8c71efd8a8e3f0dab4360067f97e3ad15b146b0"),
					h2b("90f93e1fcc15cc4885b4b306eb0d7ebc5fb0df9be3485f521543fbec6141c11ceace40797d26cde13c68d084985c7c2f"),
					h2b("9681840cfbbdf49c3dc8878dcff5b862a705de97a21a8146fe6a8a5d4af566411e62830517e6b19c12a6f494d3565a34"),
					h2b("a03838eed202d88c6448061e7fc865ef5de7a1461febe31af88f638f63b562f977cfa179cbec99ecdd98d5457447b6df"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("10f32dff205b8c189b434cb2aa47cd9936d8e8ec21407f23c2497702e4dd8f2e"),
					2:  h2b("38de81ffe2190c4b46f5e2b988ecf628c3ff294f38017a9e3df139303c902856"),
					3:  h2b("04405af09bf7e3e56dfe78af82d1dd7c25927d4940644508a91e739b8bb36c3e"),
					4:  h2b("5804af1506bfe4542da23b39e9e23db783bcef08019a279561503c0deee2d0c5"),
					5:  h2b("0179857e62b9b952395569b64f63ff12f0d56c43659dfbe9bd3bf7b8f20966a7"),
					6:  h2b("59e40b38eee542d6a3c65509ec07ffedfb7175acd65786291234cd8b6ccca545"),
					7:  h2b("455c9a13c91eb4aeabbd73482af9d15804e42a32dc181b23d76c887da5363f5f"),
					8:  h2b("1e017d54e16b8c52c1162179e2da00e7c589f293fd2d4e81e70de9abb89d00e9"),
					9:  h2b("0dcb89f3ccbfb3f05a9512840867ef37872c4ce0debe2907dd306a5e4accb5db"),
					10: h2b("15deff8e1112a72bd039ebae74f8b35b5157f81da6d8a6e794a56df17b6d1a0c"),
				},
				ProofCommitment: h2b("b1bbb6930180ebb2ff21b73e21424e3f8589eb2f05795a4d1f4ff02141c2a94251d8dc930c8d81e9d8f89dd0bcd9ba1a"),
				ProofResponse:   h2b("1a9c328b04602833a22024784b3c52964472b6baad9898ec6886f47deed9d0df"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("321e0c9e9c922ac10b96837029bcb6841e207e507451b4837815bf316c99be82"),
					h2b("6585a7deeb4be155d517cd096106df6084f2a8831bced6d1853e49dc8083ba41"),
					h2b("1bf31897a7a584dad927e734c5201df16ca08549ce5b43e17eae7d8f05601279"),
					h2b("1e9d72197e73bdf090a70d7672f645597658acc7cda645c6858e4c64cce2674a"),
					h2b("058b6be4ccec3034cbf6c2a50b9524bd1d48ff0c34e603420720061d82289e6f"),
					h2b("6d1ce0e0840e5561ec2997d6aef5547e43f4fd1185b1151782f63b95566beb68"),
					h2b("0046fa7a16fdcc20943351d2e3013f5029631cc8ede23660576038f6447bcdd9"),
				},
				BlindFactor: h2b("1b644cf789d2bc94561416db05f70c25adf9bcebd1f0227e07b118e79f6502a3"),
				DlogR:       h2b("7260d3e1a13c3f39882a9721796b42e4fdd996873f0ff7d7537fd6daab52d72a"),
				Commitment:  h2b("eb7cf510a40b40a90c422243e2511ece5bacd8159846646cf018968dd0da4520"),
				DeCommitment: [][]byte{
					h2b("a9f4a935e416454e1b3b6d8e6cf9a4c8cdbf0767f01aec0fb2852dd9c590ac27df2d6e54155b117aae48d9eb907b4f1f"),
					h2b("81c13f36c0251a7bc21ec5970674400ee3697c3c08db76083f550681a3d53550e0c9129657592586417d69f38d2f637f"),
					h2b("84104eb5f0c87a8a710cab1935796841e08a09b1d781c54c2b8ec1c411af4183bbfdcc53b911e77a8f71125267a43713"),
					h2b("8cb0ceed1e7f7a7b2bcbdff003f0f90364084803f4a66a1ca0db06bc99ce481fb506275606924536038f12844059919b"),
					h2b("815c5446692db5a82e786da2365c786760c9927b45b0d0941e76bbdfcdf58d80bfc8d20fb09df44c5c9f35db493cca61"),
					h2b("b83f3efcb2f52f79374c298d2a7b141a06fb1997e30688f30ce4919e921981a2b729c9325e4ae9e350e2e5e5972c3f2c"),
					h2b("ad27b17bb7b46563a13a76dc8e45a457a66597d46ec0ab7ac58c3dcef080cb8a739179caebe0ba66498d028b324d8438"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("5d4837c7c2b4a609305d41634d2201b0693229c5d49eabb8e3074dacdc704a34"),
					2:  h2b("228740b83b3909e4bcef5ec7958766bb7298ff7fcfda1f2b5294be84f92f8144"),
					3:  h2b("501cfe844536b4dff00097b545033c6a229da84ec81b7659f3d3c9c5f63e3b26"),
					4:  h2b("3dbb9e2733f8f8cfd44066c03b50dd397ead7441070970fdd80db3dc926d3ca4"),
					5:  h2b("2d2701d0c9539634936200e1e2c6ebdce21b9e5208221fe5b4bbbf98210cedf6"),
					6:  h2b("5f1bb5a6c33bebb2cc71da393a1f1100abf8c11bcf4a90d06223d4d62f360321"),
					7:  h2b("07fb0d40bff5469d540d9bd93e7c1a76ac2ff89f4b9ab13f18955bd0c54516a4"),
					8:  h2b("37974bcaef59bc192b8722ee0c5f13fc4194753ec656a4366c464d064489367c"),
					9:  h2b("49c5ce56d419495f14c0207fda00f11fa2d5c3a55e3b470507d075d1e1256368"),
					10: h2b("19c7de1f788b71b829d822accf2290d4275b15ec8ef198f5254ef199b8250288"),
				},
				ProofCommitment: h2b("8828925fe282a230de490dcef5d71d5a826b152ad09226a0ed22257bfa7969c54caf5cf80a95635e0362b42d82fe4605"),
				ProofResponse:   h2b("49cdf0b1122e05e023ec14babf58a7d512cce25a3f06631633a202b8781c4f8a"),
			},
			3: {
				Coefficients: [][]byte{
					h2b("0c3bbf3c94c6214bf44373d31c9d9841753086c08685da1bb57da8bd7e5c74b9"),
					h2b("044bb8d01b66dbfdcf3322053f9d9c822416607997ef4020e942eaf4eee6255b"),
					h2b("366c55c87922aaac273687ece94bd39a54b046d7e3c409eef5e6a2010d43063a"),
					h2b("38d776f3d64a8b046e0824f3c01d9f5424dc8d8e06af3467c125d52321dbb2ac"),
					h2b("6e090269c4f9c108cdf27fda69299f5703ff87d6af033b5d84c49b891aed5127"),
					h2b("39d27bc25fab696f6c726be74e2acc2c0021cd6d439bf6fecb810e382390c42e"),
					h2b("4f0e2a50b67c2ba163e1314457838256968c410a2e923fbda18a76a0ac60c064"),
				},
				BlindFactor: h2b("e68c8322393a99dc018f5c2599b9a9e542e6cc6a2905c469695ce575fedbac76"),
				DlogR:       h2b("347b578e0638d16f780e2751ad6ba59ed1b84f3b7986f441b9ff0e2e9fcd74e6"),
				Commitment:  h2b("91f710397a08ef431f0a5fd0e7bfe97973cccf7ea36e149efc46613d907e23a3"),
				DeCommitment: [][]byte{
					h2b("a9089f5c6221bfd006570a373105fa258b4cff735a5f31da4107070629b8aad3452cafeef10d05e73c2be59017e5d9bd"),
					h2b("82b7559a9eeb55bf53331195e876a0de4afe9b3946f1c8babe38d64f2bb8695d2ff68411abbe679015d963918a4c7930"),
					h2b("ab9827d75fc2b61011a3606e7793d50e5b09869012527e68e9b7cc7aa605462fea3a59f3a1308109001011bde3dd2df4"),
					h2b("b095d8a82e5486ba0bae102de1edd9b80d027db98465730454ff5a559d80fecdba8e0351737c7175c3441fb537bd8d3b"),
					h2b("a4b98688ef58e13f7e0e93fe2047998d26e48354ed732446e018b98fd302f34bb920245d82e6835135bd122c0e1b8b96"),
					h2b("9642d1bec4a4ee83f22e97bc6a57982a72fb02820fb6d241df4013035d5847cf4a7c9fac36604ebaa59bb66a5e959090"),
					h2b("991ff06813ff0717c7e08c7c4c0435c4d64b20863d69331935ada379858d44e52bd8086f8218b0b1d8617dcabd8a3a2a"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("1aebf74c5de3113b5d4dd7a6f7970d7bb24865e52a1eb6b0479d2b3b874028b0"),
					2:  h2b("596628b6454fdc5f037348527835a93bf301fc3283b3eac187d5d7d5d9301e97"),
					3:  h2b("67a9d6f4cf5e551e6b5edcd936693e22f0ce1b0a34cb5d11a6611e587b9ab3f5"),
					4:  h2b("08a25fb8a13aec8165815429d71a8ec6cc0c96477c5f1c4c72a73fdbc17a29c3"),
					5:  h2b("3081a1188ed309c809bcb1b9ab8620684f259bff64858cd3955017d5fd795417"),
					6:  h2b("72b47ba8c7fae42fc2aa52ababfd8ecd6173470aee9f89d41fc18b050635611c"),
					7:  h2b("0e4ba5abbeaab3689c18fc03682ce195e5d5ac8b3df23995710b9c1f8a9cb767"),
					8:  h2b("55348abcf63fc451d4a1a9bbbcbaaaf0623d3058c0b4c90670422653366aeba4"),
					9:  h2b("478c741cceaea2ea1841ee30a0e1bf4e8cd5982757af969b1c443da1a6c1cd89"),
					10: h2b("54b08b55a8f37fc1fe70bae549e15ec98a4401c27c36a85970f134fb2edf8c30"),
				},
				ProofCommitment: h2b("b7704b55cb1d093f9d1f26c408102c3382df58284538346f0b0410e7bf947acc3a2b03c24eafa00914df3009da02447e"),
				ProofResponse:   h2b("18cede58ea541879eca2c6a942f0200836b38ef01df649bcbb5b22e7ff9b8fdd"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("25d41a05e2ab51446aab73e10d243b2e0711f68032513f65486f533f9b27ced5"),
					h2b("259d225ffd34f0606aa2c8033be814670e3e2af4a38ace56830b451d399ba473"),
					h2b("5d5703bc0fa17e9c4b749c16e4445e8d6fd8891fd3042c1d267265373255509a"),
					h2b("0d5ac4a96bcfa6bdbd74f09b53d9ab7d832762f4d60de48e0abe51a9b6ff4ed9"),
					h2b("0a713b78ac95fc75050b6cfbec34129b0f829f3686804d7af3bca531da9a455d"),
					h2b("3ce9b7abe716a1201f2b09324d749b906517989c4f9b247d14e5720cfc014354"),
					h2b("0162933473a1acac973997db8d56d55acd6694c50be0baf0818512ab86e067b7"),
				},
				BlindFactor: h2b("a1bdd82cf7a524531ccf073c74dfefa7118c1b3495247854e06dd91837ca00a1"),
				DlogR:       h2b("53b502820600ab57249b7254813f5c01fd9d61220dc04a6191050a6ed340b01b"),
				Commitment:  h2b("22844e74c20b6ee1c1db238ebd71665c1f4c16bba6fcb8c46993fe9c42322ca6"),
				DeCommitment: [][]byte{
					h2b("b259850b10869fea69f0bd49d0e107979843635f7c09fd80ed5e89286e6917cdeb8b0b1e4ec9fb0b6602751c8047663f"),
					h2b("ae84bce486839eb6ae12c7e2d74b8caa11b12bc0c79003562c6411b74b0a8400a39f6b0bea2fd1f6e79779118dbac3d0"),
					h2b("aad48b6689a867a94f3ec1ff077015b2797298765b3bef661a3a7c9313d39be0ab2a85bd38e8b990af0ae3307870d2e8"),
					h2b("a1ea2623b5bafd9f28cb6bbc5b678fb9ed6f04601dccc1b8f1fe469b766a79c7239ae4550cec09921932984a5183c3a3"),
					h2b("8c850e6d3ea720e0d807fc84100b1745ebc9e3f8c3bd89059d6b3cfb6db1ba2643ecac6544105a78fb7db610d6f49c2d"),
					h2b("8e016bc5bae324f78101e8e84c288ab1365a4e5f877c62af7d36973d5b35e00658e2a373b1d4bbd73395b22a09c65649"),
					h2b("977cb40e1cc84c12b37d328ec85f255d24d536f75a46a572e844cbb63ff0b97fd69e34f4fe487303759b90979c038cf8"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("17053c7e0f64b6b03334269034e62d1ba2d5921b60ed935186d2792a1b940321"),
					2:  h2b("0fe85d5a27f3ba19907ca5489530ec0f1726607cc9a6348177ff3f5b71957ee3"),
					3:  h2b("27cb6afed0d7db2f1a7582e788aa331208964d3d0e5d1c8816170095cef21528"),
					4:  h2b("6af00b293d9c9df937387e58ed60695ae94aa6a370d213eb0eef46450c8d3b0a"),
					5:  h2b("335e86609153192eceb9ff3201a122a6739c9764a80b50293f4a84f724d8d572"),
					6:  h2b("0b2dc0f2f8e8b21ab6c823b1a82226e2da031e8a1dd0c23074555b430569b987"),
					7:  h2b("352ec245625ba9673d6767bd964508e76ba5315a944f694f7368532ab7afdfbd"),
					8:  h2b("634453b570f45da1d1cf552cccc7bc6df2fc6f6233d4429a480e25eee0d2498e"),
					9:  h2b("26790fba84e5c62e98cd6557514053487acac49700a078c1d84e825498ae99dc"),
					10: h2b("18d3e145d65526abec8e13eb2863e1d37962f999b8d6d45ebf3d555d97fc5ffc"),
				},
				ProofCommitment: h2b("a43e668c01a329960ba7122e9c5d99940f5e2f801cc4caddbb95280cbbbb73a5a11d58205992f7c2d6e7f4e3c709ec67"),
				ProofResponse:   h2b("62c64dfd54c5e2efc046d965fe89992c2a22080f10b7b0bff2e6882749a8d2f8"),
			},
			5: {
				Coefficients: [][]byte{
					h2b("6c4b57f66952aa94eeefe9c8b97933d5f797955b04b37d56d00f8ccdbbf18eab"),
					h2b("52d8aea46b3bb01f272f0a3a2cfd7e1de411cd15d00deaae1e87abc1f07d4f9d"),
					h2b("4597831eb2da78fe466b4a1f016e4284b8d6b1b00dc285d2ec3495b653201bf4"),
					h2b("3a62f125eb1b1d9f52d3ddb754fff2eb33ac8f5e867f5cb5ac7dc1a1037cc457"),
					h2b("2a9b36dcb6f59b1afaf949ea1caf2d6d53df5d06c26c3de2859507716c5ed09f"),
					h2b("2979971f76843958adc7ef01a20db11b3b0afebeee8ddb29afea5c08448cf682"),
					h2b("5ab7044bda9b0cba987889433b0b576ad4ac5fbfea25fb300ffe84886c9db526"),
				},
				BlindFactor: h2b("8e6115e751b9585e22dc16febbe146155a6dfe6c52f9177cb4657b937a21aba9"),
				DlogR:       h2b("12578469c48458bd2baf595d3943eef7453921787e960e2a588bd8d0de3467b4"),
				Commitment:  h2b("11fbaf68022348b508891aa533ffca7b1bf0f5c14ee2f9c025fc029b50814977"),
				DeCommitment: [][]byte{
					h2b("8ac71e72058e564c1fd2953d1e94de0a4e7ddb7546fc7f9133abd78a129de99ad9f5f9abd37ac3bd607d4c4602edc416"),
					h2b("91de7eb3235336aec65fb2b965bfaf2b4f9b8d0b5e5c264543a60cc4212f7677680a330814b255198c9cb2f5809d77d7"),
					h2b("b0a0ce70f308d503ede3440b811df7f7458d0a7f47de1694ce95a6112a7093a3075c970d32d1270551f50099d15088d7"),
					h2b("95cdabadf2f41a4a1ed815a89f98ea6a6247dcfcabb9c85ba5c7f1f0ddb60b9eca47510291474f82b16ffea44af312b1"),
					h2b("8a95326188a6cf49b2e186bbd535cf86802cf5b9be739724a3a6b3b7d6e69dbacaa6915c09f1f2a2e626748b598d7669"),
					h2b("aae8f1c96e96a315afbf2788aa684280314ed381e67b74c4b4211691250d2c2104ed875d429197cfb9673f120355bc85"),
					h2b("b25895544f4a7422312603a03c928846c054b20c0ccf42eb82dfa1142af3bd14e75d5327c806e6f913f3dcacdee42d53"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("1e33afdad422dd5f23b07de81025bd41dccccef90429eecdccc777ed20953ad6"),
					2:  h2b("178c02e9efefb5c33f22c1b5b49f84db2e9b589b924f0f5a781c62b9844be3d1"),
					3:  h2b("1b1a105a6623e4cd484ed9b4a0cdc88fde4e0657dfbcd8d576276176809a0469"),
					4:  h2b("19ae17a080017065b56faa48be9e21b7e7b535f6dc10a76d53a82742fd169aa0"),
					5:  h2b("65f48029e728fca898548ba36ddacaceb66ed8da6cc4b0cc71740ee006606222"),
					6:  h2b("47a5fd3f0e64672028a790f8c4714dd636f58700327e30a2ebb9e56020a08763"),
					7:  h2b("5d3a05f03616934512ce835dfcbfbba3805a4b45dcbc9a887d1a76b3159ad372"),
					8:  h2b("317835e1f5a01be8be1b87b85ae563afdf95fcb60c08fd495d95dc213e4b506e"),
					9:  h2b("1379a9788a06b32d8ad2505d49bbecd8cf42f519c284c78a1d4d8ca2481166b1"),
					10: h2b("662d878279a4d9d2079e31eb1efd26bef3a288a162f089d17b1b2f1e756872a2"),
				},
				ProofCommitment: h2b("a38e22e05e9fbaff0f49eb2e171396f26d5aad4577f7116ee2d4ef07a94bb28e27dac899138e5fc6e316101ed7aba8e8"),
				ProofResponse:   h2b("3bba4c65ce6c34efb3cea0485cf957ebf1f4a3492e1343502a819b686c29d00d"),
			},
			6: {
				Coefficients: [][]byte{
					h2b("093dbe7e668cf76e8eedc257fa85ae303494206306acc667c1d82e947c445596"),
					h2b("4d3b332ad08f874d0f03d6e41c03306c7c08f7b5262fc7135a1212bc0b834c9c"),
					h2b("397f8d1b12cb6bf471624ca36c7b49865e12c7a89dc699c51b3662591b462f00"),
					h2b("55df2a647048f9e00adbdfb06745116018ce56d100a225c257d9db0772d92791"),
					h2b("1f7578a30e50526c0e56fac19f29175f7188da8e84fe35daebaa9437c91a5c07"),
					h2b("50ea6595332852509104a6453c965b040f65756b775cc9dd5cc1440b5c612a95"),
					h2b("340061a6de5a30d9ad7dcf0f64a3afc0a7050cd432710316feaaecbfd104d62b"),
				},
				BlindFactor: h2b("1dbde427ff6cdb5fdb42f979ae1ef130c03382ae4e118136c8aded80995e2f03"),
				DlogR:       h2b("52eb47d7125c837d86b243f948d235c75d5cac136e34ec8d72418fe5daf418c9"),
				Commitment:  h2b("336fbfdfe449d7d88176bc094a3c6dea732f7550947f88c25cb73ee7d62e2f1e"),
				DeCommitment: [][]byte{
					h2b("8511e42c962171bedddbb6b0188306b993b52c8ed1a1960776e70eb59d8626f1a62d3e4645bbd8e941637aa2c1239304"),
					h2b("a2e976998cedc0cee86718552344cf96b7d38f4420e7f793012cdaebb7f94fa52aa58d6a0f2effeb9033cc283cae42a6"),
					h2b("ab39ff820e50bc598f285c5222b20637382c77df822ebebd3577947436f04b982bbf814d210e2d51392e0e75b3734733"),
					h2b("acdf26578feacf0744efb9576a65c3fb568503d61b46a61b5a02ef643297b389278cd43905f505643f341773a4a664aa"),
					h2b("b939213f4e34d71726b29d465ca2bfba4dc901392594d9c3d21e10ffa9e5fe37b33b7472017c263ca15b7f042574b13f"),
					h2b("afc7af9b7c1ee2db9bcd34ad50250a52747981492cbdc8b8d68c75c30af5b1548c0365ed6c914f406adc5b79541832df"),
					h2b("b09385333859cb8890e25851a806aee050d1bb9c7be2f79a688f1e4f3f44aa35a6c4164a9e89fc54fc93667b18b8cc6d"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("2e6ef30e5d2b424dcd5bad8e0dc6d3975438a756fa163bd4d61143b70c675587"),
					2:  h2b("51fe7e284e5dda0c605fd6c964748571427a6f9ba75e490f9f31aac8f62d84e6"),
					3:  h2b("0fc33ae506524f630bab7512f6790ca3755e2e761f24190eda1aa27e398576ad"),
					4:  h2b("4715546efd4f8686cc4ae4f7f45f3aae536852263718f65d326dcdb501675ccf"),
					5:  h2b("03b78f85c2525034fe83b1080301137b0774edc32dddf00d5ef5c1009b62d8e6"),
					6:  h2b("2de5e944b1d8ac19d3f10bf2ecdcde96529a664a2eb3808e7ec86ce44cc7249d"),
					7:  h2b("6bafef6e1867061b19176b54a8194f388bdda884b52c6e7eb72364dc056591b8"),
					8:  h2b("11e68206a78b2d4c17db83f09865039aa77f480ed0bc87641202031aefee61d2"),
					9:  h2b("6d0b124e43ac14d96381972bd2e64b784f6faf6b482d0c519d6d6905dfe7f5cb"),
					10: h2b("3b8a9b3d3d938fd8f423e135b6da055f2994d8ef9d14588acb865c7b9d4054d4"),
				},
				ProofCommitment: h2b("afee2319814310578a5586d2f5ef12e879b6d9ac4cb3b42801e716df4245593a2942c381254f49cad307c2df93f7f3a5"),
				ProofResponse:   h2b("142a4a55a1ff4f83335b36411e2690bfe4af55010c784386e0b5fa418b868d91"),
			},
			7: {
				Coefficients: [][]byte{
					h2b("6c37681b97889aecb32b031666c7da730cccda81eb0b4e20a505d2de4a2e6720"),
					h2b("2ba6ecda3ae0eb9a4b36c1f092236a7adad8e79e534275b32e237f6b521e1b17"),
					h2b("6124a3d81deb44bfb39cd522275cea513218db9de68a68055424d8da7aa616dc"),
					h2b("454b23f00cf0d881273785e7b0dcc972a1e62cb27011430a51192f1c75928f0a"),
					h2b("4d4cbbe4628146ee0f4a9fffc16c5c455e5fe2ffdd9abec904b0e9b6c3e3ab91"),
					h2b("2bb367ad51c9198f8bf24363ba670e65c086c65d9a967584c85cc2507bbe31ba"),
					h2b("2c95f61aeaeff53856056364a28422b1960c7f2e1dc1b70774d69b930b831e23"),
				},
				BlindFactor: h2b("69e544b1128d29131a5475b5f972f72c1d7926afb6fa2c16195082e7787fdc8a"),
				DlogR:       h2b("10290e9020f1fe90eaaea784fd8769189f8bf45fee2a33f40d434568764810b2"),
				Commitment:  h2b("1316028c45d82b7ed1bc5a6bd9b4473dcc62cb6437054365d5d2c594fd43b43d"),
				DeCommitment: [][]byte{
					h2b("a22b626af5c1f352882d221d18baa1a1f220579a527885fb53adec588bbce56c8cebb53cc2af6aa9d2f4bb0e9acb9a2c"),
					h2b("8e3fd4aa3cdfcc12d177fd05142e50237ce5a8c15e6b8893ff5ffbaa15f3cd047d4db93ba737bd2c736fc8449a62bc64"),
					h2b("985b7535d071b5622ea4ed1b7665d0453318b01fe127072e165aa3c5e30e79f7f2b77be30d1a03df376964560d00b9e7"),
					h2b("8fc6a5a2c12ca45a438f8249c7baa3c61a4ff419f158471786cbc8c1a26fcc944e657f3e17e08b7961a9216f55caf183"),
					h2b("914623adb2cc179c7c1931daec11a22ee51a3743b52c03f9a9d85e0f8286581a6ea371e06496ac5a99ad4ff0a7181f9f"),
					h2b("844593fb7bfa238f9cde8fa8263726731ccb8efd1c2b7326e665ffd7398e840471fa38dbec37eeb87f2d0aae09502dc0"),
					h2b("897779d9dbd9baf5ef57741a64ff016a404565736257a0453ac0f795dd6ea46b1ea5a3d656f312964fa3b8389649063b"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("142d991df60a045cfd9106b8c8f525f921a162f02ae2ea3cba4ba1ded7aa2387"),
					2:  h2b("133ce9d7c5a252fdde32368f8f0fc8866ed0c0574989d52e66f7787a1c5fe9e5"),
					3:  h2b("5451a381347fe9e76c30f18b1c02acafd24157946bd07624158b55872fc7e886"),
					4:  h2b("0e4bd140e8279cb4e5a9fb186718de5e44de35fbc3142bd473d95745b000a432"),
					5:  h2b("42f7d4e33c9aa478a5c1367b829a8becaf6150e642864205e374b5dcb4e22a99"),
					6:  h2b("50939b1121f65758024f34184cb4e4b2e56fab8d0721753085b84c8c23d6060d"),
					7:  h2b("5c036d78d28c5c24450f817a451cc588484a0b8072235c20e362b1766473f2a4"),
					8:  h2b("1137a915f78813048c910372bd257ace3283c77bf62c2bb02fc7dc1775e354ad"),
					9:  h2b("10c27f2582be0d2362122abde925684cc11ea8c096eddc8427985a5b65017086"),
					10: h2b("278c418c20a636b56a590ac655f09e9cb96a8acf1b8006e38b3e1461234c63c6"),
				},
				ProofCommitment: h2b("928fd48978a48c45f1a6ebd7c0126d25ec26a39c4d7d7e7188fe3b215412e87cad93838d8fba385cfaa98cc5dd7613df"),
				ProofResponse:   h2b("2d3214f7a5bf894dda2754e7da99fd04fb6fdcfc1835b55fad4bbd68106b67d9"),
			},
			8: {
				Coefficients: [][]byte{
					h2b("049f7e4bae32401759158e53c72131b416149b25dbd82e8a4c42f9e3fb295120"),
					h2b("39409abcdee5165ddfd9ec61bbf592c5725a5eadaf3de6a3eb859ff8704c5f27"),
					h2b("60479ac3ec30daf61a151bb17985d402b20e4de080fb8f5678fbc735baf67cd9"),
					h2b("666472d98d037049dbeb3deab0563eea828615289864c31deadb4cb8911b88b1"),
					h2b("3a6e1e354063f15d4b25e5fd1f4f8702fe4e7fa6401d592143321de5671fb239"),
					h2b("1534483ae7307b196ded917635aa3527e465b2e7e13e88a37ad8cfb408e37e5f"),
					h2b("4c2ff72a52a68c4b5d00b453e0e7c7a0f92517a5642611026c55ce24e743cb80"),
				},
				BlindFactor: h2b("4fab8d4264817ecb7d6e5494ef8eaec5f9b55436ba667b9b3db1280ee58c1681"),
				DlogR:       h2b("183572957cc6c90105beff11f07ee23c547734528a5f9b78a57e1a34e244c9f0"),
				Commitment:  h2b("93dba508221f91ea653f63091ed1e057fde03944d9a2b49034226332887bf09e"),
				DeCommitment: [][]byte{
					h2b("8b037b6c362ca6b9470fa9508b42680d1fa5ae9a756780d4c2d6a285969a9d72ba044aaa2e4e3356aba103987ff7a025"),
					h2b("86606c04d70a0dc2a48210caf04819d69fb7179476b7cfcfd453380234040cb52d36cc08cd7f6a0bc70f028a2404a788"),
					h2b("aa642f2a65d8ead5318e6a7068d12b7615a83bcb5c9f1a2e425ac16ef0098cefddd6480bbe52851174b632fa33085260"),
					h2b("a170f34c6a826413e3c1925ba26faa8d367b3598edabe255852055816dff3f8b9db351fd704c4583f74d851e14d0fdcb"),
					h2b("a4c8b2c7c0fdb02c69b4b64cbf24e9585fa437d1f9aaa86e3bd00de54ceb0e9f3c8ef60470f21ac0ea511003151965f5"),
					h2b("abd11ce1149e073f645caa01688a2365e725df3a8fc4e0744b3016b5d0fa2fa31cd601f43a8adb2df0fddb86e66aed14"),
					h2b("b548a8f502c1fa093fa3829a2451800c46b7d4e7fda143b6168e76a0553764b534df9ae6960d8da1451185d85d754152"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("44958e4703ae229eab567800c5eed3229da3bb0729fd466cc600698c0eceb1e6"),
					2:  h2b("2d98a02336951ceea891dea906ab12a3a7758d383501257d01c71ac4afd61787"),
					3:  h2b("18559daf9a9a9e8695cc64dfb023e2a1e888c00bd0e94d9773cf990dcda71242"),
					4:  h2b("02a4f15c0d8787920c2068eef68b29a0571095fa4e0287f833282b4c9d087084"),
					5:  h2b("6b84f79de9ea9471b43bf4415d692fbc82c54b22418e3b994aa7a85cd77252bd"),
					6:  h2b("1861db6f66e875d98764c7a175fb76f0a1f37652fc931c41bb34154c690c9555"),
					7:  h2b("462997bbcf16dc70c9c13cf8e446b68531d8d9dfb0aebe302f5d0d2c8d5990d6"),
					8:  h2b("54898c00311551a5237125856cff21025edc7fd8431e32826049f2995a8d302a"),
					9:  h2b("50e5dd2c10ce8147a625559553c5a55653d1ea0dcdc5e4372bfbead5bb905d0f"),
					10: h2b("6a9e419c85d92fbf583d4d1a35d19abb3dfb3da0ce5ccde15ce3a291d8b0c2ab"),
				},
				ProofCommitment: h2b("a313aa20d00c42801c2e8d5e6b530c97c87b46417e1681ec80a4375bbcc446b15afeb079c7ec14fbb34510989aad621d"),
				ProofResponse:   h2b("54020d473e0acb1183681c6206e9492170144feba696dadc2b6bd06d620f8d97"),
			},
			9: {
				Coefficients: [][]byte{
					h2b("373541d3fbb8f566289dbd95f06553bf93d01f50f8ff9c2cc92e1060bb8b8167"),
					h2b("4635ff8265648950794d2210d3f03a74381719a53ac3500bb7d1490f4524baa6"),
					h2b("463961c4dbb68f51d449d0bdfbc0640c971e88328db64a3e9361756cf114963f"),
					h2b("207a531983fc88cba2ec0058d85752deaf4956c57db07720bd8db1e26ac06f7c"),
					h2b("3d87c29e02aa9c0bf446d6208c2159f4f90744adc0d2d24f1e0c0a9a5126f778"),
					h2b("6ef4ff1ffead05c337b5483ff6e42640a00e4ef1175d0540fee02ced65e87b82"),
					h2b("3f331b2e7c7b7581b76bd90b64b286c9383023443c8e79c51e8d58e56d5e0fef"),
				},
				BlindFactor: h2b("278b1a6bf2f2fa963c41d3efd8614d993819c92b095bfdf13e20f878815d6270"),
				DlogR:       h2b("222c8c131e5639aa246f89f0cc7fb94a4378d521c19c08bbaa82b2c39d77a673"),
				Commitment:  h2b("d4e420fb772f141affd30bcdb19c496a2eaaf15ee496533ac48beb001407fa55"),
				DeCommitment: [][]byte{
					h2b("889db2a6a7389a1e7f5368768d66f422025d114e1136d9f6b2e35155daaee2e3d49151a9c91edb09e204dd0e0aa20bd7"),
					h2b("b96e54636c3186a03f5809e73e537a88130fe46a376e9656d013c20d1d0c889909e758d4d56321bac668bb242b1f10d8"),
					h2b("8ecc52f5f4a0521a43d8501fa27f5a0606e71962d8cc80b637b719e3a40597139169b5e81b47d01606ca00629c10590f"),
					h2b("a0d8d6ca5955a39d3b75204ea513c837d419d9afc7e0ebe95fdc01932aa48b60428481e84ac1cb2f64397d7d8c044331"),
					h2b("a718d52501fecca57e712e0d7968991b4262adedc72662de0b906484b6fa9bce66ade94d9feebd171d8560042b60c7f0"),
					h2b("abfa1b730b674ee89d9abb74521f0a2391525bf19c07a7db1cdf8df36cdd74e139b6b97b33ff157433bb8e4f3ee10703"),
					h2b("a59722da841c55e5f445f950fffa331f4f96006d19e1361c41bdcf65da6a47ec83b8d63087e67f6083840141048437e1"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("001835d4982db9042fa14809599dec08949e3ec553ee8ef10d68113080f2c4ad"),
					2:  h2b("29f7f03a73313eaf97b871d412aa87d74dc6e2afafd29cedd2e08843872daebf"),
					3:  h2b("231beb360b7ccd220c2754b0a5ba4f94b531abef8c70c5b58e1cf77420dba1e9"),
					4:  h2b("35d299a31824d6fbc278cc8ac92a72b1ab039cbc6fbaa6a9a443b7f93b6811bc"),
					5:  h2b("3d62cadf86ad11635c740315c7f16de693acf891aa3c33bbb60a1b2221646c2b"),
					6:  h2b("6687e75581a5260c05b8ff3f86d9ff4e0dedbf786ebc7c18d82c91d9d56ed154"),
					7:  h2b("5c0bcb1f32e216603335783fd97f083026e3a4983d1e9ac6b970dafa05a599ed"),
					8:  h2b("0afdd20a64a2bcca6b06d47c6a7544a35650fc22a077572bb632386697a7ac5c"),
					9:  h2b("2ed077af5b2976fb47e1f421101b78edf426108e405e0585d977aaf5cd21a076"),
					10: h2b("61b946fb16b786c3b7df8f445de7dc8d8498ab124580db55cb94352a00e7b1e0"),
				},
				ProofCommitment: h2b("9682c52d869298eb4bae7ddf0f213827ca177fe9c719336c3ab545ccb63f70c5250aaec8a36f021f8f104f1eef8a02fe"),
				ProofResponse:   h2b("37b862893d62d8ea6ebe2d48e1ce84ee76669646c48323e333df9f2f5a3d7dfa"),
			},
			10: {
				Coefficients: [][]byte{
					h2b("57fed304f17d9bc428eb4baca47f01b1ea1128cdf9c9ea789099d81d135a7d13"),
					h2b("2f57ab45bed40606864bf258d57998b6917fb699751bf6f68c7d96312d3ac70e"),
					h2b("2fa026e45604506625e89a8dbb801f11610248828e616ccc4e04dabc4a6e729a"),
					h2b("4c96d84f302af061ce0fc3d852bf018bf1a9318d084fe7f37da2c27250cfeea8"),
					h2b("6c8fe4147f102aa99b42093ec27b7658331cfe47c819fc46191e2ccd61a8e84f"),
					h2b("4806c132f15c2225813ed424af025ea7b64ca54e0557bf9ca8aed9d7c751ea90"),
					h2b("3a06b7f64254aee899c6deff14a881bd4a2adcdddf4c307a4abf69cc144dd8ba"),
				},
				BlindFactor: h2b("649c13b329e73c32bc4476069e450d5d43fcbec71e457382ee2f9fd54d65a3fb"),
				DlogR:       h2b("5260911ae752b5b8e95f637b0c422cc52c784c99cf2899c1530c5105eed40c31"),
				Commitment:  h2b("24c0361f5bf6e8b5bc1d569ad14d079dce49de033f8b51d710dc0752a42cb447"),
				DeCommitment: [][]byte{
					h2b("91ee90abd429b5738459c8472b4f157e7087e25d0ad3420b31352c791d70a46023ba48517c95aa7497b4c7c2be92adeb"),
					h2b("8763aaf2a072a20c66031e1dcafc29212058ee54dbf9f9c86975a01ae3f0310eeaf47208e1b100cdc1c2593d5364c7ef"),
					h2b("9686f9dc349a864132a9991d9e974a2e0e797a15a97484c751d1dd33b79c956b2d69ab3c07aa4c57a5ca27d9322aab18"),
					h2b("a1cd85bede9b160bc2ad939942fef439e14e4dbeec18d128e660ce3c32dedecbc1b6905598f78e66e888e408f97d1ac5"),
					h2b("b8641e889b63a2d27bb286abee8eb9795c7e2926db81f73932b4069c14c4de5c95ebb6ec7a3df883767465321a9713e0"),
					h2b("afd3a83fda686fa1b1309aecf7ed0b7dc76535adad1f3082459db0e4ceea3e41e8d889411155dfd7eae53badc4ce7335"),
					h2b("a6d6c4a1a78e60125a6e3cba3b0d2734139cc5a6c68e9b17e4dcd3577102c2d975713c1da3993a6fe9b7c02a24070ee2"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("22743d6f42cbe9298c8ff8ade7d6b1adb2da49deb25bb28ff54b7bf2191c50f8"),
					2:  h2b("2ec9e774c91d61c18964474433241c011beffd476025ae072656fe22364b4ffc"),
					3:  h2b("1bb574888107e0ee72898e43ceedd78ac2e461728162d7c120020ab8a9bc5451"),
					4:  h2b("2cd778803302b7551c3641516c236e5ad792058e17611cffdbbf729f7146904a"),
					5:  h2b("2b3a4742fe5be216fc7e04e8a1beb2e822edcab23d9c9682f4052f857b8559d9"),
					6:  h2b("4239796c9356f7d58efc849f1f41923d32ade0dfc96b2afc827bf721f48b7fc9"),
					7:  h2b("3964eb4bcf89e13d61944b644701a2f30bd8bfbef245bae438885ddcad8828ac"),
					8:  h2b("36dec982dfc5ca1230d685f573ae591eae6b973302a359a1b22d88d49f5d3b84"),
					9:  h2b("52a2bf55eddb5b0ae0b89f56af46cf1173a65fb111706c1403496f4a88275221"),
					10: h2b("567f3ca1d712b455eb453f67f864aafc42367662c41cc0737f2aab6ca3b73537"),
				},
				ProofCommitment: h2b("a439c7dbfd3e1edd4173e59944eb08faae50dab4be7840b07b5b573a5861aac1bd0d4c31a3965d69b4beebd15c098298"),
				ProofResponse:   h2b("6e8ee8354c4b29848b85d9bce1a88ceee3be97c06b88e68ef760380059715873"),
			},
		},
	}
}

func TestSuiteThirteenOperators() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("8133a79b4509213e9baf557b72f112d98276b95b1e72d44a9e26604ca6dd6a56ee6376a456dde1e9610ab153ec34c153"),
		SecretShares: map[types.OperatorID][]byte{
			1:  h2b("1a8b38f084a6fe5369493c01c242e44fc5474f98ffe3de0883f93cd2d30f3a1c"),
			2:  h2b("6bfe03936c70f08a967d771cc217ba6eedb1dcb05f64b9a29260670ddbb78190"),
			3:  h2b("404cdb9e786f62d8f1ce47e641738eb80ebc3681204e0bc74ebd26b1a7cc3cd4"),
			4:  h2b("44fd392e739b4f334533ca666d1c968201ef267c9554bc63f105a02a2f33e534"),
			5:  h2b("594b45f4a24579f4ce01715ebf309755eb86be8da90e75462223cb448b246832"),
			6:  h2b("122e5cd1726bf3b9748cebed55cbc583f85e7c579369f12c6f145c7da6dfbdd1"),
			7:  h2b("4256f3ff32cfcaf70984ab97485c8c2977478ebbe110b6a2d9fbc4c789328e55"),
			8:  h2b("5d9c5126f89b30f25c3e3739c94b791078f60f3fd058389a2185ad3be42dd9ef"),
			9:  h2b("4af2e122fc4c575a356082588c43e79dcc9af69c1d029c557d518e7171f27062"),
			10: h2b("6c2a5a2399bb78bade2e613c353b3d7f3e7e2cb86a110636b49d3da77ec2e205"),
			11: h2b("1f7ad422a61eab26e133bcc9a4dc6c436cfac9d68e9dd4ad9acc8449d9ee6f0c"),
			12: h2b("1fe712cdd921d1c31be92b58ea9fc36bfc15682c206ea25023db10a140855585"),
			13: h2b("457bd7a93c86295d3c3524df03ec0612a69d4624acec420a693449002e24b9cf"),
		},
		SharePublicKeys: map[types.OperatorID][]byte{
			1:  h2b("aeab88b11aefb5b450e0f30eb8105f3cc418cf69bae4c2651c50726fdc97c2f1478b3c7b7f23e59ecc07d80ca1aa7d1c"),
			2:  h2b("9933e66c66cde885cc98015ae3306033615f9a3b696a8972474b9028798e5905fdee55424c449d7638c684e74c8d26b8"),
			3:  h2b("886202280b9f9432173eb4019ad2230ef70f13aa1933b9f64ffee00dd22246989978f248cfde76093a1bf7e51dd0e920"),
			4:  h2b("a738b16e037b0b385b4e176e3a067f6247437b6a40faf453d52c7f1a716117c683e85e60688246f3d1829c33fbc405b5"),
			5:  h2b("8c1bb2098333faa172d352e130f4966b18c87ec7a443c6ad43ecb80deaa73e5d8bf64329f23037e7a7755eda47292cbd"),
			6:  h2b("948fa29e211c242563e2bfbc862fae99e44527fbe658076bf220252eec14e719f5d84cfaefd87b5a93fa8ef7398f40ea"),
			7:  h2b("a5e8cdea8e9b9a0f47200610f9d866ef740061c02c1f812cfd45e71963c7d47e414a687f8d654364b7c0cd2caa7e4eeb"),
			8:  h2b("839d1d578fbf334a7dd9b26ff1872b7779dc6a39dcba6d02612f8b06a2566b09202030b90db027519b08f26f75b8c1ee"),
			9:  h2b("82455da3416b4890f07cf034efa64dea3aaf166fe017154a0f1cdabb4b7471c61e588d6d23f8befe6d78260921dca60b"),
			10: h2b("8ca66923ae11310ba4d6511e0377c3bad5bb5daa93c66dd6eeff20684d611d6b6aace0493526a4736e9006a37fd45b78"),
			11: h2b("82290b906f2859657106a80d065703037991cca4d7a1d98812a41e8e449fa4b6bbc2ecd07c3cfcea004dc997449df20b"),
			12: h2b("b1915b7cffd26ad081e1076ed57e918864ad8c22892be3e321b1927794f1cc7795f3b89ee7b35bb5543a120f116863e0"),
			13: h2b("85b7d282d8a52bc21f704023e4dca55e518cbae3cc55146c329c1fbce47d5eced668d6985841096253a37771acf3928d"),
		},
		PartyData: map[types.OperatorID]*DkgPartyData{
			1: {
				Coefficients: [][]byte{
					h2b("4cfe717ea858335f9ef0a8a343bf47a1bafd361b7805c011eea9e41ede18f062"),
					h2b("59a7cef35960c686441080e9cf7f57935769f927c2c751ed40ba60e7c9319e15"),
					h2b("0ce8eb8a5cf456028be3a77e80a36b062ed0720f8c47200626095865ecb350c9"),
					h2b("63497f596d89ca4e015fd37fc295f006df56e3cc4747153144bb1b4726d45b63"),
					h2b("13010dc14fa983e52ca640af7b61679b4af59037ad193bebd634cf58dc6d1bf6"),
					h2b("3eddbf211fc97e605d648ac9efbcb2d6556ab0ac9667d818b1be5893ef3b0172"),
					h2b("168d8ea7e00ef19da3505428191666612303352b76734503e8e35d9a663a7756"),
					h2b("1844c8e60971bd55c844a7bed863f2a03aa7f4729dbfb7856d74c2dbbc3e22fa"),
					h2b("2839dabeaa4e07cf75f2e8b98a7dafc1e5b744b7425214c546ef5be57021dab4"),
				},
				BlindFactor: h2b("a3ff9c91a47446ff4793937690d45dc1dcc79ff4ac906cfab8178333a5307de8"),
				DlogR:       h2b("19899e212f108b83696342570c18a691b5f43dd12d3ee0dd3d60c1826dbc28e1"),
				Commitment:  h2b("21c29dfab217f4dec1f57bc2f827fea461126dae0d7603fdf19a2cc531ce24a7"),
				DeCommitment: [][]byte{
					h2b("903fdbf5287e829666771b283bbc789c925b1c294be7931d1c294084ca81cf8c0e281eecad21b6869b667615e375a743"),
					h2b("9041875c00fecb11d4d5e96480abaffd384331e4b0b0988faa4ad541751af5437503f093bdcdec952e5ecbd030f365c1"),
					h2b("a3bda79f36e054bc488fddccbf94945eb105c0b64360c314d6530c3eb3f42319351b49915af6676e2d89e28f1fbd7a2c"),
					h2b("ad4e3f44162da5123dfb7198258b33405f02695ee7c0132e765d25b450f03c987deb395193a8d42e7374405868cb33f9"),
					h2b("b5c431f2f57c7905092d2de9c8fe26aaf5d4b7170aeea4b09d10e1a1622e5aae2f0e3892bd1a737a260802bfba1d6a57"),
					h2b("b29941835af99910455f580ce4ff39ee33a9a91f804b0d8714fcdea0e29e87d72f39d92e6bed37454a43f21138c426b7"),
					h2b("ac5975d32ba2f8a6409ac200ee7756ffd97a09842fbd915db6cc4f40825546d6b4b4264ee71554e29db180ffe133b053"),
					h2b("ae3d81b47fb908e2b06469613493726aac0abd7b234ce4d90f06de5f3d497efea2fb2c2191214100c9a8e19cf2b0af33"),
					h2b("8c08ab6ab959b07426983abed6161e9267a8a3bbc891ea974be3e4ccbc7c1d73f597ed99af7df28655c3f7785fd29588"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("63fab48b52a05b664229cc8d20a895670918484fa866588cbf635cff1914cd0c"),
					2:  h2b("40e9b639b40e81c0576c6eb3b4a9269eed37cdb0d73bacf7abc9a44fd7a83e4b"),
					3:  h2b("30a2d6e2df74a25732579f9716a0ae47e1b8910b7d2c803a05c1297b8ee2321f"),
					4:  h2b("4e29c704eb38b3257c027a36387c5467c974d63db9f68fda7197d9fce734801e"),
					5:  h2b("0a5953ce5440fae333b189a2cc9c62c0897b75c855bad2b615db49232f95396e"),
					6:  h2b("5db1aea24d71cc03c47438e8405c0a26f4f14362cfc0ebaec949dd005a419e44"),
					7:  h2b("30fc01a6cf6f7173010d16de4b1f0f97afe01f39a90bf7a3143dd7b870a2a84f"),
					8:  h2b("1cc51ec4ad82c630c640c06f471e70d600bdf166f23daf9e1dfaf35e3fe0a0ff"),
					9:  h2b("3299c96a2a74bcb8bc6cd97641fededc2ef6c5dffe159232e15125e6d3c0c1b7"),
					10: h2b("163f16b7f960c6d6a001854869112d155fcc1efbf2fc9891471101b129756382"),
					11: h2b("016341add036f5c498999ec29252db99a3c53fd0cfa201c407d806585816cb27"),
					12: h2b("3b289e05020fde0a13dac711777df24adee906f75022985379c514ac408825cf"),
					13: h2b("031a2ac0880462fc86f5aa196d02e28c735f31b6fa98484cb1af17c8aa99d1d1"),
				},
				ProofCommitment: h2b("85d2cb4d6292e2f7b264c86ba5e4aebad877f4764961704a43363bccb2d51aec5971466f87acc477834ee3e45e560cd1"),
				ProofResponse:   h2b("2f3db29b6a461bbac3e8fbb4a0d9f8cc0dc08bcd59316ed45c8e526e05229c8e"),
			},
			2: {
				Coefficients: [][]byte{
					h2b("3375993c89e6763f10dc80e68d15dd6d4448c9a5d282ab322606fc18d575e5ad"),
					h2b("394bf6323f0267c3c794ae444fe3fac81f489ad033f6e255117ae3d589f5d27c"),
					h2b("08d4ace189ca0f550d7aec930911dfd4b958b2b73de17435fbeea51963fd6779"),
					h2b("6e142357489b598f1b421e81994b40b1e0fb7e7b752be1e3d8f7013e9e427e8c"),
					h2b("4d1fb2ebe311e6dc89cd4d0226c745c0ea51c12b83445e190d836b3db4abbfeb"),
					h2b("3f0d382e8c59a1067531ec1d290fabfb057e9d87322ccc0481d54a885d7546cb"),
					h2b("2be7a2107533abf6971e4f53acc9f9aae4edc90f8d7076722f68e3e5d14a80e0"),
					h2b("3332f5d7e33c14adb652f6654a00b28ec853e3845882386b3a74a3902850119c"),
					h2b("2e658219289ec783c3ba7d48e46e418cc8734cc5e029cf5c5ba7fed4661146c4"),
				},
				BlindFactor: h2b("8cd4aff31b851c830f2f4e24b343fc8fd3841ed778d2cac1ed772f7a26f73d4b"),
				DlogR:       h2b("52b9ca7c10158be78f2234cda8cf567f04274e2d0c8f6380fda9154f074b982c"),
				Commitment:  h2b("d715efa79321f697adebdab919b0be1a35845be8daa718d29c14b0d384f1da3d"),
				DeCommitment: [][]byte{
					h2b("b1a32612f268ca6bdbef2c5c8f98e514473936d24648ebcf1a5c3f5c71f3fc122ba63c1b33b4cfa5f95916cc819648cf"),
					h2b("86da56e029194aaa42231139784223687f600ef9066ed8f0ae9283bfaf0a1c18fd6ce63724a9cd7d79d444f4c1f95212"),
					h2b("95e8b4b22663b0dd32820418fb4e3aa6ec8975de4e2292757f61ba4093499da4f42709a389d34b3c9baa3bf51f5b8269"),
					h2b("a5b07a80fef3cdc63e9e8877e5da163c71465621a77bb778a29ecfbae425ce17303c3a9e96d78eca22fd1c99688294f4"),
					h2b("914da647e2a9d7ca37b3391d225633f26f29e1539404c85265e461ece8f8d939ad866cb546af1901250b3205f29037fc"),
					h2b("af16cd218056746a125ddadb88ec8878d0be0ba43ff97db7a50933f63974f3091062c2abf05cce0eed4e105275788807"),
					h2b("80b1b8f4aa3edbec1784486b333c90a6090037cff2fd5e7e2a736947eb5001f81bdbad5a4a63a532ca24a67f12baa7f4"),
					h2b("8e1d74a503020d087eb1c4a332ee0ffe0b0500c4e44c8db41e3d74c6f3f856564bea51256a1551b759044ceb57ebb52e"),
					h2b("af7eabe508e6c902cabe6ab020be11004650de37abd5509f76621322ac864819c1fc51442b2cb0a201c37118f09ce4ce"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("2da0c776e55261d14471d64083df782914745da9351b1bfc6145c25ad3787e20"),
					2:  h2b("464d4630f6a60c6398eb811c03eccda0a6c05af91da6d154cfd8fdcff0bfbe1d"),
					3:  h2b("0a9331b457b1967c037ba3794e74a5ee96c07c102f8a76840c9fa086cae34818"),
					4:  h2b("335dfbbc060aa9d0848affd9800eb70d1a05b7090d54d10bf729e4f041d6d54e"),
					5:  h2b("5d286e4c7bb4651b06cd9f330f2dc67473919cbf981df4dd1ce55aebb87e3946"),
					6:  h2b("5885965acf39b21d451af4a48a180c4715ddbc7bcd2c04c053dc18cc2cc82b0a"),
					7:  h2b("56fb6e7ac5227e8496d4a36fb9fb0be845f7649745a20fbbded6f893834b2afd"),
					8:  h2b("16ec73b1fdd2173de9051a026526939bfe0a60fe7c69e59d4dfab56b03693335"),
					9:  h2b("40cf0328058cb789c3d5e2ecd861b27b6a9910d31efc96372a3fa92da208373e"),
					10: h2b("715015926d13758aa19449f052c3a12b79413c8ace4ad061556bd5c059f6cf0f"),
					11: h2b("16240d01c2a5d423433caa031464bfd4c3e3fa82d55a031de37ed1af721dc00b"),
					12: h2b("0094a72ecd86d7ad47554e79cd54d1039265a41a7bc700c30dc71b6e32a86dcd"),
					13: h2b("3139031b15716d6d7aa25f40a49496835d3ae5a5bae56d7cbb2042702b579370"),
				},
				ProofCommitment: h2b("8a553a6cf341e65a7ff8b634697cfead71f3b25e0cf544c1b61a79f97a1bf7edd940035a72736eb094280de7539b56fc"),
				ProofResponse:   h2b("20fc9191fcf77e1252c8eb896f3e37f1a6cf6913c4391ebc368f8bcd1a9775b7"),
			},
			3: {
				Coefficients: [][]byte{
					h2b("71dbb1ab10317331f2b64d0a3e8a4c67fb9e7261ac977fc778999f01691f64b7"),
					h2b("047414daf1ac75e2eb724709d6b44b703d2106546193e1a53e233a3c0fdb1d0a"),
					h2b("4825965c6e8711e996dac1fde6e43ba7b1d264f680a804b4a8c82ae73e8229dd"),
					h2b("6af0f82b529bf53426d0d8560188c1e3b82f591c93f7385dec63778f43a83c99"),
					h2b("2ebedff4f6042110009988fb74468127f652ba36ceea8d0c9945d4e8fdf00801"),
					h2b("0ff58f15d5026946e6f7201ed512ab52b345460d413a76489cb8a147ef2091dd"),
					h2b("6662454faec679db4922e03bbb3a5dead589201599269b0bb59b373c0b0eaf16"),
					h2b("2aac8f0cdb583834012b69c15a8697641213b725af3cf538aee78ae5063ec671"),
					h2b("498bbcebf57e854858c8d070eec1ab77c773d10ab05dcf3ad52cbf6f357bb311"),
				},
				BlindFactor: h2b("a1449485366b5119b182be2d74cda6d801fbf0a99b8990e8c83634db6f7d7712"),
				DlogR:       h2b("58e427e5fb0508a8a61aa82022426be71dbe56cd5e5d977241e7f16f3a523ff4"),
				Commitment:  h2b("85c67914127efd9758e34abcd67ba8b5caf34fb00d8e9b99b57cdb29081fd8b3"),
				DeCommitment: [][]byte{
					h2b("8814e0f2b46247052a80e682e7df02a5f20ea938f62dc22adf4c766bed47d104433d387587de83b1fb04fc60674c1cd8"),
					h2b("852c2285d32a37ccb3268e5ecd3c7a9d6fa1780afbe6fa44a60927072794f313bd02d8d3ef6f71ecd8040d3ad3838c0c"),
					h2b("a20ffff4f13ec3d57131a63b81e679f2d183775f97470600165aec3adebddda2036778a2122f443e13d92896c90cfe76"),
					h2b("9325954098a44f9e7cb47cb1ba0c96f0b2049093d7431261c3b7211234055139b12232a056e123856232ce14b82293a1"),
					h2b("a50d054aeda81b282ded3af01286e0b17074d95921ed18a7d769f867f5503ad213a24cacddc949c6ac779ed81fcc2a9e"),
					h2b("ad5d276c9ed734a7725a66652340e831fe7691cc18d1d79101d6f596383c964bc251367c6e1a4b3cf29166885044b236"),
					h2b("b8dcfb77ecdc7710e9b81ab145ef551351bd69a7f08d27ebd931faf4f7f662a5b0f4a2835f48942de33edffd42684e39"),
					h2b("914fb221b109cab1be76a865fa302fa971b841a1432a94bccec1bc32c65cc0c899df2a6661405b1be110cb4ae606d6ba"),
					h2b("8d1b86af65dd103207c5d7b828310a6373ae2bebfd6f8eaeeb341600d3a4265ff2144fae0ff75bce2cd47a87d54abf1f"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("72feb814672ebcc0599491d02500028fac734f472bb79157bb9673792efeaaa9"),
					2:  h2b("60a9084a77c5ef567d6d68df81a4d467cad3046f6fd829ea2de0a4f8c1f4f498"),
					3:  h2b("6dd79021f0eb59fd0f06e603e714ef6f94196e78ba8671961dca14f0d141a33b"),
					4:  h2b("2df65d9a499da1020af0bc44c018612cc5847bdbd02249c1947e24003caae89c"),
					5:  h2b("4f0d665422e49bad46cc0efec120909e65167504f120bcb9cb3395e185f4da86"),
					6:  h2b("022a5e3eae9ce69212cd9cdd06f992bd9b1e02f1e7eac72f1558160555f78022"),
					7:  h2b("7134f457b85cad1a1dfb4d0926e91b4fa6683fe775709bf0496d768fa68ba1d3"),
					8:  h2b("5316e9e96280d534250ba50125340f1bd6b9e0efc1cebc89ba1812f5bd30b3c0"),
					9:  h2b("2c31ac8ec9a4977f838a476685421aecef67cf5dd0e0f73ad8a50e8abcf8c553"),
					10: h2b("3bfcc3a59bb06218e4ac99b0bd38763ff9a60c8155fed50da21676e130f1f143"),
					11: h2b("56de207e3607ef17ad60c810fb05d681c884c2109cb5f8b9ff8aa1078aef5b87"),
					12: h2b("1abe689086d1c92079c1ddee0dd3dc36c9e0a54d9767a9274a9b680e2d3e5c0f"),
					13: h2b("565cfeb5f96100632fd5f1512fde5072edd962bb7f93306b2e194293307f06ba"),
				},
				ProofCommitment: h2b("936e92b359d84cfdbc30ca333f1b8cc88190fc745d253bffc82dda0fe1b9ba57d83fe28ae6326861811352a358c03409"),
				ProofResponse:   h2b("684571f684ac4d74f29a011b9b7fe54057302cea954f36214f3770825567f329"),
			},
			4: {
				Coefficients: [][]byte{
					h2b("41f2250f498ecc7b75250a05eea39bb03ed129dceb3979743b752479e379bc6c"),
					h2b("59deb3baca005d68b62168037c7b8e4053bfa6ee2924c7bbb0d96f3bdf7a3c58"),
					h2b("1e19f2f8283de49cc2983df91a4c7dba7cdaf458dabb8c4e982bf7e14e8c406d"),
					h2b("65294ebc4b187494299c56998a053bc0e0ee48b60450b07339d810905583a03d"),
					h2b("2e6ce0605038ff0acb06ed6412bfe75f68b0575f40aa6866b951f2c8c709332d"),
					h2b("3e2771fe38f033a500156d33717852c5cb8eddfcdc22948cb04870bdeca7a123"),
					h2b("50a23b4a038b672448cca3d4b3006e91ee81aa231d1c9561178628d5bc90da16"),
					h2b("473719da38b6e06dbf8d143cb9b6c6711ceb536288b4e24f9f26cfc556b097b6"),
					h2b("1972201e7334809272e55f48ed7df8ae620c5869c301e56ecfe4d72da8ba6a53"),
				},
				BlindFactor: h2b("74d466fd9e2fa605f26fc9c7b81323da35c116fca22cef911e5ed5f617942762"),
				DlogR:       h2b("2fca19d127fc8071d48bbca8c5bf4d5dce3f1d971da041ee219306075bbfb385"),
				Commitment:  h2b("a65596120b020cd5b2a6f788af4078018775e51517031b086a52005df04d38cb"),
				DeCommitment: [][]byte{
					h2b("a4f2b028b0cd40fd2c1739cd1cb8bf012f729696cf29a69917f7c903f13b52db986f11cb373b4ae1cc2a052cb662f8e8"),
					h2b("b0efe39f791fa0a6c57d0fa160aa2826cbc427a3291660fb5511601f6cdeee0600a6f32116838d7e4008b72406cac4d0"),
					h2b("9061664f0c5e20be011d9e44b18e4144693f14967be844bfbcf089fec888c02ec4dd27e7e7c99fcbe3326297a8c413d6"),
					h2b("b5729e00d15001263c596e12176968a859741d6c1790ab2d805fc36767f97c8b36e5ef2be5578dfad224d3cf022aceef"),
					h2b("84d8b10e533c98d63b506a689c02779b47c3e436cdc475f478250737d8a684bec7a830be3f04828b3dcf80e149754531"),
					h2b("95e1996f8220e35dd12098321d4e76ca94508b3cdf5fc00c83d765e246702859455d1ba0b90274eb44dce4eea43360af"),
					h2b("a0cbf93d4b5357db50c131e2dd1f28340c522b0eb9bb31194a5cd4c9391ddcca9e02741e41ca6ca4777223615145c2a7"),
					h2b("8721b12c8b091366860a1469c5ff407ce6e52c767d8e90cb7e0fd19beb27ea35e622c8242e27aaf3f52ff26713e5caca"),
					h2b("b8339a56249333afcf584e06dbd68c867c0a22254ef66b4592da43b94132f6f4d8fbbe425bb1516e2b8882078e59a33c"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("6d3d44d3190f88c890ef186dc756eb2d431c091979116808ae7ecf7ad6b089d9"),
					2:  h2b("049db9f2cba279002bb9a13411ce2acf53b2f1abfec9397bc48ef1dbc5304293"),
					3:  h2b("12df47e7921c1312c4ad70b6721d20f3c9e8b3a6a818f7e66e91d86dc77bc53f"),
					4:  h2b("5de3263dde13618838df6176562ba06b643106101f4c1e2959010ba9f0a7ef6f"),
					5:  h2b("5cbf533d4255ca2a65c5f958827451d70d253cc6fe0f9787109ef038039b728f"),
					6:  h2b("26dec2a9d435f6e7ef59ab8885574dff8c69747b90dc4ab5d629d11063ab77ff"),
					7:  h2b("2a87bf5d5fde6f907bafd61bc17627a688555e8005cfba39f17a67015eb29820"),
					8:  h2b("51e4b6af8eeb77e1bb390aefe997c48e34af87ff031d0b4e9ee4ab3c9f29960c"),
					9:  h2b("2e8736369a6fcfa233affb90df5776a958ba61f2849cd41afc147c7045333fa1"),
					10: h2b("7094315ccb548a636289f019cc7e9f7f020e84efc366556ca40257d2d5f6bf72"),
					11: h2b("14b767fa15a9e78f6a3ad0dcd53190b19f917979cdd16fb7f3ff2176e10d9c12"),
					12: h2b("0a6791bf1c49c142ea3243fe1f5c1e82acf87fbc5964ea733e38b0f7ec348e3f"),
					13: h2b("1cd190f3733648363eeebec4f0ca013b6778feb53beaf090787b91a7e5d8452c"),
				},
				ProofCommitment: h2b("a54fbc07d4a29a2414b1ec1fbdf357dde6bd4d0a4f3080c9f1a67b1af2b6780ec715952367051390af3abafc351d49ee"),
				ProofResponse:   h2b("222f7413f2cacabf2d2c3b4a91cb96c1e3eae641b809027109753ccbb2f7861d"),
			},
			5: {
				Coefficients: [][]byte{
					h2b("366bda1b3886f4f198194a7a223d5db0183cc3c2d8660fc8af6f54da34ac0aa4"),
					h2b("0cdca37f40e0da134eb66fabe9b0091682af3067895f0ba2fa173b05c3894f31"),
					h2b("2941b0e87c280a5fe8cabc3f14a486b983152c8b4fcc43b6ccb5e0ac92bc5d33"),
					h2b("3e52c089d966213fe238813b449e9c3b3f94a9a80faa8c548690a74d6466c209"),
					h2b("1fbbcec68fabded4398bdc6581b17889fde3bfe07f06559a66059344c6fbac7d"),
					h2b("6ead43904395f5c3afef698fe363b8d4548128e320c0d090a83a1566220c31fd"),
					h2b("64a60e378d4cc7a9913210d7001dc79ff572505095013f545437d6233a525a35"),
					h2b("03ee1beece6ce2f8a74db90d91d11145705fdac53b21bb806fb389270c8939cc"),
					h2b("100c2c44eaeb498766b476f99b21f9b94f413c4cf835bc26725d4a857e803b96"),
				},
				BlindFactor: h2b("686a6d591ce87f496d71f3a7efb808dc4836689a0392577cc2028ea61018509d"),
				DlogR:       h2b("6bdecab2956ab2c4d091d228338b2eb62e9aa4df5ddc82c1cba7a68195c223e5"),
				Commitment:  h2b("5ded6ed313f19fc4e02a1c742be216d4f5c4f4b76c77f67e4157ae7bc91bf0b0"),
				DeCommitment: [][]byte{
					h2b("950b62fbff5a0d4e3fbe39e20df82572b04a99b9cc20e04924f9e56ef168c14f6cd355c6bb446dbb185e1604194ef3e6"),
					h2b("a6e63a66f9f5febc2b2de0102e50021cb875707e2ccca2e4f60364357111aa0d9aa06535693c07bbc1112633ad1a5ed2"),
					h2b("8bb36692ae159d1904ef4ef48ee356bb0a65645d45f79b0027891f31a92b9a144734b02972c5d68e23a7ad5b576e3d7c"),
					h2b("8640c9cf00e67bb972530d4d5f277b450cc526685b8e9b0aaa41e84884c1a3d27573adf873a0a0ab3ebbfc4860c2d9b4"),
					h2b("a92ebdcb74b188f90aab5620a1a2156af9beef4d2ae13cd736a1efbe54feeafe4512b2a6ca1f1459bb6895c069d97335"),
					h2b("9919d207f49424c5b41d0316e2ff074c83cbc3f0dbee0a7d46c5fdb3edb9c532e90174275ad214913da4f9c90f947d5d"),
					h2b("aef04bfd2280acf9f5d354257506749f7bb25ab526cae0d65e9f9a2f46f00e054a8cf901fb9ddcc6d562f97a9bdd92a3"),
					h2b("b69d822f002bc6e88b8f346169c00a7068a7a4f0538ed7671ea9622833a41d24622db52b0f19dc731ca17082de18a181"),
					h2b("8d40eb6082b1fcbb79c4288d2f3403c334d43905c26dae345e210451980f98672c64f60310d5115e521454673b984697"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("561d61d56c044b8da0d4f65bda7105a869d52e7b2960b4a041556a579dbc271f"),
					2:  h2b("4fa42d833cd8e213f23d0edbe8ac90f7e5f4e08eb7f1e090b79b0b6f34963e42"),
					3:  h2b("6f2c2baab735c047212073b6a2461eab65d7772a10b2924916f9e372c73d7751"),
					4:  h2b("167a5bbf9e80e6a264c0707a97025975ed65ce2a895477379f6ea96a64682618"),
					5:  h2b("45963956dec1dcc7c1b1cd2afceefa57544099431f4ba8cef2cc47463c4190e8"),
					6:  h2b("4f15970185e6cc64ddf1e90097f76d804c7e336b9c0c17c25595a0d56fd05b16"),
					7:  h2b("5085240b3fac697262ccabd850754eebccb706705a7bac738b726cce3cf1e6dc"),
					8:  h2b("5be75dec57fef463dff333245abab950dc31393d3c8bad2435ecc982d983b481"),
					9:  h2b("167591dd583ce56f009ea9cd491918c4ee03092092c41b8a017e6de334fc77fc"),
					10: h2b("02703d8cdd6b6f0af36793f791ed386706da1f7eea9ba494524d8e92bc5f5173"),
					11: h2b("63dbdc97d563b511c1eba0f90974384c017b4a9aeb907e678752c606211c481e"),
					12: h2b("3039ed761af4fe8b65adcff1bfd7ff4edfb5bd8c9de5a9ad4ffa85b00819da57"),
					13: h2b("72382c5f92c6abf677abfd53443a329d69c450dda8e78c7ef5adcd606acc28d7"),
				},
				ProofCommitment: h2b("99df50e307baaaa2686c460288976e5941eb0cc7485d6a3924ced8112600c312ef42601a3a1b1aea74c0faab5a814c01"),
				ProofResponse:   h2b("0c46485de110cbaa2f3b4cc115931202d97761062b54cb2cc281ebed7ec0e485"),
			},
			6: {
				Coefficients: [][]byte{
					h2b("37264fb32b78fce0fef91d552be5a1b4e4b9327656a437232343a652cf6aadf1"),
					h2b("5c250dc3f2a4c37bb18f269f13406c83d4fdf47b98e45765ed1871b525e9e510"),
					h2b("1e5abf6604f746a4c469ba196922edbb49fa2c7d87ce317891f763d6f983aea9"),
					h2b("6b4fa043f46804ee9d193030d93ec0e439b26bc779840b6398ee5b67effe0d2f"),
					h2b("26d50f002da73ff7f2b283011a46e6fcfb43b90fa23d1d6a5c1dc8b724d1a672"),
					h2b("6ea6d3ebf3140d6a82195ff5f4dea7fccf16e3303ffc15d29a2da3a205bcb38e"),
					h2b("5da970a39427060c6aef62e722b21c53fec0f2e0ca7d0d28a5b2596b201ba700"),
					h2b("393fb4d1f19c1e500f897e50c1047c59e58ce26a8ccb3edf6392d246532a60ba"),
					h2b("70335c4462e3132bac01b7ffd6fb7b0827e1fe7fc1a6c485272f07fc22b26187"),
				},
				BlindFactor: h2b("36ffee6042f276b1d1f59bd822b35057991628139f5895a7c91c7df4b389fe98"),
				DlogR:       h2b("3e23290b83897bc1c454a264126696882d357e7e0e9ce1dbdc2cf73da7799719"),
				Commitment:  h2b("130227e79d8c2545a6075f45ddb61b4f6d1ed24ae3534a98d20b6cac8c74fe36"),
				DeCommitment: [][]byte{
					h2b("ad28e6a76331bf6a8422bc3085c218d2cbfee849825d6f89d3a944c02ee18bdd89c7fe92251d300507308bdf9067bf7d"),
					h2b("988749ac55896dffa5c426d066e5da1f66ac4bcfc373c9bfbaf66c3c83e565480757cec4bc29575dd0c1ebbc973b9c56"),
					h2b("8b4e989d4b8e3c346eccee7b238ce609dba918723f5644907543cf0975843f6ddbcaee2d98a250b25f54a0f6efa5b3b4"),
					h2b("ae315bf50296fb5a2f4def7ce2d07d59965b756988eed7e2d2564d505af19b391f91b99219e52969e903acc0465666dc"),
					h2b("8739a04e4b3fd55fb11872deb5c6896daa2d130af4a8fee39e1a51123abaafa094724ded22abb1968850c60dbd89ba90"),
					h2b("b19f274bb00441b8d363770438a0d182938df744e7bd048146e88b31464c027af4a80ae0431cd1fccb69fb850ebd0dc5"),
					h2b("a0893d256ee93cbde95a2e393b8f90a19c692424f1029eaceab07c6925261f6c59e3a47f2e671cc7b5d84800c1b03f09"),
					h2b("94da589584838da23ae7aca9f812ff74b8e385d1fd765a7f82d87756efdfd9adc669e59f6e11abb95153437eb2ddafcc"),
					h2b("84cbc6a6e09db8e39be21dd3948be9d12bc474677bef0786a14ce443c799bbffaaf297004b730ca68a5add62f350381c"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("01fc35d4272da12979f69a3d11944f681d7c572fec0ce735620177539f5d1214"),
					2:  h2b("0c665ad7431c212475c1c078c7d9d1cca9fa87114560dc82795d6fb8d46a1774"),
					3:  h2b("2265ffe9eb84f84cfbb84a2daaa2c43d2509f0637cdad3cf68a8bcf596854754"),
					4:  h2b("09d5921609a7bf79f7bfff11d3b9a54e08494a1c116d2914e92f5604776979ca"),
					5:  h2b("3a0d0bd4fa6e3fb2810a3301b598f6aab059c6da5be36f1500f4c47ef766a0fb"),
					6:  h2b("53f0ab66b573e0687707d7abaa981b7a8ce798757a8481d4fbe86cfbda25df50"),
					7:  h2b("66e0c2f7cc4c6c1beff1b0c6475d4152e40781f46a7a09c82c10ec243a340a10"),
					8:  h2b("52742b8d7d038fd00faa1f24de57bfc5ff9b91e01a651d98409855cd0922afc6"),
					9:  h2b("0da07df3d3457bd791b5ed0daee58b7778c31367556c1072730d5b1a690faffb"),
					10: h2b("36651f62ad544fd837953b84607949f458e4590c7a193ad3399794820b6367d5"),
					11: h2b("43491000d079b715ccfeeda2acaf0bf4c6a283de04695fd99c47644e7f956d38"),
					12: h2b("109052574c5616c38dec95393b60948e79dea6c906496923bb141dfb1bc7c829"),
					13: h2b("5292df60d5e79c38a08cf0c949015be64d519da7bde78f578185569ee808902d"),
				},
				ProofCommitment: h2b("b1f5c82e3835dc24324c75d8d78b6b04892806f0bb053ab7770a929a97b1baff06b8200cbf442b99bc8eaabff7516349"),
				ProofResponse:   h2b("588afe381209ae232aa315607f02eaa016d29904044000782a3fe77938d00d6b"),
			},
			7: {
				Coefficients: [][]byte{
					h2b("738740680884ed50582f2222c6bd1f1c93e9e49b3f181e7ca33bd2d3259baa38"),
					h2b("4685a0544ab6c862a0f7e817c1a6aa3bc3f8c1203219abd7cb4c6583d4e5953b"),
					h2b("3176136f88d89fdb577b4e6d5f27f1a04211982d3a02f6f08522ae84bf8ef87f"),
					h2b("41f980f1ae3ce7e6dbc700326ea3a438d8c75aacdbf802d9d072dc4e57a87d22"),
					h2b("35fffd53d934b7c9c7b52b14f70330a808d3a60866dc29c4a71c513f585e9f51"),
					h2b("6b73f739ed331dabaa28375f1e8eed70dc120569499430aabd07916f595d9854"),
					h2b("123c4b7b4af41a5c423f590241d09865117c115e728e9fe1ab7bba6f3cf8c566"),
					h2b("5a69c87e9a43bd7c0a7183354f5986088a2e9e521eaaf83c29935bbb2354a7f3"),
					h2b("6908ac2a53665fffe7892a28322463b051d1566c67ee3cd9a76217602c4013cf"),
				},
				BlindFactor: h2b("6a0eeb316be2e778751e08a3189f263cb1d701cf78e35067e4b3952b937e21de"),
				DlogR:       h2b("02345b1c9c6dc550443fafc8117ec6859b8744a41982937653b3368a7829f085"),
				Commitment:  h2b("3e6307678a71ebc47a813ce7bdd9c2739293b4396eab18b502d136651f29eb86"),
				DeCommitment: [][]byte{
					h2b("86391539e8911c08955d1302beea4446cc5b197268de7a325b9fd813df60df702bdd75944338f8032a273dd2568f6237"),
					h2b("a1af4356657e0bbe32f25a1fa1418ef72a3d7d415be38a32869ce1a7656297cd3b818b040071f6878de86304f2780987"),
					h2b("b3abbbcbadc18db875c787a934fae7060e4d3fd0e17e877b04b23e57982abf9d3f8bc3131c3fb8a7669ee5f38b2770f9"),
					h2b("af2e87bcd6f2868c3206c6d5fa96217d1a9aa5b53585bf67839262a694cd8ae24f47aa8c242f588e247754874f72181e"),
					h2b("9781284dea61f93bc6915b43dbc6c93e0c24e8dcdef5a64422ef6de30365cf27ae9be26469a55a987666f2e3c46d5a39"),
					h2b("b6201a712f880e18cfc233002989658970afa4aa718af27d55c64f044a313ffb5edab781fc08aa263f2d645079d15118"),
					h2b("b58c5e9b312b29c3d56239e60be183741ec4506d27aabfd3c4de6f65300c7d863696a55f96fdf60ac697a7b5676c3ef3"),
					h2b("8579b3b97c2a88218eaf70400e12b60619a40e386ade81e175ec2a0414b2afd9c77e68228fdff8cb5de587842d0759b1"),
					h2b("879286bf73276d991330ae9886e3bc3f150e354402f35cfb885ef4eb1c732ee82a54dcce27d7a172e3a9eaccaad85c5a"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("60fae52fb943d859d25f8985fee6c74da269161530cd278aa4b2d36850026ddc"),
					2:  h2b("3114e892c91c5fa377ee3a12337b5245a9d2a731f53bf808ef615751651cbfc7"),
					3:  h2b("56237bfe4dcbd13273b94dd0c2dca2e2fea1e8e88112840057cae9c3f3bcce78"),
					4:  h2b("3f51fe1478f8ecead101f0cf63c22f1fa14311499134782c23dd6b20fb626fda"),
					5:  h2b("15bfb503c70d2cf61461d810e2649204a5035634c29d5a0cb22b1df18ae7c0e0"),
					6:  h2b("4c3503ca1658169d7cba44fd211464b0eaff6418e94078a3850dab3e54b8edb9"),
					7:  h2b("4d5179c9e287072513863324dcf56c6ad5449fec78fae6c32fa7ec9cc5e71ca2"),
					8:  h2b("3f8ef21ee1a02af2e37d9640d2ba35fe18e3d575d03003057c1e84329995c12f"),
					9:  h2b("47e72860dd37516e74e20d8df64d7860833dd8953059a8bcfbc184f24e70e99c"),
					10: h2b("427c7dc3c0b6e76d7841b22408b19b2749fe26542d441c520b43a72c500a6397"),
					11: h2b("2504a7a4e359bd5e514acd967eeba76012de14455845c0af2e7dbf83122ae1bb"),
					12: h2b("4ac89c0348aa8eddc900ce027ae3337fc4e82e91e9540e90809f683ec85396bf"),
					13: h2b("6a0f532170340f29f7922c876ed6fd3d6d4e1ae924670db8cf1006eacfdc0728"),
				},
				ProofCommitment: h2b("acd2b9a22a98cb205e286b6a4741d836dcc1d515adaf0f21e0d2c094916f51e431e27a082885248223e2bcdb1dc01c98"),
				ProofResponse:   h2b("54b26244563950163570c4b8e21fea560ce0fb4fa9565820d7d73daa9014416d"),
			},
			8: {
				Coefficients: [][]byte{
					h2b("720bc816dfd1e19152b3837c3ec1045d946eec77371255ed3e91f277799b2a84"),
					h2b("4bc320e262863e434cec90d809ef9d8b704a9f6fd59f0a4d23a92939ba2af977"),
					h2b("6ef200d9eb5338bcbe74fff8dbfce6ec4a8af9cbf405608f7f672b022b4dd8ff"),
					h2b("0f86b44ab1107151bc56179d030a6420fbdd0ef6ff5c76e09ff031510d553a41"),
					h2b("188fef3ce3e80fe56a20cc4254fb0f72d87795de23e4e9babdb12f6f68a2372f"),
					h2b("2e42a5a98e848bfc958090a496bc6f777ba3095a490113e9f0a23bccacd9ee8c"),
					h2b("1778943a7676bb647bb37047ee59049f83fe8a5a4965e9ae73cdd463cbb8baf8"),
					h2b("08618b59b725d1180fcdd406721ca36964ccb8fbc14ef77724f1477e626d8bd6"),
					h2b("5b5f5dd0ad5c8e12e0276213e4c00afeb1efe69a7ef9e7656cd8f37685fadb0a"),
				},
				BlindFactor: h2b("203fa9c7b7e8fd0fe8a9fcd7afb0119d207745c7c1c35a7850f30bd6fb2898e9"),
				DlogR:       h2b("6fb4840092b238cdc3651797595d729caf0733d1c578f2f4567612c509c5ff20"),
				Commitment:  h2b("b315446fd0ce53ff7faf8592b6e8af106efb6f9092ab0ada0560a0b1655169e3"),
				DeCommitment: [][]byte{
					h2b("835650ca0f5c1c18f1dab1035b1fba463071ffcb73ca228166d0162fbff6955e720945f112f9aec2a381d1b6d317f128"),
					h2b("b82e070972e6edb1dc2764ee8463ee481acb5c654607b926777911c8ac15d2a52676ff52fde91cd921acf2a355e7feb5"),
					h2b("a80453e29941081dbf1c31c0093d54bdeb6b243cc02e7b88b9551770527c265420d330e6d968dec40dab9262870bf25a"),
					h2b("b804bb66c87b94a28d28c642bb2c6aa16a8623761cf35530d85adc1db97d20b073e7b9ffd620b23c070029946709a660"),
					h2b("a16420f53640f34080f27dca161a05e33f29e9d88f9757b9d0b5487072bf45e0ade6bfe6a96be46608f5a5f6cb050499"),
					h2b("8e7e66ac763b612f335d16d8b21076411f9610f9017d9347623fe68161fb72c254fdc444c7e2948dd7ea2298043e410a"),
					h2b("afa05f4e1c1358ff28b5fea7b3c5efdd8b4dafa1f485368081bb397ce0ee8bfd119830fc6582a7b67b7890559bf20de1"),
					h2b("adb83cb738705be6123f2bcbecc4c8f9375b07aa65640f8f4f8e2bd2a93189842589269ebd4b0fea0e075fc41d1893ce"),
					h2b("8033bd863e3c3b7a0ba1945aecb90d16130066435cd51c4f60698c78b1b1bef4d7f9bc56255ecf6561bc4017eeacb8b3"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("2e9d131c85ab8b33b8cdcf13321dbed2eb00cdc6f6ae8dde357df29d36067eca"),
					2:  h2b("26df73824ff2e0dfa988f4081b9c506c03df8377082d8647b7693ba94b0349ef"),
					3:  h2b("1e6ef47318c8175998199048a8b9b10f797240b1b9d7a3a7742945b48514704c"),
					4:  h2b("2799b628e7ebbe2a9d400e21dd050ebf3567ad38ed2f8ba21c716cc04012bc02"),
					5:  h2b("0bd150f093e1daf124260859440184f27afb1a1be0ad5ea323a71655692ec9af"),
					6:  h2b("3520f13440485fecfcd94de7bf58e64145ea40868ec2224ba1a8d92b95ee24c7"),
					7:  h2b("3dc6671f35af0ca7430a129ff97049aa3ead1533cf8994eaccfadafb68fc2d40"),
					8:  h2b("1b5d89bb1e20972b8a47b21bac72aa50859593faf5f844c1ff94801d53698b30"),
					9:  h2b("0229f508e970b5dcf5bd59229e419d52d5adc36c2cfb6484b255703a96b71fc5"),
					10: h2b("4def154bde5887454a37018f0cb0fec6e35169466474991454f1dadc41cb8de7"),
					11: h2b("3be23febdbeb68c40b26ebfb7ca7fbf2ae103d36298ca31861efce3ab8b49fab"),
					12: h2b("5a01d3ae641f2ce7829e8a66ccfa6e518807e82a514941a088195a342cd7f8a8"),
					13: h2b("2ff078689fdf76fcc8b14d82bdd6bf818ead12e7e80ab1b73c9120ff40f8affe"),
				},
				ProofCommitment: h2b("8e567cfd116dcc34d3c727df8f797d1405280b225d15e7eacbb2e76800abef2c17649890e363d094b05645ce6f38bf88"),
				ProofResponse:   h2b("5a9fa377f418359e23d503236300f82dc6382272bb81ce442568a6286d6e58c7"),
			},
			9: {
				Coefficients: [][]byte{
					h2b("714f7eeece84bf7d44bd7f2f62fa60e71137a7537ddeb7033a7548ccae9c2e94"),
					h2b("5e9fd262e3580cc99da2ad76629b4e0f650716c767e975a0ffbbd87fbef27b63"),
					h2b("6c3510c459fe96d0d35cdf00e53f3f6cbf1677ed9bef23609ac12c943a9b86e4"),
					h2b("6ed6dc7c96154fde0b1f15ab0582d6a1d1f4023ab565b443ddc3524407503f10"),
					h2b("46bf2a461c9965ff88f256b7421c1d6032f0b1428e2da0cfefbcb20e68fa4f12"),
					h2b("5b39eec4c2a4e1daae23c069fcf76d0698cbc4dbbd5e01c4e422df2097c9e1cd"),
					h2b("35ba6ce45f9dbde8006be879f6afee5f0eb234b91ee212fa71abc0cae9b0f371"),
					h2b("1d4a0047eaad03dfaad56313554072a4edac54cdfc0e9c862525ebdb725ce632"),
					h2b("71ea39f351a96f5a1f63b9ca9a986f18c08da95ad5a22cc9a6d1913e3fafb370"),
				},
				BlindFactor: h2b("549c9417c8272038a7a4fa975323cb3264d1c4a0e12772c0b3d00a342843f175"),
				DlogR:       h2b("0da340a5469cec0b08b3164e029f716bc72cf8c46f957703b02392ae329c32ab"),
				Commitment:  h2b("9331acc2dc29e2914c79fd4e8834607a32e5ca6dd3a51be046aff3bdc73be374"),
				DeCommitment: [][]byte{
					h2b("820a4f955828420ff9b0c8e5d5cf20fea5c789fe11f46fe24a98e771ceeba8a63e797f5e9698546f5fde0904bc0fb08d"),
					h2b("b274b87c84c39770ebef8eea234ee2bc15b8a5463e87f8cda60f5043b0983fb0b1bcbdca1ef4e11b1166ccaa0be27b52"),
					h2b("90156fc7dd123894f49ab3a7706fa3a277faee773036f7994ebb090e93a7fe5841e49ef91feb34e9d5c257036649bf4f"),
					h2b("b1e0d04ed48ae526868ca78d087446186ee6075ee9742273f291f16eaf8abc73f20743a027c66e30f8b3610264ebe584"),
					h2b("98fe2df2f8870c74205320cc6960adb6d5f477d76689d9df35fd676efab4de05f1c1a4c0203107ef2641229ed37b80c9"),
					h2b("a7e719d32c111b6c242dece7af9abc21849da7755488956d29847bb2764ee8a08014b07df19ee5836e5b6e69d4de49e5"),
					h2b("97b7a70284a08af63ffecac60f511998543aa03d71c2ea26b99a4fce89b225027094fd3cedd6170ffe0cea2d9b927ea4"),
					h2b("aa7ba86d1df8fa08a3a61016fd59e5723b78b2daa517483049d8af8abe06cead7d95105d647626435b8f4b85bc5d4f33"),
					h2b("b447e0a9909276c46dfb662b94c1c55aa84c1a466ab5cd7e67225f67cf02babf0efd9ea3fa19c599ea9e1f1edd17a3db"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("5a5112ca23723c408f3c2d9a9c290f689980093173455b2dc4386f3e4bfc2dd7"),
					2:  h2b("39117c6adab9ea6ef3a67390efbcb21055fa79ec5b87cc41e2aaa35424b5c800"),
					3:  h2b("56f3420965ca578eae47b02198a3a84085a79c7ab42f782be3e4aad8bac3d55a"),
					4:  h2b("6b0098c8f91babce65e08238c59b188111586e729aba9a7366bd77f046370d2b"),
					5:  h2b("6f21f7303e1006e495616604334b221537f64282ec3d0afc6dcd4a34ae550ebd"),
					6:  h2b("53d14aea775d33fbcfa51d0c2929ea180b135705bfaf141963d054d0a97ebff9"),
					7:  h2b("2147a28ffac805fca775415aaa68b7c92de8a49e925aa8032fffc73b9ff2c351"),
					8:  h2b("27778ef16094509d0fcbee14922415a0efc4f2eec25618b667d47d2e866593ee"),
					9:  h2b("4c43e816c37f8d8cab9b7a0a54eb36c07d099f8c3775541b9a83015d7e4da452"),
					10: h2b("20ff9e2089ddc4fe8e503b2067fb1275ce195d9c5f890940a15cf56ea2a6a8e0"),
					11: h2b("6272ff9a0145ee785e5c838947c3198349f7ee93e73e9e50cf20589b00d404f6"),
					12: h2b("11aec84af3aebd13d48523849cecba9373d0de95dd8311edc7198d67562f1e73"),
					13: h2b("6b4731f6ae4dfb4f03cab5b279a21e0e5e97a9833384c54ab3d965aec3b328e5"),
				},
				ProofCommitment: h2b("b46c36fc93cf406512b51f272841c43a05004d1779a6693d3e851be81582a7d60e67905b0ba8df1a6b3f4d312b5c24b5"),
				ProofResponse:   h2b("3a5e2ece7e2af1d8c861a1bdf508640c6ffed5e334b996364dc2fc91502a7e11"),
			},
			10: {
				Coefficients: [][]byte{
					h2b("43ef0eec4faf4b0adfe7e9b60d26082e10d2260c61eaef7f2f67f75560c2e591"),
					h2b("48bc5c5899fd1ed54b6f0029a2be4dda555a6722e9e36db63d1df17aac1e393e"),
					h2b("5d57e9c08722a4d1ecdb33bd19e916cdfa71e64092e0c71e2e11fbdad6bfabe1"),
					h2b("6a5b9f4edc069ac7b39200a07e177197620b9c21275ec38270c46fb27ec04a54"),
					h2b("3149c7a67a72626eac9060bd21904ee74c30f2beac8297462881e1081f2081ca"),
					h2b("3e6d44ea0dc21caef1896d248a99d1ab04454517e4f51baca9a8404125aae467"),
					h2b("0464e33158ca2a247f087e0dbf05f1ed066982d0a1efdb821dce5c0543976bfd"),
					h2b("6a0d043b4ca4a0259c1982239f2b9b9b60de87251d79962456b371d3f4d8aae3"),
					h2b("351e885d4568d8582349cd2c90b6f75fcfe674ef73ff6dcb3c4414cbaebc3bd7"),
				},
				BlindFactor: h2b("e5315adb4993adf3115e70ae21613b2244242cd8415519ded9b66abe1b0fa31b"),
				DlogR:       h2b("228ed2dd245b6d3cf52e98a27e1c766205d8f75dc86134dc914f38b82b92e393"),
				Commitment:  h2b("28c818531569d89b066c595742ca8f00ab9f6d7fbb7b0f1d532df333dffc1d67"),
				DeCommitment: [][]byte{
					h2b("a448308867a5b9a64fcab5fd141dcf157576f7d80de4d15360b3a30cd9f5b3cb66dba5f601139c245957862b89eee62e"),
					h2b("85eb9a6e7a8433bfd88cf42266de2083c472145f273d4393b89b2ca8e83662819f9c9f9e6f85913bb700dddce13387e9"),
					h2b("8a543b2817396b77ad80761f44d1442ccaa244c114364c38e196a330d28c42302b33648522ef7014fecb3842ce2a9169"),
					h2b("8e82602666105da0eb4b32f50233aee9152cb20d78db19c5bde6f069d37afc4259604c0cea5ee954d1b463003ddd5613"),
					h2b("aca65f8b4fd35eeab7fbb775c37553dbaa4acdb91741c1b3987136ffeab84da4db4e2df83932d4a5d89f9c90273ac90e"),
					h2b("ae4bca4926d713684df3490ec36b4a0af204f3fa0d10c9042f8359cc9977e70901074d9e0af23a3e883cdf5d9b93431d"),
					h2b("8695616e387b1ca2260d105eca064aca2a057c51003a042a09f7833a077d5441d64761e8c8064df59bb98716399c4d0e"),
					h2b("b5d2e9ffc498a2ab814165a2f092aabceb129a12f5ae3d98fd32752ef81810468c5665da3b7915886608ba08736d8408"),
					h2b("8f079b47c31ec202cbd696b622903e397f9e825ed66dc23305bd399b7a6a15e11f83e60b56f5e2b3a86550457ea9d375"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("24022c0eefce58d0a8288154b2ce4bcda79a923dcaf6ae3f8e4c58508e58cde7"),
					2:  h2b("1dbd58526c93091638f4cfa2dfebfabea4408c3d5ce7906fb69a17fbbfd14a60"),
					3:  h2b("2b11b928883e8503445b4c970d764c7754f30502531bb9a862632508620a433a"),
					4:  h2b("27b506ceaca4a7e51db448b4a8035ae9c503ccb05d7d44c72fff27c1d6ab6046"),
					5:  h2b("18277fbad9fa64b33743e3a0f6d4b38aae483e095943307aa0ab1353ee0a2335"),
					6:  h2b("6dd07ca496d64c0749df9f4b9c9cd0f09d0156d5b974a2702745c57f8e38da96"),
					7:  h2b("685e78702e97e5a89b85ec7fbe9af6e50453c6aa6a10f0b1cdc35841c4be6ab8"),
					8:  h2b("604a5f500749c2c248eec534f0aaed1172e64f0315165a6edabc610c4f456809"),
					9:  h2b("296f4038b0075f937d9676ff522f6dc928f984aacd45ef47c349f17876617cc2"),
					10: h2b("334749b13509749bf5401ff8ea39c90faaf6702a3ec201d4982ac568e45a9225"),
					11: h2b("3f26e01ec4e67e1ef64ecea6c521f5684c24d001a24287f5fff9a2d4f0bdb3cd"),
					12: h2b("00383904d4ea67f030544cfe8e49aa606dbe2110cbebb52ab308878ea84421fa"),
					13: h2b("0ed34b018dd2cf005886054251cdd276cfe72ac5cb5bccef503fd98da771860b"),
				},
				ProofCommitment: h2b("aa461f3feb9048ade5773686ae5322d9e206d11181a56d290a25bf5454ab053451aa99874df87e440f941bf5a4c0b0f8"),
				ProofResponse:   h2b("372ba40fc28bbd94300ac37de3b3376f6375d61a2afe1b79c83499036be3a51a"),
			},
			11: {
				Coefficients: [][]byte{
					h2b("595007c847b934e50313ef4e209fad0039b75a9e0f1befdd1a1de81a3697335f"),
					h2b("025a99f1a823bbcf964b0dd138af2b02c70b23a004738fa2845b26bf5e27e580"),
					h2b("388c25277d3213f2eb2650712bd4dfac7548ed6f26c18aa1c0f5b224f25c805a"),
					h2b("1ab96afa9b16819ae34a5c8dabe9830085a6a6c666f6bf7a8fe5c5814975c4b2"),
					h2b("72a8cc99b5f47e369a441844ff0d86078e9861393d6baa62c5568f84806ebdde"),
					h2b("3fe49d0dd51662896adbd3ae401c281920138198bccf22bb32539b713d01955f"),
					h2b("40b4066864e7782f33f29cdb103412df2c2a3500238495d4ee7b1803ea73b111"),
					h2b("1ad457a9d5fead2e5724bac5cdcd34d97caab9cc9d42fb150837bc940b389d8f"),
					h2b("274ff759b7c09aab6e867345cc0ef081da689031a940e0392e4e116a21063ed9"),
				},
				BlindFactor: h2b("8801f1ad7471bd9225f2bae6b06d9ccbf517bf9b60445a5dcae9a3e197aa07c0"),
				DlogR:       h2b("35bfd92707aa153ff3e76a94f43548176bcb70bfb27688c85c991ff6e220dac3"),
				Commitment:  h2b("87940726bd4d42f76faf224ee1ce26fafb5efe0d72436b84ec389b7385b36178"),
				DeCommitment: [][]byte{
					h2b("85ffb477b22b2a0909303e4d4dd8fc0cb689c32a2e05b92b8dc6bd625a52eb9a0ff95296a499d7a6d60fb611a82c2172"),
					h2b("99d3ab129855517c1902bd4eca35ee6968cc93710a1e6a4e47590c7e368a8d8a97b27ad4f3a454dcf71080ae588e51af"),
					h2b("8279df013fccfac76c85d824f637a8f6989c4ed01e065b479648125155bc930fc4d4f7bee8240b32254111152273d895"),
					h2b("968e433425873cfedbbb8f3e6aec4b1de4bab768433bede1c7902427e5a22af6fb04b4f4486c30d00baf33b2efee1700"),
					h2b("8be7bef34fd484c4fb9fb49b8498248669bab85f1d63a27c9a27e85e8826737637f76b206004f36b0552d0093f084945"),
					h2b("87ec647fb83ac03bbd53ce6d1cb643fce7fdb8cceab9b93b4290c474fcd6590c8b0514882c6a3d96f9a448b43cc22256"),
					h2b("ac68024c8a8b2d57a50ecd6b3246bab3825dfa57584ffb63b6a47337808751359da441e637350749d36f18e07888cc0f"),
					h2b("961b0882231f7b89e5bab27b6d6c5e5b381a8359036e9b93781f62eac69a8e21a4bd0f64e315deb529d7d9e186ba4464"),
					h2b("893e6e8f552a791017132f9f5e6c2ff3db72171626d4e39e227932e0303b9badb39d515243544eb588d16d98a061897b"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("149f53a2df6131ea99a600d7f3bfc0f5dea4e438059197e10bff977ba4b43e9d"),
					2:  h2b("19e42c777e7d908f06dbb4ba716c4fc857d77240201924bc4e6b4685ee9f9319"),
					3:  h2b("2e699d440610305b58f2edcb92f0b1e89fffa62d79ff4bf9e24a2faf0273ea9e"),
					4:  h2b("533d70633e9884e805b21c9d3a72bdee268c732396074fcf8b64237f48d0b3a9"),
					5:  h2b("211517e21b489a4bbd387a5795848de437afc6168414d6443ae4aa355064f841"),
					6:  h2b("1663a195835744be55db4dfb1d0f5f6d2a30715aa03e64c774aaeaf0ff1eecc4"),
					7:  h2b("17669748884b34ef18fc9a8037d4e9c06294d9162f06aec913b58227e6721bb1"),
					8:  h2b("48ebe653d6439209db3850bc148a445726779fb3dbea7530d9b271c8afd06fa4"),
					9:  h2b("5291b1433d95140a05386b124bc7c80b69006a9fe1252c0f3d725641ee1ca8bc"),
					10: h2b("46ab3d34218a90c5a196053c3004041aae4ec13061080336f3b50c4af88716ca"),
					11: h2b("0986dfdea41b376f727bddbee6d33f8e483fb8d63ae98de2b57cf415caefbdb1"),
					12: h2b("5b5fd39affa5b3899a60aadaa95bff2e7c4fc3446730f4e2ceccff3d3f935584"),
					13: h2b("0033455bbbaf472e515ccfe58ecb2f1aa4407069a7b27247086cb7904d72e3f3"),
				},
				ProofCommitment: h2b("963d1ceacc6c963d679ddf43f0591274e278bdfe69433e73b2723a0d8c560375d09396dffbe1b27db308a4faa0c0e9fa"),
				ProofResponse:   h2b("0921ebae29a2cbd706e54b7920d3c188ce3663f8144a50afcb23b5912ac102ad"),
			},
			12: {
				Coefficients: [][]byte{
					h2b("141448c1cc37e06e82285da8c3d4e1c335ae27fba209895a99a82fb684b0d928"),
					h2b("1d50dc3df73f21d2444dcdfca43966aad54ed9f103189e20a8c3f2e2e9232407"),
					h2b("2b143160531f6b294f7df077a0b490f0d89809ef53c075ab703ff490ac6e5f56"),
					h2b("1f37fe2ec86abd81f51611bdfea48d33c8bc7e17f8f1cf994b94178ecdecf770"),
					h2b("3ce98baed396b50d1839d786342e3b97dd6c4b04494f4419b398fb7ab1f3dc55"),
					h2b("13d7e23955e1636f0f96c1ec31cf7aab2bc56e3c05c4463a5804c1e6f54548ad"),
					h2b("3d31f4107e12787692299aff66b15e5d5b1a9e89e23d206d7e0a12c057f099f4"),
					h2b("1beeb654dcd5985e220ca1ce0b972d189fffec793423f2b516b53b7108e322ee"),
					h2b("2c2ba5a92cd026f752ca75dca94fa8743483353ca7241f6e8240d65800dc5bb3"),
				},
				BlindFactor: h2b("64a2822fd4b24d71d2de2e4b19feccf1b060c562d0e4cb61f35f9f6efb695906"),
				DlogR:       h2b("33eb9c361ee84f17788c97b412441b9042d9b3d20f1e9d7ec8498f77476b78c4"),
				Commitment:  h2b("a08459ee1d8f6ec5b28154d190e403bd79d1abc0e103db93292d1000172de961"),
				DeCommitment: [][]byte{
					h2b("876ca867f0522c5993d04fc4fd10ab239015b7fad82bfe7b14b70f0794c7a1dd3f9e8ddea11da2d725777fc733353f65"),
					h2b("8eddce0b6e06854129c07ec2e0ce3dd457d9073ee41974dcc407766972efbdec0057e0a3e3eb90e9c2d2d6aca9bc62e9"),
					h2b("a849a1fdad776d5c044d790eb7a6bd551efde2a9b99e176d26967b8c73b2d62778d59ee1049f85e0ec6fc7f4b013e3a4"),
					h2b("a92f9af4dd75b1fa43b9de43abd2d5e2a17753296879e01568fbddba1364f26bcf52b988975746d111c92f3a39e7b0cc"),
					h2b("b09bbc7d90b192b2107e56bca6da4a3a27142d4738d183bfaa59531d29faa276b2916c63edad006b1d6034d34dd9adff"),
					h2b("ad32f7e1d34d9549295161652613615531dad535a7fd194d2c24f7bebb1186d7ccce4ee6840f6dec427da1e651567758"),
					h2b("a700384b1cec38551dc2b67d423786e342c1fabca1bb2cccae7667d4d96231b7ff9254247a2102bdf1d0332d42139b63"),
					h2b("8b12680cdd1c820346ae5888cf3767c374b704f1d2d660df9c598c7ede2e076c2c5491eee560cd843371aef376f13927"),
					h2b("88618117404372e942fe01f7087fc38845da749e53e6627ed8c8cfcc5170fed2bcd2db5ccc6446ea490113d8e9345b4f"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("69e3c3df3cf680a3d367c9e775b9a0b53da5bb6dfe7071a720de10a5f118918a"),
					2:  h2b("2abfeb8ebcf6a8bfbe639d3c02dc1f7a407e35bdd2d0571b60f15a0e8a12db4a"),
					3:  h2b("23fa7dff3a83ccfc9a0b4242bec550ac106ff388f46ab332076565148f02531d"),
					4:  h2b("177bc48e23080ed2a18efbe304eb317db6d533433eb9f87f982101d1a26d88ee"),
					5:  h2b("582ae3da674c10881b59978b106aadbf0e345012d3541e0823f42f2ff99435c7"),
					6:  h2b("59a9161f81eadbc954043ff1dec04e7a32c388a68ebead4465de6e27bf6dc670"),
					7:  h2b("36f2e91d81e24845a9a3b293cb185a2cf16512f51ee754b974a5c8522345436a"),
					8:  h2b("209e0a5ede28561f7d95b12e1f3f83beceea0e5743e706441a2e02681d910351"),
					9:  h2b("219d6a8a4467075db100b8d6f199544d1c99fa6edf65694dc5ded3bc880ef93e"),
					10: h2b("313c5d266eb096b9b74aa7f6edf4fa6a37df560992c6424bf3475c50b898d53a"),
					11: h2b("546c7b8db5c1afbab0fed43b1db5e00b5e9139f9c8161610c20c73ad608d9150"),
					12: h2b("5bed3326bf2b11a368cfe9653a621efecce8f4ad1d76f499600ecb3405186a49"),
					13: h2b("67f5a59bf6f5c16f407d21d1b2cd027774c06c1f46936c28ca953f23f707a79e"),
				},
				ProofCommitment: h2b("a5dba27b79f6c006d95a0924040d67160d492429d77bc62cb44fbee74d1ebda5c06da26a19a0fc15165d8e25ad64b747"),
				ProofResponse:   h2b("0b4ebad5ca7668e69ac3aae0cffebc1394d30e507e279c231e8a7a90e5209e3a"),
			},
			13: {
				Coefficients: [][]byte{
					h2b("6850d49dfd9ccd907564adb329fbf7fdb6f1e9cadd0f0ff3fa929ce55d0a1e5f"),
					h2b("711a777b751ff8802d236a3f4fcb18060e149d3798e08b90369a9d17bd3230d5"),
					h2b("12d3fce8b07ceaf9a07487d65c65e7ad2d1ff1090825f5a898981f8a0eaca95a"),
					h2b("08aeec6c22c0e44da7f7a6d87bede82ba1aa03a9a669a97e22151f1a7f713341"),
					h2b("31174770695268f90768c077b2ce35bf43d5b65078a852dd7a4189a2cc5a9703"),
					h2b("0f615ec65ab925a55f445ad327e3ceb41856d70a4c8de0bb2d0fd194e85fe93d"),
					h2b("60e643dfdad370a8cc4adf1a5a29c68ca3130cee1dae3442db38b608faeb1fe9"),
					h2b("5e7f2a445c1cab77e60335671495137aba15a2069f22ec02c3aec21f2f668b4d"),
					h2b("3e8367d27d2802fa1085360833f3b2d8e5c8a4c257d7cd55ca3d771626287177"),
				},
				BlindFactor: h2b("2e37a7e3f2dee7176ca5a768285d39da4b56714274cab101874a4f28b56b0187"),
				DlogR:       h2b("15f9c4b8a599a0266bfb862f03bbf9ea961af7c27c6e25c6b68ba1884e5d8bcc"),
				Commitment:  h2b("9bdd6c811312c4a8e8a745e8d445d7f50d0fe170a0ac40475a5378995118a5b6"),
				DeCommitment: [][]byte{
					h2b("af897b227a3d2294a70e9c7cc50e97259d4e006b85a4aeeb73e48c30ff75ceb817b99028d1317f9606e528a90389274a"),
					h2b("96d27d536df80e326abd76cd8ba0555012ea4d2e8bb5594967d6d2b77b7409bfaa72c4d8cb6c7f4b16ef1bbaaf0ea1b8"),
					h2b("b90d7a6bcf7c2acdf4dfcf273faeebf3c64968c60975f9e1a44fda23c234f4e157e462b1072857efaf39e3a54cc8771a"),
					h2b("925bfbf55d8214936853767f0e17d46091f621ad91d72a1f7fb17d41479f1896dbc6bd407afc02d22813f7e39c6b950a"),
					h2b("838b18f2669be9835f88413ec57dee2bc1595088d11b47f1ad63e4c8bd74cc6cdf020de6daba8e9c3a4bdd7fde808914"),
					h2b("8ed920e568aef5d8fa8be91cd23308944f2ac166ecb9af78aa96cc560aff3867e1843ef6004df605d4208b0e55edac08"),
					h2b("88aca7a3d8268c08b0f8a6a1749768520100f4f4ac4dd985e223267fe06f5a695fb111f4b10e0076304403ecc7d45ef7"),
					h2b("91b328a19e148defc40b918a1f18f7161ad9328ca034e404c43d7a3423d50425aa238241ad03b5211fb5e720ad323112"),
					h2b("8b884f32414ffacfef9cbe2aacb6e249d8004bb3464598b8a16df21499e78714962da36b6023e026ef3a062c53cd0c0d"),
				},
				Shares: map[types.OperatorID][]byte{
					1:  h2b("6399144f17a84df0478d4c55a8f8111ae3f7ccbafe64ebe2fc50c31bad8ec8b8"),
					2:  h2b("73b2b85832a1fbe981cf22c862a8dd8b0db5b18955c38ffe03e9c40d75d06dd3"),
					3:  h2b("616434758007319d0d54f5870ca80c16a0b312fcd2c114c634723ac530710671"),
					4:  h2b("6a8367ec43e9a6c43e98f0e07f98fa1b6b193708f8724fe817961539b96e41fd"),
					5:  h2b("5bd0f87398086902f9d01443309f828f21f46b2810988114dcac2a1c0b698c03"),
					6:  h2b("0862d38750cb41353d79c05d5de079a100df68c346f675b5889879f03b5161c0"),
					7:  h2b("2fa5001c5194d9198fa2990ac426dce951fb53bf1f430ef171ba8c617b94b92b"),
					8:  h2b("41ed25c2647db443f1236d2dd92d8ce692e307728864026233e9d02ed1d79c63"),
					9:  h2b("686c0513509c7d4521a5a8a7cc6981f8a24de0d89f9f62921946595e0bee7e18"),
					10: h2b("5a2cb29e4d01b0debed68cedb183734b73e22aac06d3b55dc568ceeb62b42cf6"),
					11: h2b("4c56d99f3c15143dbd9a9eeea49363a8cdb35ab080c38310c140ce781b1cadc7"),
					12: h2b("53df60599b0448cc15a269b256bb84aee555997556cd2e0356b581f657b5a57f"),
					13: h2b("147dc6dcc4a25ec8388c66cc4485de1a7c91d840db7f474efbe098ac2e935e03"),
				},
				ProofCommitment: h2b("abf6ecd23517069bff84256213458ac759262c8f7f501f6639717b243848d05debb0d9b23eae70ce5772e362ac055a8e"),
				ProofResponse:   h2b("6cca3e626992c8cc76590840ccb2abdff17ed68ba7d6414ed46f028cef5508bb"),
			},
		},
	}
}
