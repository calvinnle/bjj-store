<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="$emit('close')"></div>

      <!-- Modal -->
      <div class="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-lg">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">Update Order Status</h3>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div v-if="order">
          <!-- Order Info -->
          <div class="mb-6 p-4 bg-gray-50 rounded-lg">
            <p class="font-medium">{{ order.order_number }}</p>
            <p class="text-sm text-gray-600">{{ order.guest_email }}</p>
            <p class="text-sm text-gray-600">Current Status: 
              <span :class="getStatusClass(order.status)" class="inline-block px-2 py-1 rounded-full text-xs font-medium ml-1">
                {{ formatStatus(order.status) }}
              </span>
            </p>
          </div>

          <!-- Status Selection -->
          <form @submit.prevent="handleSubmit">
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-3">New Status</label>
              <div class="space-y-2">
                <label v-for="status in statusOptions" :key="status.value" class="flex items-center">
                  <input
                    v-model="selectedStatus"
                    type="radio"
                    :value="status.value"
                    class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                  />
                  <span class="ml-3 flex items-center">
                    <span :class="getStatusClass(status.value)" class="inline-block px-2 py-1 rounded-full text-xs font-medium mr-2">
                      {{ status.label }}
                    </span>
                    <span class="text-sm text-gray-600">{{ status.description }}</span>
                  </span>
                </label>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex justify-end space-x-4">
              <button
                type="button"
                @click="$emit('close')"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="!selectedStatus || loading"
                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span v-if="loading">Updating...</span>
                <span v-else>Update Status</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Order } from '@/types'

interface Props {
  order: Order | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  update: [orderId: number, status: string]
}>()

// Component state
const selectedStatus = ref('')
const loading = ref(false)

// Status options
const statusOptions = [
  { value: 'pending', label: 'Pending', description: 'Order received, payment pending' },
  { value: 'paid', label: 'Paid', description: 'Payment confirmed, preparing for shipment' },
  { value: 'shipped', label: 'Shipped', description: 'Order shipped to customer' },
  { value: 'delivered', label: 'Delivered', description: 'Order delivered to customer' },
  { value: 'cancelled', label: 'Cancelled', description: 'Order cancelled' },
]

// Methods
const handleSubmit = async () => {
  if (!props.order || !selectedStatus.value) return

  loading.value = true

  try {
    emit('update', props.order.id, selectedStatus.value)
  } catch (error) {
    console.error('Failed to update status:', error)
  } finally {
    loading.value = false
  }
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