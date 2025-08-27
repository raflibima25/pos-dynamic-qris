'use client'

import { useState } from 'react'
import { CartItem } from '@/types'
import { CartSummary } from '@/store/cart'
import { useCartStore } from '@/store/cart'
import { useTransactionStore } from '@/store/transaction'
import { XMarkIcon } from '@heroicons/react/24/outline'

interface CheckoutModalProps {
  isOpen: boolean
  onClose: () => void
  items: CartItem[]
  summary: CartSummary
  notes: string
}

export function CheckoutModal({ isOpen, onClose, items, summary, notes }: CheckoutModalProps) {
  const [step, setStep] = useState<'confirm' | 'processing' | 'success' | 'error'>('confirm')
  const [transactionId, setTransactionId] = useState<string | null>(null)
  const { clearCart } = useCartStore()
  const { createTransaction, loading, error } = useTransactionStore()

  if (!isOpen) return null

  const handleConfirmOrder = async () => {
    setStep('processing')
    
    try {
      const transaction = await createTransaction(items, notes)
      
      if (transaction) {
        setTransactionId(transaction.id)
        setStep('success')
        // Clear cart after successful transaction
        clearCart()
      } else {
        setStep('error')
      }
    } catch (err) {
      setStep('error')
    }
  }

  const handleClose = () => {
    if (step !== 'processing') {
      onClose()
      // Reset step when closing
      setStep('confirm')
    }
  }

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        {/* Backdrop */}
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={handleClose} />
        
        {/* Modal */}
        <div className="relative w-full max-w-md bg-white rounded-lg shadow-xl">
          {step === 'confirm' && (
            <ConfirmStep
              items={items}
              summary={summary}
              notes={notes}
              onConfirm={handleConfirmOrder}
              onClose={handleClose}
            />
          )}
          
          {step === 'processing' && (
            <ProcessingStep />
          )}
          
          {step === 'success' && (
            <SuccessStep
              transactionId={transactionId}
              total={summary.total}
              onClose={handleClose}
            />
          )}
          
          {step === 'error' && (
            <ErrorStep
              error={error}
              onRetry={handleConfirmOrder}
              onClose={handleClose}
            />
          )}
        </div>
      </div>
    </div>
  )
}

interface ConfirmStepProps {
  items: CartItem[]
  summary: CartSummary
  notes: string
  onConfirm: () => void
  onClose: () => void
}

function ConfirmStep({ items, summary, notes, onConfirm, onClose }: ConfirmStepProps) {
  return (
    <>
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b">
        <h3 className="text-lg font-semibold text-gray-900">Confirm Order</h3>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600"
        >
          <XMarkIcon className="w-6 h-6" />
        </button>
      </div>

      {/* Content */}
      <div className="p-6 space-y-4">
        {/* Items Review */}
        <div>
          <h4 className="font-medium text-gray-900 mb-3">Order Items</h4>
          <div className="space-y-2 max-h-40 overflow-y-auto">
            {items.map(item => (
              <div key={item.product.id} className="flex justify-between text-sm">
                <span className="text-gray-900">
                  {item.quantity}x {item.product.name}
                </span>
                <span className="text-gray-600">
                  ${(item.quantity * item.product.price).toFixed(2)}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Notes */}
        {notes && (
          <div>
            <h4 className="font-medium text-gray-900 mb-2">Notes</h4>
            <p className="text-sm text-gray-600 bg-gray-50 p-3 rounded">
              {notes}
            </p>
          </div>
        )}

        {/* Total */}
        <div className="border-t pt-4">
          <div className="flex justify-between text-lg font-semibold">
            <span>Total</span>
            <span>${summary.total.toFixed(2)}</span>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="flex gap-3 p-6 border-t">
        <button
          onClick={onClose}
          className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Cancel
        </button>
        <button
          onClick={onConfirm}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          Create Transaction
        </button>
      </div>
    </>
  )
}

function ProcessingStep() {
  return (
    <div className="p-8 text-center">
      <div className="w-16 h-16 mx-auto mb-4 bg-blue-100 rounded-full flex items-center justify-center">
        <svg className="w-8 h-8 text-blue-500 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      <h3 className="text-lg font-semibold text-gray-900 mb-2">Processing Transaction</h3>
      <p className="text-gray-600">Please wait while we create your transaction...</p>
    </div>
  )
}

interface SuccessStepProps {
  transactionId: string | null
  total: number
  onClose: () => void
}

function SuccessStep({ transactionId, total, onClose }: SuccessStepProps) {
  return (
    <>
      <div className="p-8 text-center">
        <div className="w-16 h-16 mx-auto mb-4 bg-green-100 rounded-full flex items-center justify-center">
          <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">Transaction Created!</h3>
        <p className="text-gray-600 mb-4">
          Transaction #{transactionId} has been created successfully.
        </p>
        <p className="text-2xl font-bold text-green-600 mb-6">
          ${total.toFixed(2)}
        </p>
        <p className="text-sm text-gray-500 mb-6">
          You can now proceed to payment or generate a QRIS code.
        </p>
      </div>

      <div className="flex gap-3 p-6 border-t">
        <button
          onClick={onClose}
          className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Start New Order
        </button>
        <button
          onClick={() => {
            // TODO: Navigate to payment/QRIS generation
            onClose()
          }}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          Generate QRIS
        </button>
      </div>
    </>
  )
}

interface ErrorStepProps {
  error: string | null
  onRetry: () => void
  onClose: () => void
}

function ErrorStep({ error, onRetry, onClose }: ErrorStepProps) {
  return (
    <>
      <div className="p-8 text-center">
        <div className="w-16 h-16 mx-auto mb-4 bg-red-100 rounded-full flex items-center justify-center">
          <svg className="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">Transaction Failed</h3>
        <p className="text-gray-600 mb-4">
          {error || 'An error occurred while creating the transaction.'}
        </p>
      </div>

      <div className="flex gap-3 p-6 border-t">
        <button
          onClick={onClose}
          className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Cancel
        </button>
        <button
          onClick={onRetry}
          className="flex-1 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
        >
          Try Again
        </button>
      </div>
    </>
  )
}