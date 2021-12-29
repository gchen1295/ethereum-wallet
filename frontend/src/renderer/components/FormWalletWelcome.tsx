import * as React from 'react';
import {
  Flex,
  Heading,
  Button,
  Text,
  Textarea,
  Input,
  IconButton,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogBody,
  AlertDialogContent,
  useToast,
  useColorModeValue,
} from "@chakra-ui/react"

import { ArrowLeftIcon } from '@chakra-ui/icons'
import { FocusableElement } from '@chakra-ui/utils';


export const FormWalletWelcome: React.FC<{
  getSecret: () => Promise<string>;
  handleWalletForm: (secret: string, passphrase: string) => void;
  formLoading: boolean;
}> = (props) => {
  const [titleText, setTitleText] = React.useState("")
  const [description, setDescription] = React.useState("")
  const [secret, setSecret] = React.useState("")
  const [passphrase, setPassphrase] = React.useState("")
  const [confirmPassphrase, setConfirmPassphrase] = React.useState("")
  const [currentStep, setCurrentStep] = React.useState(0)
  const [mode, setMode] = React.useState(0)

  const [isOpen, setIsOpen] = React.useState(false)
  const onClose = () => setIsOpen(false)
  const cancelRef = React.useRef() as React.RefObject<FocusableElement>;

  const toast = useToast()
  const bgColor = useColorModeValue("gray.100", "gray.700")

  React.useEffect(() => {
    setStepValues()
  }, [currentStep])

  const setStepValues = () => {
    switch (currentStep) {
      case 0:
        setTitleText("Welcome!\nLet's setup a new wallet.")
        setDescription("Create or import an new wallet to get started.")
        break;
      case 1:
        if (mode === 0) {
          setTitleText("Import Wallet")
          setDescription("Import your wallet using your secret seed phrase.")
        } else {
          setTitleText("Secret Seed Phrase")
          setDescription("Copy and save your secret seed somewhere safe! There is no way to recover your account if you lose this.")
        }

        break;
      case 2:
        setTitleText("Password")
        setDescription("Set a password for your wallet.")
        break;

      default:
        break;
    }
  }

  const handleNextStep = () => {
    setCurrentStep(currentStep + 1)
  }

  const handlePreviousStep = () => {
    switch (currentStep - 1) {
      case 0:
        setSecret("")
        break
      case 1:
        setPassphrase("")
        setConfirmPassphrase("")
        break
    }
    setCurrentStep(currentStep - 1)
  }

  const handleConfirmSecret = (e: React.MouseEvent) => {
    if (secret === "") {
      toast({
        title: "Enter a Seed Phrase",
        description: "A seed phrase must be provided to continue.",
        status: "error",
        duration: 5000,
        isClosable: true,
        position: "bottom-right"
      })

      return
    }

    setCurrentStep(currentStep + 1)
  }


  const getStepElements = () => {
    switch (currentStep) {
      case 0:
        return (
          <React.Fragment>
            <Button colorScheme="facebook" mb={3} onClick={handleCreateWallet}>Create Wallet</Button>
            <Button colorScheme="facebook" onClick={handleImportWallet}>Import Wallet</Button>
          </React.Fragment>
        )
      case 1:
        if (mode == 0) {
          return (
            <React.Fragment>
              <Textarea value={secret} mb={3} height="120px" background="gray.50" readOnly cursor="default" color={fontColor}/>
              <Button colorScheme="red" onClick={() => setIsOpen(true)}>Confirm</Button>
            </React.Fragment >
          )
        } else {
          return (
            <React.Fragment>
              <Textarea value={secret} mb={3} height="120px" background="gray.50" onChange={(e) => { if (secret.length < 500) setSecret(e.target.value); }} color={fontColor}/>
              <Button colorScheme="facebook" onClick={handleConfirmSecret}>Import</Button>
            </React.Fragment>
          )
        }
      case 2:
        return (
          <React.Fragment>
            <Input value={passphrase} mb={3} background="gray.50" onChange={(e) => { if (passphrase.length < 40) setPassphrase(e.target.value) }} type="password" placeholder="Password" w="350px"/>
            <Input value={confirmPassphrase} mb={3} background="gray.50" onChange={(e) => { if (confirmPassphrase.length < 40) setConfirmPassphrase(e.target.value) }} type="password" placeholder="Confirm Password" w="350px"/>
            <Button colorScheme="facebook" isLoading={props.formLoading} onClick={() => {
              if (passphrase !== confirmPassphrase) {
                toast({
                  title: "Password mismatch",
                  description: "Passwords don't match.",
                  status: "error",
                  duration: 5000,
                  isClosable: true,
                  position: "bottom-right"
                })
                return
              }
              try {
                console.log(passphrase)
                props.handleWalletForm(secret, passphrase)
              } catch(e) {
                toast({
                  title: "Internal Error",
                  description: e.toString(),
                  status: "error",
                  duration: 5000,
                  isClosable: true,
                  position: "bottom-right"
                })
              }
              
            }} >Confirm</Button>
          </React.Fragment>
        )
    }
  }


  const handleCreateWallet = async () => {
    try {
      setMode(0)
      const secret = await props.getSecret()
      setSecret(secret)
      handleNextStep()
    }catch(e) {
      toast({
        title: "Internal Error",
        description: e.toString(),
        status: "error",
        duration: 5000,
        isClosable: true,
        position: "bottom-right"
      })
    }
  };

  const handleImportWallet = () => {
    setMode(1)
    handleNextStep()
  };

  const boxColor = useColorModeValue("rgba(0,0,0,0.08)", "rgba(255,255,255,0.08)")
  const fontColor = useColorModeValue("", "black")

  return (
    <Flex direction="column" maxWidth="60vw" borderWidth={1} borderRadius={8} boxShadow="lg"  p={12} rounded={6} bg={boxColor}>
      {
        currentStep > 0
          ? <IconButton aria-label="Search database" icon={<ArrowLeftIcon />} onClick={handlePreviousStep} position="fixed" marginLeft="-45px" marginTop={currentStep == 1 ? "130px" : "90px"} cursor="pointer" bg="rgba(0,0,0,0.0)"/>
          : <></>
      }
      <Heading mb={6} fontSize="1.6rem">{titleText}</Heading>
      <Text mb={6} fontSize="md">{description}</Text>
      {
        getStepElements()
      }
      <AlertDialog
        motionPreset="slideInBottom"
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
        isCentered
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Save Your Secret!
            </AlertDialogHeader>

            <AlertDialogBody>
              Make sure you have your secret saved. This is your last warning!
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef as React.LegacyRef<HTMLButtonElement>} onClick={onClose}>
                Cancel
              </Button>
              <Button colorScheme="red" onClick={() => {
                onClose()
                handleNextStep()
              }} ml={3}>
                Confirm
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </Flex>
  )
}


export default FormWalletWelcome