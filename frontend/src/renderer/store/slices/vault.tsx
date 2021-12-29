import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface VaultState {
  accounts: string[]
}

const initialState = {
  accounts: [],
} as VaultState

const vaultSlice = createSlice({
  name: "vault",
  initialState,
  reducers: {
    setAccounts(state, action: PayloadAction<string[]>) {
      state.accounts = action.payload
    },
  },
})

export const { setAccounts } = vaultSlice.actions

export default vaultSlice.reducer