package testutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/vss"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/ethereum/go-ethereum/crypto"
)

type DkgPartyDataSet struct {
	PublicKey       []byte
	SecretShares    map[types.OperatorID][]byte
	SharePublicKeys map[types.OperatorID][]byte
	PartyData       map[types.OperatorID]*DkgPartyData
}

func (s DkgPartyDataSet) VkVec() [][]byte {
	out := make([][]byte, len(s.SharePublicKeys))
	for i := 0; i < len(s.SharePublicKeys); i++ {
		out[i] = s.SharePublicKeys[types.OperatorID(i+1)]
	}
	return out
}

func (s DkgPartyDataSet) R1(operatorId types.OperatorID) *keygen.ParsedMessage {
	return &keygen.ParsedMessage{
		Header: &base.MessageHeader{
			MsgType: int32(dkg.ProtocolMsgType),
			Sender:  uint32(operatorId),
		},
		Body: &keygen.KeygenMsgBody{
			Round1: &keygen.Round1Msg{
				Commitment: s.PartyData[operatorId].Commitment,
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R2(operatorId types.OperatorID) *keygen.ParsedMessage {
	return &keygen.ParsedMessage{
		Header: &base.MessageHeader{
			MsgType: int32(dkg.ProtocolMsgType),
			Sender:  uint32(operatorId),
		},
		Body: &keygen.KeygenMsgBody{
			Round2: &keygen.Round2Msg{
				DeCommitment: s.PartyData[operatorId].DeCommitment,
				BlindFactor:  s.PartyData[operatorId].BlindFactor,
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R3(operatorId types.OperatorID, receiver types.OperatorID) *keygen.ParsedMessage {
	return &keygen.ParsedMessage{
		Header: &base.MessageHeader{
			MsgType:  int32(dkg.ProtocolMsgType),
			Sender:   uint32(operatorId),
			Receiver: uint32(receiver),
		},
		Body: &keygen.KeygenMsgBody{
			Round3: &keygen.Round3Msg{
				Share: s.PartyData[operatorId].Shares[receiver],
			},
		},
		Signature: nil,
	}
}

func (s DkgPartyDataSet) R4(operatorId types.OperatorID) *keygen.ParsedMessage {
	return &keygen.ParsedMessage{
		Header: &base.MessageHeader{
			MsgType: int32(dkg.ProtocolMsgType),
			Sender:  uint32(operatorId),
		},
		Body: &keygen.KeygenMsgBody{
			Round4: &keygen.Round4Msg{
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

var BaseInstance = func() dkg.Protocol {
	state := keygen.EmptyKeygen(2, 4)
	partyData := TestSuiteThreeOfFourSmallValues().PartyData[1]
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

var SignDKGMsg = func(sk *ecdsa.PrivateKey, msg *keygen.ParsedMessage) *keygen.ParsedMessage {
	domain := types.PrimusTestnet
	sigType := types.DKGSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig, _ := crypto.Sign(r, sk)

	return &keygen.ParsedMessage{
		Header:    msg.Header,
		Body:      msg.Body,
		Signature: sig,
	}
}

func TestSuiteThreeOfFour() DkgPartyDataSet {
	return DkgPartyDataSet{
		PublicKey: h2b("8adbbb94ab3b4741e651e20255ad33e73483d0c83181b3aedad5fec9d648e952bfd4baeef8236781ce00300d17ae31ad"),
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
				Commitment:  h2b("c242ab540de0ccf6107ddc1a3b494ec4f93621d4b09a443a025f83b90ce8b462"),
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
				Commitment:  h2b("ae24241d1b463462c2acb4c0f37c1708b5ff4f80bf50998e97842afe1a3c1eb1"),
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
				Commitment:  h2b("68997f5b0cceff6e1ea3b1b483ace3778dae3db2e25ebd5f819b0be56eabce12"),
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
				Commitment:  h2b("924699f92ba71c39903f491fdb08a13c1dbe5eeda83fdf164110cca7e03ae054"),
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