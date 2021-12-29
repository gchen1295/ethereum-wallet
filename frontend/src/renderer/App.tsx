import { Route, Switch, useHistory } from 'react-router-dom'
import * as React from 'react'
import {
  Grid,
  GridItem,
  chakra,
  Box
} from '@chakra-ui/react'
import '@/renderer/styles/globals.css'
import { Sidebar } from '@components/Sidebar'
import { Header } from '@components/Header'

import Landing from './pages/index'
import Login from './pages/login'
import Home from './pages/hub/home'
import Ether from './pages/hub/eth'
import Wallet from './pages/hub/wallet'
import Contract from './pages/hub/contracts'
import OpenSea from './pages/hub/os'
import { listenWallets } from '@/renderer/handlers/vault'


export const MyApp = () => {
  React.useEffect(()=> {
    listenWallets()
  }, [])
  return (
    <Box>
      <Header />
      <Sidebar />
      <chakra.div className="draggable" />
      <Switch>
        <Route path={"/login"} component={Login} />
        <Route path={"/ether"} component={Ether} />
        <Route path={"/wallet"} component={Wallet} />
        <Route path={"/contracts"} component={Contract} />
        <Route path={"/os"} component={OpenSea} />
        <Route path={"/"} component={Home} />
      </Switch>
    </Box>
  )
}


export default MyApp