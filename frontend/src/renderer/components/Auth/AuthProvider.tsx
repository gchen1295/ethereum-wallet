import * as React from 'react';
import { useEffect } from 'react'
import { getAuth, onAuthStateChanged } from 'firebase/auth'
import { store, setAuth, logout } from '@store/store'

interface AuthProps {
  children: JSX.Element | JSX.Element[]
}

export const AuthProvider: React.FC<AuthProps> = ({ children }) => {
  const auth = getAuth();
  useEffect(() => {
    return onAuthStateChanged(auth, (user) => {
      if (!user) {
        store.dispatch(logout())
        return
      }

      store.dispatch(setAuth({
        name: user.displayName ?? '',
        email: user.email ?? '',
        claim: user.refreshToken ?? '',
        uid: user.uid ?? '',
        photo: user.photoURL ?? '',
      }))
    })
  }, [])

  return <>{children}</>
}

export default AuthProvider