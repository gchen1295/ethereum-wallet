/* eslint-disable */
import Long from "long";
import {
  makeGenericClientConstructor,
  ChannelCredentials,
  ChannelOptions,
  UntypedServiceImplementation,
  handleUnaryCall,
  handleServerStreamingCall,
  Client,
  ClientUnaryCall,
  Metadata,
  CallOptions,
  ClientReadableStream,
  ServiceError,
} from "@grpc/grpc-js";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "proto";

export enum StatusLevel {
  Log = 0,
  Warning = 1,
  Error = 2,
  UNRECOGNIZED = -1,
}

export function statusLevelFromJSON(object: any): StatusLevel {
  switch (object) {
    case 0:
    case "Log":
      return StatusLevel.Log;
    case 1:
    case "Warning":
      return StatusLevel.Warning;
    case 2:
    case "Error":
      return StatusLevel.Error;
    case -1:
    case "UNRECOGNIZED":
    default:
      return StatusLevel.UNRECOGNIZED;
  }
}

export function statusLevelToJSON(object: StatusLevel): string {
  switch (object) {
    case StatusLevel.Log:
      return "Log";
    case StatusLevel.Warning:
      return "Warning";
    case StatusLevel.Error:
      return "Error";
    default:
      return "UNKNOWN";
  }
}

export enum ReponseStatus {
  OK = 0,
  UNRECOGNIZED = -1,
}

export function reponseStatusFromJSON(object: any): ReponseStatus {
  switch (object) {
    case 0:
    case "OK":
      return ReponseStatus.OK;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ReponseStatus.UNRECOGNIZED;
  }
}

