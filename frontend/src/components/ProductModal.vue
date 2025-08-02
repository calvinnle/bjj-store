<template>
  <!-- Modal overlay that covers the whole screen -->
  <div class="modal-backdrop">
    <div class="modal-content">
        <!-- Modal header with title and close button -->
        <div class="flex-between">
          <h3 class="text-lg font-semibold text-gray-900">
            {{ isEdit ? 'Edit Product' : 'Add New Product' }}
          </h3>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Product form with simple spacing -->
        <form @submit.prevent="handleSubmit" style="margin-top: 1.5rem;">
          <!-- Product Name -->
          <div class="form-group">
            <label class="form-label">Product Name *</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="form-input"
              :class="{ 'form-input-error': errors.name }"
            />
            <p v-if="errors.name" class="form-error-text">{{ errors.name }}</p>
          </div>

          <!-- Description -->
          <div class="form-group">
            <label class="form-label">Description *</label>
            <textarea
              v-model="form.description"
              rows="4"
              required
              class="form-input"
              :class="{ 'form-input-error': errors.description }"
            ></textarea>
            <p v-if="errors.description" class="form-error-text">{{ errors.description }}</p>
          </div>

          <!-- Price and Category side by side -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Price *</label>
              <input
                v-model.number="form.price"
                type="number"
                step="0.01"
                min="0"
                required
                class="form-input"
                :class="{ 'form-input-error': errors.price }"
              />
              <p v-if="errors.price" class="form-error-text">{{ errors.price }}</p>
            </div>
            <div class="form-group">
              <label class="form-label">Category *</label>
              <input
                v-model="form.category"
                type="text"
                required
                class="form-input"
                :class="{ 'form-input-error': errors.category }"
                placeholder="e.g., gi, rashguard, belt"
              />
              <p v-if="errors.category" class="form-error-text">{{ errors.category }}</p>
            </div>
          </div>

          <!-- Stock and Image URL side by side -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Stock *</label>
              <input
                v-model.number="form.stock"
                type="number"
                min="0"
                required
                class="form-input"
                :class="{ 'form-input-error': errors.stock }"
              />
              <p v-if="errors.stock" class="form-error-text">{{ errors.stock }}</p>
            </div>
            <div class="form-group">
              <label class="form-label">Product Image URL</label>
              <input
                v-model="form.image_url"
                type="url"
                class="form-input"
                placeholder="https://example.com/image.jpg"
              />
              <p style="color: #6b7280; font-size: 0.875rem; margin-top: 0.25rem;">Enter a direct link to the product image</p>
              
              <!-- Image Preview -->
              <div v-if="form.image_url" style="margin-top: 1rem;">
                <img 
                  :src="form.image_url" 
                  alt="Product preview" 
                  class="image-preview"
                />
              </div>
            </div>
          </div>

          <!-- Size Options -->
          <div class="form-group">
            <label class="form-label">Size Options</label>
            <input
              v-model="form.size_options"
              type="text"
              class="form-input"
              placeholder="A1, A2, A3, A4 (comma separated)"
            />
            <p style="color: #6b7280; font-size: 0.875rem; margin-top: 0.25rem;">Enter sizes separated by commas (e.g., A1, A2, A3, A4)</p>
          </div>

          <!-- Action buttons -->
          <div class="border-top flex-end">
            <button
              type="button"
              @click="$emit('close')"
              class="btn btn-secondary"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="btn btn-primary"
              :class="{ 'btn-disabled': loading }"
            >
              <span v-if="loading">{{ isEdit ? 'Updating...' : 'Creating...' }}</span>
              <span v-else>{{ isEdit ? 'Update Product' : 'Create Product' }}</span>
            </button>
          </div>
        </form>
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