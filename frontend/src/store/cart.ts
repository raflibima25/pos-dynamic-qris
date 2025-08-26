import { create } from 'zustand'
import { CartItem, Product } from '@/types'

interface CartState {
  items: CartItem[]
  total: number
  addItem: (product: Product, quantity?: number) => void
  removeItem: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clearCart: () => void
  getItemCount: () => number
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  total: 0,

  addItem: (product: Product, quantity = 1) => {
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
      
      const total = newItems.reduce(
        (sum, item) => sum + (item.product.price * item.quantity), 
        0
      )
      
      return { items: newItems, total }
    })
  },

  removeItem: (productId: string) => {
    set((state) => {
      const newItems = state.items.filter(item => item.product.id !== productId)
      const total = newItems.reduce(
        (sum, item) => sum + (item.product.price * item.quantity), 
        0
      )
      
      return { items: newItems, total }
    })
  },

  updateQuantity: (productId: string, quantity: number) => {
    if (quantity <= 0) {
      get().removeItem(productId)
      return
    }

    set((state) => {
      const newItems = state.items.map(item =>
        item.product.id === productId
          ? { ...item, quantity }
          : item
      )
      
      const total = newItems.reduce(
        (sum, item) => sum + (item.product.price * item.quantity), 
        0
      )
      
      return { items: newItems, total }
    })
  },

  clearCart: () => {
    set({ items: [], total: 0 })
  },

  getItemCount: () => {
    const state = get()
    return state.items.reduce((count, item) => count + item.quantity, 0)
  },
}))