import { create } from 'zustand'
import { CartItem, Product } from '@/types'

export interface CartSummary {
  subtotal: number
  tax: number
  discount: number
  total: number
  itemCount: number
}

interface CartState {
  items: CartItem[]
  summary: CartSummary
  discountPercent: number
  notes: string
  
  // Actions
  addItem: (product: Product, quantity?: number) => void
  removeItem: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clearCart: () => void
  applyDiscount: (percent: number) => void
  setNotes: (notes: string) => void
  
  // Computed properties
  getItemCount: () => number
  getItem: (productId: string) => CartItem | undefined
  isEmpty: () => boolean
}

// Helper function to calculate summary
const calculateSummary = (items: CartItem[], discountPercent: number = 0): CartSummary => {
  const subtotal = items.reduce((sum, item) => sum + (item.product.price * item.quantity), 0)
  const tax = subtotal * 0.1 // 10% tax
  const discount = (subtotal * discountPercent) / 100
  const total = subtotal + tax - discount
  const itemCount = items.reduce((count, item) => count + item.quantity, 0)
  
  return {
    subtotal,
    tax,
    discount,
    total,
    itemCount
  }
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  summary: {
    subtotal: 0,
    tax: 0,
    discount: 0,
    total: 0,
    itemCount: 0
  },
  discountPercent: 0,
  notes: '',

  addItem: (product: Product, quantity = 1) => {
    const { discountPercent } = get()
    set((state) => {
      const existingItem = state.items.find(item => item.product.id === product.id)
      
      let newItems: CartItem[]
      
      if (existingItem) {
        // Update quantity if item already exists
        newItems = state.items.map(item =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        )
      } else {
        // Add new item
        newItems = [...state.items, { product, quantity }]
      }
      
      return { 
        items: newItems, 
        summary: calculateSummary(newItems, discountPercent)
      }
    })
  },

  removeItem: (productId: string) => {
    const { discountPercent } = get()
    set((state) => {
      const newItems = state.items.filter(item => item.product.id !== productId)
      return { 
        items: newItems, 
        summary: calculateSummary(newItems, discountPercent)
      }
    })
  },

  updateQuantity: (productId: string, quantity: number) => {
    if (quantity <= 0) {
      get().removeItem(productId)
      return
    }

    const { discountPercent } = get()
    set((state) => {
      const newItems = state.items.map(item =>
        item.product.id === productId
          ? { ...item, quantity }
          : item
      )
      
      return { 
        items: newItems, 
        summary: calculateSummary(newItems, discountPercent)
      }
    })
  },

  clearCart: () => {
    set({ 
      items: [], 
      summary: {
        subtotal: 0,
        tax: 0,
        discount: 0,
        total: 0,
        itemCount: 0
      },
      discountPercent: 0,
      notes: ''
    })
  },

  applyDiscount: (percent: number) => {
    const { items } = get()
    set({
      discountPercent: percent,
      summary: calculateSummary(items, percent)
    })
  },

  setNotes: (notes: string) => {
    set({ notes })
  },

  getItemCount: () => {
    const state = get()
    return state.items.reduce((count, item) => count + item.quantity, 0)
  },

  getItem: (productId: string) => {
    const state = get()
    return state.items.find(item => item.product.id === productId)
  },

  isEmpty: () => {
    const state = get()
    return state.items.length === 0
  },
}))