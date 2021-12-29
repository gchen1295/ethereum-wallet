import * as React from 'react';
import { useHistory } from 'react-router-dom'
import { useEffect } from 'react'
import { useAppSelector } from '@store/store'

interface AuthGuardProps {
  children: JSX.Element
}

export const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
  const history = useHistory()
  const user = useAppSelector(state => state.auth.user)

  useEffect(() => {
    if (!user) history.push('/')
  }, [history, user])

  if (user)
    return <>{children}</>

  return null
}

export default AuthGuard