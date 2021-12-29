import {
  createSlice,
  PayloadAction,
  createAsyncThunk,
} from '@reduxjs/toolkit'

interface CommonState {
  showLogin: boolean
  sidebarOpen: boolean
  version: string
}

const initialState = {
  showLogin: false,
  sidebarOpen: false,
  version: "v0.0.0"
} as CommonState

export const commonSlice = createSlice({
  name: "common",
  initialState,
  reducers: {
    toggleLogin(state, action: PayloadAction<boolean>) {
      state.showLogin = action.payload
    },
    toggleSidebar(state) {
      state.sidebarOpen = !state.sidebarOpen
    }
  },
})

export const {
  toggleSidebar,
  toggleLogin
} = commonSlice.actions

export default commonSlice.reducer