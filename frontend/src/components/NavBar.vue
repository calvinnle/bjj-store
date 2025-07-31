<template>
  <nav class="bg-white shadow-lg sticky top-0 z-50 border-b border-gray-200">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <router-link to="/" class="flex items-center space-x-2">
          <div class="bg-blue-600 text-white p-2 rounded-lg">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
              <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <span class="text-xl font-bold text-gray-900">BJJ Store</span>
        </router-link>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-8">
          <router-link 
            to="/" 
            class="text-gray-700 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 font-medium': $route.name === 'home' }"
          >
            Home
          </router-link>
          
          <router-link 
            to="/products" 
            class="text-gray-700 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 font-medium': $route.name === 'products' }"
          >
            Products
          </router-link>
          
          <router-link 
            to="/track-order" 
            class="text-gray-700 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 font-medium': $route.name === 'track-order' }"
          >
            Track Order
          </router-link>

          <!-- Cart Button -->
          <button 
            @click="toggleCart"
            class="relative p-2 text-gray-700 hover:text-blue-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                    d="M3 3h2l.4 2M7 13h10l4-8H5.4m1.6 8L3 3H1m6 10v6a2 2 0 002 2h4a2 2 0 002-2v-6m-6 0a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z"/>
            </svg>
            <!-- Cart Badge -->
            <span 
              v-if="cart_store.totalItems > 0"
              class="absolute -top-1 -right-1 bg-blue-600 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
            >
              {{ cart_store.totalItems }}
            </span>
          </button>

          <!-- Admin Link -->
          <router-link 
            to="/admin" 
            class="bg-gray-100 text-gray-700 px-3 py-2 rounded-lg hover:bg-gray-200 transition-colors text-sm"
          >
            Admin
          </router-link>
        </div>

        <!-- Mobile Menu Button -->
        <button 
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="md:hidden p-2 text-gray-700 hover:text-blue-600"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Mobile Menu -->
      <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200">
        <div class="py-4 space-y-2">
          <router-link 
            to="/" 
            @click="mobileMenuOpen = false"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 bg-gray-50 font-medium': $route.name === 'home' }"
          >
            Home
          </router-link>
          
          <router-link 
            to="/products" 
            @click="mobileMenuOpen = false"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 bg-gray-50 font-medium': $route.name === 'products' }"
          >
            Products
          </router-link>
          
          <router-link 
            to="/track-order" 
            @click="mobileMenuOpen = false"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors"
            :class="{ 'text-blue-600 bg-gray-50 font-medium': $route.name === 'track-order' }"
          >
            Track Order
          </router-link>

          <button 
            @click="toggleCart(); mobileMenuOpen = false"
            class="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors flex items-center justify-between"
          >
            <span>Cart</span>
            <span 
              v-if="cart_store.totalItems > 0"
              class="bg-blue-600 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
            >
              {{ cart_store.totalItems }}
            </span>
          </button>

          <router-link 
            to="/admin" 
            @click="mobileMenuOpen = false"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-blue-600 transition-colors"
          >
            Admin
          </router-link>
        </div>
      </div>
    </div>

    <!-- Cart Sidebar -->
    <div v-if="cartOpen" class="fixed inset-0 z-50 overflow-hidden">
      <!-- Backdrop -->
      <div 
        @click="cartOpen = false" 
        class="absolute inset-0 bg-black bg-opacity-50 transition-opacity"
      ></div>
      
      <!-- Cart Panel -->
      <div class="absolute right-0 top-0 h-full w-full max-w-md bg-white shadow-xl">
        <div class="flex flex-col h-full">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b">
            <h2 class="text-lg font-semibold">Shopping Cart</h2>
            <button 
              @click="cartOpen = false"
              class="p-2 text-gray-400 hover:text-gray-600"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Cart Items -->
          <div class="flex-1 overflow-y-auto p-4">
            <div v-if="!cart_store.hasItems" class="text-center py-8">
              <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M3 3h2l.4 2M7 13h10l4-8H5.4m1.6 8L3 3H1m6 10v6a2 2 0 002 2h4a2 2 0 002-2v-6m-6 0a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z"/>
              </svg>
              <p class="text-gray-500">Your cart is empty</p>
            </div>

            <div v-else class="space-y-4">
              <div 
                v-for="item in cart_store.itemsWithProducts" 
                :key="`${item.product_id}-${item.size}`"
                class="flex items-center space-x-4 p-4 border rounded-lg"
              >
                <!-- Product Image -->
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

                <!-- Product Info -->
                <div class="flex-1 min-w-0">
                  <h4 class="text-sm font-medium text-gray-900 truncate">
                    {{ item.product?.name || 'Loading...' }}
                  </h4>
                  <p class="text-sm text-gray-500">Size: {{ item.size }}</p>
                  <p class="text-sm font-medium text-blue-600">${{ item.price.toFixed(2) }}</p>
                </div>

                <!-- Quantity Controls -->
                <div class="flex items-center space-x-2">
                  <button 
                    @click="cart_store.updateQuantity(item.product_id, item.size, item.quantity - 1)"
                    class="w-8 h-8 border border-gray-300 rounded flex items-center justify-center hover:bg-gray-50"
                  >
                    -
                  </button>
                  <span class="w-8 text-center text-sm">{{ item.quantity }}</span>
                  <button 
                    @click="cart_store.updateQuantity(item.product_id, item.size, item.quantity + 1)"
                    class="w-8 h-8 border border-gray-300 rounded flex items-center justify-center hover:bg-gray-50"
                  >
                    +
                  </button>
                </div>

                <!-- Remove Button -->
                <button 
                  @click="cart_store.removeItem(item.product_id, item.size)"
                  class="p-1 text-red-500 hover:text-red-700"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div v-if="cart_store.hasItems" class="border-t p-4 space-y-4">
            <div class="flex justify-between items-center text-lg font-semibold">
              <span>Total:</span>
              <span>${{ cart_store.totalPrice.toFixed(2) }}</span>
            </div>
            <router-link 
              to="/checkout" 
              @click="cartOpen = false"
              class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-semibold hover:bg-blue-700 transition-colors block text-center"
            >
              Checkout
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useCartStore } from '@/stores/cart_store'

const cart_store = useCartStore()

// Component state
const mobileMenuOpen = ref(false)
const cartOpen = ref(false)

// Methods
const toggleCart = () => {
  cartOpen.value = !cartOpen.value
}

// Load cart from storage on mount
onMounted(() => {
  cart_store.loadFromStorage()
})
</script>