import * as React from 'react';
import { FieldValues, useForm } from "react-hook-form";
import {
  chakra,
  Button,
  useColorModeValue,
  ButtonGroup,
  useToast
} from "@chakra-ui/react"
import { useAppDispatch, } from '@store/store'

interface VerifyFormProps {
  handleClose: () => void
}

export const VerifyForm: React.FC<VerifyFormProps> = ({ handleClose }) => {
  const {
    handleSubmit,
    register,
    setError,
    formState: { errors }
  } = useForm<{
    verifyCode: string
  }>();

  const hoverBgColor = useColorModeValue("rgba(0,0,0,.1)", "rgba(255,255,255,.1)")
  const toast = useToast()
  const dispatch = useAppDispatch()

  function onSubmit(values: FieldValues) {

  }

  return (
    <chakra.form onSubmit={handleSubmit(onSubmit)}>
      {/* <FormControl isInvalid={errors.verifyCode !== undefined}>
        <FormLabel htmlFor="verifyCode">Verification Code</FormLabel>
        <Input
          id="verifyCode"
          borderColor="gray.500"
          placeholder="xxxxxx"
          _hover={{ backgroundColor: hoverBgColor }}
          _focus={{ outline: "none" }}
          inputMode="text"
          {...register("verifyCode", {
            required: "Input your verification code",
            pattern: {
              value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
              message: "invalid email address"
            },
          })}
        />
        <FormErrorMessage>
          {errors.verifyCode && errors.verifyCode.message}
        </FormErrorMessage>
      </FormControl> */}

      <ButtonGroup
        spacing={6}
        mt={6}
        variant="outline"
      >
        <Button
          // isLoading={isLoading}
          type="submit"
          width="150px"
          borderColor="gray.500"
          _hover={{ backgroundColor: hoverBgColor }}
        >
          Resend
        </Button>
        <Button
          type="button"
          width="150px"
          onClick={handleClose}
          borderColor="gray.500"
          _hover={{ backgroundColor: hoverBgColor }}
        >
          Close
        </Button>
      </ButtonGroup >
    </chakra.form>
  )
}