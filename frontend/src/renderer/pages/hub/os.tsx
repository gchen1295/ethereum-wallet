import * as React from 'react';
import styles from '@styles/hub/Home.module.css'
import {
  Grid,
  GridItem,
  useColorModeValue,
  Flex,
} from '@chakra-ui/react'

const HubOpensea = () => {
  const bgColor = useColorModeValue("var(--color-light)", "")

  return (
    <Grid
      className={styles.container}
      background={bgColor}
      h="100vh"
      gridTemplateColumns="repeat(3, 1fr)"
      gridTemplateRows="75px 1fr 1fr 60px"
      gridTemplateAreas={`
        'header header header'
        'main main main'
        'main main main'
        'footer footer footer'
      `}
      p="0"
    >
      <GridItem gridArea="main" minHeight="100%">
        <Flex alignItems="center" minHeight="calc(100vh - 135px)" width="100%">
          <Flex width="100%" justifyContent="center">
            Opensea Content Here
          </Flex>
        </Flex>
      </GridItem>
      <GridItem gridArea="footer">
      </GridItem>
    </Grid >
  )

}

export default HubOpensea
