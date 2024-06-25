package messages

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MarshalJustificationsWithoutFullData tests marshalling justifications without full data.
func MarshalJustificationsWithoutFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndFullData(ks.Shares[1], 1, 2, nil),
	}

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRoundAndFullData(ks.Shares[1], types.OperatorID(1), 1, nil),
	}

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs))

	r, err := hex.DecodeString("38977ca5244c072549564f2fd90a88fd7e3cf5124ebd3c77154c042ea3ed8aa0")
	if err != nil {
		panic(err)
	}

	b, err := hex.DecodeString("97b88bdd65b68dc9f35bcaa236a41c2e279e0d60d1f3b88f43b5184ce29b2bcde2ff79f4e9877e2655eb269f103a7ed7070469f2fc0f53db4676f85f7602d42ae8715f2a9e8fc08e75ee09ecb213e0e94d01d95a5214d165569ada5f8f9434a96c000000740000005402000001000000000000000000000000000000000000000000000002000000000000004c000000be956fb7df4ef37531682d588320084fc914c3f0fed335263e5b44062e6c29b4000000000000000050000000180100000102030404000000b0a85408a6085a2dc0d571f37a6d687f1b76bb152f8764fb6111048d91642a656e66bb0e260981b752980cfff4d2c0ff0a79e5b1877d8f5cd62f4d5c59c13695d672f2b4b83c5596b9ec3838d5eea245bb6086676989fefcd65fbb6bf93678806c00000074000000c400000001000000000000000300000000000000000000000000000002000000000000004c000000e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85500000000000000005000000050000000010203040400000094f75ab09524ef2d4fff17e641437abcc9aef029962b4298c13ee7914f88e10463f27edb5d12f19c211df499a10f3852126184cc2acd150f087f298bd4ff42d917f94910e82fd8a6974400626c99ea2dacac83ead8cafd1b913359b23c5734a36c00000074000000c400000001000000000000000100000000000000000000000000000001000000000000004c000000e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b8550000000000000000500000005000000001020304010203040506070809010203040506070809010203040506070809")
	if err != nil {
		panic(err)
	}

	return &tests.MsgSpecTest{
		Name: "marshal justifications",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][32]byte{
			*(*[32]byte)(r),
		},
	}
}
