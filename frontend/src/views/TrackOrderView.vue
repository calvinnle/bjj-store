<template>
  <div class="track-order">
    <!-- Navigation -->
    <nav class="bg-gray-50 py-4">
      <div class="container mx-auto px-4">
        <div class="flex items-center space-x-2 text-sm text-gray-600">
          <router-link to="/" class="hover:text-blue-600">Home</router-link>
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-gray-900">Track Order</span>
        </div>
      </div>
    </nav>

    <div class="container mx-auto px-4 py-8">
      <div class="max-w-2xl mx-auto">
        <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">Track Your Order</h1>
        
        <!-- Search Methods -->
        <div class="bg-white rounded-lg shadow-sm border p-6 mb-8">
          <div class="space-y-6">
            <!-- Search by Order Number -->
            <div>
              <h2 class="text-lg font-semibold text-gray-900 mb-4">Search by Order Number</h2>
              <div class="flex gap-4">
                <input
                  v-model="orderNumber"
                  type="text"
                  placeholder="Enter your order number (e.g., BJJ-1234567890)"
                  class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  @keyup.enter="searchByOrderNumber"
                />
                <button
                  @click="searchByOrderNumber"
                  :disabled="!orderNumber || loadingOrder"
                  class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span v-if="loadingOrder">Searching...</span>
                  <span v-else>Track</span>
                </button>
              </div>
            </div>

            <!-- Divider -->
            <div class="relative">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-gray-200"></div>
              </div>
              <div class="relative flex justify-center text-sm">
                <span class="bg-white px-4 text-gray-500">or</span>
              </div>
            </div>

            <!-- Search by Email -->
            <div>
              <h2 class="text-lg font-semibold text-gray-900 mb-4">Search by Email</h2>
              <div class="flex gap-4">
                <input
                  v-model="email"
                  type="email"
                  placeholder="Enter your email address"
                  class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  @keyup.enter="searchByEmail"
                />
                <button
                  @click="searchByEmail"
                  :disabled="!email || loadingOrders"
                  class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span v-if="loadingOrders">Searching...</span>
                  <span v-else>Search</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Error Messages -->
        <div v-if="errorMessage" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
          {{ errorMessage }}
        </div>

        <!-- Single Order Result -->
        <div v-if="singleOrder" class="bg-white rounded-lg shadow-sm border p-6 mb-6">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-semibold text-gray-900">Order Details</h2>
            <span :class="getStatusClass(singleOrder.status)" class="px-3 py-1 rounded-full text-sm font-medium">
              {{ formatStatus(singleOrder.status) }}
            </span>
          </div>

          <!-- Order Info Grid -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Order Number</h3>
              <p class="text-lg font-semibold text-gray-900">{{ singleOrder.order_number }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Order Date</h3>
              <p class="text-lg text-gray-900">{{ formatDate(singleOrder.created_at) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Total Amount</h3>
              <p class="text-lg font-semibold text-gray-900">${{ singleOrder.total_amount.toFixed(2) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Email</h3>
              <p class="text-lg text-gray-900">{{ singleOrder.guest_email }}</p>
            </div>
          </div>

          <!-- Progress Steps -->
          <div class="mb-6">
            <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-4">Order Progress</h3>
            <div class="flex items-center space-x-4">
              <div class="flex items-center">
                <div :class="getStepClass('pending', singleOrder.status)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <svg v-if="isStepCompleted('pending', singleOrder.status)" class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span v-else class="text-xs font-medium">1</span>
                </div>
                <span class="ml-2 text-sm text-gray-700">Ordered</span>
              </div>
              <div class="flex-1 h-0.5 bg-gray-200"></div>
              <div class="flex items-center">
                <div :class="getStepClass('paid', singleOrder.status)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <svg v-if="isStepCompleted('paid', singleOrder.status)" class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span v-else class="text-xs font-medium">2</span>
                </div>
                <span class="ml-2 text-sm text-gray-700">Paid</span>
              </div>
              <div class="flex-1 h-0.5 bg-gray-200"></div>
              <div class="flex items-center">
                <div :class="getStepClass('shipped', singleOrder.status)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <svg v-if="isStepCompleted('shipped', singleOrder.status)" class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span v-else class="text-xs font-medium">3</span>
                </div>
                <span class="ml-2 text-sm text-gray-700">Shipped</span>
              </div>
              <div class="flex-1 h-0.5 bg-gray-200"></div>
              <div class="flex items-center">
                <div :class="getStepClass('delivered', singleOrder.status)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <svg v-if="isStepCompleted('delivered', singleOrder.status)" class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span v-else class="text-xs font-medium">4</span>
                </div>
                <span class="ml-2 text-sm text-gray-700">Delivered</span>
              </div>
            </div>
          </div>

          <!-- Shipping Address -->
          <div class="border-t pt-6">
            <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-3">Shipping Address</h3>
            <div class="text-gray-700">
              <p>{{ singleOrder.shipping_address.first_name }} {{ singleOrder.shipping_address.last_name }}</p>
              <p>{{ singleOrder.shipping_address.address1 }}</p>
              <p v-if="singleOrder.shipping_address.address2">{{ singleOrder.shipping_address.address2 }}</p>
              <p>{{ singleOrder.shipping_address.city }}, {{ singleOrder.shipping_address.state }} {{ singleOrder.shipping_address.zip_code }}</p>
              <p>{{ singleOrder.shipping_address.country }}</p>
            </div>
          </div>
        </div>

        <!-- Multiple Orders Result -->
        <div v-if="orders.length > 0" class="space-y-4">
          <h2 class="text-xl font-semibold text-gray-900">Your Orders</h2>
          <div 
            v-for="order in orders" 
            :key="order.id"
            class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow cursor-pointer"
            @click="viewOrderDetails(order)"
          >
            <div class="flex items-center justify-between">
              <div>
                <h3 class="font-semibold text-gray-900">{{ order.order_number }}</h3>
                <p class="text-sm text-gray-500">{{ formatDate(order.created_at) }}</p>
                <p class="text-sm font-medium text-gray-900">${{ order.total_amount.toFixed(2) }}</p>
              </div>
              <div class="text-right">
                <span :class="getStatusClass(order.status)" class="inline-block px-3 py-1 rounded-full text-sm font-medium">
                  {{ formatStatus(order.status) }}
                </span>
                <p class="text-sm text-gray-500 mt-2">{{ order.items?.length || 0 }} items</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { orderService } from '@/services/orders'
import type { Order } from '@/types'

// Component state
const orderNumber = ref('')
const email = ref('')
const singleOrder = ref<Order | null>(null)
const orders = ref<Order[]>([])
const loadingOrder = ref(false)
const loadingOrders = ref(false)
const errorMessage = ref('')

// Methods
const searchByOrderNumber = async () => {
  if (!orderNumber.value.trim()) return

  loadingOrder.value = true
  errorMessage.value = ''
  singleOrder.value = null
  orders.value = []

  try {
    singleOrder.value = await orderService.trackOrder(orderNumber.value.trim())
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || 'Order not found. Please check your order number.'
    console.error('Error tracking order:', error)
  } finally {
    loadingOrder.value = false
  }
}

const searchByEmail = async () => {
  if (!email.value.trim()) return

  loadingOrders.value = true
  errorMessage.value = ''
  singleOrder.value = null
  orders.value = []

  try {
    orders.value = await orderService.getOrdersByEmail(email.value.trim())
    if (orders.value.length === 0) {
      errorMessage.value = 'No orders found for this email address.'
    }
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || 'Failed to fetch orders. Please try again.'
    console.error('Error fetching orders:', error)
  } finally {
    loadingOrders.value = false
  }
}

const viewOrderDetails = (order: Order) => {
  singleOrder.value = order
  orders.value = []
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
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

const getStepClass = (step: string, currentStatus: string) => {
  const statusOrder = ['pending', 'paid', 'shipped', 'delivered']
  const stepIndex = statusOrder.indexOf(step)
  const currentIndex = statusOrder.indexOf(currentStatus)
  
  if (currentStatus === 'cancelled') {
    return 'bg-red-500 text-white'
  }
  
  if (currentIndex >= stepIndex) {
    return 'bg-blue-600 text-white'
  }
  
  return 'bg-gray-200 text-gray-600'
}

const isStepCompleted = (step: string, currentStatus: string) => {
  const statusOrder = ['pending', 'paid', 'shipped', 'delivered']
  const stepIndex = statusOrder.indexOf(step)
  const currentIndex = statusOrder.indexOf(currentStatus)
  
  return currentIndex >= stepIndex && currentStatus !== 'cancelled'
}
</script>