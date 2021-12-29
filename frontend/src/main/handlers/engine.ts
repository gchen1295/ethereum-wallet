import { connectivityState, credentials, ServiceError } from '@grpc/grpc-js'
import { execFile } from 'child_process'
import { app } from 'electron'
import { createWriteStream, existsSync, mkdirSync } from 'fs'
import path from 'path'

import {
  VaultHandlerClient,
  EngineHandlerClient,
  Empty,
  EngineStatus,
} from '../proto/paws'
import { initVault, initVaultChannels} from './vault'
import { statuses } from '../index'


let state: connectivityState
const grpcUrl = 'localhost:17529'
const grpcCredentials = credentials.createInsecure()
const grpcConfig = {
  'grpc.max_send_message_length': 1024 * 1024 * 1024,
  'grpc.max_receive_message_length': 1024 * 1024 * 1024,
}

export const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export const VaultHandlerService = new VaultHandlerClient(grpcUrl, grpcCredentials, grpcConfig)
export const EngineHandlerService = new EngineHandlerClient(grpcUrl, grpcCredentials, grpcConfig)

export const spawnEngine = async (enginePath: string, token: string, licenseId: string) => {
  const engine = execFile(enginePath, [`-token=${token}`, `-licenseid=${licenseId}`])

  if (!app.isPackaged) {
    engine.stdout?.on('data', (data) => {
      console.log(data.toString())
    })
  }

  if (!existsSync(path.join(app.getPath('userData'), 'logs'))) {
    mkdirSync(path.join(app.getPath('userData'), 'logs'))
  }

  const launchTime = Date.now()

  const logStream = createWriteStream(path.join(app.getPath('userData'), `logs/${launchTime}.log`))

  engine.stderr?.pipe(logStream)

  engine.stderr?.on('data', (data) => {
    const message = data.toString() as string

    if (!app.isPackaged) {
      console.log(message)
    }

    // if (message.includes('panic: ')) {
    //   const panicMessage = message.match('panic: (.*)\\n') as RegExpMatchArray
    //   mainWindow.ipcEmitter("notice-event", {
    //     title: 'Engine Crashed',
    //     message: `${panicMessage[1]} please restart.`,
    //   })
    // }
  })
}

const engineReady = async () =>
  // eslint-disable-next-line
  new Promise<void>(async (resolve, reject) => {
    while (state !== connectivityState.READY) {
      state = EngineHandlerService.getChannel().getConnectivityState(true)
      // eslint-disable-next-line
      await sleep(100)
    }

    resolve()
  })


export const listenEngineStatus = async () => {
  const stream = EngineHandlerService.listen(Empty)

  stream.on('data', (status: EngineStatus) => {
    statuses.engineStatus = status
  })

  stream.on('error', () => { })
}

export const initEngine = async (
  enginePath: string,
  token: string,
  licenseId: string,
) => {
  await spawnEngine(enginePath, token, licenseId)
  await engineReady()

  listenEngineStatus()
  initVaultChannels()

  initVault()
}