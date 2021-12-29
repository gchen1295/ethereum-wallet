import {createSlice,PayloadAction} from '@reduxjs/toolkit'

interface AuthState {
  user: User | null
  initializing: boolean
  vksp: string
}

export interface User {
  name: string
  email: string
  uid: string
  claim: string
  photo: string
}

const initialState = {
  user: null,
  vksp: '',
} as AuthState

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    logout(state) {
      state.user = null
    },
    setAuth(state, action: PayloadAction<User | null>) {
      if (!action.payload) {
        state.user = null
      } else {
        state.user = { ...action.payload }
      }
    },
    setInitializing(state, action: PayloadAction<boolean>) {
      state.initializing = action.payload
    },
    setVksp(state, action: PayloadAction<string>) {
      state.vksp = action.payload
    }
  },
})


export const {
  logout,
  setAuth,
  setInitializing,
} = authSlice.actions

export default authSlice.reducer