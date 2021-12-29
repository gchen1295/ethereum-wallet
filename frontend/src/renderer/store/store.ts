import { configureStore } from '@reduxjs/toolkit'
import { combineReducers } from 'redux'
import { TypedUseSelectorHook, useSelector as _useSelector } from 'react-redux'
import authReducer from '@store/slices/auth'
import commonReducer from '@store/slices/common'
import vaultReducer from '@store/slices/vault'
import { authApi } from '@services/auth'

export * from '@store/slices/auth'
export * from '@store/slices/common'

const rootReducer = combineReducers({
  auth: authReducer,
  [authApi.reducerPath]: authApi.reducer,
  common: commonReducer,
  vault: vaultReducer,
})

export const store = configureStore({
  reducer: rootReducer
})

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
export const useSelector: TypedUseSelectorHook<RootState> = _useSelector
export * from './hooks'
export default store