import * as React from 'react';
import { useState } from 'react'
import {
  Flex,
  Box,
  Heading,
  useColorModeValue,
} from "@chakra-ui/react"

import { LoginForm } from './LoginForm'
import { RegisterForm } from './RegisterForm'
import { VerifyForm } from './VerifyForm'

export enum AuthMode {
  Login = 'LOGIN',
  Register = 'REGISTER',
  Verify = 'VERIFY',
  None = 'NONE'
}

interface AuthCardProps {
  mode?: AuthMode
}

export const AuthCard: React.FC<AuthCardProps> = ({ mode }) => {
  const [authMode, setAuthMode] = useState(mode ?? AuthMode.Login)

  const HeaderText = () => {
    switch (mode) {
      case AuthMode.Login:
        return "Login"
      case AuthMode.Register:
        return "Register"
      case AuthMode.Verify:
        return "Verify Email"
      default:
        return "Login"
    }
  }

  const bgColor = useColorModeValue("var(--color-light-primary)", "var(--color-dark-secondary)")

  if (authMode === AuthMode.None) {
    return null
  }

  return (
    <Flex
      width="500px"
      height={authMode !== AuthMode.Register ? "420px" : "520px"}
      flexDir="column"
      justifyContent="center"
      alignItems="center"
      background={bgColor}
      borderRadius="15px"
      boxShadow="0px 0px 6px 0 rgba(0, 0, 0, .2)"
    >
      <Box>
        <Heading>{HeaderText()}</Heading>
      </Box>
      <br />
      <Box width="65%">
        {
          authMode === AuthMode.Register
            ? <RegisterForm handleCancel={() => setAuthMode(AuthMode.Login)} />
            : authMode === AuthMode.Verify
              ? <VerifyForm handleClose={() =>  setAuthMode(AuthMode.None) } />
              : <LoginForm handleSignup={() => setAuthMode(AuthMode.Register)} />
        }
      </Box>
    </Flex >
  )
}

export default AuthCard