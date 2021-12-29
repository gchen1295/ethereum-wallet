import * as React from 'react';
import { useHistory } from 'react-router-dom'
import {
  Flex,
  Divider,
  Heading,
  Text,
  Avatar,
  IconButton,
  useColorMode,
  useColorModeValue,
} from "@chakra-ui/react"

import {
  FaHome,
  FaWallet,
  FaEthereum,
  FaFileContract,
  FaChevronRight,
  FaChevronLeft,
} from 'react-icons/fa'

import { OpenSeaIcon } from '@assets/Opensea'
import { SidebarItem } from './SidebarItem'
import { useAppDispatch, useAppSelector } from '@store/store';
import { toggleSidebar } from '@store/store';


export const Sidebar: React.FC = () => {
  const history = useHistory()
  const dispatch = useAppDispatch()
  const sidebarOpen = useAppSelector(state => state.common.sidebarOpen)
  const { colorMode } = useColorMode()
  const border = useColorModeValue(".5px solid rgba(0, 0, 0, 0.1)", ".5px solid rgba(255, 255, 255, 0.1)")

  const [location, setLocation] = React.useState(history.location.pathname)

  React.useEffect(() => {
    const unlisten = history.listen((l) => {
      setLocation(l.pathname)
    });

    return unlisten
  }, [])

  return (
    <Flex
      minHeight="80%"
      width={sidebarOpen ? "200px" : "75px"}
      boxShadow="0 1px 5px rgba(0, 0, 0, 0.3)"
      borderRight={border}
      borderTop={border}
      pos="fixed"
      borderRightRadius={sidebarOpen ? "30px" : "15px"}
      flexDir="column"
      justifyContent="space-between"
      background={colorMode == "dark" ? "" : "#FAF9E0"}
      p="0"
      top="100px"
      className="noselect"
    >
      <Flex
        p="5%"
        flexDir="column"
        alignItems={sidebarOpen ? "flex-start" : "center"}
        as="nav"
      >
        <Flex
          p={3}
          justifyContent="space-around"
          alignItems="center"
          w="100%"
        >
          <Heading as="h3" size="sm" display={sidebarOpen ? "flex" : "none"}>Menu</Heading>
          <IconButton
            aria-label="Open/close sidebar"
            icon={sidebarOpen ? <FaChevronLeft /> : <FaChevronRight />}
            background="none"
            _hover={{ background: "none" }}
            onClick={() => {
              dispatch(toggleSidebar())
            }}
            _focus={{ outline: "none" }}
          />
        </Flex>

        <Divider />
        <SidebarItem
          icon={FaHome}
          title="Home"
          active={location === "/"}
          handleItemClick={() => { history.push("/") }}
        />
        <SidebarItem
          icon={FaWallet}
          title="Wallet"
          active={location === "/wallet"}
          handleItemClick={() => { history.push("/wallet") }}
        />
        <SidebarItem
          icon={FaEthereum}
          title="Ethereum"
          active={location === "/ether"}
          handleItemClick={() => { history.push("/ether") }}
        />
        <SidebarItem
          icon={FaFileContract}
          title="Contracts"
          active={location === "/contracts"}
          handleItemClick={() => { history.push("/contracts") }}
        />
        <SidebarItem
          icon={OpenSeaIcon}
          title="OpenSea"
          active={location === "/os"}
          handleItemClick={() => { history.push("/os") }}
        />
      </Flex>
      <Flex
        p="5%"
        flexDir="column"
        w="100%"
        alignItems={sidebarOpen ? "flex-start" : "center"}
        mb={4}
      >
        <Divider />
        <Flex mt={4} mb={sidebarOpen ? -4 : 0} align="center">
          <Avatar size="sm" />
          <Flex flexDir="column" ml={4} display={sidebarOpen ? "flex" : "none"}>
            <Heading as="h3" size="sm">Woof#6969</Heading>
            <Text color="gray">Admin</Text>
          </Flex>
        </Flex>
      </Flex>
    </Flex>

  )
}

export default Sidebar