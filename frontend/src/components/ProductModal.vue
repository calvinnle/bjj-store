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
              <label class="block text-sm font-medium text-gray-700 mb-2">Product Image</label>
              
              <!-- Image Upload Area -->
              <div class="space-y-4">
                <!-- Current Image Preview -->
                <div v-if="imagePreview || form.image_url" class="relative">
                  <img 
                    :src="imagePreview || form.image_url" 
                    alt="Product preview" 
                    class="w-32 h-32 object-cover rounded-lg border border-gray-300"
                  />
                  <button
                    type="button"
                    @click="clearImage"
                    class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center hover:bg-red-600 transition-colors"
                  >
                    ×
                  </button>
                </div>

                <!-- Upload Options -->
                <div class="space-y-3">
                  <!-- File Upload -->
                  <div>
                    <label class="block text-sm font-medium text-gray-600 mb-2">Upload from device:</label>
                    <div class="flex items-center space-x-3">
                      <input
                        ref="fileInput"
                        type="file"
                        accept="image/*"
                        @change="handleFileUpload"
                        class="hidden"
                      />
                      <button
                        type="button"
                        @click="$refs.fileInput?.click()"
                        class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center space-x-2"
                      >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                        </svg>
                        <span>Choose Image</span>
                      </button>
                      <span v-if="selectedFile" class="text-sm text-gray-600">{{ selectedFile.name }}</span>
                    </div>
                  </div>

                  <!-- Or URL Input -->
                  <div>
                    <label class="block text-sm font-medium text-gray-600 mb-2">Or enter image URL:</label>
                    <input
                      v-model="form.image_url"
                      type="url"
                      @input="handleUrlInput"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      placeholder="https://example.com/image.jpg"
                    />
                  </div>
                </div>

                <!-- Upload Progress -->
                <div v-if="uploadProgress > 0 && uploadProgress < 100" class="space-y-2">
                  <div class="flex justify-between text-sm text-gray-600">
                    <span>Uploading...</span>
                    <span>{{ uploadProgress }}%</span>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2">
                    <div class="bg-blue-600 h-2 rounded-full transition-all duration-300" :style="`width: ${uploadProgress}%`"></div>
                  </div>
                </div>

                <!-- Upload Success -->
                <div v-if="uploadProgress === 100" class="flex items-center space-x-2 text-green-600 text-sm">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span>Image uploaded successfully!</span>
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
import { uploadService } from '@/services/upload'

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

// File upload state
const selectedFile = ref<File | null>(null)
const imagePreview = ref('')
const uploadProgress = ref(0)
const fileInput = ref<HTMLInputElement>()

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

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return

  // Validate file type
  if (!file.type.startsWith('image/')) {
    alert('Please select an image file')
    return
  }

  // Validate file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    alert('Image size must be less than 5MB')
    return
  }

  selectedFile.value = file
  
  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
    form.value.image_url = '' // Clear URL when file is selected
  }
  reader.readAsDataURL(file)
}

const handleUrlInput = () => {
  if (form.value.image_url) {
    selectedFile.value = null
    imagePreview.value = ''
  }
}

const clearImage = () => {
  selectedFile.value = null
  imagePreview.value = ''
  form.value.image_url = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const uploadImage = async (): Promise<string> => {
  if (!selectedFile.value) return form.value.image_url

  uploadProgress.value = 1 // Start progress

  try {
    // Simulate progress for better UX
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 10
      }
    }, 100)

    const imageUrl = await uploadService.uploadImage(selectedFile.value)

    clearInterval(progressInterval)
    uploadProgress.value = 100
    
    // Keep success message for a few seconds
    setTimeout(() => {
      if (uploadProgress.value === 100) {
        uploadProgress.value = 0
      }
    }, 3000)
    
    return imageUrl
  } catch (error) {
    uploadProgress.value = 0
    console.error('Image upload error:', error)
    throw error
  }
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true

  try {
    // Upload image if file is selected
    let imageUrl = form.value.image_url
    if (selectedFile.value) {
      imageUrl = await uploadImage()
    }

    const productData = {
      name: form.value.name.trim(),
      description: form.value.description.trim(),
      price: form.value.price,
      category: form.value.category.trim(),
      stock: form.value.stock,
      image_url: imageUrl.trim(),
      size_options: form.value.size_options.trim()
    }

    emit('save', productData)
  } catch (error: any) {
    console.error('Form submission error:', error)
    let errorMessage = 'Failed to save product. Please try again.'
    if (error instanceof Error) {
      errorMessage = `Failed to save product: ${error.message}`
    }
    alert(errorMessage)
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
    // Clear upload state when switching products
    clearImage()
  }
}, { immediate: true })

// Watch for modal close to reset upload state
watch(() => props.product, (newProduct, oldProduct) => {
  if (!newProduct && oldProduct) {
    // Modal is closing, reset upload state
    clearImage()
    uploadProgress.value = 0
  }
})

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