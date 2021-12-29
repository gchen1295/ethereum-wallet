import * as React from 'react';
import styles from '@styles/hub/Home.module.css'
import {
  Grid,
  GridItem,
  useColorModeValue,
  Flex,
  Text,
  Table,
  TableCaption,
  Thead,
  Tr,
  Th,
  Tfoot,
  Td,
  Tbody,
  ButtonGroup,
  Button,
  Container,
  Box
} from '@chakra-ui/react'

import { useEtherBalance, useEthers, } from '@usedapp/core'
import { utils } from 'ethers'
import store, {
  toggleLogin,
  useAppDispatch,
  useAppSelector,
} from '@store/store'
import {createWallet} from '@/renderer/handlers/vault'
import {FormWalletWelcome} from '@components/FormWalletWelcome'
import { generateMnemonic, createHDWallet, getWallets } from '@/renderer/handlers/vault';
import { setAccounts } from '@/renderer/store/slices/vault';

const HubWallet = () => {
  const bgColor = useColorModeValue("var(--color-light)", "")
  const { account } = useEthers()
  const balance = useEtherBalance(account)
  const accounts = useAppSelector(state => state.vault.accounts)
  const vksp = useAppSelector(state => state.auth.vksp)

  function handleNewWallet() {
    if (!vksp) {
      // show password dialog
      createWallet("supersecretpassword123")
      return
    }

    createWallet(vksp)
  }

  const [formLoading, setFormLoading] = React.useState(false)


  return (
    <Grid
      className={styles.container}
      background={bgColor}
      h="100vh"
      gridTemplateColumns="100px 1fr 1fr 1fr"
      gridTemplateRows="75px 60px 1fr 1fr 40px"
      gridTemplateAreas={`
        'header header header header'
        'side thead thead thead'
        'side main main main'
        'side main main main'
        'footer footer  footer footer'
      `}
      p="0"
    >
      {accounts.length > 0 
      ? <>
          <GridItem gridArea="thead" minHeight="100%" w="100%" justifyContent="end" p={3}>
            <ButtonGroup variant='outline' minHeight="100%" spacing='6' justifyContent="end" w="90%" p={3}>
              <Button colorScheme='blue' onClick={handleNewWallet}>New</Button>
              <Button colorScheme='blue'>Import</Button>
            </ButtonGroup>
          </GridItem>
          <GridItem 
          gridArea="main" 
          height="100%" p={3} 
          justifyContent="start"  
          overflow="auto" 
          css={{
            '::-webkit-scrollbar': {width: '12px'},
            '::-webkit-scrollbar-track': {width: '6px',},
            '::-webkit-scrollbar-thumb': {borderRadius: '24px',}, 
            }}>
            <Table variant='simple' size={"md"}>
              <Thead position="sticky" top="0" bg="white">
                <Tr>
                  <Th>Account</Th>
                  <Th>Address</Th>
                  <Th>Balance</Th>
                </Tr>
              </Thead>
              <Tbody overflowY="auto" h="100%" w="100%" mt="50px">
                {accounts.map(acc => {
                  return(
                    <Tr>
                      <Td>Account</Td>
                      <Td>{acc}</Td>
                      <Td>0.3028E</Td>
                    </Tr>
                  )
                })}
              </Tbody> 
            </Table>
          </GridItem></>
      : <>
      <GridItem gridArea="main" height="100%" p={3} >
        <Flex height="100%" w="100%" justifyContent="center" alignItems="center" mt="-10">
          <FormWalletWelcome getSecret={generateMnemonic} handleWalletForm={(secret, passphrase)=>{
            setFormLoading(true)
            console.log("Passphrase: ", passphrase)
            createHDWallet(secret, passphrase)
            .then(() => getWallets())
            .then(accounts => {
              store.dispatch(setAccounts(accounts))
              setFormLoading(false)
            })
            
          }}
          formLoading={formLoading}
          />
        </Flex>
      </GridItem>
      </>}
      <GridItem gridArea="footer">
      </GridItem>
    </Grid >
  )
}

export default HubWallet
