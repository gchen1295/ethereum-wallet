import { extendTheme, ThemeConfig } from "@chakra-ui/react"

const config: ThemeConfig = {
  initialColorMode: "dark",
  useSystemColorMode: true,
}

const theme = extendTheme({
  config: {
    styles: {
      global: {
        body: {
          innerWidth: "100%",
          outerWidth: "100%"
        },
      },
    },
    ...config,
  },
})

export default theme