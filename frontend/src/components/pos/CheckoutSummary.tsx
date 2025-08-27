'use client'

import { useState } from 'react'
import { CartSummary } from '@/store/cart'
import { useCartStore } from '@/store/cart'
import { useTransactionStore } from '@/store/transaction'
import { CheckoutModal } from './CheckoutModal'

interface CheckoutSummaryProps {
  summary: CartSummary
}

export function CheckoutSummary({ summary }: CheckoutSummaryProps) {
  const [isCheckoutModalOpen, setIsCheckoutModalOpen] = useState(false)
  const { items, notes, applyDiscount, setNotes, discountPercent } = useCartStore()
  const { createTransaction, loading } = useTransactionStore()

  const handleApplyDiscount = (discount: number) => {
    applyDiscount(discount)
  }

  const handleProceedToCheckout = () => {
    setIsCheckoutModalOpen(true)
  }

  return (
    <>
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-900">Order Summary</h3>

        {/* Order Notes */}
        <div>
          <label htmlFor="notes" className="block text-sm font-medium text-gray-700 mb-1">
            Notes (Optional)
          </label>
          <textarea
            id="notes"
            rows={2}
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Add any special instructions..."
          />
        </div>

        {/* Discount */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Discount
          </label>
          <div className="flex gap-2">
            {[5, 10, 15, 20].map(discount => (
              <button
                key={discount}
                onClick={() => handleApplyDiscount(discount)}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                  discountPercent === discount
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {discount}%
              </button>
            ))}
            <button
              onClick={() => handleApplyDiscount(0)}
              className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                discountPercent === 0
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              None
            </button>
          </div>
        </div>

        {/* Price Breakdown */}
        <div className="space-y-2 pt-4 border-t border-gray-200">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Subtotal ({summary.itemCount} items)</span>
            <span className="text-gray-900">${summary.subtotal.toFixed(2)}</span>
          </div>
          
          {summary.discount > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Discount ({discountPercent}%)</span>
              <span className="text-red-600">-${summary.discount.toFixed(2)}</span>
            </div>
          )}
          
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Tax (10%)</span>
            <span className="text-gray-900">${summary.tax.toFixed(2)}</span>
          </div>
          
          <div className="flex justify-between text-lg font-semibold pt-2 border-t border-gray-200">
            <span>Total</span>
            <span>${summary.total.toFixed(2)}</span>
          </div>
        </div>

        {/* Checkout Button */}
        <button
          onClick={handleProceedToCheckout}
          disabled={items.length === 0 || loading}
          className={`w-full py-3 px-4 rounded-lg font-medium transition-colors ${
            items.length === 0 || loading
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-blue-500 hover:bg-blue-600 text-white'
          }`}
        >
          {loading ? (
            <div className="flex items-center justify-center gap-2">
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Processing...
            </div>
          ) : (
            `Proceed to Checkout • $${summary.total.toFixed(2)}`
          )}
        </button>
      </div>

      {/* Checkout Modal */}
      {isCheckoutModalOpen && (
        <CheckoutModal
          isOpen={isCheckoutModalOpen}
          onClose={() => setIsCheckoutModalOpen(false)}
          items={items}
          summary={summary}
          notes={notes}
        />
      )}
    </>
  )
}