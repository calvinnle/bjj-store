<template>
  <div class="product-detail">
    <!-- Navigation Breadcrumb -->
    <nav class="bg-gray-50 py-4">
      <div class="container mx-auto px-4">
        <div class="flex items-center space-x-2 text-sm text-gray-600">
          <router-link to="/" class="hover:text-blue-600">Home</router-link>
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-gray-900">{{ products_store.currentProduct?.name || 'Product' }}</span>
        </div>
      </div>
    </nav>

    <!-- Loading State -->
    <div v-if="products_store.loading" class="container mx-auto px-4 py-16">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading product...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="products_store.error" class="container mx-auto px-4 py-16">
      <div class="text-center">
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded max-w-md mx-auto">
          {{ products_store.error }}
        </div>
        <button 
          @click="loadProduct" 
          class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Try Again
        </button>
      </div>
    </div>

    <!-- Product Content -->
    <div v-else-if="product" class="container mx-auto px-4 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Product Image -->
        <div class="space-y-4">
          <div class="aspect-square bg-gray-200 rounded-lg overflow-hidden">
            <img 
              v-if="product.image_url" 
              :src="product.image_url" 
              :alt="product.name"
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
              <svg class="w-24 h-24" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Product Info -->
        <div class="space-y-6">
          <!-- Title and Price -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <span class="inline-block bg-gray-100 text-gray-800 text-sm px-3 py-1 rounded-full uppercase tracking-wide">
                {{ product.category }}
              </span>
              <span :class="product.stock > 0 ? 'text-green-600' : 'text-red-600'" class="text-sm font-medium">
                {{ product.stock > 0 ? `${product.stock} in stock` : 'Out of stock' }}
              </span>
            </div>
            <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ product.name }}</h1>
            <p class="text-4xl font-bold text-blue-600">${{ product.price.toFixed(2) }}</p>
          </div>

          <!-- Description -->
          <div>
            <h3 class="text-lg font-semibold mb-2">Description</h3>
            <p class="text-gray-700 leading-relaxed">{{ product.description }}</p>
          </div>

          <!-- Size Selection -->
          <div v-if="sizeOptions.length > 0">
            <h3 class="text-lg font-semibold mb-3">Size</h3>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="size in sizeOptions"
                :key="size"
                @click="selectedSize = size"
                :class="[
                  'px-4 py-2 border rounded-lg font-medium transition-colors',
                  selectedSize === size 
                    ? 'border-blue-600 bg-blue-600 text-white' 
                    : 'border-gray-300 text-gray-700 hover:border-blue-600'
                ]"
              >
                {{ size }}
              </button>
            </div>
          </div>

          <!-- Quantity -->
          <div>
            <h3 class="text-lg font-semibold mb-3">Quantity</h3>
            <div class="flex items-center space-x-4">
              <button 
                @click="decreaseQuantity"
                :disabled="quantity <= 1"
                class="w-10 h-10 border border-gray-300 rounded-lg flex items-center justify-center hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                -
              </button>
              <span class="text-xl font-semibold w-8 text-center">{{ quantity }}</span>
              <button 
                @click="increaseQuantity"
                :disabled="quantity >= product.stock"
                class="w-10 h-10 border border-gray-300 rounded-lg flex items-center justify-center hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                +
              </button>
            </div>
          </div>

          <!-- Add to Cart -->
          <div class="space-y-4">
            <button
              @click="addToCart"
              :disabled="!canAddToCart"
              class="w-full bg-blue-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ product.stock > 0 ? 'Add to Cart' : 'Out of Stock' }}
            </button>
            
            <!-- Success Message -->
            <div v-if="showAddedMessage" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
              Added to cart successfully!
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Product Not Found -->
    <div v-else class="container mx-auto px-4 py-16">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">Product Not Found</h1>
        <p class="text-gray-600 mb-8">The product you're looking for doesn't exist.</p>
        <router-link to="/" class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700">
          Back to Home
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useProductStore } from '@/stores/products_store'
import { useCartStore } from '@/stores/cart_store'

const route = useRoute()
const products_store = useProductStore()
const cart_store = useCartStore()

// Component state
const selectedSize = ref('')
const quantity = ref(1)
const showAddedMessage = ref(false)

// Computed
const product = computed(() => products_store.currentProduct)
const sizeOptions = computed(() => {
  if (!product.value?.size_options) return []
  return typeof product.value.size_options === 'string' 
    ? product.value.size_options.split(',').map(s => s.trim())
    : product.value.size_options
})

const canAddToCart = computed(() => {
  return product.value && 
         product.value.stock > 0 && 
         quantity.value > 0 && 
         quantity.value <= product.value.stock &&
         (sizeOptions.value.length === 0 || selectedSize.value)
})

// Methods
const loadProduct = async () => {
  const productId = Number(route.params.id)
  if (productId) {
    await products_store.fetchProduct(productId)
    // Set default size if available
    if (sizeOptions.value.length > 0) {
      selectedSize.value = sizeOptions.value[0]
    }
  }
}

const increaseQuantity = () => {
  if (product.value && quantity.value < product.value.stock) {
    quantity.value++
  }
}

const decreaseQuantity = () => {
  if (quantity.value > 1) {
    quantity.value--
  }
}

const addToCart = () => {
  if (!canAddToCart.value || !product.value) return

  cart_store.addItem({
    product_id: product.value.id,
    quantity: quantity.value,
    size: selectedSize.value,
    price: product.value.price
  })

  // Show success message
  showAddedMessage.value = true
  setTimeout(() => {
    showAddedMessage.value = false
  }, 3000)

  // Reset quantity to 1
  quantity.value = 1
}

// Watch for route changes
watch(() => route.params.id, () => {
  if (route.params.id) {
    loadProduct()
  }
})

// Load product on mount
onMounted(() => {
  loadProduct()
})
</script>