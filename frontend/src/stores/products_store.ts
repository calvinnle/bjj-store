import { defineStore } from 'pinia'
import { productService } from '@/services/products'
import type { Product } from '@/types'

export const useProductStore = defineStore('products', {
  state: () => ({
    // Product data
    products: [] as Product[],
    currentProduct: null as Product | null,

    // UI states
    loading: false,
    error: null as string | null,

    // Filters
    searchQuery: '',
    selectedCategory: '',

    // Available categories (will be populated from products)
    categories: [] as string[],
  }),

  getters: {
    // Get products filtered by search and category
    filteredProducts: (state) => {
      let filtered = state.products

      if (state.selectedCategory) {
        filtered = filtered.filter((product) => product.category === state.selectedCategory)
      }

      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase()
        filtered = filtered.filter(
          (product) =>
            product.name.toLowerCase().includes(query) ||
            product.description.toLowerCase().includes(query),
        )
      }

      return filtered
    },

    // Get unique categories from products
    availableCategories: (state) => {
      const cats = [...new Set(state.products.map((p) => p.category))]
      return cats.filter(Boolean) // Remove empty categories
    },

    // Check if products are loaded
    hasProducts: (state) => state.products.length > 0,
  },

  actions: {
    // Fetch all products
    async fetchProducts() {
      this.loading = true
      this.error = null

      try {
        this.products = await productService.getProducts()
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to load products'
        console.error('Error fetching products:', error)
      } finally {
        this.loading = false
      }
    },

    // Fetch single product
    async fetchProduct(id: number) {
      this.loading = true
      this.error = null

      try {
        this.currentProduct = await productService.getProduct(id)
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Product not found'
        console.error('Error fetching product:', error)
      } finally {
        this.loading = false
      }
    },

    // Search products
    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    // Filter by category
    setCategory(category: string) {
      this.selectedCategory = category
    },

    // Clear filters
    clearFilters() {
      this.searchQuery = ''
      this.selectedCategory = ''
    },

    // Admin actions (create, update, delete)
    async createProduct(productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) {
      try {
        const newProduct = await productService.createProduct(productData)
        this.products.push(newProduct)
        return newProduct
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to create product'
        throw error
      }
    },

    async updateProduct(id: number, productData: Partial<Product>) {
      try {
        const updatedProduct = await productService.updateProduct(id, productData)
        const index = this.products.findIndex((p) => p.id === id)
        if (index !== -1) {
          this.products[index] = updatedProduct
        }
        return updatedProduct
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to update product'
        throw error
      }
    },

    async deleteProduct(id: number) {
      try {
        await productService.deleteProduct(id)
        this.products = this.products.filter((p) => p.id !== id)
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to delete product'
        throw error
      }
    },
  },
})
