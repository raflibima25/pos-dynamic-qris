import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '@/types'
import { api } from '@/lib/api'

interface AuthState {
  user: User | null
  token: string | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  checkAuth: () => Promise<void>
  setUser: (user: User) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isLoading: false,
      isAuthenticated: false,

      login: async (email: string, password: string) => {
        set({ isLoading: true })
        try {
          const response = await api.login(email, password)
          const { user, token } = response.data
          
          set({ 
            user, 
            token, 
            isAuthenticated: true, 
            isLoading: false 
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: () => {
        api.removeToken()
        set({ 
          user: null, 
          token: null, 
          isAuthenticated: false 
        })
      },

      checkAuth: async () => {
        const token = get().token
        if (!token) {
          set({ isAuthenticated: false })
          return
        }

        try {
          api.setToken(token)
          const response = await api.getCurrentUser()
          set({ 
            user: response.data, 
            isAuthenticated: true 
          })
        } catch (error) {
          // Token is invalid, clear auth state
          set({ 
            user: null, 
            token: null, 
            isAuthenticated: false 
          })
          api.removeToken()
        }
      },

      setUser: (user: User) => {
        set({ user })
      },
    }),
    {
      name: 'auth-store',
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token,
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
)