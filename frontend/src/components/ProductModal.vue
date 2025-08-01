<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="$emit('close')"></div>

      <!-- Modal -->
      <div class="inline-block w-full max-w-2xl p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-lg">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">
            {{ isEdit ? 'Edit Product' : 'Add New Product' }}
          </h3>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Product Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Product Name *</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.name }"
            />
            <p v-if="errors.name" class="text-red-500 text-sm mt-1">{{ errors.name }}</p>
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Description *</label>
            <textarea
              v-model="form.description"
              rows="4"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.description }"
            ></textarea>
            <p v-if="errors.description" class="text-red-500 text-sm mt-1">{{ errors.description }}</p>
          </div>

          <!-- Price and Category -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Price *</label>
              <input
                v-model.number="form.price"
                type="number"
                step="0.01"
                min="0"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.price }"
              />
              <p v-if="errors.price" class="text-red-500 text-sm mt-1">{{ errors.price }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Category *</label>
              <input
                v-model="form.category"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.category }"
                placeholder="e.g., gi, rashguard, belt"
              />
              <p v-if="errors.category" class="text-red-500 text-sm mt-1">{{ errors.category }}</p>
            </div>
          </div>

          <!-- Stock and Image URL -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Stock *</label>
              <input
                v-model.number="form.stock"
                type="number"
                min="0"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.stock }"
              />
              <p v-if="errors.stock" class="text-red-500 text-sm mt-1">{{ errors.stock }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Product Image URL</label>
              
              <!-- Image URL Input -->
              <div class="space-y-4">
                <input
                  v-model="form.image_url"
                  type="url"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="https://example.com/image.jpg"
                />
                <p class="text-sm text-gray-500">Enter a direct link to the product image</p>
                
                <!-- Image Preview -->
                <div v-if="form.image_url" class="relative">
                  <img 
                    :src="form.image_url" 
                    alt="Product preview" 
                    class="w-32 h-32 object-cover rounded-lg border border-gray-300"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Size Options -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Size Options</label>
            <input
              v-model="form.size_options"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="A1, A2, A3, A4 (comma separated)"
            />
            <p class="text-sm text-gray-500 mt-1">Enter sizes separated by commas (e.g., A1, A2, A3, A4)</p>
          </div>

          <!-- Actions -->
          <div class="flex justify-end space-x-4 pt-6 border-t">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading">{{ isEdit ? 'Updating...' : 'Creating...' }}</span>
              <span v-else>{{ isEdit ? 'Update Product' : 'Create Product' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { Product } from '@/types'

interface Props {
  product?: Product | null
  isEdit?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  product: null,
  isEdit: false
})

const emit = defineEmits<{
  close: []
  save: [data: any]
}>()

// Form data
const form = ref({
  name: '',
  description: '',
  price: 0,
  category: '',
  stock: 0,
  image_url: '',
  size_options: ''
})

// State
const imageError = ref(false)

// Form validation
const errors = ref({} as Record<string, string>)
const loading = ref(false)

// Methods
const validateForm = () => {
  errors.value = {}

  if (!form.value.name.trim()) {
    errors.value.name = 'Product name is required'
  }

  if (!form.value.description.trim()) {
    errors.value.description = 'Description is required'
  }

  if (form.value.price <= 0) {
    errors.value.price = 'Price must be greater than 0'
  }

  if (!form.value.category.trim()) {
    errors.value.category = 'Category is required'
  }

  if (form.value.stock < 0) {
    errors.value.stock = 'Stock cannot be negative'
  }

  return Object.keys(errors.value).length === 0
}


const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true

  try {
    const productData = {
      name: form.value.name.trim(),
      description: form.value.description.trim(),
      price: form.value.price,
      category: form.value.category.trim(),
      stock: form.value.stock,
      image_url: form.value.image_url.trim(),
      size_options: form.value.size_options.trim()
    }

    emit('save', productData)
  } catch (error: any) {
    console.error('Form submission error:', error)
    alert('Failed to save product. Please try again.')
  } finally {
    loading.value = false
  }
}

// Watch for product changes to populate form
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    form.value = {
      name: newProduct.name || '',
      description: newProduct.description || '',
      price: newProduct.price || 0,
      category: newProduct.category || '',
      stock: newProduct.stock || 0,
      image_url: newProduct.image_url || '',
      size_options: typeof newProduct.size_options === 'string' 
        ? newProduct.size_options 
        : Array.isArray(newProduct.size_options) 
          ? newProduct.size_options.join(', ')
          : ''
    }
    imageError.value = false
  }
}, { immediate: true })

// Initialize form on mount
onMounted(() => {
  if (props.product) {
    form.value = {
      name: props.product.name || '',
      description: props.product.description || '',
      price: props.product.price || 0,
      category: props.product.category || '',
      stock: props.product.stock || 0,
      image_url: props.product.image_url || '',
      size_options: typeof props.product.size_options === 'string' 
        ? props.product.size_options 
        : Array.isArray(props.product.size_options) 
          ? props.product.size_options.join(', ')
          : ''
    }
  }
})
</script>