import * as React from 'react';
import { useForm } from "react-hook-form";
import {
  chakra,
  Button,
  Input,
  FormErrorMessage,
  FormLabel,
  FormControl,
  useColorModeValue,
  ButtonGroup,
  useToast
} from "@chakra-ui/react"

import { getAuth, signInWithEmailAndPassword } from 'firebase/auth'
import { LoginRequest } from '@interfaces/auth'
import { useAppDispatch, useAppSelector } from '@store/hooks';
import { setAuth, setInitializing } from '@store/store';
import { useHistory } from 'react-router-dom';


interface LoginFormProps {
  handleSignup: () => void
}

export const LoginForm: React.FC<LoginFormProps> = ({ handleSignup }) => {
  const {
    handleSubmit,
    register,
    formState: { errors }
  } = useForm<LoginRequest>();

  const toast = useToast()
  const dispatch = useAppDispatch()
  const history = useHistory()
  const auth = getAuth()

  const isInitializing = useAppSelector(state => state.auth.initializing)
  const hoverBgColor = useColorModeValue("rgba(0,0,0,.1)", "rgba(255,255,255,.1)")

  function onSubmit(values: LoginRequest) {
    dispatch(setInitializing(true))
    signInWithEmailAndPassword(auth, values.email, values.password)
      .then(credentials => {
        dispatch(setAuth({
          name: credentials.user.displayName ?? "",
          email: credentials.user.email ?? "",
          uid: credentials.user.uid ?? "",
          claim: credentials.user.refreshToken,
          photo: credentials.user.photoURL ?? "",
        }))

        history.push("/home")
      })
      .catch((e) => {
        toast({
          id: 'login-error',
          title: "Login Failed!",
          variant: "solid",
          status: "error",
          description: 'User authentication failed.',
          isClosable: true,
          position: 'top'
        })
      })
      .finally(() => dispatch(setInitializing(false)))
  }

  return (
    <chakra.form
      onSubmit={handleSubmit(onSubmit)}
    >
      <FormControl isInvalid={errors.email !== undefined}>
        <FormLabel htmlFor="email">Email</FormLabel>
        <Input
          id="email"
          borderColor="gray.500"
          placeholder="Email"
          _hover={{ backgroundColor: hoverBgColor }}
          _focus={{ outline: "none" }}
          inputMode="email"
          {...register("email", {
            required: "Email is required",
            pattern: {
              value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
              message: "invalid email address"
            },
          })}
        />
        <FormErrorMessage>
          {errors.email && errors.email.message}
        </FormErrorMessage>
      </FormControl>
      <FormControl mt={2} isInvalid={errors.password !== undefined}>
        <FormLabel htmlFor="password">Password</FormLabel>
        <Input
          type="password"
          borderColor="gray.500"
          id="password"
          placeholder="Password"
          _hover={{ backgroundColor: hoverBgColor }}
          _focus={{ outline: "none" }}
          {...register("password", {
            required: "This is required",
            minLength: { value: 4, message: "Minimum length should be 4" }
          })}
        />
        <FormErrorMessage>
          {errors.password && errors.password.message}
        </FormErrorMessage>
      </FormControl>
      <ButtonGroup
        spacing={6}
        mt={4}
        variant="outline"
      >
        <Button
          isLoading={isInitializing}
          type="submit"
          width="150px"
          borderColor="gray.500"
          _hover={{ backgroundColor: hoverBgColor }}
        >
          Submit
        </Button>
        <Button
          type="button"
          width="150px"
          onClick={handleSignup}
          borderColor="gray.500"
          _hover={{ backgroundColor: hoverBgColor }}
        >
          Register
        </Button>
      </ButtonGroup >
    </chakra.form >
  )
}

export default LoginForm