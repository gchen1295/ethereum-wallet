import * as React from 'react';
import { FieldValues, useForm } from "react-hook-form";
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
import { validatePassword } from '@libs/utils'
import { getAuth, createUserWithEmailAndPassword, sendEmailVerification } from 'firebase/auth'
import { useAppDispatch, useAppSelector } from '@store/hooks';
import { setInitializing } from '@store/store';

interface RegisterFormProps {
  handleCancel: () => void
}

export const RegisterForm: React.FC<RegisterFormProps> = ({ handleCancel }) => {
  const {
    handleSubmit,
    register,
    setError,
    formState: { errors }
  } = useForm<{
    email: string
    password: string
    passwordConfirm: string
  }>();

  const hoverBgColor = useColorModeValue("rgba(0,0,0,.1)", "rgba(255,255,255,.1)")
  const toast = useToast()
  const dispatch = useAppDispatch()
  const auth = getAuth()

  const isInitializing = useAppSelector(state => state.auth.initializing)
  function onSubmit(values: FieldValues) {
    const { email, password, passwordConfirm } = values
    if (password !== passwordConfirm) {
      setError("passwordConfirm", { message: "Passwords must match" })
      return
    }

    dispatch(setInitializing(true))
    createUserWithEmailAndPassword(auth, email, password)
      .then((credentials) => {
        credentials.user?.emailVerified && sendEmailVerification(credentials.user)
      })
      .catch(() => {
        if (!toast.isActive("register-error")) {
          toast({
            id: "register-error",
            title: "Registration failed!",
            status: "error",
            description: "Failed to register user.",
            position: "top",
            isClosable: true,
          })
        }
      })
      .finally(() => dispatch(setInitializing(false)))
  }

  return (
    <chakra.form onSubmit={handleSubmit(onSubmit)}>
      <FormControl isInvalid={errors.email !== undefined} isRequired>
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
      <FormControl mt={2} isInvalid={errors.password !== undefined} isRequired>
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
            validate: value => validatePassword(value) || "Password strength too weak",
          })}
        />
        <FormErrorMessage>
          {errors.password && errors.password.message}
        </FormErrorMessage>
      </FormControl>
      <FormControl mt={2} isInvalid={errors.passwordConfirm !== undefined} isRequired>
        <FormLabel htmlFor="password">Confirm Password</FormLabel>
        <Input
          type="password"
          borderColor="gray.500"
          id="passwordConfirm"
          placeholder="Confirm Password"
          _hover={{ backgroundColor: hoverBgColor }}
          _focus={{ outline: "none" }}
          {...register("passwordConfirm", {
            required: "This is required",
          })}
        />
        <FormErrorMessage>
          {errors.passwordConfirm && errors.passwordConfirm.message}
        </FormErrorMessage>
      </FormControl>
      <ButtonGroup
        spacing={6}
        mt={6}
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
          onClick={handleCancel}
          borderColor="gray.500"
          _hover={{ backgroundColor: hoverBgColor }}
        >
          Cancel
        </Button>
      </ButtonGroup >
    </chakra.form>
  )
}