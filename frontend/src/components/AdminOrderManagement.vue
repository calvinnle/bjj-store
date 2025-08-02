<template>
  <div class="admin-orders">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-900">Order Management</h2>
      <div class="flex items-center space-x-4">
        <div class="flex items-center text-sm text-gray-500">
          <div class="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></div>
          Auto-refresh: 30s
        </div>
        <button
          @click="loadOrders"
          class="bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700"
          :disabled="loading"
        >
          <span v-if="loading">Loading...</span>
          <span v-else>Refresh Now</span>
        </button>
      </div>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white p-4 rounded-lg shadow-sm border mb-6">
      <div class="flex flex-col md:flex-row gap-4">
        <div class="flex-1">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by order number or email..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
        <div class="md:w-48">
          <select
            v-model="selectedStatus"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="paid">Paid</option>
            <option value="shipped">Shipped</option>
            <option value="delivered">Delivered</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Orders Table -->
    <div class="bg-white rounded-lg shadow-sm border overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="p-8 text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading orders...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="p-8 text-center">
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {{ error }}
        </div>
        <button
          @click="loadOrders"
          class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Try Again
        </button>
      </div>

      <!-- Orders Table -->
      <div v-else>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Order
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Customer
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Total
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th v-if="auth_store.canManageOrders" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="order in paginatedOrders" :key="order.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ order.order_number }}</div>
                    <div class="text-sm text-gray-500">{{ order.items?.length || 0 }} items</div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm font-medium text-gray-900">
                      {{ order.shipping_address.first_name }} {{ order.shipping_address.last_name }}
                    </div>
                    <div class="text-sm text-gray-500">{{ order.guest_email }}</div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(order.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  ${{ order.total_amount.toFixed(2) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getStatusClass(order.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ formatStatus(order.status) }}
                  </span>
                </td>
                <td v-if="auth_store.canManageOrders" class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                  <button
                    @click="viewOrder(order)"
                    class="text-blue-600 hover:text-blue-900"
                  >
                    View
                  </button>
                  <button
                    @click="updateOrderStatus(order)"
                    class="text-green-600 hover:text-green-900"
                  >
                    Update Status
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- No Orders -->
        <div v-if="filteredOrders.length === 0 && !loading" class="p-8 text-center text-gray-500">
          No orders found
        </div>
        
        <!-- Pagination (shows when there are orders) -->
        <div v-if="filteredOrders.length > 0" class="p-4 border-t">
          <div v-if="totalPages <= 1" class="pagination">
            <span class="page-item active">
              <span class="page-link">Page 1</span>
            </span>
          </div>
          <Paginate
            v-else
            v-model="currentPage"
            :page-count="totalPages"
            :page-range="5"
            :margin-pages="2"
            :click-handler="changePage"
            :prev-text="'Previous'"
            :next-text="'Next'"
            :container-class="'pagination'"
            :active-class="'active'"
            :disabled-class="'disabled'"
          />
        </div>
      </div>
    </div>

    <!-- Order Details Modal -->
    <OrderDetailsModal
      v-if="showOrderModal"
      :order="selectedOrder"
      @close="showOrderModal = false"
    />

    <!-- Update Status Modal -->
    <UpdateStatusModal
      v-if="showStatusModal"
      :order="selectedOrder"
      @close="showStatusModal = false"
      @update="handleStatusUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth_store'
import { orderService } from '@/services/orders'
import OrderDetailsModal from '@/components/OrderDetailsModal.vue'
import UpdateStatusModal from '@/components/UpdateStatusModal.vue'
import Paginate from 'vuejs-paginate-next'
import type { Order } from '@/types'

const auth_store = useAuthStore()

// Component state
const orders = ref<Order[]>([])
const loading = ref(false)
const error = ref('')
const searchQuery = ref('')
const selectedStatus = ref('')
const showOrderModal = ref(false)
const showStatusModal = ref(false)

// Auto-refresh
let refreshInterval: number | null = null
const selectedOrder = ref<Order | null>(null)

// Pagination
const currentPage = ref(1)
const itemsPerPage = 4

// Computed
const filteredOrders = computed(() => {
  let filtered = orders.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(order =>
      order.order_number.toLowerCase().includes(query) ||
      order.guest_email.toLowerCase().includes(query) ||
      `${order.shipping_address.first_name} ${order.shipping_address.last_name}`.toLowerCase().includes(query)
    )
  }

  if (selectedStatus.value) {
    filtered = filtered.filter(order => order.status === selectedStatus.value)
  }

  return filtered.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

// Paginated orders
const paginatedOrders = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredOrders.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredOrders.value.length / itemsPerPage)
})

// Methods
const loadOrders = async () => {
  loading.value = true
  error.value = ''

  try {
    orders.value = await orderService.getAllOrders()
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to load orders'
    console.error('Error loading orders:', err)
  } finally {
    loading.value = false
  }
}

const viewOrder = (order: Order) => {
  selectedOrder.value = order
  showOrderModal.value = true
}

const updateOrderStatus = (order: Order) => {
  selectedOrder.value = order
  showStatusModal.value = true
}

const changePage = (pageNum: number) => {
  currentPage.value = pageNum
}

const handleStatusUpdate = async (orderId: number, newStatus: string) => {
  try {
    await orderService.updateOrderStatus(orderId, newStatus)

    // Update local order
    const orderIndex = orders.value.findIndex(o => o.id === orderId)
    if (orderIndex !== -1) {
      orders.value[orderIndex].status = newStatus as any
    }

    showStatusModal.value = false
  } catch (error) {
    console.error('Failed to update order status:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
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

// Reset to page 1 when filters change
watch([searchQuery, selectedStatus], () => {
  currentPage.value = 1
})

// Load orders on mount
onMounted(() => {
  loadOrders()

  // Set up auto-refresh every 30 seconds
  refreshInterval = setInterval(() => {
    loadOrders()
  }, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
