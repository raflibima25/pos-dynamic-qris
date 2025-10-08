import { create } from 'zustand'
import { Transaction, TransactionStatus, CartItem } from '@/types'
import { api } from '@/lib/api'

interface TransactionState {
  currentTransaction: Transaction | null
  transactions: Transaction[]
  loading: boolean
  error: string | null
  
  // Actions
  createTransaction: (items: CartItem[], notes?: string) => Promise<Transaction | null>
  getTransaction: (id: string) => Promise<Transaction | null>
  listTransactions: (filters?: {
    status?: TransactionStatus
    dateFrom?: string
    dateTo?: string
    limit?: number
    offset?: number
  }) => Promise<void>
  addItemToTransaction: (transactionId: string, productId: string, quantity: number) => Promise<Transaction | null>
  removeItemFromTransaction: (transactionId: string, productId: string) => Promise<Transaction | null>
  updateItemQuantity: (transactionId: string, productId: string, quantity: number) => Promise<Transaction | null>
  cancelTransaction: (transactionId: string) => Promise<void>
  clearError: () => void
  setCurrentTransaction: (transaction: Transaction | null) => void
}

export const useTransactionStore = create<TransactionState>((set, get) => ({
  currentTransaction: null,
  transactions: [],
  loading: false,
  error: null,

  createTransaction: async (items: CartItem[], notes?: string) => {
    set({ loading: true, error: null })

    try {
      const transactionItems = items.map(item => ({
        product_id: item.product.id,
        quantity: item.quantity
      }))

      console.log('Sending transaction request:', { items: transactionItems, notes })

      const response = await api.post('/transactions', {
        items: transactionItems,
        notes: notes || ''
      })

      console.log('Full API response:', response)

      // With fetch API (not axios), response is already parsed JSON
      // Backend returns: { success: true, message: "...", data: {...} }
      // So response.data contains the transaction object
      const transaction = response.data

      console.log('Extracted transaction:', transaction)

      if (!transaction || !transaction.id) {
        console.error('Invalid transaction structure:', transaction)
        throw new Error('Invalid transaction response from server')
      }

      set({
        currentTransaction: transaction,
        transactions: [transaction, ...get().transactions],
        loading: false
      })

      console.log('Transaction stored successfully:', transaction.id)
      return transaction
    } catch (error: any) {
      console.error('Transaction creation error details:', error)
      console.error('Error response:', error.response)
      const errorMessage = error.message || error.response?.data?.message || 'Failed to create transaction'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  getTransaction: async (id: string) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.get(`/transactions/${id}`)
      const transaction = response.data.data
      
      set({ 
        currentTransaction: transaction,
        loading: false 
      })
      
      return transaction
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to get transaction'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  listTransactions: async (filters = {}) => {
    set({ loading: true, error: null })
    
    try {
      const params = new URLSearchParams()
      
      if (filters.status) params.append('status', filters.status)
      if (filters.dateFrom) params.append('date_from', filters.dateFrom)
      if (filters.dateTo) params.append('date_to', filters.dateTo)
      if (filters.limit) params.append('limit', filters.limit.toString())
      if (filters.offset) params.append('offset', filters.offset.toString())

      const response = await api.get(`/transactions?${params.toString()}`)
      const transactions = response.data.data
      
      set({ 
        transactions,
        loading: false 
      })
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to list transactions'
      set({ error: errorMessage, loading: false })
    }
  },

  addItemToTransaction: async (transactionId: string, productId: string, quantity: number) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.post(`/transactions/${transactionId}/items`, {
        product_id: productId,
        quantity
      })

      const updatedTransaction = response.data.data
      
      set(state => ({
        currentTransaction: state.currentTransaction?.id === transactionId ? updatedTransaction : state.currentTransaction,
        transactions: state.transactions.map(t => 
          t.id === transactionId ? updatedTransaction : t
        ),
        loading: false
      }))
      
      return updatedTransaction
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to add item to transaction'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  removeItemFromTransaction: async (transactionId: string, productId: string) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.delete(`/transactions/${transactionId}/items/${productId}`)
      const updatedTransaction = response.data.data
      
      set(state => ({
        currentTransaction: state.currentTransaction?.id === transactionId ? updatedTransaction : state.currentTransaction,
        transactions: state.transactions.map(t => 
          t.id === transactionId ? updatedTransaction : t
        ),
        loading: false
      }))
      
      return updatedTransaction
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to remove item from transaction'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  updateItemQuantity: async (transactionId: string, productId: string, quantity: number) => {
    set({ loading: true, error: null })
    
    try {
      const response = await api.put(`/transactions/${transactionId}/items/${productId}`, {
        quantity
      })

      const updatedTransaction = response.data.data
      
      set(state => ({
        currentTransaction: state.currentTransaction?.id === transactionId ? updatedTransaction : state.currentTransaction,
        transactions: state.transactions.map(t => 
          t.id === transactionId ? updatedTransaction : t
        ),
        loading: false
      }))
      
      return updatedTransaction
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to update item quantity'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  cancelTransaction: async (transactionId: string) => {
    set({ loading: true, error: null })
    
    try {
      await api.put(`/transactions/${transactionId}/cancel`)
      
      set(state => ({
        currentTransaction: state.currentTransaction?.id === transactionId ? null : state.currentTransaction,
        transactions: state.transactions.map(t => 
          t.id === transactionId ? { ...t, status: 'cancelled' as TransactionStatus } : t
        ),
        loading: false
      }))
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to cancel transaction'
      set({ error: errorMessage, loading: false })
    }
  },

  clearError: () => {
    set({ error: null })
  },

  setCurrentTransaction: (transaction: Transaction | null) => {
    set({ currentTransaction: transaction })
  }
}))