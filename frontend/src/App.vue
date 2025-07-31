<template>
  <div id="app">
    <!-- Navigation (only show on non-admin routes) -->
    <NavBar v-if="!isAdminRoute" />
    
    <!-- Main Content -->
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { useAuthStore } from '@/stores/auth_store'
import { useCartStore } from '@/stores/cart_store'

const route = useRoute()
const auth_store = useAuthStore()
const cart_store = useCartStore()

// Computed
const isAdminRoute = computed(() => {
  return route.path.startsWith('/admin')
})

// Initialize stores on app mount
onMounted(() => {
  // Initialize auth state
  auth_store.initializeAuth()
  
  // Load cart from localStorage
  cart_store.loadFromStorage()
})
</script>

<style>
#app {
  min-height: 100vh;
  background-color: #f9fafb;
}
</style>
