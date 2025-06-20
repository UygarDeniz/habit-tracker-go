import { createContext, useContext, useEffect, useState } from 'react'

import type { User } from '../types/types'

export interface UserContextType {
  user: User | null
  setUser: React.Dispatch<React.SetStateAction<User | null>>
  isAuthenticated: boolean
  isLoading: boolean
  logout: () => Promise<void>
}

const UserContext = createContext<UserContextType | undefined>(undefined)

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await fetch('/api/auth/session', {
          credentials: 'include',
        })
        if (!res.ok) {
          throw new Error('Failed to fetch user')
        }
        const data = await res.json()
        setUser({
          accessToken: data.access_token,
          ...data.user,
        })
      } catch (error) {
        console.log(error)
      } finally {
        setIsLoading(false)
      }
    }
    fetchUser()
  }, [])

  const logout = async () => {
    try {
      await fetch('/api/auth/logout', { method: 'POST' })
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      setUser(null)
    }
  }

  const isAuthenticated = !!user
  console.log(isAuthenticated)
  return (
    <UserContext.Provider
      value={{ user, setUser, isAuthenticated, isLoading, logout }}
    >
      {children}
    </UserContext.Provider>
  )
}

export const useUser = () => {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useUser must be used within an UserProvider')
  }
  return context
}
