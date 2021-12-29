/**
 * This file will automatically be loaded by webpack and run in the "renderer" context.
 * To learn more about the differences between the "main" and the "renderer" context in
 * Electron, visit:
 *
 * https://electronjs.org/docs/tutorial/application-architecture#main-and-renderer-processes
 *
 * By default, Node.js integration in this file is disabled. When enabling Node.js integration
 * in a renderer process, please be aware of potential security implications. You can read
 * more about security risks here:
 *
 * https://electronjs.org/docs/tutorial/security
 *
 * To enable Node.js integration in this file, open up `main.js` and enable the `nodeIntegration`
 * flag:
 *
 * ```
 *  // Create the browser window.
 *  mainWindow = new BrowserWindow({
 *    width: 800,
 *    height: 600,
 *    webPreferences: {
 *      nodeIntegration: true
 *    }
 *  });
 * ```
 */

import { Provider } from 'react-redux'
import ReactDOM from 'react-dom'
import * as React from 'react'
import {
  ChakraProvider,
  CSSReset,
  ColorModeScript,
} from '@chakra-ui/react'
import { ChainId, DAppProvider, Config } from '@usedapp/core'
import { Router } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import { store } from '@store/store'
import { firebaseConfig } from '@libs/firebase'
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";

import theme from '@/renderer/theme'
import '@/renderer/styles/globals.css'

declare global {
  interface Window {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    [x: string]: any;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    require: any;
  }
}

const config: Config = {
  readOnlyChainId: ChainId.Mainnet,
  readOnlyUrls: {
    [ChainId.Mainnet]: "http://127.0.0.1:8545/",
  },
  pollingInterval: 10000,
  supportedChains: [ChainId.Mainnet, 31337]
}

import App from '@/renderer/App'
import './index.css';


// Initialize Firebase
const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);

ReactDOM.render(
  <React.StrictMode>
    <DAppProvider config={config}>
      <ChakraProvider theme={theme}>
        <Provider store={store}>
          <Router history={createBrowserHistory()}>
            <CSSReset />
            <App />
            <ColorModeScript initialColorMode={theme.config.initialColorMode} />
          </Router>
        </Provider>
      </ChakraProvider>
    </DAppProvider>
  </React.StrictMode>,
  document.getElementById('root') as HTMLDivElement
)


console.log('👋 This message is being logged by "renderer.js", included via webpack');
