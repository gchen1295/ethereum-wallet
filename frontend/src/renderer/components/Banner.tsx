import * as React from 'react';
import {
  Flex,
  Box,
  Spacer,
  Heading,
  IconButton,
  useColorMode,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Button
} from "@chakra-ui/react"

export const Banner: React.FC = () => {

  return (
    <Flex
      width="100%"
      height="100%"
      justifyContent="center"
    >
      <Box>
        <Heading>HUB Banner</Heading>
      </Box>
    </Flex>
  )
}

export default Banner