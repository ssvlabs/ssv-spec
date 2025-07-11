package types

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// RSA private key cannot be unmarshalled directly as it uses big.Int
func (s *OperatorSigner) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		SSVOperatorSK string     `json:"SSVOperatorSK"`
		OperatorID    OperatorID `json:"OperatorID"`
	}{
		SSVOperatorSK: fmt.Sprintf("N:%s,E:%s,D:%s", s.SSVOperatorSK.N.String(), strconv.Itoa(s.SSVOperatorSK.E), s.SSVOperatorSK.D.String()),
		OperatorID:    s.OperatorID,
	})
}

func (s *OperatorSigner) UnmarshalJSON(data []byte) error {
	aux := &struct {
		SSVOperatorSK string     `json:"SSVOperatorSK"`
		OperatorID    OperatorID `json:"OperatorID"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the private key string format "N:x,E:y,D:z"
	var nStr, eStr, dStr string
	parts := strings.Split(aux.SSVOperatorSK, ",")
	if len(parts) != 3 {
		return fmt.Errorf("invalid private key format: expected 3 parts, got %d", len(parts))
	}
	fmt.Sscanf(parts[0], "N:%s", &nStr)
	fmt.Sscanf(parts[1], "E:%s", &eStr)
	fmt.Sscanf(parts[2], "D:%s", &dStr)

	// Convert strings back to big integers
	n := new(big.Int)
	n.SetString(nStr, 10)

	e, err := strconv.Atoi(eStr)
	if err != nil {
		return err
	}

	d := new(big.Int)
	d.SetString(dStr, 10)

	// Reconstruct the RSA private key
	s.SSVOperatorSK = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: e,
		},
		D: d,
	}
	s.OperatorID = aux.OperatorID

	return nil
}
