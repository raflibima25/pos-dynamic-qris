import { create } from 'zustand'
import { Payment } from '@/types'
import { api } from '@/lib/api'

interface PaymentState {
  currentPayment: Payment | null
  loading: boolean
  error: string | null
  
  // Actions
  generateQRIS: (transactionId: string, amount: number) => Promise<Payment | null>
  getPaymentStatus: (transactionId: string) => Promise<Payment | null>
  refreshQRIS: (transactionId: string) => Promise<Payment | null>
  clearError: () => void
  setCurrentPayment: (payment: Payment | null) => void
}

export const usePaymentStore = create<PaymentState>((set) => ({
  currentPayment: null,
  loading: false,
  error: null,

  generateQRIS: async (transactionId: string, amount: number) => {
    set({ loading: true, error: null })

    try {
      const response = await api.post('/qris/generate', {
        transaction_id: transactionId,
        amount: amount,
        expiry_minutes: 10 // Default 10 minutes
      })

      // With fetch API, response.data contains the payment object
      const payment = response.data
      set({
        currentPayment: payment,
        loading: false
      })

      return payment
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to generate QRIS'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  getPaymentStatus: async (transactionId: string) => {
    set({ loading: true, error: null })

    try {
      const response = await api.get(`/payments/${transactionId}/status`)
      // With fetch API, response.data contains the payment status
      const paymentStatus = response.data

      set({
        currentPayment: paymentStatus,
        loading: false
      })

      return paymentStatus
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to get payment status'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  refreshQRIS: async (transactionId: string) => {
    set({ loading: true, error: null })

    try {
      const response = await api.post(`/qris/${transactionId}/refresh`)

      // With fetch API, response.data contains the payment object
      const payment = response.data
      set({
        currentPayment: payment,
        loading: false
      })

      return payment
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to refresh QRIS'
      set({ error: errorMessage, loading: false })
      return null
    }
  },

  clearError: () => {
    set({ error: null })
  },

  setCurrentPayment: (payment: Payment | null) => {
    set({ currentPayment: payment })
  }
}))