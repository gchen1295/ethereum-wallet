package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func deriveAddress(path accounts.DerivationPath, root *hdkeychain.ExtendedKey) (common.Address, error) {
	publicKeyECDSA, err := derivePublicKey(path, root)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func derivePublicKey(path accounts.DerivationPath, root *hdkeychain.ExtendedKey) (*ecdsa.PublicKey, error) {
	privKey, err := derivePrivateKey(path, root)
	if err != nil {
		return nil, err
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

func derivePrivateKey(path accounts.DerivationPath, root *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	var err error
	derivedKey := root

	for _, bit := range path {
		derivedKey, err = derivedKey.Derive(bit)
		if err != nil {
			return nil, err
		}
	}

	privKey, err := derivedKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return privKey.ToECDSA(), nil
}

func GenerateMnemonic() (string, error) {
	m, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(m)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func encryptFile(filename string, data []byte, passphrase string) error {
	f, _ := os.Create(filename)
	defer f.Close()
	d, err := encrypt(data, passphrase)
	if err != nil {
		return err
	}
	f.Write(d)
	return nil
}

func decryptFile(filename string, passphrase string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return decrypt(data, passphrase)
}
