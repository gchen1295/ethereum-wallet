import * as React from 'react';
import {
  Flex,
  Box,
  Text,
  Spacer,
  Heading,
  IconButton,
  useColorMode,
  useColorModeValue,
  CloseButton,
  Link,
  Badge,
} from "@chakra-ui/react"
import {
  MoonIcon,
  SunIcon,
} from '@chakra-ui/icons'
import {
  toggleLogin,
  useAppDispatch,
  useAppSelector,
} from '@store/store'
import { useHistory } from 'react-router-dom'

// TODO move Brand title to props
export const Header: React.FC = () => {
  const { toggleColorMode } = useColorMode();
  const colorModeIcon = useColorModeValue(<MoonIcon />, <SunIcon />)
  const btnHoverBgColor = useColorModeValue("rgba(0,0,0,.1)", "rgba(255,255,255,.1)")
  const border = useColorModeValue(".5px solid rgba(0, 0, 0, 0.1)", ".5px solid rgba(255, 255, 255, 0.1)")
  const boxShadow = useColorModeValue("0 1px 6px rgba(0, 0, 0, 0.1)", "")

  const dispatch = useAppDispatch();
  const history = useHistory()
  const version = useAppSelector(state => state.common.version)

  const user = {
    id: ''
  }

  return (
    <>
      <Flex
        height="70px"
        width="100%"
        alignItems="center"
        position="fixed"
        top="0"
        p="2"
        className="noselect"
      >
        <Box >
          <Link
            _hover={{ textDecor: "none" }}
            onClick={() => {
              dispatch(toggleLogin(false))
              history.push("/")
            }}
          >
            <Heading size="lg" ml={2}>🐾 Paws</Heading>
          </Link>
        </Box>
        <Spacer />
        <Box mr="5" mt="3">
          <Badge variant='subtle' colorScheme='green' p="2" borderRadius="10px">
            {version}
          </Badge>
          <IconButton
            aria-label="Toggle Light/Dark mode"
            icon={colorModeIcon}
            size="md"
            background="transparent"
            onClick={toggleColorMode}
            _focus={{ outline: "none" }}
            _hover={{}}
          />
        </Box>
        <Box mt="-5">
          <CloseButton
            aria-label={"Exit application"}
            size="md"
            background="transparent"
            onClick={() => window.ipcRenderer.send("quit")}
            _focus={{ outline: "none" }}
            _hover={{ backgroundColor: btnHoverBgColor }}
          />
        </Box>
      </Flex>
    </>
  )
}

export default Header
