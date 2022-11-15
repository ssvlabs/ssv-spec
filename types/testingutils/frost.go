package testingutils

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"encoding/hex"
	"math/big"
	mrand "math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var (
	SessionPKs map[types.OperatorID]string = map[types.OperatorID]string{
		1: "036ff75a45bb43f1190f89838326ed4f2e090293184e56ff4a01a1a6db548fbae6",
		2: "038680ce08d663c436ddb98265dd26a0c775bf4728ab5ae385671eeb5b87ab08e7",
		3: "0204470b016f243d34ff27d8c869c3b8012612232390d8d3259bc40bf4dc3c4551",
		4: "0328893f709ce7ad1ee70f393cf5ba152fc11043043f0a0acb1591923ebea52dbd",
	}

	Round1 = map[types.OperatorID]struct {
		Commitments []string
		Shares      map[uint32]string
		ProofS      string
		ProofR      string
	}{
		1: {
			Commitments: []string{
				"b1fc06dbfe90a494bbda98ce51663eeec74134b3faf1f667f872ffff9b7a7747e31fab29b8d2d7e30e5b8dfe26f4d552",
				"821696e875f72e1b85be294972c02787f839626077298791c4e7c553c6762433c2b7cfeba151752316ca926fa9ba236c",
				"8f66b1662e494ad37239040812743d7d68e256f57a59beb310c93e2df17b41e51d534168a0abc93ebb56785524f0586a",
			},
			ProofS: "65e597c5000b9b6cc2953b7c44b93e5ca5511f7c2cd81c7680bede9bbb9828b6",
			ProofR: "4429058e0da3ca520b3d311dd9eae36c132151f99c804278f59eb3ddc128407e",
			Shares: map[uint32]string{
				2: "042a60c3ef5efbcc802221c346ac73e997c449a66ecb9dbe4f63108c67678e62af7e151cc8b8289d1abb43ce49361490ed889778d44d6cb623ae169c300881d6124cb9ce563eaa36f0b96eac28a033e582f35a7eca9d0b73d585428a647f6614e42b54946930066fcfb71d4c767b85b1b146e6cedc5ad5f743b1c8080f776460e6",
				3: "04ce72a9cb7c24930acc3ee568e9999254d425e99e71d72e0500dd81d1818ea303ae8b7b0c075420403a9a7c3f5ebe8b2df2e702601aa49932b5b771b14b2bb3a0730dc4c7f9c1422f14b7198361b631e1e04c89301dca6796bc4b417c5888e1820b2906d96571550425655c101a16357be1bc2c162ca7ba0062fe4890fddae895",
				4: "0446fe513897e4ba460aeefc76c5e8ac2f60f7fd0e93797c406c93dbca02fd7f55e79e9b322b9641b1a77d11adf4f38feb132270b3d28463c5e3be517c1d80ae4a5d262ab28924fd6187cde8490979e7a0e3ff7896ea7fb9a17bd6095e4afe913967655f9a00999f9fbd7c16e262a3686802b2817e5ca216aa2029c420042e5fd1",
			},
		},
		2: {
			Commitments: []string{
				"af8a4e775d7d80b5f3b74eb1c2f997b9d218ac414797b6724fbc970cccf75adfc298af8a90f45e9ba9eae28bf949fdea",
				"81087c05df25266a244fdcf31542bf8dabd6a7ccfc54e229e51fe40b25f10a862aa480b23b09a0d1755864c2f0223fd0",
				"ab31eb44022999aee599b6e6dfb49accd0985fa584dbf458405fa541dafeb3e48eff70249cba0c1f4a1cb80163c7e15a",
			},
			ProofS: "458e8c8e76a908fb6de9c47d61759d6dd77731ef014890ab9532695fdcb96cbf",
			ProofR: "34044356f2931faa528e262d33052a899498bf594a35593a75e0c1f4571d7640",
			Shares: map[uint32]string{
				1: "04e483c6b63c6e6c0ef9677062cfd38030bebe25c2610aa574da04992bad387ae58c1a081a64bd28a9b16c0604a5d4ed93e6a1cb23a71b49479f7593acdf4f24479813a42970682b4645163aafd14b01d05f9d11d348949829d71f907759551dc934c03b392f7a77d6aeae96fae550e15bb4d24e7dd16f7109b06fb4633851727e",
				3: "04e94fa027890da52caf76d15463e5534259d4e80f9bc4ad6c0f770d8c8468be5ee58ecf3731320a7ac33d656856a9cd6cf93a1a29161e819dfec389081e25edf062c9dd07d07e98a70cf8e0c16112ef2a8015b1146786ef7d0764abf8fba9a3d38f137e48f597eb531dda88f0edc68f3c160f9d543a5d5a215ebcb30b2542a78b",
				4: "04765affeae2264ce133b3c5908173865524a97d8adcfd624193b5c094539b2bcfbaaa989c07342ac6e44a70d2abb193c300392fe8ed95edd8e0d9715451f996af51627b703670f2465883c33c62b3cc6544342dae70057c8c92c3eac325aedc9715944dabcc8c77846b00d662ff4c20a6cf21dd00c2b9573116fba72703328016",
			},
		},
		3: {
			Commitments: []string{
				"a2e4badfdc21372375a741669676f51bfea8a9e21dbd16286217d5931f690c8aa9c06bf0138c648d61b1acee7bb68198",
				"8deb67f209dc09235c4627f750b2a95ace958f97baafcd95560517f39e4a0d3bec31cddb1fc52f9769c3df44d15e8a63",
				"87addd83d2a222ca610102e401da2a2070d18c32ca1b089cdc683f45c3149fbeea3312451d2c658ff41ea5330e98319e",
			},
			ProofS: "455bfc61aff5b8f90eb32a2444ed72fac45cf69533ac2d343b9beac2ac8eb08e",
			ProofR: "018330f5d32492fbc2a3e65dd80c33f63122d743affff9133cbf448b66f3d5a1",
			Shares: map[uint32]string{
				1: "04a93822a49b25a5f0c2694596cb0f307ed0eefb35ada9f749b827df128a2155bbdfd027e225513acab4a8acfbcd329c41cf6c00def8c37b10ca0d7be25a92fcc19b872a79069f1a6515e66a893fe6e0c59b1cea6fd93dba3bf900d22290defd68c6a4c9794ab054147ad48fa3edb0ee63cfe649164d1eec46e110436b86fa6f4a",
				2: "045234ca3f8e31bf2cc0f0350c68db67ea16a3cf9f1ac439dbefa486ea5c003592d8b057b8370237bd7560ac84b02c9c0d6bf90352c2fd482aa68282a16b88f05aad6900642af14353a88b7464410803468dbebe83e611c84a74f22226355ab39ed86f90a5ba76620fdffaae9b00b64b14293a7f09a323b5b2fcc127832b599662",
				4: "04757e324022b39f7b27ef4d32dc38d09978623bbfc4ed121529e67e3f3a684ae5f8b3cdcbb983767fdd3d26b6f714cb349d81ff8d8b3b5fbb659b4d03234c94d878501438f98e3ed493250693624397a3ca7ef827c9285b4a4af1b56f321f52cd59f825dcd5db095875b23c420f25cea6928d2cdfdfb53f25b96219ebd63726d7",
			},
		},
		4: {
			Commitments: []string{
				"a764e878ade532a2ec3fa8f8813fddd2f33a3318ec2652354d41d90e198f3c17b4b8735d4d5984c5f5df650d082d78c9",
				"95146fe1402bce6dd575e313294d7b02ef5d79eeaeb26ce3da3e02eb93d93a46a7b4fe1db0c842b3e2de8503c0c47e28",
				"967597e453bb617e77dd2341c3ace40a1beb4795962c720c67fec55bfba4ae628b2171b33ab3ea2148c6c287a97e0b34",
			},
			ProofS: "259cf2af7c9662ad2415889901c368ec25aa4abf177b8c01d82635c85ef7cc41",
			ProofR: "498842cfe23f756abb34b50cf866126ebf68aac5af49b28860eca15ea9768be7",
			Shares: map[uint32]string{
				1: "0446fe513897e4ba460aeefc76c5e8ac2f60f7fd0e93797c406c93dbca02fd7f55e79e9b322b9641b1a77d11adf4f38feb132270b3d28463c5e3be517c1d80ae4a5d262ab28924fd6187cde8490979e7a0fee327110c77430f348c8aaaa8dda7dfccc0bf452f7998aca584f2a5e2ca43b2996603124a3f50e12ca6c128199e1543",
				2: "041270eb0c929301038afac519236be54fa97f0a428c2d73b46cefc89e198eb8b0058ae7b1ab2a5a778d873b88d07eb44aa944ccf9397672b9078c83f4c7c0bf2ab5d2f8c17070fc16e3d796ccc5cc29becc2e9b260461383ec4a789a6273ed37ac61247539c72fb11f4ea954e9080b98d34e5b9e60d4b208af59d1239dafb81dc",
				3: "04ca7e6fb4b57c8d847cc30f943045f970bb2b0ecfd926ac7d4275dfccb53c62a3eedd9ab9091699992f4073d9daea3dccde2b239e7be5a2f9e42a8601ac21ad1b9be9b36d34d09c5c22eee19a2957436c52b960f37e77a8d8bf184296ab6d7b637776914aaf9e67f02f26bfbb674aaac2ea7419582fe6df68bf7e6ec4a1e163e4",
			},
		},
	}

	Round2 = map[types.OperatorID]struct {
		Vk      string
		Sk      string
		VkShare string
		SkShare string
	}{
		1: {

			Vk:      "871ac12101b6a8dc6edac502d776a10c01fceedf727b10acb40c6f9de0c977c63879ccb895db7d92ad3d4e77144b64ac",
			Sk:      "65f80c8ae56d2385612e206d53a3db6083b86ffb7ea093ddcf40cb003a60014f",
			VkShare: "a20f5933a6a97e7cd2a0ea65e8e55493bccc60031065429dc1184d72bbbb916491da73cc649760b3d739647077944236",
			SkShare: "6c190c998b847e71713336f884badd4b10ba7cfb488e8bb6c2c95e92c6406a0a",
		},
		2: {

			Vk:      "871ac12101b6a8dc6edac502d776a10c01fceedf727b10acb40c6f9de0c977c63879ccb895db7d92ad3d4e77144b64ac",
			Sk:      "65f80c8ae56d2385612e206d53a3db6083b86ffb7ea093ddcf40cb003a60014f",
			VkShare: "8412c202a39d5b68f3e56c34d6ec721d69f9ef9fc12560c76bc53b1a07118f79ce70850c3f0773b45cdf996270c8b52d",
			SkShare: "64425aa3e6d73fae40f81c2a761503a86ccfc3faaef47e3dc3b00f1513fb1ffd",
		},
		3: {

			Vk:      "871ac12101b6a8dc6edac502d776a10c01fceedf727b10acb40c6f9de0c977c63879ccb895db7d92ad3d4e77144b64ac",
			Sk:      "65f80c8ae56d2385612e206d53a3db6083b86ffb7ea093ddcf40cb003a60014f",
			VkShare: "8a819950175eb638494ac7a037362469f8716d13800a58cb2109b8a7b128ff6f3fbf5fe2b507377d41dee71bf2d906f3",
			SkShare: "4e73f6a9f765673bd07cd00327b24e7897f844f9b1d26b72d1f4dc8723902328",
		},
		4: {

			Vk:      "871ac12101b6a8dc6edac502d776a10c01fceedf727b10acb40c6f9de0c977c63879ccb895db7d92ad3d4e77144b64ac",
			Sk:      "65f80c8ae56d2385612e206d53a3db6083b86ffb7ea093ddcf40cb003a60014f",
			VkShare: "925aeda7183eecdc8e9d27ebc32a32bfad15683a36cb2bfa3ca5f450ad92917d46a35a2d4b66dd976defd0ba6306e78e",
			SkShare: "2aade0abbd2ef51a1fc152829992bdbb9233fff851285355ed97c6e8f4ff738b",
		},
	}
)

