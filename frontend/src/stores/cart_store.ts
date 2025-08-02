import { defineStore } from 'pinia'
import type { CartItem, Product } from '@/types'
import { productService } from '@/services/products'

interface CartItemWithProduct extends CartItem {
  product?: Product
}

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [] as CartItemWithProduct[],
    loading: false,
    error: null as string | null,
  }),

  getters: {
    // Total number of items in cart
    totalItems: (state) => {
      return state.items.reduce((total, item) => total + item.quantity, 0)
    },

    // Total price of all items
    totalPrice: (state) => {
      return state.items.reduce((total, item) => total + item.price * item.quantity, 0)
    },

    // Check if cart has items
    hasItems: (state) => state.items.length > 0,

    // Get cart items with product details
    itemsWithProducts: (state) => {
      return state.items.filter((item) => item.product)
    },
  },

  actions: {
    // Add item to cart
    addItem(cartItem: CartItem) {
      // Check if item already exists (same product and size)
      const existingItemIndex = this.items.findIndex(
        (item) => item.product_id === cartItem.product_id && item.size === cartItem.size,
      )

      if (existingItemIndex !== -1) {
        // Update quantity of existing item
        this.items[existingItemIndex].quantity += cartItem.quantity
      } else {
        // Add new item to cart
        this.items.push({ ...cartItem })
      }

      // Fetch product details for the new item
      this.fetchProductDetails(cartItem.product_id)

      // Save to localStorage
      this.saveToStorage()
    },

    // Remove item from cart
    removeItem(productId: number, size: string) {
      this.items = this.items.filter(
        (item) => !(item.product_id === productId && item.size === size),
      )
      this.saveToStorage()
    },

    // Update item quantity
    updateQuantity(productId: number, size: string, quantity: number) {
      const item = this.items.find((item) => item.product_id === productId && item.size === size)

      if (item) {
        if (quantity <= 0) {
          this.removeItem(productId, size)
        } else {
          item.quantity = quantity
          this.saveToStorage()
        }
      }
    },

    // Clear entire cart
    clearCart() {
      this.items = []
      this.saveToStorage()
    },

    // Fetch product details for cart items
    async fetchProductDetails(productId: number) {
      try {
        const product = await productService.getProduct(productId)
        const item = this.items.find((item) => item.product_id === productId)
        if (item) {
          item.product = product
        }
      } catch (error) {
        console.error('Error fetching product details:', error)
      }
    },

    // Load all product details for cart items
    async loadAllProductDetails() {
      this.loading = true
      this.error = null

      try {
        const productPromises = this.items
          .filter((item) => !item.product)
          .map((item) => this.fetchProductDetails(item.product_id))

        await Promise.all(productPromises)
      } catch (error: any) {
        this.error = 'Failed to load product details'
        console.error('Error loading product details:', error)
      } finally {
        this.loading = false
      }
    },

    // Save cart to localStorage
    saveToStorage() {
      try {
        localStorage.setItem('bjj_store_cart', JSON.stringify(this.items))
      } catch (error) {
        console.error('Error saving cart to localStorage:', error)
      }
    },

    // Load cart from localStorage
    loadFromStorage() {
      try {
        const savedCart = localStorage.getItem('bjj_store_cart')
        if (savedCart) {
          this.items = JSON.parse(savedCart)
          // Load product details for all items
          this.loadAllProductDetails()
        }
      } catch (error) {
        console.error('Error loading cart from localStorage:', error)
        this.items = []
      }
    },

    // Get cart data for checkout
    getCheckoutData() {
      return this.items.map((item) => ({
        product_id: item.product_id,
        quantity: item.quantity,
        size: item.size,
        price: item.price,
      }))
    },
  },
})
