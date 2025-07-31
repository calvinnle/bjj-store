<template>
  <div class="order-confirmation">
    <div class="container mx-auto px-4 py-16">
      <!-- Success Icon -->
      <div class="text-center mb-8">
        <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Order Confirmed!</h1>
        <p class="text-gray-600 mb-8">Thank you for your purchase. We'll send you updates about your order.</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading order details...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8">
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded max-w-md mx-auto mb-4">
          {{ error }}
        </div>
        <button 
          @click="loadOrder" 
          class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Try Again
        </button>
      </div>

      <!-- Order Details -->
      <div v-else-if="order" class="max-w-2xl mx-auto">
        <!-- Order Info -->
        <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Order Number</h3>
              <p class="text-lg font-semibold text-gray-900">{{ order.order_number }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Order Date</h3>
              <p class="text-lg text-gray-900">{{ formatDate(order.created_at) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Status</h3>
              <span :class="getStatusClass(order.status)" class="inline-block px-3 py-1 rounded-full text-sm font-medium">
                {{ formatStatus(order.status) }}
              </span>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">Total</h3>
              <p class="text-lg font-semibold text-gray-900">${{ order.total_amount.toFixed(2) }}</p>
            </div>
          </div>
        </div>

        <!-- Shipping Address -->
        <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Shipping Address</h3>
          <div class="text-gray-700">
            <p>{{ order.shipping_address.first_name }} {{ order.shipping_address.last_name }}</p>
            <p>{{ order.shipping_address.address1 }}</p>
            <p v-if="order.shipping_address.address2">{{ order.shipping_address.address2 }}</p>
            <p>{{ order.shipping_address.city }}, {{ order.shipping_address.state }} {{ order.shipping_address.zip_code }}</p>
            <p>{{ order.shipping_address.country }}</p>
          </div>
        </div>

        <!-- Order Items -->
        <div class="bg-white rounded-lg shadow-sm border p-6 mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Order Items</h3>
          <div class="space-y-4">
            <div 
              v-for="item in order.items" 
              :key="item.id"
              class="flex items-center space-x-4 py-4 border-b border-gray-200 last:border-b-0"
            >
              <div class="w-16 h-16 bg-gray-200 rounded-lg overflow-hidden flex-shrink-0">
                <img 
                  v-if="item.product?.image_url" 
                  :src="item.product.image_url" 
                  :alt="item.product.name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                  <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
              <div class="flex-1">
                <h4 class="font-medium text-gray-900">{{ item.product?.name }}</h4>
                <p class="text-sm text-gray-500">Size: {{ item.size }}</p>
                <p class="text-sm text-gray-500">Quantity: {{ item.quantity }}</p>
              </div>
              <div class="text-right">
                <p class="font-medium text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</p>
                <p class="text-sm text-gray-500">${{ item.price.toFixed(2) }} each</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="text-center space-y-4">
          <p class="text-gray-600">
            We'll send order updates to <strong>{{ order.guest_email }}</strong>
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <router-link 
              to="/track-order" 
              class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors"
            >
              Track Your Order
            </router-link>
            <router-link 
              to="/" 
              class="bg-gray-100 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Continue Shopping
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { orderService } from '@/services/orders'
import type { Order } from '@/types'

const route = useRoute()

// Component state
const order = ref<Order | null>(null)
const loading = ref(false)
const error = ref('')

// Methods
const loadOrder = async () => {
  const orderNumber = route.params.orderNumber as string
  if (!orderNumber) {
    error.value = 'Order number is required'
    return
  }

  loading.value = true
  error.value = ''

  try {
    order.value = await orderService.trackOrder(orderNumber)
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Order not found'
    console.error('Error loading order:', err)
  } finally {
    loading.value = false
  }
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

onMounted(() => {
  loadOrder()
})
</script>