export function reponseStatusToJSON(object: ReponseStatus): string {
  switch (object) {
    case ReponseStatus.OK:
      return "OK";
    default:
      return "UNKNOWN";
  }
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

const baseEmpty: object = {};

export const Empty = {
  encode(_: Empty, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Empty {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseEmpty } as Empty;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): Empty {
    const message = { ...baseEmpty } as Empty;
    return message;
  },

  toJSON(_: Empty): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Empty>, I>>(_: I): Empty {
    const message = { ...baseEmpty } as Empty;
    return message;
  },
};

const baseNotification: object = { message: "", level: 0 };

export const Notification = {
  encode(
    message: Notification,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.message !== "") {
      writer.uint32(10).string(message.message);
    }
    if (message.level !== 0) {
      writer.uint32(16).int32(message.level);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Notification {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseNotification } as Notification;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.message = reader.string();
          break;
        case 2:
          message.level = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Notification {
    const message = { ...baseNotification } as Notification;
    message.message =
      object.message !== undefined && object.message !== null
        ? String(object.message)
        : "";
    message.level =
      object.level !== undefined && object.level !== null
        ? statusLevelFromJSON(object.level)
        : 0;
    return message;
  },

  toJSON(message: Notification): unknown {
    const obj: any = {};
    message.message !== undefined && (obj.message = message.message);
    message.level !== undefined &&
      (obj.level = statusLevelToJSON(message.level));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Notification>, I>>(
    object: I
  ): Notification {
    const message = { ...baseNotification } as Notification;
    message.message = object.message ?? "";
    message.level = object.level ?? 0;
    return message;
  },
};

const baseEngineStatus: object = { connected: false, version: "", hwid: "" };

export const EngineStatus = {
  encode(
    message: EngineStatus,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.connected === true) {
      writer.uint32(8).bool(message.connected);
    }
    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }
    if (message.hwid !== "") {
      writer.uint32(26).string(message.hwid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EngineStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseEngineStatus } as EngineStatus;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.connected = reader.bool();
          break;
        case 2:
          message.version = reader.string();
          break;
        case 3:
          message.hwid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EngineStatus {
    const message = { ...baseEngineStatus } as EngineStatus;
    message.connected =
      object.connected !== undefined && object.connected !== null
        ? Boolean(object.connected)
        : false;
    message.version =
      object.version !== undefined && object.version !== null
        ? String(object.version)
        : "";
    message.hwid =
      object.hwid !== undefined && object.hwid !== null
        ? String(object.hwid)
        : "";
    return message;
  },

  toJSON(message: EngineStatus): unknown {
    const obj: any = {};
    message.connected !== undefined && (obj.connected = message.connected);
    message.version !== undefined && (obj.version = message.version);
    message.hwid !== undefined && (obj.hwid = message.hwid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EngineStatus>, I>>(
    object: I
  ): EngineStatus {
    const message = { ...baseEngineStatus } as EngineStatus;
    message.connected = object.connected ?? false;
    message.version = object.version ?? "";
    message.hwid = object.hwid ?? "";
    return message;
  },
};

const baseWallet: object = { address: "" };

export const Wallet = {
  encode(
    message: Wallet,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Wallet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseWallet } as Wallet;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Wallet {
    const message = { ...baseWallet } as Wallet;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : "";
    return message;
  },

  toJSON(message: Wallet): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Wallet>, I>>(object: I): Wallet {
    const message = { ...baseWallet } as Wallet;
    message.address = object.address ?? "";
    return message;
  },
};

const baseKeystoreOptions: object = {};

export const KeystoreOptions = {
  encode(
    message: KeystoreOptions,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.passphrase !== undefined) {
      writer.uint32(10).string(message.passphrase);
    }
    if (message.mnemonic !== undefined) {
      writer.uint32(18).string(message.mnemonic);
    }
    if (message.address !== undefined) {
      writer.uint32(26).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): KeystoreOptions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseKeystoreOptions } as KeystoreOptions;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.passphrase = reader.string();
          break;
        case 2:
          message.mnemonic = reader.string();
          break;
        case 3:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): KeystoreOptions {
    const message = { ...baseKeystoreOptions } as KeystoreOptions;
    message.passphrase =
      object.passphrase !== undefined && object.passphrase !== null
        ? String(object.passphrase)
        : undefined;
    message.mnemonic =
      object.mnemonic !== undefined && object.mnemonic !== null
        ? String(object.mnemonic)
        : undefined;
    message.address =
      object.address !== undefined && object.address !== null
        ? String(object.address)
        : undefined;
    return message;
  },

  toJSON(message: KeystoreOptions): unknown {
    const obj: any = {};
    message.passphrase !== undefined && (obj.passphrase = message.passphrase);
    message.mnemonic !== undefined && (obj.mnemonic = message.mnemonic);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<KeystoreOptions>, I>>(
    object: I
  ): KeystoreOptions {
    const message = { ...baseKeystoreOptions } as KeystoreOptions;
    message.passphrase = object.passphrase ?? undefined;
    message.mnemonic = object.mnemonic ?? undefined;
    message.address = object.address ?? undefined;
    return message;
  },
};

const baseKeystoreResponse: object = { accounts: "" };

export const KeystoreResponse = {
  encode(
    message: KeystoreResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.accounts) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): KeystoreResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseKeystoreResponse } as KeystoreResponse;
    message.accounts = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accounts.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): KeystoreResponse {
    const message = { ...baseKeystoreResponse } as KeystoreResponse;
    message.accounts = (object.accounts ?? []).map((e: any) => String(e));
    return message;
  },

  toJSON(message: KeystoreResponse): unknown {
    const obj: any = {};
    if (message.accounts) {
      obj.accounts = message.accounts.map((e) => e);
    } else {
      obj.accounts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<KeystoreResponse>, I>>(
    object: I
  ): KeystoreResponse {
    const message = { ...baseKeystoreResponse } as KeystoreResponse;
    message.accounts = object.accounts?.map((e) => e) || [];
    return message;
  },
};

const baseMnemonicResponse: object = { mnemonic: "" };

export const MnemonicResponse = {
  encode(
    message: MnemonicResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.mnemonic !== "") {
      writer.uint32(10).string(message.mnemonic);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MnemonicResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMnemonicResponse } as MnemonicResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mnemonic = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MnemonicResponse {
    const message = { ...baseMnemonicResponse } as MnemonicResponse;
    message.mnemonic =
      object.mnemonic !== undefined && object.mnemonic !== null
        ? String(object.mnemonic)
        : "";
    return message;
  },

  toJSON(message: MnemonicResponse): unknown {
    const obj: any = {};
    message.mnemonic !== undefined && (obj.mnemonic = message.mnemonic);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MnemonicResponse>, I>>(
    object: I
  ): MnemonicResponse {
    const message = { ...baseMnemonicResponse } as MnemonicResponse;
    message.mnemonic = object.mnemonic ?? "";
    return message;
  },
};

const baseGenericResponse: object = { status: 0, message: "" };

export const GenericResponse = {
  encode(
    message: GenericResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenericResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenericResponse } as GenericResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenericResponse {
    const message = { ...baseGenericResponse } as GenericResponse;
    message.status =
      object.status !== undefined && object.status !== null
        ? reponseStatusFromJSON(object.status)
        : 0;
    message.message =
      object.message !== undefined && object.message !== null
        ? String(object.message)
        : "";
    return message;
  },

  toJSON(message: GenericResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = reponseStatusToJSON(message.status));
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenericResponse>, I>>(
    object: I
  ): GenericResponse {
    const message = { ...baseGenericResponse } as GenericResponse;
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    return message;
  },
};

export const EngineHandlerService = {
  init: {
    path: "/proto.EngineHandler/Init",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: Empty) =>
      Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: Empty) =>
      Buffer.from(Empty.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Empty.decode(value),
  },
  notify: {
    path: "/proto.EngineHandler/Notify",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: Empty) =>
      Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: Notification) =>
      Buffer.from(Notification.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Notification.decode(value),
  },
  listen: {
    path: "/proto.EngineHandler/Listen",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: Empty) =>
      Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: EngineStatus) =>
      Buffer.from(EngineStatus.encode(value).finish()),
    responseDeserialize: (value: Buffer) => EngineStatus.decode(value),
  },
} as const;

