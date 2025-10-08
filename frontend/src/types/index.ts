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
  user_id: string
  total_amount: number
  tax_amount: number
  discount: number
  status: TransactionStatus
  notes?: string
  created_at: string
  updated_at: string
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

export type PaymentStatus = 'pending' | 'success' | 'failed' | 'expired' | 'cancelled'
export type PaymentMethod = 'qris'

export interface Payment {
  id: string
  transaction_id: string
  amount: number
  method: PaymentMethod
  status: PaymentStatus
  external_id?: string
  paid_at?: string
  expires_at: string
  created_at: string
  updated_at: string
  qr_code?: QRISCode
}

export interface QRISCode {
  id: string
  transaction_id: string
  payment_id: string
  qr_code: string      // QRIS EMVCo string for QR generation
  url?: string         // Midtrans simulator URL for testing
  expires_at: string
  created_at: string
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: any
}