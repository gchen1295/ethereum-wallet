package engine

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"

	"encoding/hex"
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

// Vault is a vault that manages access to wallets.
type Vault struct {
	sync.RWMutex
	// idx of derivation path
	nonce int

	ks *keystore.KeyStore
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
	if err = encryptFile(APPLICATION_CONFIG_PATH+"/vault/root.txt", []byte(masterKey.String()), string(passphrase)); err != nil {
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
	decrypted, err := decryptFile(APPLICATION_CONFIG_PATH+"/vault/root.txt", string(passphrase))
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

	mk, err := v.getRoot([]byte("1"))
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

// UTILS
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
