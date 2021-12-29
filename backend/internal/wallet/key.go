package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKey returns the private key of an account.
func (wallet *Wallet) PrivateKey(account accounts.Account) (*ecdsa.PrivateKey, error) {
	path, err := accounts.ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}

	return wallet.derivePrivateKey(path)
}

// PublicKey returns the public key of an account.
func (wallet *Wallet) PublicKey(account accounts.Account) (*ecdsa.PublicKey, error) {
	path, err := accounts.ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}

	return wallet.derivePublicKey(path)
}

// Address returns the address of an account.
func (wallet *Wallet) Address(account accounts.Account) (common.Address, error) {
	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*publicKey), nil
}

func (wallet *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	wallet.mux.RLock()
	derivedKey := wallet.masterKey
	wallet.mux.RUnlock()

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

func (wallet *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	privKey, err := wallet.derivePrivateKey(path)
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

func (wallet *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := wallet.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}
