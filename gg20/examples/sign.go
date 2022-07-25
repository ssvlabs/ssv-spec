package main

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func init() {
	_ = bls.Init(bls.BLS12_381)
	_ = bls.SetETHmode(bls.EthModeDraft07)
}

var skFromHex = func(str string) *bls.SecretKey {
	types.InitBLS()
	ret := &bls.SecretKey{}
	if err := ret.DeserializeHexStr(str); err != nil {
		panic(err.Error())
	}
	return ret
}

func main() {
	shares := map[types.OperatorID]*bls.SecretKey{
		1: skFromHex("5f4711a796c1116b5118ec35279fb64d551d9b38813d2939954dd2df5160d3d9"),
		2: skFromHex("48e4c0a38e90f9352d1d09489446443ebd17b1904f4f0002fe894c2c3f62457a"),
		3: skFromHex("65dc7c179f68347cf12f86e1c51e54e8aeeed579d4c715082bb8a0382c1a8153"),
		4: skFromHex("42409cb09fa945fa6a168cf8b0861045d6e562f211a70c4a1cdbcf0417898763"),
	}
	bytes, _ := hex.DecodeString("1fa0068233c6c0ffedd8fb1c6dea0fd13d67d5a558b803acecf15601e36dd8d9")
	for id, key := range shares {
		sig := key.SignByte(bytes)
		fmt.Printf("sig %v: %x\n", id, sig.Serialize())
	}
}
