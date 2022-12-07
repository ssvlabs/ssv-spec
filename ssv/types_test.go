package ssv

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestBLSVerify(t *testing.T) {
	var wg sync.WaitGroup

	types.InitBLS()
	sk := bls.SecretKey{}
	sk.SetByCSPRNG()
	pk := sk.GetPublicKey()
	sign := sk.SignByte([]byte{1, 2, 3, 4})

	start := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(pk bls.PublicKey) {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				require.True(t, sign.VerifyByte(&pk, []byte{1, 2, 3, 4}))
			}
		}(*pk)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("took %s\n", elapsed)
}
