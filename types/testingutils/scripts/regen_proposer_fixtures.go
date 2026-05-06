//go:build ignore

// regen_proposer_fixtures rewrites the BLS-typed JSON values inside the
// capellaBlock and denebBlockContents literals of beacon_node_consts.go to
// valid compressed BLS12-381 points. BLS-typed positions are discovered by
// JSON-decoding each literal into its go-eth2-client struct and walking the
// known BLS fields via internal/proposerbls; substitution is then a scoped
// text replacement of the discovered values inside the original literal. KZG
// commitments, KZG proofs, transactions, blobs, and other 48/96-byte non-BLS
// fields are never touched because they are not on the BLS-field walk.
//
// Run from types/testingutils/: go run scripts/regen_proposer_fixtures.go
package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types/testingutils/internal/proposerbls"
)

const (
	fixtureFilePath = "beacon_node_consts.go"
	validatorSKHex  = "3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f"
	sigMessage      = "ssv-spec proposer fixture"
)

func init() {
	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(err)
	}
	if err := bls.SetETHmode(bls.EthModeDraft07); err != nil {
		panic(err)
	}
}

func main() {
	var sk bls.SecretKey
	if err := sk.SetHexString(validatorSKHex); err != nil {
		log.Fatalf("parse validator SK: %v", err)
	}
	validPubkey := "0x" + sk.GetPublicKey().SerializeToHexStr()
	validSig := "0x" + sk.SignByte([]byte(sigMessage)).SerializeToHexStr()

	src, err := os.ReadFile(fixtureFilePath)
	if err != nil {
		log.Fatalf("read %s: %v", fixtureFilePath, err)
	}

	out := bytes.Clone(src)
	for _, name := range []string{"capellaBlock", "denebBlockContents"} {
		rewriteLiteral(out, name, validPubkey, validSig)
	}

	if err := os.WriteFile(fixtureFilePath, out, 0o644); err != nil {
		log.Fatalf("write %s: %v", fixtureFilePath, err)
	}
	fmt.Printf("rewrote BLS placeholders in %s\n", fixtureFilePath)
}

// rewriteLiteral substitutes BLS values in-place inside the backtick literal
// `var <name> = []byte(` ... `)`. Substitution preserves length, so surrounding
// file offsets are unchanged.
func rewriteLiteral(out []byte, name, validPubkey, validSig string) {
	prefix := []byte("var " + name + " = []byte(`")
	start := bytes.Index(out, prefix)
	if start < 0 {
		log.Fatalf("literal %q not found", name)
	}
	contentStart := start + len(prefix)
	rel := bytes.Index(out[contentStart:], []byte("`)"))
	if rel < 0 {
		log.Fatalf("literal %q has no closing backtick", name)
	}
	contentEnd := contentStart + rel

	literal := out[contentStart:contentEnd]
	pubkeys, sigs := discoverBLSValues(name, literal)

	rewritten := bytes.Clone(literal)
	for v := range pubkeys {
		rewritten = bytes.ReplaceAll(rewritten, []byte(`"`+v+`"`), []byte(`"`+validPubkey+`"`))
	}
	for v := range sigs {
		rewritten = bytes.ReplaceAll(rewritten, []byte(`"`+v+`"`), []byte(`"`+validSig+`"`))
	}
	if len(rewritten) != contentEnd-contentStart {
		log.Fatalf("%s: substitution changed literal length", name)
	}
	copy(out[contentStart:contentEnd], rewritten)
}

func discoverBLSValues(name string, literal []byte) (pubkeys, sigs map[string]struct{}) {
	pubkeys = map[string]struct{}{}
	sigs = map[string]struct{}{}
	visitor := proposerbls.Visitor{
		Pubkey:    func(_ string, b []byte) { pubkeys["0x"+hex.EncodeToString(b)] = struct{}{} },
		Signature: func(_ string, b []byte) { sigs["0x"+hex.EncodeToString(b)] = struct{}{} },
	}

	switch name {
	case "capellaBlock":
		var blk capella.BeaconBlock
		if err := json.Unmarshal(literal, &blk); err != nil {
			log.Fatalf("decode %s: %v", name, err)
		}
		proposerbls.WalkCapella(blk.Body, visitor)
	case "denebBlockContents":
		var bc apiv1deneb.BlockContents
		if err := json.Unmarshal(literal, &bc); err != nil {
			log.Fatalf("decode %s: %v", name, err)
		}
		proposerbls.WalkDeneb(bc.Block.Body, visitor)
	default:
		log.Fatalf("unknown literal %q", name)
	}
	return
}