func SignDKGMsg2(sk *ecdsa.PrivateKey, opID types.OperatorID, msg *dkg.Message) *dkg.SignedMessage {
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

type TestOutcome struct {
	KeygenOutcome TestKeygenOutcome
	BlameOutcome  TestBlameOutcome
}

type TestKeygenOutcome struct {
	ValidatorPK     string
	Share           map[uint32]string
	OperatorPubKeys map[uint32]string
}

func (o TestKeygenOutcome) ToKeygenOutcomeMap(threshold uint64, operators []uint32) map[uint32]*dkg.KeyGenOutput {
	m := make(map[uint32]*dkg.KeyGenOutput)

	opPublicKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range operators {

		pk := &bls.PublicKey{}
		_ = pk.DeserializeHexStr(o.OperatorPubKeys[operatorID])
		opPublicKeys[types.OperatorID(operatorID)] = pk

		share := o.Share[operatorID]
		sk := &bls.SecretKey{}
		_ = sk.DeserializeHexStr(share)

		vk, _ := hex.DecodeString(o.ValidatorPK)

		m[operatorID] = &dkg.KeyGenOutput{
			Share:           sk,
			ValidatorPK:     vk,
			OperatorPubKeys: opPublicKeys,
			Threshold:       threshold,
		}
	}

	return m
}

func ResetRandSeed() {
	src := mrand.NewSource(1)
	src.Seed(12345)
	crand.Reader = mrand.New(src)
}

func GetRandRequestID() dkg.RequestID {
	requestID := dkg.RequestID{}
	for i := range requestID {
		rndInt, _ := crand.Int(crand.Reader, big.NewInt(255))
		if len(rndInt.Bytes()) == 0 {
			requestID[i] = 0
		} else {
			requestID[i] = rndInt.Bytes()[0]
		}
	}
	return requestID
}

type TestBlameOutcome struct {
	Valid        bool
	BlameMessage []byte
}
