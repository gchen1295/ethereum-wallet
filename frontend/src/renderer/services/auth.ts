import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { LoginResponse, LoginRequest } from '@interfaces/auth'
import { RootState } from '@store/store'


export const authApi = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: '/api',
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as RootState).auth.user?.claim
      if (token) {
        headers.set('authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (credentials) => ({
        url: 'auth',
        method: 'POST',
        body: credentials
      })
    }),
  }),
})

export const { useLoginMutation } = authApi