<template>
  <div class="home">
    <!-- Hero Section -->
    <section class="bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white py-24 relative overflow-hidden">
      <!-- Background pattern -->
      <div class="absolute inset-0 bg-black/20"></div>
      <div class="absolute inset-0 bg-gradient-to-r from-blue-600/10 to-purple-600/10"></div>
      
      <div class="container mx-auto px-4 text-center relative z-10">
        <h1 class="text-6xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
          BJJ Store
        </h1>
        <p class="text-xl md:text-2xl mb-8 text-gray-300 max-w-2xl mx-auto">
          Premium Brazilian Jiu-Jitsu Gear & Equipment for Champions
        </p>
        <button
          @click="scrollToProducts"
          class="bg-blue-600 hover:bg-blue-700 px-10 py-4 rounded-lg font-semibold text-lg transition-all transform hover:scale-105 hover:shadow-xl"
        >
          Shop Now
        </button>
      </div>
    </section>

    <!-- Products Section -->
    <section id="products" class="py-20 bg-white">
      <div class="container mx-auto px-4">
        <div class="text-center mb-16">
          <h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">Featured Products</h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            Discover our handpicked selection of premium BJJ gear trusted by athletes worldwide
          </p>
        </div>

        <!-- Loading State -->
        <div v-if="products_store.loading" class="text-center py-8">
          <div
            class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
          ></div>
          <p class="mt-2 text-gray-600">Loading products...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="products_store.error" class="text-center py-8">
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {{ products_store.error }}
          </div>
          <button
            @click="products_store.fetchProducts()"
            class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Try Again
          </button>
        </div>

        <!-- Products Grid -->
        <div
          v-else-if="products_store.hasProducts"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8"
        >
          <div
            v-for="product in products_store.products"
            :key="product.id"
            class="group bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2"
          >
            <!-- Product Image -->
            <div class="aspect-square bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center relative overflow-hidden">
              <img
                v-if="product.image_url"
                :src="product.image_url"
                :alt="product.name"
                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
              />
              <div v-else class="text-gray-400">
                <svg class="w-20 h-20" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    fill-rule="evenodd"
                    d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
              <!-- Overlay on hover -->
              <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-all duration-300"></div>
            </div>

            <!-- Product Info -->
            <div class="p-6">
              <!-- Category Badge -->
              <div class="mb-3">
                <span
                  class="inline-block bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full uppercase tracking-wide font-medium"
                >
                  {{ product.category }}
                </span>
              </div>

              <h3 class="font-bold text-xl mb-3 text-gray-900 group-hover:text-blue-600 transition-colors">{{ product.name }}</h3>
              <p class="text-gray-600 text-sm mb-4 line-clamp-2 leading-relaxed">{{ product.description }}</p>

              <!-- Price and Stock -->
              <div class="flex items-center justify-between mb-6">
                <span class="text-3xl font-bold text-gray-900">
                  ${{ product.price.toFixed(2) }}
                </span>
                <span
                  :class="product.stock > 0 ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50'"
                  class="text-sm px-3 py-1 rounded-full font-medium"
                >
                  {{ product.stock > 0 ? `${product.stock} in stock` : 'Out of stock' }}
                </span>
              </div>

              <!-- Action Buttons -->
              <div class="flex gap-3">
                <router-link
                  :to="`/products/${product.id}`"
                  class="flex-1 bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-all text-center font-semibold hover:shadow-lg"
                >
                  View Details
                </router-link>
                <button
                  v-if="product.stock > 0"
                  @click="addToCart(product)"
                  class="bg-green-600 text-white py-3 px-4 rounded-lg hover:bg-green-700 transition-all font-semibold hover:shadow-lg"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4m1.6 8L3 3H1m6 10v6a2 2 0 002 2h4a2 2 0 002-2v-6m-6 0a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- No Products -->
        <div v-else class="text-center py-8">
          <p class="text-gray-600">No products available</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useProductStore } from '@/stores/products_store'
import { useCartStore } from '@/stores/cart_store'
import type { Product } from '@/types'

// Store with your naming convention
const products_store = useProductStore()
const cart_store = useCartStore()

// Methods
const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

const addToCart = (product: Product) => {
  // Add to cart with default size (first available size or empty string)
  const sizeOptions = typeof product.size_options === 'string' 
    ? product.size_options.split(',').map(s => s.trim())
    : product.size_options || []
  
  cart_store.addItem({
    product_id: product.id,
    quantity: 1,
    size: sizeOptions.length > 0 ? sizeOptions[0] : '',
    price: product.price
  })
  
  alert(`Added ${product.name} to cart!`)
}

// Load products when component mounts
onMounted(() => {
  products_store.fetchProducts()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
