import { defineStore } from 'pinia'
import { authService } from '@/services/auth'
import type { AdminUser } from '@/types'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as AdminUser | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.user,
    isLoading: (state) => state.loading,
    hasError: (state) => !!state.error,
    userRole: (state) => state.user?.role || null,
    canManageProducts: (state) => {
      if (!state.user) return false
      return ['super_admin', 'inventory'].includes(state.user.role)
    },
    canManageOrders: (state) => {
      if (!state.user) return false
      return ['super_admin', 'order_manager'].includes(state.user.role)
    },
    canViewOrders: (state) => {
      if (!state.user) return false
      return ['super_admin', 'order_manager', 'viewer'].includes(state.user.role)
    },
    canViewProducts: (state) => {
      if (!state.user) return false
      return ['super_admin', 'inventory', 'viewer'].includes(state.user.role)
    },
  },

  actions: {
    // Login
    async login(email: string, password: string) {
      this.loading = true
      this.error = null

      try {
        const response = await authService.login(email, password)

        // Store token
        authService.setToken(response.token)

        // Set user
        this.user = response.admin

        return response
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Login failed'
        throw error
      } finally {
        this.loading = false
      }
    },

    // Logout
    async logout() {
      this.loading = true

      try {
        if (authService.isAuthenticated()) {
          await authService.logout()
        }
      } catch (error) {
        // Continue with logout even if API call fails
        console.error('Logout API call failed:', error)
      } finally {
        // Clear local state regardless
        authService.removeToken()
        this.user = null
        this.error = null
        this.loading = false
      }
    },

    // Load user profile
    async loadProfile() {
      if (!authService.isAuthenticated()) {
        return
      }

      this.loading = true
      this.error = null

      try {
        this.user = await authService.getProfile()
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to load profile'

        // If unauthorized, clear auth state
        if (error.response?.status === 401) {
          authService.removeToken()
          this.user = null
        }
      } finally {
        this.loading = false
      }
    },

    // Initialize auth state from localStorage
    async initializeAuth() {
      if (authService.isAuthenticated()) {
        await this.loadProfile()
      }
    },

    // Clear error
    clearError() {
      this.error = null
    },
  },
})
