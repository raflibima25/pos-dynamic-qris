'use client'

import { useState, useEffect } from 'react'
import { Payment, PaymentStatus } from '@/types'
import { usePaymentStore } from '@/store/payment'

interface PaymentStatusMonitorProps {
  transactionId: string
  onPaymentSuccess?: () => void
  onPaymentFailed?: () => void
  pollInterval?: number // in milliseconds
}

export function PaymentStatusMonitor({ 
  transactionId, 
  onPaymentSuccess, 
  onPaymentFailed,
  pollInterval = 5000 
}: PaymentStatusMonitorProps) {
  const [paymentStatus, setPaymentStatus] = useState<PaymentStatus>('pending')
  const [lastChecked, setLastChecked] = useState<Date>(new Date())
  const { getPaymentStatus, loading, error } = usePaymentStore()

  useEffect(() => {
    let isMounted = true
    let intervalId: NodeJS.Timeout | null = null
    let hasCallbackFired = false

    const checkPaymentStatus = async () => {
      try {
        const response = await getPaymentStatus(transactionId)
        if (!isMounted) return

        if (response) {
          // Handle PaymentStatusResponse type
          const status = response.status || (response as any).payment_status
          setPaymentStatus(status)
          setLastChecked(new Date())

          // Handle status changes (only fire callbacks once)
          if (!hasCallbackFired) {
            switch (status) {
              case 'success':
                hasCallbackFired = true
                // Clear interval when payment succeeds
                if (intervalId) clearInterval(intervalId)
                onPaymentSuccess?.()
                break
              case 'failed':
              case 'cancelled':
              case 'expired':
                hasCallbackFired = true
                // Clear interval when payment fails
                if (intervalId) clearInterval(intervalId)
                onPaymentFailed?.()
                break
            }
          }
        }
      } catch (err) {
        console.error('Payment status check error:', err)
      }
    }

    // Check immediately
    checkPaymentStatus()

    // Set up polling - ONLY poll if status is pending
    intervalId = setInterval(() => {
      // Stop polling if callback has fired
      if (hasCallbackFired && intervalId) {
        clearInterval(intervalId)
        return
      }
      checkPaymentStatus()
    }, pollInterval)

    return () => {
      isMounted = false
      if (intervalId) {
        clearInterval(intervalId)
        intervalId = null
      }
    }
    // IMPORTANT: Only transactionId and pollInterval as dependencies
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [transactionId, pollInterval])
  
  const getStatusColor = () => {
    switch (paymentStatus) {
      case 'success': return 'bg-green-500'
      case 'failed': return 'bg-red-500'
      case 'expired': return 'bg-orange-500'
      case 'cancelled': return 'bg-gray-500'
      default: return 'bg-blue-500'
    }
  }
  
  const getStatusText = () => {
    switch (paymentStatus) {
      case 'success': return 'Payment Successful'
      case 'failed': return 'Payment Failed'
      case 'expired': return 'Payment Expired'
      case 'cancelled': return 'Payment Cancelled'
      default: return 'Waiting for Payment'
    }
  }
  
  return (
    <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
      <div className="flex items-center gap-3">
        <div className={`w-3 h-3 rounded-full ${getStatusColor()} ${paymentStatus === 'pending' ? 'animate-pulse' : ''}`}></div>
        <div>
          <div className="font-medium text-gray-900">{getStatusText()}</div>
          <div className="text-sm text-gray-500">
            Last checked: {lastChecked.toLocaleTimeString()}
          </div>
        </div>
      </div>
      
      {loading && (
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Checking...
        </div>
      )}
      
      {error && (
        <div className="text-sm text-red-600">
          Error checking status
        </div>
      )}
    </div>
  )
}