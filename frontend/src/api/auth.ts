import { apiClient } from './client'
import type { LoginRequest, LoginResponse, User } from '@/types'

export const authApi = {
  login: (data: LoginRequest): Promise<LoginResponse> => {
    return apiClient.post('/auth/login', data)
  },

  logout: (): Promise<void> => {
    return apiClient.post('/auth/logout')
  },

  getCurrentUser: (): Promise<User> => {
    return apiClient.get('/auth/me')
  },

  refreshToken: (): Promise<{ token: string }> => {
    return apiClient.post('/auth/refresh')
  }
}