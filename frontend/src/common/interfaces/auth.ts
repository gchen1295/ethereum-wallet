export interface User {
  email: string
}

export interface LoginResponse {
  success?: boolean
  error?: string
}

export interface LoginRequest {
  email: string
  password: string
}
