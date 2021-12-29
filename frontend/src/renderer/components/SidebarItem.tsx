import * as React from 'react';
import {
  Flex,
  Text,
  Menu,
  MenuButton,
  Link,
  Icon,
} from "@chakra-ui/react"
import { useAppSelector } from '@store/store'

interface SidebarProps {
  icon: any;
  title: string;
  active?: boolean;
  handleItemClick?: () => void;
}

export const SidebarItem: React.FC<SidebarProps> = ({ icon, title, active, handleItemClick }) => {
  const sidebarOpen = useAppSelector(state => state.common.sidebarOpen)
  return (
    <Flex
      mt={30}
      flexDir="column"
      w="100%"
      alignItems={sidebarOpen ? "flex-start" : "center"}
    >
      <Menu placement="right">
        <Link
          backgroundColor={active ? "rgba(0, 0, 0, 0.2)" : ""}
          p={3}
          borderRadius={8}
          _hover={{ textDecor: 'none', backgroundColor: 'rgba(0, 0, 0, 0.2)' }}
          w={sidebarOpen ? "100%" : ""}
          onClick={handleItemClick}
        >
          <MenuButton w="100%">
            <Flex>
              <Icon as={icon} fontSize="xl" />
              <Text ml={5} display={sidebarOpen ? "flex" : "none"}>{title}</Text>
            </Flex>
          </MenuButton>
        </Link>
      </Menu>
    </Flex>
  )
}

export default SidebarItem