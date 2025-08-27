'use client'

import { useEffect } from 'react'
import { useAuthStore } from '@/store/auth'
import { api } from '@/lib/api'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { checkAuth, token, isAuthenticated } = useAuthStore()

  useEffect(() => {
    // Initialize API token from store if available
    if (token) {
      api.setToken(token)
    }
    
    // Check authentication status
    if (token && !isAuthenticated) {
      checkAuth()
    }
  }, [token, isAuthenticated, checkAuth])

  return <>{children}</>
}