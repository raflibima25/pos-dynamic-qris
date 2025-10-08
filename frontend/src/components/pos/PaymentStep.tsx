'use client'

import { useState, useEffect, useRef } from 'react'
import { QRISDisplay } from '@/components/pos/QRISDisplay'
import { PaymentStatusMonitor } from '@/components/pos/PaymentStatusMonitor'
import { usePaymentStore } from '@/store/payment'
import { useTransactionStore } from '@/store/transaction'
import { formatRupiah } from '@/lib/currency'
import { QRISCode } from '@/types'

interface PaymentStepProps {
  transactionId: string
  total: number
  onPaymentComplete: () => void
  onBack: () => void
}

export function PaymentStep({ transactionId, total, onPaymentComplete, onBack }: PaymentStepProps) {
  const [qrCode, setQrCode] = useState<QRISCode | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)
  const hasFetchedRef = useRef(false) // Track if we already fetched
  const { generateQRIS } = usePaymentStore()
  const { getTransaction } = useTransactionStore()

  useEffect(() => {
    let isMounted = true

    // Reset fetch flag when transactionId changes
    hasFetchedRef.current = false

    const fetchQRIS = async () => {
      // Only fetch once per transaction
      if (hasFetchedRef.current) {
        console.log('Already fetched QRIS for this transaction, skipping')
        return
      }

      hasFetchedRef.current = true

      try {
        setLoading(true)
        console.log('Generating QRIS for transaction:', transactionId)
        const response = await generateQRIS(transactionId, total)
        console.log('QRIS API response:', response)

        if (isMounted) {
          if (response && response.qr_code) {
            console.log('Setting QR code:', response.qr_code)
            setQrCode(response.qr_code)
            console.log('QRIS generated successfully')
          } else {
            console.error('No qr_code in response:', response)
            setError('Invalid response from server')
          }
        }
      } catch (err) {
        if (isMounted) {
          setError('Failed to generate QRIS code')
          console.error('QRIS generation error:', err)
        }
        hasFetchedRef.current = false // Reset on error to allow retry
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    fetchQRIS()

    return () => {
      isMounted = false
    }
    // Only depend on transactionId
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [transactionId])

  const handleRefreshQRIS = async () => {
    try {
      setLoading(true)
      setError(null)
      // Use refreshQRIS if available, otherwise fall back to generateQRIS
      const { refreshQRIS } = usePaymentStore.getState()
      const response = refreshQRIS ? await refreshQRIS(transactionId) : await generateQRIS(transactionId, total)
      if (response && response.qr_code) {
        setQrCode(response.qr_code)
      }
    } catch (err) {
      setError('Failed to refresh QRIS code')
      console.error('QRIS refresh error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handlePaymentSuccess = () => {
    // Refresh transaction data to get updated status
    getTransaction(transactionId)
    onPaymentComplete()
  }

  const handlePaymentFailed = () => {
    // Handle payment failure if needed
  }

  if (loading) {
    return (
      <div className="p-8 text-center">
        <div className="w-16 h-16 mx-auto mb-4 bg-blue-100 rounded-full flex items-center justify-center">
          <svg className="w-8 h-8 text-blue-500 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">Generating Payment Code</h3>
        <p className="text-gray-600">Please wait while we create your QRIS code...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-8 text-center">
        <div className="w-16 h-16 mx-auto mb-4 bg-red-100 rounded-full flex items-center justify-center">
          <svg className="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">Failed to Generate QRIS</h3>
        <p className="text-gray-600 mb-4">{error}</p>
        <button
          onClick={handleRefreshQRIS}
          className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          Try Again
        </button>
      </div>
    )
  }

  return (
    <>
      <div className="p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Complete Payment</h3>
        <p className="text-gray-600 mb-6">
          Scan the QR code below with your mobile payment app to complete the payment of {formatRupiah(total)}.
        </p>
        
        {qrCode && (
          <>
            <div className="mb-6">
              <QRISDisplay 
                qrCode={qrCode} 
                onRefresh={handleRefreshQRIS}
                size={256}
              />
            </div>
            
            <div className="mb-6">
              <PaymentStatusMonitor 
                transactionId={transactionId}
                onPaymentSuccess={handlePaymentSuccess}
                onPaymentFailed={handlePaymentFailed}
                pollInterval={5000}
              />
            </div>
          </>
        )}
      </div>

      <div className="flex gap-3 p-6 border-t">
        <button
          onClick={onBack}
          className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Back
        </button>
        <button
          onClick={onPaymentComplete}
          className="flex-1 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
        >
          Payment Completed
        </button>
      </div>
    </>
  )
}