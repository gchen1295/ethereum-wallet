import { VaultChannel } from '../../common/interfaces/vault'

import { store } from '@store/store'
import { setAccounts } from '@store/slices/vault'
import { IpcMainEvent, ipcRenderer } from 'electron'

export const listenWallets = () => {
  window.ipcRenderer.receive(VaultChannel.Wallets, (accounts: string[]) => {
    store.dispatch(setAccounts(accounts));
  });
}

export const createWallet = (passphrase: string) => {
  return window.ipcRenderer.invoke(VaultChannel.CreateWallet, passphrase);
}

export const importWallet = (address: string, passphrase: string) => {
  return window.ipcRenderer.invoke(VaultChannel.ImportWallet, address, passphrase);
}

export const removeWallet = (address: string, passphrase: string) => {
  return window.ipcRenderer.invoke(VaultChannel.RemoveWallet, address, passphrase);
}

export const generateMnemonic = async (): Promise<string> => {
  return window.ipcRenderer.invoke(VaultChannel.GenerateMnemonic)
}

export const getWallets = async (): Promise<string[]> => {
  return await window.ipcRenderer.invoke(VaultChannel.GetWallets)
}

export const createHDWallet = async (mnemonic: string, passphrase: string): Promise<void> => {
  return window.ipcRenderer.invoke(VaultChannel.CreateHDWallet, mnemonic, passphrase);
}