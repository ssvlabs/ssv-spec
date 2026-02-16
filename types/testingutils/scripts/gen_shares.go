package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/herumi/bls-eth-go-binary/bls"
)

func init() {
	// Initialize BLS with the ETH2 curve
	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(err)
	}
	if err := bls.SetETHmode(bls.EthModeDraft07); err != nil {
		panic(err)
	}
}

func main() {

	count := 3000

	output, err := os.Create("shares.json")
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}

	// Generate the key sets
	for i := 1; i <= count; i++ {
		entry := generateKeySetEntry(i)
		_, err = fmt.Fprint(output, entry)
		if err != nil {
			fmt.Printf("Failed to write to output file: %v\n", err)
			os.Exit(1)
		}
	}

	err = output.Close()
	if err != nil {
		fmt.Printf("Failed to close output file: %v\n", err)
		os.Exit(1)
	}
}

func generateKeySetEntry(index int) string {
	// Create a random validator secret key
	var validatorSK bls.SecretKey
	validatorSK.SetByCSPRNG()

	// Get the public key
	validatorPK := validatorSK.GetPublicKey()

	// Create shares using Shamir's secret sharing
	// For threshold t=3 out of n=4, we need a polynomial of degree t-1 = 2
	shareCount := 4
	threshold := 3

	// Create a polynomial with the master secret as the constant term
	// The polynomial is: f(x) = a0 + a1*x + a2*x^2 where a0 = validatorSK
	msk := make([]bls.SecretKey, threshold)
	msk[0] = validatorSK
	for i := 1; i < threshold; i++ {
		msk[i].SetByCSPRNG()
	}

	// Generate shares for each operator (IDs 1, 2, 3, 4)
	shares := make(map[int]*bls.SecretKey)
	for opID := 1; opID <= shareCount; opID++ {
		var share bls.SecretKey
		// Evaluate polynomial at point opID
		var id bls.ID
		if err := id.SetDecString(fmt.Sprintf("%d", opID)); err != nil {
			panic(err)
		}
		if err := share.Set(msk, &id); err != nil {
			panic(err)
		}
		shares[opID] = &share
	}

	// Format the output
	return formatKeySetEntry(index, &validatorSK, validatorPK, shares)
}

func formatKeySetEntry(index int, validatorSK *bls.SecretKey, validatorPK *bls.PublicKey, shares map[int]*bls.SecretKey) string {
	// Format validator SK as hex (without 0x prefix)
	skHex := validatorSK.SerializeToHexStr()

	// Format validator PK as hex (without 0x prefix)
	pkHex := hex.EncodeToString(validatorPK.Serialize())

	// Entry string
	result := fmt.Sprintf("%d: {\n", index)
	result += fmt.Sprintf("\tValidatorSK:      skFromHex(\"%s\"),\n", skHex)
	result += fmt.Sprintf("\tValidatorPK:      pkFromHex(\"%s\"),\n", pkHex)
	result += "\tShareCount:       4,\n"
	result += "\tThreshold:        3,\n"
	result += "\tPartialThreshold: 2,\n"
	result += "\tShares: map[types.OperatorID]*bls.SecretKey{\n"

	// Add shares
	for opID := 1; opID <= 4; opID++ {
		shareHex := shares[opID].SerializeToHexStr()
		result += fmt.Sprintf("\t\t%d: skFromHex(\"%s\"),\n", opID, shareHex)
	}

	result += "\t},\n"
	result += "\tOperatorKeys: TestingOperatorKeys4Map,\n"
	result += "\tDKGOperators: TestingDKGOperators4Map,\n"
	result += "},\n"

	return result
}
