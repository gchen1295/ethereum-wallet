import { VaultHandlerService } from './engine'
import { KeystoreOptions, Empty, KeystoreResponse } from '../proto/paws'
import { ipcMain, IpcMainEvent } from 'electron'
import { VaultChannel } from '../../common/interfaces/vault'
import { mainWindow } from '..'
export const initVault = () => {
  const opts: KeystoreOptions = {}

  VaultHandlerService.init(opts, (err, res) => {
    if (err) {
      console.log(err)
      return
    }

    if (res.accounts)
      mainWindow.webContents.send(VaultChannel.Wallets, res.accounts)
  })
}

export const listenAccounts = () => {
  const stream = VaultHandlerService.listenWallets(Empty)

  stream.on('data', (res: KeystoreResponse) => {
    if (res.accounts)
      mainWindow.webContents.send(VaultChannel.Wallets, res.accounts)
  })

  stream.on('error', (e) => { console.log(e) })
}
export const createWallet = (passphrase: string): Promise<string[]> => {
  const opts: KeystoreOptions = { passphrase }

  return new Promise((resolve, reject) => {
    VaultHandlerService.createWallet(opts, (err, res) => {
      if (err) reject(err)
      if (res?.accounts) resolve(res.accounts)
    })
  })
}


export const importWallet = (address: string, passphrase: string): Promise<string[]> => {
  const opts: KeystoreOptions = { passphrase, address }

  return new Promise((resolve, reject) => {
    VaultHandlerService.importWallet(opts, (err, res) => {
      if (err) reject(err)
      if (res?.accounts) resolve(res.accounts)
    })
  })
}

export const removeWallet = (address: string, passphrase: string): Promise<string[]> => {
  const opts: KeystoreOptions = { passphrase, address }

  return new Promise((resolve, reject) => {
    VaultHandlerService.deleteWallet(opts, (err, res) => {
      if (err) reject(err)
      if (res?.accounts) resolve(res.accounts)
    })
  })
}

export const generateMnemonic = (): Promise<string> => {
  return new Promise((resolve, reject) => {
    VaultHandlerService.generateMnemonic({}, (err, res) => {
      if (err) reject(err)
      if (res?.mnemonic) resolve(res.mnemonic)
    })
  })
}

export const createHDWallet = (mnemonic: string, passphrase: string): Promise<string[]> => {
  const opts: KeystoreOptions = { mnemonic, passphrase }
  console.log(opts)
  return new Promise((resolve, reject) => {
    VaultHandlerService.createHDWallet(opts, (err, res) => {
      if (err) reject(err)
      if (res?.accounts) resolve(res.accounts)
    })
  })
}

export const getWallets = (): Promise<string[]> => {
  const opts: KeystoreOptions = {}
  return new Promise((resolve, reject) => {
    VaultHandlerService.getWallets(opts, (err, res) => {
      if (err) reject(err)
      if (res?.accounts) resolve(res.accounts)
    })
  })
}


export const initVaultChannels = () => {
  listenAccounts()
  newHandler<string[]>(VaultChannel.GetWallets, getWallets)
  newHandler<string>(VaultChannel.GenerateMnemonic, generateMnemonic)
  newHandler<string[]>(VaultChannel.CreateHDWallet, createHDWallet)
  newHandler<string[]>(VaultChannel.CreateWallet, createWallet)
  newHandler<string[]>(VaultChannel.ImportWallet, importWallet)
  newHandler<string[]>(VaultChannel.RemoveWallet, removeWallet)
}
export const newListener = (e: string, fn: (...args: any[]) => void) => {
  ipcMain.on(e, (e, ...args) => { fn(...args) })
}

export function newHandler<T>(e: string, fn: (...args: any[]) => Promise<T>): void {
  ipcMain.handle(e, async (e, ...args) => {
    if (args) return fn(...args)
    return fn()
  })
}
