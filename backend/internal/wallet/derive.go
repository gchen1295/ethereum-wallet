package wallet

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

// Derive attempts to explicitly derive a hierarchical deterministic account at
// the specified derivation path. If requested, the derived account will be added
// to the wallet's tracked account list.
func (wallet *Wallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	wallet.mux.RLock()
	address, err := wallet.deriveAddress(path)
	wallet.mux.RUnlock()
	if err != nil {
		return accounts.Account{}, err
	}

	account := accounts.Account{
		Address: address,
		URL: accounts.URL{
			Path: path.String(),
		},
	}

	if pin {
		wallet.mux.Lock()
		defer wallet.mux.Unlock()

		if _, found := wallet.paths[address]; !found {
			wallet.accounts = append(wallet.accounts, account)
			wallet.paths[address] = path
			wallet.path = path
		}
	}
	
	return account, nil
}

// SelfDerive sets a base account derivation path from which the wallet attempts
// to discover non zero accounts and automatically add them to list of tracked
// accounts.
//
// Note, self derivation will increment the last component of the specified path
// opposed to decending into a child path to allow discovering accounts starting
// from non zero components.
//
// Some hardware wallets switched derivation paths through their evolution, so
// this method supports providing multiple bases to discover old user accounts
// too. Only the last base will be used to derive the next empty account.
//
// You can disable automatic account discovery by calling SelfDerive with a nil
// chain state reader.
func (wallet *Wallet) SelfDerive(base []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	wallet.mux.Lock()
	wallet.chainReader = chain
	wallet.mux.Unlock()

	if chain == nil {
		return
	}

	if len(base) == 0 {
		dPath, _ := accounts.ParseDerivationPath(DEFAULT_BASE_DERIVATION_PATH)
		base = append(base, dPath)
	}

	var currentPath accounts.DerivationPath
	var foundAccounts []accounts.Account
	var foundPaths map[common.Address]accounts.DerivationPath = make(map[common.Address]accounts.DerivationPath)
	// derive the first accounts
	for i, p := range base {

		iterate := accounts.DefaultIterator(p)
		for {
			currentPath = iterate()

			add, err := wallet.deriveAddress(currentPath)
			if err != nil {
				return
			}

			nonce, err := chain.NonceAt(wallet.ctx, add, nil)
			if err != nil {
				return
			}

			balance, err := chain.BalanceAt(wallet.ctx, add, nil)
			if err != nil {
				return
			}

			if nonce == 0 && balance.Cmp(big.NewInt(0)) != 1 {
				if i == len(base)-1 && wallet.accountMax > currentPath[len(currentPath)-1] {

					foundAccounts = append(foundAccounts, accounts.Account{
						Address: add,
						URL: accounts.URL{
							Path: currentPath.String(),
						},
					})
					foundPaths[add] = currentPath
					// check if wallet path is set and check up to wallet path
					continue
				}

				currentPath[len(currentPath)-1]--
				break
			}

			foundAccounts = append(foundAccounts, accounts.Account{
				Address: add,
				URL: accounts.URL{
					Path: currentPath.String(),
				},
			})
			foundPaths[add] = currentPath
		}
	}

	wallet.mux.Lock()
	wallet.paths = foundPaths
	wallet.accounts = foundAccounts
	wallet.path = currentPath
	wallet.mux.Unlock()
	go wallet.selfDerive()
}

// selfDerive is a background job runs automatic wallet discovery.
func (wallet *Wallet) selfDerive() {
	// return if already running
	if wallet._selfDerive {
		return
	}
	wallet.mux.Lock()
	wallet._selfDerive = true
	wallet.mux.Unlock()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for range ticker.C {
		wallet.mux.Lock()
		if wallet.chainReader == nil {
			wallet._selfDerive = false
			wallet.mux.Unlock()
			return
		}

		nextPath, _ := accounts.ParseDerivationPath(DEFAULT_BASE_DERIVATION_PATH)
		wallet.mux.Unlock()
		// get next wallet address
		nextPath[len(nextPath)-1]++
		nextAddress, err := wallet.deriveAddress(nextPath)
		if err != nil {
			log.Fatalf("selfDerivation failed with error: %q\n", err)
		}

		// read nonce from last block
		nonce, err := wallet.chainReader.NonceAt(wallet.ctx, nextAddress, nil)
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		balance, err := wallet.chainReader.BalanceAt(wallet.ctx, nextAddress, nil)
		if err != nil {
			return
		}

		// if there were txns
		if nonce > 0 || balance.Cmp(big.NewInt(0)) == 1 {
			wallet.mux.Lock()
			wallet.path = nextPath
			wallet.paths[nextAddress] = nextPath
			wallet.accounts = append(wallet.accounts, accounts.Account{
				Address: nextAddress,
				URL: accounts.URL{
					Path: nextPath.String(),
				},
			})
			wallet.mux.Unlock()
		}
	}

	wallet.mux.Lock()
	wallet._selfDerive = false
	wallet.mux.Unlock()
}
