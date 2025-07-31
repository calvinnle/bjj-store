import api from './api'
import type { Order, CartItem, Address, OrderResponse } from '@/types'

interface CreateOrderRequest {
  guest_email: string
  shipping_address: Address
  items: CartItem[]
}

export const orderService = {
  // Create new order
  async createOrder(orderData: CreateOrderRequest): Promise<OrderResponse> {
    const response = await api.post('/api/orders', orderData)
    return response.data
  },

  // Track order by order number
  async trackOrder(orderNumber: string): Promise<Order> {
    const response = await api.get(`/api/orders/track/${orderNumber}`)
    return response.data.order
  },

  // Get orders by email
  async getOrdersByEmail(email: string): Promise<Order[]> {
    const response = await api.get(`/api/orders/email/${email}`)
    return response.data.orders
  },

  // Admin: Get all orders (requires JWT)
  async getAllOrders(): Promise<Order[]> {
    const response = await api.get('/api/admin/orders')
    return response.data.orders
  },

  // Admin: Update order status (requires JWT)
  async updateOrderStatus(orderId: number, status: string): Promise<Order> {
    const response = await api.put(`/api/admin/orders/${orderId}/status`, { status })
    return response.data.order
  },
}