import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authApi } from '@/api'
import type { User, LoginRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const isLoading = ref(false)
  const isAuthInitialized = ref(false)

  const isLoggedIn = computed(() => {
    return !!token.value && !!user.value
  })
  const isAdmin = computed(() => user.value?.role === 'admin')

  const login = async (credentials: LoginRequest) => {
    try {
      isLoading.value = true
      const response = await authApi.login(credentials)
      
      token.value = response.token
      user.value = response.user
      localStorage.setItem('token', response.token)
      
      return response
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
    }
  }

  const getCurrentUser = async () => {
    if (!token.value) return null
    
    try {
      const currentUser = await authApi.getCurrentUser()
      user.value = currentUser
      return currentUser
    } catch (error) {
      console.error('Get current user error:', error)
      logout()
      return null
    }
  }

  const initAuth = async () => {
    if (token.value) {
      await getCurrentUser()
    }
    isAuthInitialized.value = true
  }

  return {
    user,
    token,
    isLoading,
    isAuthInitialized,
    isLoggedIn,
    isAdmin,
    login,
    logout,
    getCurrentUser,
    initAuth
  }
})