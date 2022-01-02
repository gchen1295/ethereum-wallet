package vault

import (
	"crypto/ecdsa"

	"errors"

	"strconv"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

const DEFAULT_BASE_DERIVATION_PATH = "m/44'/60'/0'/0/"

// Vault is a secure vault that manages access to wallets.
type Vault struct {
	sync.RWMutex
	// idx of derivation path
	nonce int

	configPath string
	ks         *keystore.KeyStore
}

// NewVault is a wrapper around the geth keystore
func NewVault(configPath string) *Vault {
	ks := keystore.NewKeyStore(configPath+"/vault", keystore.StandardScryptN, keystore.StandardScryptP)
	return &Vault{
		nonce: len(ks.Wallets()) - 1,
		ks:    ks,
	}
}

// ImportHDWallet setups and loads the HD Wallet that is used to derive the rest of our wallets.
func (v *Vault) ImportHDWallet(mnemonic, passphrase []byte) error {
	// start deriving accounts from the given seed phrase
	seed := bip39.NewSeed(string(mnemonic), "")
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}
	seed = nil
	// load first account from wallet
	accPk, err := v.deriveAccount(masterKey)
	if err != nil {
		return err
	}

	// encrypt and store mk
	// pk, err := masterKey.ECPrivKey()
	// if err != nil {
	// 	return err
	// }
	if err = encryptFile(v.configPath+"/vault/root.txt", []byte(masterKey.String()), string(passphrase)); err != nil {
		return err
	}

	// remove from memory
	masterKey = nil

	// // import into keystore
	// _, err = v.ks.ImportECDSA(pk.ToECDSA(), string(passphrase))
	// if err != nil {
	// 	return err
	// }

	// import first account into keystore
	_, err = v.ks.ImportECDSA(accPk, string(passphrase))
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) getRoot(passphrase []byte) (*hdkeychain.ExtendedKey, error) {
	decrypted, err := decryptFile(v.configPath+"/vault/root.txt", string(passphrase))
	if err != nil {
		return nil, err
	}

	mk, err := hdkeychain.NewKeyFromString(string(decrypted))
	if err != nil {
		return nil, err
	}

	return mk, nil
}

// UpdatePassphrase updates the password for all accounts in the keystore.
func (v *Vault) UpdatePassphrase(old, new []byte) error {
	for _, acc := range v.ks.Accounts() {
		err := v.ks.Update(acc, string(old), string(new))
		if err != nil {
			return err
		}
	}

	return nil
}

// Accounts returns the underlying keystores accounts
func (v *Vault) Accounts() []accounts.Account {
	v.RLock()
	defer v.RUnlock()
	accounts := v.ks.Accounts()
	if len(accounts) == 0 {
		return accounts
	}
	return accounts[1:]
}

// Addresses returns the hex-encoded public addresses of all accounts in the keystore
func (v *Vault) Addresses() []string {
	v.RLock()
	defer v.RUnlock()
	return v.addresses()
}

func (v *Vault) addresses() []string {
	accs := v.ks.Accounts()

	addresses := []string{}
	for _, acc := range accs {
		addresses = append(addresses, acc.Address.Hex())
	}

	return addresses
}

// CreateWallet creates a new wallet and adds it to the keystore with the given passphrase.
func (v *Vault) CreateWallet(passphrase []byte) ([]string, error) {
	v.Lock()
	defer v.Unlock()

	ksAccounts := v.ks.Accounts()
	if len(ksAccounts) == 0 {
		return nil, errors.New("keystore not configured")
	}

	mk, err := v.getRoot([]byte(passphrase))
	if err != nil {
		return nil, err
	}

	// derive account
	newAcc, err := v.deriveAccount(mk)
	if err != nil {
		return nil, err
	}

	// import first account into keystore
	_, err = v.ks.ImportECDSA(newAcc, "1")
	if err != nil {
		return nil, err
	}

	return v.addresses(), nil
}

func (v *Vault) ImportWallet(pk, passphrase []byte) ([]string, error) {
	privKey, err := crypto.HexToECDSA(string(pk))
	if err != nil {
		return nil, err
	}

	v.Lock()
	defer v.Unlock()
	_, err = v.ks.ImportECDSA(privKey, string(passphrase))
	if err != nil {
		return nil, err
	}

	return v.addresses(), nil
}

func (v *Vault) RemoveWallet(address, passphrase []byte) ([]string, error) {
	v.Lock()
	defer v.Unlock()

	err := v.ks.Delete(accounts.Account{
		Address: common.BytesToAddress(address),
	}, string(passphrase))
	if err != nil {
		return nil, err
	}

	return v.addresses(), nil
}

func (v *Vault) deriveAccount(root *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	// Get metamask default derivation path
	v.nonce++
	path, _ := accounts.ParseDerivationPath(DEFAULT_BASE_DERIVATION_PATH + strconv.Itoa(v.nonce))

	privKey, err := derivePrivateKey(path, root)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
