<template>
  <!-- Modal overlay -->
  <div class="modal-backdrop">
    <div class="modal-content" style="max-width: 400px;">
        <!-- Modal header -->
        <div class="flex-between">
          <h3 class="text-lg font-semibold text-gray-900">Update Order Status</h3>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div v-if="order">
          <!-- Order Info -->
          <div style="margin: 1.5rem 0; padding: 1rem; background-color: #f9fafb; border-radius: 6px;">
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
            <div class="form-group">
              <label class="form-label">New Status</label>
              <div style="display: flex; flex-direction: column; gap: 0.5rem;">
                <label v-for="status in statusOptions" :key="status.value" style="display: flex; align-items: center; cursor: pointer;">
                  <input
                    v-model="selectedStatus"
                    type="radio"
                    :value="status.value"
                    style="margin-right: 0.75rem;"
                  />
                  <span style="display: flex; align-items: center;">
                    <span :class="getStatusClass(status.value)" class="inline-block px-2 py-1 rounded-full text-xs font-medium mr-2">
                      {{ status.label }}
                    </span>
                    <span style="font-size: 0.875rem; color: #6b7280;">{{ status.description }}</span>
                  </span>
                </label>
              </div>
            </div>

            <!-- Action buttons -->
            <div class="flex-end">
              <button
                type="button"
                @click="$emit('close')"
                class="btn btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="!selectedStatus || loading"
                class="btn btn-primary"
                :class="{ 'btn-disabled': !selectedStatus || loading }"
              >
                <span v-if="loading">Updating...</span>
                <span v-else>Update Status</span>
              </button>
            </div>
          </form>
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