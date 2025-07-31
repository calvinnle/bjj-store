import api from './api'
import type { Product } from '@/types'

export const productService = {
  // Customer API calls (public endpoints)

  // GET /api/products - Browse all products
  async getProducts(category?: string, search?: string): Promise<Product[]> {
    const params = new URLSearchParams()
    if (category) params.append('category', category)
    if (search) params.append('search', search)

    const response = await api.get(`/api/products?${params}`)
    return response.data
  },

  // GET /api/products/:id - Get single product
  async getProduct(id: number): Promise<Product> {
    const response = await api.get(`/api/products/${id}`)
    return response.data
  },

  // Admin API calls (JWT protected endpoints)

  // POST /api/admin/products - Create new product
  async createProduct(
    product: Omit<Product, 'id' | 'created_at' | 'updated_at'>,
  ): Promise<Product> {
    const response = await api.post('/api/admin/products', product)
    return response.data
  },

  // PUT /api/admin/products/:id - Update product
  async updateProduct(id: number, product: Partial<Product>): Promise<Product> {
    const response = await api.put(`/api/admin/products/${id}`, product)
    return response.data
  },

  // DELETE /api/admin/products/:id - Delete product
  async deleteProduct(id: number): Promise<void> {
    await api.delete(`/api/admin/products/${id}`)
  },
}
