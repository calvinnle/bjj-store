<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="$emit('close')"></div>

      <!-- Modal -->
      <div class="inline-block w-full max-w-4xl p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-lg">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">Order Details</h3>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div v-if="order" class="space-y-6">
          <!-- Order Info -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4 bg-gray-50 rounded-lg">
            <div>
              <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Order Number</h4>
              <p class="text-lg font-semibold text-gray-900">{{ order.order_number }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Date</h4>
              <p class="text-lg text-gray-900">{{ formatDate(order.created_at) }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Status</h4>
              <span :class="getStatusClass(order.status)" class="inline-block px-3 py-1 rounded-full text-sm font-medium">
                {{ formatStatus(order.status) }}
              </span>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Total</h4>
              <p class="text-lg font-semibold text-gray-900">${{ order.total_amount.toFixed(2) }}</p>
            </div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Customer Information -->
            <div class="bg-white border rounded-lg p-4">
              <h4 class="text-lg font-semibold text-gray-900 mb-4">Customer Information</h4>
              <div class="space-y-2">
                <p><span class="font-medium">Name:</span> {{ order.shipping_address.first_name }} {{ order.shipping_address.last_name }}</p>
                <p><span class="font-medium">Email:</span> {{ order.guest_email }}</p>
              </div>
            </div>

            <!-- Shipping Address -->
            <div class="bg-white border rounded-lg p-4">
              <h4 class="text-lg font-semibold text-gray-900 mb-4">Shipping Address</h4>
              <div class="text-gray-700">
                <p>{{ order.shipping_address.first_name }} {{ order.shipping_address.last_name }}</p>
                <p>{{ order.shipping_address.address1 }}</p>
                <p v-if="order.shipping_address.address2">{{ order.shipping_address.address2 }}</p>
                <p>{{ order.shipping_address.city }}, {{ order.shipping_address.state }} {{ order.shipping_address.zip_code }}</p>
                <p>{{ order.shipping_address.country }}</p>
              </div>
            </div>
          </div>

          <!-- Order Items -->
          <div class="bg-white border rounded-lg">
            <div class="p-4 border-b">
              <h4 class="text-lg font-semibold text-gray-900">Order Items</h4>
            </div>
            <div class="p-4">
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
                    <h5 class="font-medium text-gray-900">{{ item.product?.name || 'Product not found' }}</h5>
                    <p class="text-sm text-gray-500">Size: {{ item.size }}</p>
                    <p class="text-sm text-gray-500">Quantity: {{ item.quantity }}</p>
                    <p class="text-sm text-gray-500">Unit Price: ${{ item.price.toFixed(2) }}</p>
                  </div>
                  <div class="text-right">
                    <p class="font-medium text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</p>
                  </div>
                </div>
              </div>

              <!-- Order Summary -->
              <div class="border-t pt-4 mt-4">
                <div class="flex justify-between items-center text-lg font-semibold">
                  <span>Total:</span>
                  <span>${{ order.total_amount.toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Payment Information (if available) -->
          <div v-if="order.stripe_payment_id" class="bg-white border rounded-lg p-4">
            <h4 class="text-lg font-semibold text-gray-900 mb-4">Payment Information</h4>
            <p><span class="font-medium">Payment ID:</span> {{ order.stripe_payment_id }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Order } from '@/types'

interface Props {
  order: Order | null
}

defineProps<Props>()
defineEmits<{
  close: []
}>()

// Methods
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
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
</script>