export interface EngineHandlerServer extends UntypedServiceImplementation {
  init: handleUnaryCall<Empty, Empty>;
  notify: handleServerStreamingCall<Empty, Notification>;
  listen: handleServerStreamingCall<Empty, EngineStatus>;
}

export interface EngineHandlerClient extends Client {
  init(
    request: Empty,
    callback: (error: ServiceError | null, response: Empty) => void
  ): ClientUnaryCall;
  init(
    request: Empty,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Empty) => void
  ): ClientUnaryCall;
  init(
    request: Empty,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Empty) => void
  ): ClientUnaryCall;
  notify(
    request: Empty,
    options?: Partial<CallOptions>
  ): ClientReadableStream<Notification>;
  notify(
    request: Empty,
    metadata?: Metadata,
    options?: Partial<CallOptions>
  ): ClientReadableStream<Notification>;
  listen(
    request: Empty,
    options?: Partial<CallOptions>
  ): ClientReadableStream<EngineStatus>;
  listen(
    request: Empty,
    metadata?: Metadata,
    options?: Partial<CallOptions>
  ): ClientReadableStream<EngineStatus>;
}

export const EngineHandlerClient = makeGenericClientConstructor(
  EngineHandlerService,
  "proto.EngineHandler"
) as unknown as {
  new (
    address: string,
    credentials: ChannelCredentials,
    options?: Partial<ChannelOptions>
  ): EngineHandlerClient;
};

export const VaultHandlerService = {
  init: {
    path: "/proto.VaultHandler/Init",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  createWallet: {
    path: "/proto.VaultHandler/CreateWallet",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  importWallet: {
    path: "/proto.VaultHandler/ImportWallet",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  createHDWallet: {
    path: "/proto.VaultHandler/CreateHDWallet",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  deleteWallet: {
    path: "/proto.VaultHandler/DeleteWallet",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  generateMnemonic: {
    path: "/proto.VaultHandler/GenerateMnemonic",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: MnemonicResponse) =>
      Buffer.from(MnemonicResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => MnemonicResponse.decode(value),
  },
  getWallets: {
    path: "/proto.VaultHandler/GetWallets",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: KeystoreOptions) =>
      Buffer.from(KeystoreOptions.encode(value).finish()),
    requestDeserialize: (value: Buffer) => KeystoreOptions.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
  listenWallets: {
    path: "/proto.VaultHandler/ListenWallets",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: Empty) =>
      Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: KeystoreResponse) =>
      Buffer.from(KeystoreResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => KeystoreResponse.decode(value),
  },
} as const;

export interface VaultHandlerServer extends UntypedServiceImplementation {
  init: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  createWallet: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  importWallet: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  createHDWallet: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  deleteWallet: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  generateMnemonic: handleUnaryCall<KeystoreOptions, MnemonicResponse>;
  getWallets: handleUnaryCall<KeystoreOptions, KeystoreResponse>;
  listenWallets: handleServerStreamingCall<Empty, KeystoreResponse>;
}

export interface VaultHandlerClient extends Client {
  init(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  init(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  init(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createWallet(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  importWallet(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  importWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  importWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createHDWallet(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createHDWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  createHDWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  deleteWallet(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  deleteWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  deleteWallet(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  generateMnemonic(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: MnemonicResponse) => void
  ): ClientUnaryCall;
  generateMnemonic(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: MnemonicResponse) => void
  ): ClientUnaryCall;
  generateMnemonic(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: MnemonicResponse) => void
  ): ClientUnaryCall;
  getWallets(
    request: KeystoreOptions,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  getWallets(
    request: KeystoreOptions,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  getWallets(
    request: KeystoreOptions,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: KeystoreResponse) => void
  ): ClientUnaryCall;
  listenWallets(
    request: Empty,
    options?: Partial<CallOptions>
  ): ClientReadableStream<KeystoreResponse>;
  listenWallets(
    request: Empty,
    metadata?: Metadata,
    options?: Partial<CallOptions>
  ): ClientReadableStream<KeystoreResponse>;
}

export const VaultHandlerClient = makeGenericClientConstructor(
  VaultHandlerService,
  "proto.VaultHandler"
) as unknown as {
  new (
    address: string,
    credentials: ChannelCredentials,
    options?: Partial<ChannelOptions>
  ): VaultHandlerClient;
};

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;

export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
  ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & Record<
        Exclude<keyof I, KeysOfUnion<P>>,
        never
      >;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}
