import api from './api'
import type { AdminUser, LoginResponse } from '@/types'

export const authService = {
  // Admin login
  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await api.post('/api/admin/auth/login', { email, password })
    return response.data
  },

  // Admin logout
  async logout(): Promise<void> {
    await api.post('/api/admin/logout')
  },

  // Get admin profile
  async getProfile(): Promise<AdminUser> {
    const response = await api.get('/api/admin/profile')
    return response.data
  },

  // Get admin stats/dashboard data
  async getStats(): Promise<any> {
    const response = await api.get('/api/admin/stats')
    return response.data
  },

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return !!localStorage.getItem('admin_token')
  },

  // Get stored token
  getToken(): string | null {
    return localStorage.getItem('admin_token')
  },

  // Store token
  setToken(token: string): void {
    localStorage.setItem('admin_token', token)
  },

  // Remove token
  removeToken(): void {
    localStorage.removeItem('admin_token')
  },
}