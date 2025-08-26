export interface User {
  id: string
  name: string
  email: string
  role: 'admin' | 'cashier'
  is_active: boolean
}

export interface LoginResponse {
  user: User
  token: string
}

export interface Category {
  id: string
  name: string
  is_active: boolean
}

export interface Product {
  id: string
  name: string
  description: string
  price: number
  stock: number
  category_id: string
  sku: string
  is_active: boolean
  created_at: string
  updated_at: string
  category?: Category
}

export interface CartItem {
  product: Product
  quantity: number
}

export interface Transaction {
  id: string
  user_id: string
  total_amount: number
  tax_amount: number
  discount: number
  status: 'pending' | 'paid' | 'cancelled' | 'expired'
  notes?: string
  created_at: string
  updated_at: string
  items: TransactionItem[]
}

export interface TransactionItem {
  id: string
  transaction_id: string
  product_id: string
  quantity: number
  unit_price: number
  total_price: number
  product?: Product
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: any
}