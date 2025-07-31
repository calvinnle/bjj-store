<template>
  <div class="admin-dashboard">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-4">
            <router-link to="/" class="flex items-center space-x-2">
              <div class="bg-blue-600 text-white p-2 rounded-lg">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <span class="text-lg font-bold text-gray-900">BJJ Store Admin</span>
            </router-link>
          </div>

          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-600">
              Welcome, {{ auth_store.user?.email }}
              <span class="ml-1 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                {{ formatRole(auth_store.user?.role) }}
              </span>
            </span>
            <button
              @click="handleLogout"
              class="text-sm text-gray-600 hover:text-red-600 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="container mx-auto px-4 py-8">
      <!-- Navigation Tabs -->
      <div class="flex space-x-1 bg-gray-100 p-1 rounded-lg mb-8">
        <button
          v-if="auth_store.canViewProducts || auth_store.canViewOrders"
          @click="activeTab = 'dashboard'"
          :class="[
            'flex-1 py-2 px-4 text-sm font-medium rounded-md transition-colors',
            activeTab === 'dashboard' 
              ? 'bg-white text-gray-900 shadow-sm' 
              : 'text-gray-600 hover:text-gray-900'
          ]"
        >
          Dashboard
        </button>
        <button
          v-if="auth_store.canViewProducts"
          @click="activeTab = 'products'"
          :class="[
            'flex-1 py-2 px-4 text-sm font-medium rounded-md transition-colors',
            activeTab === 'products' 
              ? 'bg-white text-gray-900 shadow-sm' 
              : 'text-gray-600 hover:text-gray-900'
          ]"
        >
          Products
        </button>
        <button
          v-if="auth_store.canViewOrders"
          @click="activeTab = 'orders'"
          :class="[
            'flex-1 py-2 px-4 text-sm font-medium rounded-md transition-colors',
            activeTab === 'orders' 
              ? 'bg-white text-gray-900 shadow-sm' 
              : 'text-gray-600 hover:text-gray-900'
          ]"
        >
          Orders
        </button>
      </div>

      <!-- Dashboard Tab -->
      <div v-if="activeTab === 'dashboard'" class="space-y-6">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="bg-white p-6 rounded-lg shadow-sm border">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 rounded-lg">
                <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Total Products</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.totalProducts || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white p-6 rounded-lg shadow-sm border">
            <div class="flex items-center">
              <div class="p-2 bg-green-100 rounded-lg">
                <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Total Orders</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.totalOrders || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white p-6 rounded-lg shadow-sm border">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 rounded-lg">
                <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Revenue</p>
                <p class="text-2xl font-semibold text-gray-900">${{ (stats.totalRevenue || 0).toFixed(2) }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white p-6 rounded-lg shadow-sm border">
            <div class="flex items-center">
              <div class="p-2 bg-orange-100 rounded-lg">
                <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Pending Orders</p>
                <p class="text-2xl font-semibold text-gray-900">{{ stats.pendingOrders || 0 }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Orders -->
        <div v-if="auth_store.canViewOrders" class="bg-white rounded-lg shadow-sm border">
          <div class="p-6 border-b border-gray-200">
            <h3 class="text-lg font-semibold text-gray-900">Recent Orders</h3>
          </div>
          <div class="p-6">
            <div v-if="loadingStats" class="text-center py-4">
              <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            </div>
            <div v-else-if="recentOrders.length === 0" class="text-center py-8 text-gray-500">
              No recent orders
            </div>
            <div v-else class="space-y-4">
              <div 
                v-for="order in recentOrders" 
                :key="order.id"
                class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div>
                  <p class="font-medium text-gray-900">{{ order.order_number }}</p>
                  <p class="text-sm text-gray-500">{{ order.guest_email }}</p>
                </div>
                <div class="text-right">
                  <p class="font-medium text-gray-900">${{ order.total_amount.toFixed(2) }}</p>
                  <span :class="getStatusClass(order.status)" class="inline-block px-2 py-1 rounded-full text-xs font-medium">
                    {{ formatStatus(order.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Products Tab -->
      <div v-if="activeTab === 'products'" class="space-y-6">
        <AdminProductManagement />
      </div>

      <!-- Orders Tab -->
      <div v-if="activeTab === 'orders'" class="space-y-6">
        <AdminOrderManagement />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth_store'
import { authService } from '@/services/auth'
import AdminProductManagement from '@/components/AdminProductManagement.vue'
import AdminOrderManagement from '@/components/AdminOrderManagement.vue'
import type { Order } from '@/types'

const router = useRouter()
const auth_store = useAuthStore()

// Component state
const activeTab = ref('dashboard')
const stats = ref({
  totalProducts: 0,
  totalOrders: 0,
  totalRevenue: 0,
  pendingOrders: 0
})
const recentOrders = ref<Order[]>([])
const loadingStats = ref(false)

// Methods
const handleLogout = async () => {
  try {
    await auth_store.logout()
    router.push('/admin/login')
  } catch (error) {
    console.error('Logout failed:', error)
    // Force logout on error
    router.push('/admin/login')
  }
}

const loadStats = async () => {
  loadingStats.value = true
  try {
    const data = await authService.getStats()
    stats.value = data.stats || {}
    recentOrders.value = data.recentOrders || []
  } catch (error) {
    console.error('Failed to load stats:', error)
  } finally {
    loadingStats.value = false
  }
}

const formatRole = (role?: string) => {
  if (!role) return ''
  return role.split('_').map(word => 
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}

const formatStatus = (status: string) => {
  return status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ')
}

const getStatusClass = (status: string) => {
  const statusClasses = {
    pending: 'bg-yellow-100 text-yellow-800',
    paid: 'bg-blue-100 text-blue-800',
    shipped: 'bg-purple-100 text-purple-800',
    delivered: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800'
  }
  return statusClasses[status as keyof typeof statusClasses] || 'bg-gray-100 text-gray-800'
}

// Check authentication and load data
onMounted(async () => {
  if (!auth_store.isAuthenticated) {
    router.push('/admin/login')
    return
  }

  // Load user profile if not already loaded
  if (!auth_store.user) {
    await auth_store.loadProfile()
  }

  // Load dashboard stats
  if (auth_store.canViewProducts || auth_store.canViewOrders) {
    await loadStats()
  }

  // Set default tab based on permissions
  if (!auth_store.canViewProducts && !auth_store.canViewOrders) {
    activeTab.value = 'dashboard'
  } else if (!auth_store.canViewProducts && auth_store.canViewOrders) {
    activeTab.value = 'orders'
  } else if (auth_store.canViewProducts && !auth_store.canViewOrders) {
    activeTab.value = 'products'
  }
})
</script>