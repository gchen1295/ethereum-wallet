package wallet

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignData requests the wallet to sign the hash of the given data
// It looks up the account specified either solely via its address contained within,
// or optionally with the aid of any location metadata from the embedded URL field.
//
// If the wallet requires additional authentication to sign the request (e.g.
// a password to decrypt the account, or a PIN code to verify the transaction),
// an AuthNeededError instance will be returned, containing infos for the user
// about which fields or actions are needed. The user may retry by providing
// the needed details via SignDataWithPassphrase, or by other means (e.g. unlock
// the account in a keystore).
func (wallet *Wallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	return wallet.SignHash(account, crypto.Keccak256(data))
}

// SignDataWithPassphrase is identical to SignData, but also takes a password
// NOTE: there's a chance that an erroneous call might mistake the two strings, and
// supply password in the mimetype field, or vice versa. Thus, an implementation
// should never echo the mimetype or return the mimetype in the error-response
func (wallet *Wallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	return wallet.SignHash(account, crypto.Keccak256(data))
}

// SignText requests the wallet to sign the hash of a given piece of data, prefixed
// by the Ethereum prefix scheme
// It looks up the account specified either solely via its address contained within,
// or optionally with the aid of any location metadata from the embedded URL field.
//
// If the wallet requires additional authentication to sign the request (e.g.
// a password to decrypt the account, or a PIN code o verify the transaction),
// an AuthNeededError instance will be returned, containing infos for the user
// about which fields or actions are needed. The user may retry by providing
// the needed details via SignTextWithPassphrase, or by other means (e.g. unlock
// the account in a keystore).
//
// This method should return the signature in 'canonical' format, with v 0 or 1
func (wallet *Wallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	return wallet.SignHash(account, accounts.TextHash(text))
}

// SignTextWithPassphrase is identical to Signtext, but also takes a password
func (wallet *Wallet) SignTextWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
	return wallet.SignHash(account, accounts.TextHash(hash))
}

// SignTx requests the wallet to sign the given transaction.
//
// It looks up the account specified either solely via its address contained within,
// or optionally with the aid of any location metadata from the embedded URL field.
//
// If the wallet requires additional authentication to sign the request (e.g.
// a password to decrypt the account, or a PIN code to verify the transaction),
// an AuthNeededError instance will be returned, containing infos for the user
// about which fields or actions are needed. The user may retry by providing
// the needed details via SignTxWithPassphrase, or by other means (e.g. unlock
// the account in a keystore).
func (wallet *Wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	wallet.mux.RLock()
	path, ok := wallet.paths[account.Address]
	wallet.mux.RUnlock()
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := wallet.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(chainID)
	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}

	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

func (wallet *Wallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	wallet.mux.RLock()
	path, ok := wallet.paths[account.Address]
	wallet.mux.RUnlock()
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := wallet.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(chainID)
	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}

	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

// SignHash signs a payload hash with a private key from an address.
func (wallet *Wallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	wallet.mux.RLock()
	path, ok := wallet.paths[account.Address]
	wallet.mux.RUnlock()
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := wallet.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(accounts.TextHash(hash), privateKey)
}
