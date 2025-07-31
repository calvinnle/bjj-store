<template>
  <div class="products-page">
    <!-- Header -->
    <div class="bg-gray-900 text-white py-12">
      <div class="container mx-auto px-4">
        <h1 class="text-4xl font-bold mb-4">All Products</h1>
        <p class="text-xl text-gray-300">Browse our complete collection of BJJ gear</p>
      </div>
    </div>

    <div class="container mx-auto px-4 py-8">
      <!-- Search and Filters -->
      <div class="bg-white p-6 rounded-lg shadow-sm border mb-8">
        <div class="flex flex-col md:flex-row gap-4 mb-4">
          <!-- Search -->
          <div class="flex-1">
            <input
              v-model="products_store.searchQuery"
              type="text"
              placeholder="Search products..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          
          <!-- Category Filter -->
          <div class="md:w-48">
            <select
              v-model="products_store.selectedCategory"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Categories</option>
              <option v-for="category in products_store.availableCategories" :key="category" :value="category">
                {{ formatCategory(category) }}
              </option>
            </select>
          </div>
        </div>

        <!-- Filter Tags -->
        <div v-if="hasActiveFilters" class="flex items-center gap-2 flex-wrap">
          <span class="text-sm text-gray-600">Active filters:</span>
          <span
            v-if="products_store.searchQuery"
            class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
          >
            Search: "{{ products_store.searchQuery }}"
            <button @click="products_store.setSearchQuery('')" class="ml-2 text-blue-600 hover:text-blue-800">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </span>
          <span
            v-if="products_store.selectedCategory"
            class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-green-100 text-green-800"
          >
            Category: {{ formatCategory(products_store.selectedCategory) }}
            <button @click="products_store.setCategory('')" class="ml-2 text-green-600 hover:text-green-800">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </span>
          <button
            @click="products_store.clearFilters()"
            class="text-sm text-gray-600 hover:text-gray-800 underline"
          >
            Clear all
          </button>
        </div>

        <!-- Results Count -->
        <div class="mt-4 text-sm text-gray-600">
          Showing {{ filteredProductsCount }} of {{ products_store.products.length }} products
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="products_store.loading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading products...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="products_store.error" class="text-center py-16">
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded max-w-md mx-auto">
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
      <div v-else-if="products_store.filteredProducts.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div
          v-for="product in products_store.filteredProducts"
          :key="product.id"
          class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
        >
          <!-- Product Image -->
          <router-link :to="`/products/${product.id}`" class="block">
            <div class="aspect-square bg-gray-200 flex items-center justify-center">
              <img
                v-if="product.image_url"
                :src="product.image_url"
                :alt="product.name"
                class="w-full h-full object-cover"
              />
              <div v-else class="text-gray-400">
                <svg class="w-16 h-16" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    fill-rule="evenodd"
                    d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- Product Info -->
          <div class="p-4">
            <div class="flex items-center justify-between mb-2">
              <span class="inline-block bg-gray-100 text-gray-800 text-sm px-3 py-1 rounded-full uppercase tracking-wide">
                {{ product.category }}
              </span>
              <span
                :class="product.stock > 0 ? 'text-green-600' : 'text-red-600'"
                class="text-sm font-medium"
              >
                {{ product.stock > 0 ? `${product.stock} in stock` : 'Out of stock' }}
              </span>
            </div>

            <router-link :to="`/products/${product.id}`">
              <h3 class="font-semibold text-lg mb-2 hover:text-blue-600 transition-colors">{{ product.name }}</h3>
            </router-link>
            
            <p class="text-gray-600 text-sm mb-3 line-clamp-2">{{ product.description }}</p>

            <!-- Price -->
            <div class="flex items-center justify-between mb-4">
              <span class="text-2xl font-bold text-blue-600">${{ product.price.toFixed(2) }}</span>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2">
              <router-link
                :to="`/products/${product.id}`"
                class="flex-1 bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition-colors text-center"
              >
                View Details
              </router-link>
              <button
                v-if="product.stock > 0"
                @click="addToCart(product)"
                class="bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 transition-colors"
              >
                Add to Cart
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- No Products -->
      <div v-else class="text-center py-16">
        <svg class="mx-auto h-24 w-24 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No products found</h3>
        <p class="text-gray-600 mb-4">
          {{ hasActiveFilters ? 'Try adjusting your search or filters' : 'No products available at the moment' }}
        </p>
        <button
          v-if="hasActiveFilters"
          @click="products_store.clearFilters()"
          class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Clear Filters
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useProductStore } from '@/stores/products_store'
import { useCartStore } from '@/stores/cart_store'
import type { Product } from '@/types'

const products_store = useProductStore()
const cart_store = useCartStore()

// Computed
const hasActiveFilters = computed(() => {
  return products_store.searchQuery || products_store.selectedCategory
})

const filteredProductsCount = computed(() => {
  return products_store.filteredProducts.length
})

// Methods
const formatCategory = (category: string) => {
  return category.charAt(0).toUpperCase() + category.slice(1)
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
  
  // Show success feedback (you could use a toast notification here)
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