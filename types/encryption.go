package types

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

var keySize = 2048

// GenerateKey using rsa random generate keys
func GenerateKey() ([]byte, []byte, error) {
	// generate random private key (secret)
	sk, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to generate private key")
	}

	// convert to bytes
	skPem := PrivateKeyToPem(sk)
	pkPem, err := GetPublicKeyPem(sk)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to marshal public key")
	}
	return skPem, pkPem, nil
}

// Decrypt with secret key (base64) and bytes, return the encrypted key string
func Decrypt(sk *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, sk, cipherText)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decrypt key")
	}
	return decrypted, nil
}

// Encrypt with secret key (base64) the bytes, return the encrypted key string
func Encrypt(pk *rsa.PublicKey, plainText []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pk, plainText)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decrypt key")
	}
	return encrypted, nil
}

// PemToPrivateKey return rsa private key from pem
func PemToPrivateKey(skPem []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(skPem)
	if block == nil {
		return nil, errors.New("failed to decode PEM")
	}

	// Parse key as PKCS1
	parsedSk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse private key")
	}

	return parsedSk, nil
}

// PemToPublicKey return rsa public key from pem
func PemToPublicKey(pkPem []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pkPem)
	if block == nil {
		return nil, errors.New("failed to decode PEM")
	}
	parsedPk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse public key")
	}
	if ret, ok := parsedPk.(*rsa.PublicKey); ok {
		return ret, nil
	}
	return nil, errors.Wrap(err, "Failed to parse public key")
}

// PrivateKeyToPem converts privateKey to pem encoded
func PrivateKeyToPem(sk *rsa.PrivateKey) []byte {
	pemBytes := x509.MarshalPKCS1PrivateKey(sk)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pemBytes,
		},
	)
}

// GetPublicKeyPem get public key from private key and return []byte represent the public key
func GetPublicKeyPem(sk *rsa.PrivateKey) ([]byte, error) {
	pkBytes, err := x509.MarshalPKIXPublicKey(&sk.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal private key")
	}
	pemByte := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pkBytes,
		},
	)

	return pemByte, nil
}
