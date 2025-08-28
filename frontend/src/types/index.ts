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
  description?: string
  price: number
  stock: number
  category_id: string
  sku?: string
  image_url?: string
  is_active: boolean
  created_at: string
  updated_at: string
  category?: Category
}

export interface CartItem {
  product: Product
  quantity: number
}

export type TransactionStatus = 'pending' | 'paid' | 'cancelled' | 'expired'

export interface Transaction {
  id: string
  userId: string
  totalAmount: number
  taxAmount: number
  discount: number
  status: TransactionStatus
  notes?: string
  createdAt: string
  updatedAt: string
  items: TransactionItem[]
  user?: User
}

export interface TransactionItem {
  id: string
  transactionId: string
  productId: string
  quantity: number
  unitPrice: number
  totalPrice: number
  product?: Product
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: any
}