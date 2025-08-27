import { create } from 'zustand'
import { Product, Category } from '@/types'
import { api } from '@/lib/api'

interface ProductState {
  products: Product[]
  categories: Category[]
  loading: boolean
  error: string | null
  
  // Actions
  listProducts: (categoryId?: string) => Promise<void>
  listCategories: () => Promise<void>
  getProduct: (id: string) => Promise<Product | null>
  createProduct: (productData: Partial<Product>) => Promise<Product | null>
  updateProduct: (id: string, productData: Partial<Product>) => Promise<Product | null>
  deleteProduct: (id: string) => Promise<void>
  updateStock: (id: string, stock: number) => Promise<Product | null>
  clearError: () => void
}

export const useProductStore = create<ProductState>((set, get) => ({
  products: [],
  categories: [],
  loading: false,
  error: null,

  listProducts: async (categoryId?: string) => {
    set({ loading: true, error: null })
    
    try {
      const params = categoryId ? `?category_id=${categoryId}` : ''
      const response = await api.get(`/products${params}`)
      const products = response.data.data || []
      
      set({ products, loading: false })
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to load products'
      set({ error: errorMessage, loading: false })
    }
  },

  listCategories: async () => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.get('/categories')
      const categories = response.data.data || []
      
      set({ categories, loading: false })
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to load categories'
      set({ error: errorMessage, loading: false })
    }
  },

  getProduct: async (id: string) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.get(`/products/${id}`)
      const product = response.data.data
      
      set({ loading: false })
      return product
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to get product'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  createProduct: async (productData: Partial<Product>) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.post('/products', productData)
      const newProduct = response.data.data
      
      set(state => ({
        products: [...state.products, newProduct],
        loading: false
      }))
      
      return newProduct
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to create product'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  updateProduct: async (id: string, productData: Partial<Product>) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.put(`/products/${id}`, productData)
      const updatedProduct = response.data.data
      
      set(state => ({
        products: state.products.map(p => p.id === id ? updatedProduct : p),
        loading: false
      }))
      
      return updatedProduct
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to update product'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  deleteProduct: async (id: string) => {
    set({ loading: true, error: null })
    
    try {
      await api.delete(`/products/${id}`)
      
      set(state => ({
        products: state.products.filter(p => p.id !== id),
        loading: false
      }))
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to delete product'
      set({ error: errorMessage, loading: false })
    }
  },

  updateStock: async (id: string, stock: number) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.patch(`/products/${id}/stock`, { stock })
      const updatedProduct = response.data.data
      
      set(state => ({
        products: state.products.map(p => p.id === id ? updatedProduct : p),
        loading: false
      }))
      
      return updatedProduct
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to update stock'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  clearError: () => {
    set({ error: null })
  }
}))