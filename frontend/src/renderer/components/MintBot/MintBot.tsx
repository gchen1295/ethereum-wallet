import * as React from 'react';
import { useForm } from "react-hook-form";
import {
  chakra,
  Grid,
  Flex,
  Box,
  Heading,
  Input,
  FormErrorMessage,
  FormLabel,
  FormControl,
  useColorModeValue,
} from "@chakra-ui/react"


export const MintBot: React.FC = () => {
  const {
    handleSubmit,
    register,
    formState: { errors }
  } = useForm();


  const hoverBgColor = useColorModeValue("rgba(0,0,0,.1)", "rgba(255,255,255,.1)")

  return (
    <Flex
      w="100%"
      h="100%"
      flexDir="column"
      justifyContent="center"
      alignItems="center"
    >
      <Box>
        <Heading>Create New Mint Task</Heading>
      </Box>
      <br />
      <Box>
        <Grid>
          <FormControl isInvalid={errors.contractAddress !== undefined} w="405px">
            <FormLabel htmlFor="contractAddress">Contract Address</FormLabel>
            <Input
              id="contractAddress"
              borderColor="gray.500"
              placeholder="0x00000000000000000000"
              _hover={{ backgroundColor: hoverBgColor }}
              _focus={{ outline: "none" }}
              inputMode="text"
              {...register("contractAddress", {
                required: "Address is required",
                pattern: {
                  value: /^0x[a-zA-Z0-9]{40}$/i,
                  message: "Invalid contract address"
                },
              })}
            />
            <FormErrorMessage>
              {errors.email && errors.email.message}
            </FormErrorMessage>
          </FormControl>
          <FormControl isInvalid={errors.contractAddress !== undefined} w="405px">
            <FormLabel htmlFor="contractAddress">Contract Address</FormLabel>
            <Input
              id="contractAddress"
              borderColor="gray.500"
              placeholder="0x00000000000000000000"
              _hover={{ backgroundColor: hoverBgColor }}
              _focus={{ outline: "none" }}
              inputMode="text"
              {...register("contractAddress", {
                required: "Address is required",
                pattern: {
                  value: /^0x[a-zA-Z0-9]{40}$/i,
                  message: "Invalid contract address"
                },
              })}
            />
            <FormErrorMessage>
              {errors.email && errors.email.message}
            </FormErrorMessage>
          </FormControl>
          
        </Grid>
      </Box>
    </Flex>
  )
}