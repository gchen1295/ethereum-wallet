/* eslint-disable */
import { Observable } from "rxjs";

export const protobufPackage = "proto";

export enum StatusLevel {
  Log = 0,
  Warning = 1,
  Error = 2,
  UNRECOGNIZED = -1,
}

export enum ReponseStatus {
  OK = 0,
  UNRECOGNIZED = -1,
}

export interface Empty {}

/** Engine */
export interface Notification {
  message: string;
  level: StatusLevel;
}

export interface EngineStatus {
  connected: boolean;
  version: string;
  hwid: string;
}

/** Wallet Handler */
export interface Wallet {
  address: string;
}

export interface KeystoreOptions {
  passphrase?: string | undefined;
  mnemonic?: string | undefined;
  address?: string | undefined;
}

export interface KeystoreResponse {
  accounts: string[];
}

export interface MnemonicResponse {
  mnemonic: string;
}

/** Generics */
export interface GenericResponse {
  status: ReponseStatus;
  message: string;
}

export interface EngineHandler {
  Init(request: Empty): Promise<Empty>;
  Notify(request: Empty): Observable<Notification>;
  Listen(request: Empty): Observable<EngineStatus>;
}

export interface VaultHandler {
  Init(request: KeystoreOptions): Promise<KeystoreResponse>;
  CreateWallet(request: KeystoreOptions): Promise<KeystoreResponse>;
  ImportWallet(request: KeystoreOptions): Promise<KeystoreResponse>;
  CreateHDWallet(request: KeystoreOptions): Promise<KeystoreResponse>;
  DeleteWallet(request: KeystoreOptions): Promise<KeystoreResponse>;
  GenerateMnemonic(request: KeystoreOptions): Promise<MnemonicResponse>;
  GetWallets(request: KeystoreOptions): Promise<KeystoreResponse>;
  ListenWallets(request: Empty): Observable<KeystoreResponse>;
}
