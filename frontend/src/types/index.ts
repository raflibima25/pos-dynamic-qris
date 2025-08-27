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
  categoryId: string
  sku?: string
  image?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
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