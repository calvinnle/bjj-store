// Product type (matches your Go Product model)
export interface Product {
  id: number
  name: string
  description: string
  price: number
  category: string
  size_options: string[] | string
  stock: number
  image_url: string
  created_at: string
  updated_at: string
}

// Cart item type
export interface CartItem {
  product_id: number
  quantity: number
  size: string
  price: number
}

// Address type (matches your Go Address model)
export interface Address {
  first_name: string
  last_name: string
  address1: string
  address2?: string
  city: string
  state: string
  zip_code: string
  country: string
}

// Order type (matches your Go Order model)
export interface Order {
  id: number
  order_number: string
  guest_email: string
  shipping_address: Address
  items: OrderItem[]
  total_amount: number
  status: 'pending' | 'paid' | 'shipped' | 'delivered' | 'cancelled'
  stripe_payment_id?: string
  created_at: string
  updated_at: string
}

// Order item type
export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product: Product
  quantity: number
  price: number
  size: string
}

// Admin user type (matches your Go AdminUser model)
export interface AdminUser {
  id: number
  email: string
  role: 'super_admin' | 'inventory' | 'order_manager' | 'viewer'
  is_active: boolean
  last_login?: string
}

// API response types
export interface LoginResponse {
  success: boolean
  token: string
  admin: AdminUser
  message: string
}

export interface OrderResponse {
  order: Order
  message: string
}

export interface APIError {
  error: string
  details?: string
}
