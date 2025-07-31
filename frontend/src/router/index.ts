import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth_store'
import HomeView from '../views/HomeView.vue'
import ProductsView from '../views/ProductsView.vue'
import ProductView from '../views/ProductView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import OrderConfirmationView from '../views/OrderConfirmationView.vue'
import TrackOrderView from '../views/TrackOrderView.vue'
import AdminLoginView from '../views/AdminLoginView.vue'
import AdminDashboardView from '../views/AdminDashboardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/products',
      name: 'products',
      component: ProductsView,
    },
    {
      path: '/products/:id',
      name: 'product',
      component: ProductView,
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: CheckoutView,
    },
    {
      path: '/order-confirmation/:orderNumber',
      name: 'order-confirmation',
      component: OrderConfirmationView,
    },
    {
      path: '/track-order',
      name: 'track-order',
      component: TrackOrderView,
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: AdminLoginView,
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboardView,
      meta: { requiresAuth: true },
    },
    // Redirect /admin/* to /admin
    {
      path: '/admin/:pathMatch(.*)*',
      redirect: '/admin',
    },
  ],
})

// Auth guard
router.beforeEach(async (to) => {
  if (to.meta.requiresAuth) {
    const authStore = useAuthStore()
    
    // If we have a token but no user loaded, try to load it
    const hasToken = localStorage.getItem('admin_token')
    if (hasToken && !authStore.user && !authStore.loading) {
      try {
        await authStore.initializeAuth()
      } catch (error) {
        // If token is invalid, remove it and redirect to login
        localStorage.removeItem('admin_token')
        return {
          name: 'admin-login',
          query: { redirect: to.fullPath },
        }
      }
    }
    
    if (!authStore.isAuthenticated) {
      return {
        name: 'admin-login',
        query: { redirect: to.fullPath },
      }
    }
  }
})

export default router
