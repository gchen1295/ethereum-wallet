package utils

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// GweiToEth converts a gwei value to eth.
func GweiToEth(value *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetInt(value)
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

// EthToGwei converts an eth value to gwei.
func EthToGwei(balance *big.Float) *big.Int {
	ethVal := new(big.Float).Mul(balance, big.NewFloat(math.Pow10(18)))
	gwei := new(big.Int)
	ethVal.Int(gwei)
	return gwei
}

// PublicKey retrieves the public key given a private key.
func PublicKey(private string) (*ecdsa.PublicKey, error) {
	privKey, err := crypto.HexToECDSA(private)
	if err != nil {
		log.Fatalln(err)
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	return publicKeyECDSA, nil
}

// GeneratePrivateKey generates a private ecdsa private key.
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// SignPayload signs a payload with a private key.
func SignPayload(payload []byte, privkey *ecdsa.PrivateKey) (string, error) {
	hashedBody := crypto.Keccak256Hash(payload).Hex()
	sig, err := crypto.Sign(accounts.TextHash([]byte(hashedBody)), privkey)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(privkey.PublicKey).Hex() + ":" + hexutil.Encode(sig), nil
}
