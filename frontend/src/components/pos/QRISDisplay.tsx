'use client'

import { useState, useEffect } from 'react'
import { QRISCode } from '@/types'
import QRCode from 'qrcode'

interface QRISDisplayProps {
  qrCode: QRISCode
  onRefresh?: () => void
  size?: number
}

export function QRISDisplay({ qrCode, onRefresh, size = 256 }: QRISDisplayProps) {
  const [timeLeft, setTimeLeft] = useState<number>(0)
  const [isExpired, setIsExpired] = useState<boolean>(false)
  const [qrImageUrl, setQrImageUrl] = useState<string>('')

  // Generate QR code image from string on client-side
  useEffect(() => {
    if (qrCode.qr_code) {
      console.log('Generating QR code from string:', qrCode.qr_code.substring(0, 50) + '...')
      QRCode.toDataURL(qrCode.qr_code, {
        width: size,
        margin: 1,
        errorCorrectionLevel: 'M'
      })
        .then((url) => {
          console.log('QR code generated successfully')
          setQrImageUrl(url)
        })
        .catch((err) => {
          console.error('Failed to generate QR code:', err)
          console.error('QR string was:', qrCode.qr_code)
        })
    } else {
      console.warn('No QR code string available')
    }
  }, [qrCode.qr_code, size])

  useEffect(() => {
    // Calculate time left until expiry
    const expiryTime = new Date(qrCode.expires_at).getTime()
    const now = new Date().getTime()
    const diff = expiryTime - now

    if (diff <= 0) {
      setIsExpired(true)
      setTimeLeft(0)
    } else {
      setIsExpired(false)
      setTimeLeft(Math.floor(diff / 1000))
    }
  }, [qrCode.expires_at])

  useEffect(() => {
    if (isExpired || timeLeft <= 0) return

    const timer = setInterval(() => {
      setTimeLeft(prev => {
        if (prev <= 1) {
          setIsExpired(true)
          clearInterval(timer)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [isExpired, timeLeft])

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  return (
    <div className="flex flex-col items-center p-6 bg-white rounded-lg border border-gray-200 w-full max-w-md">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Scan to Pay</h3>

      {qrImageUrl ? (
        <div className="relative">
          <img
            src={qrImageUrl}
            alt="QRIS Code"
            width={size}
            height={size}
            className="border border-gray-300 rounded-lg"
          />

          {isExpired && (
            <div className="absolute inset-0 bg-black bg-opacity-70 flex items-center justify-center rounded-lg">
              <div className="text-center text-white">
                <div className="text-lg font-semibold mb-2">QR Code Expired</div>
                {onRefresh && (
                  <button
                    onClick={onRefresh}
                    className="px-4 py-2 bg-blue-500 hover:bg-blue-600 rounded-lg text-sm font-medium transition-colors"
                  >
                    Generate New QR Code
                  </button>
                )}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div className="flex items-center justify-center w-64 h-64 bg-gray-100 border border-gray-300 rounded-lg">
          <div className="text-gray-500 text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-2"></div>
            Generating QR Code...
          </div>
        </div>
      )}
      
      <div className="mt-4 text-center">
        <div className="text-sm text-gray-600 mb-1">
          Scan this QR code with your mobile payment app
        </div>
        
        {!isExpired ? (
          <div className="flex items-center justify-center gap-2">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span className="text-sm font-medium text-gray-900">
              Expires in: {formatTime(timeLeft)}
            </span>
          </div>
        ) : (
          <div className="text-sm font-medium text-red-600">
            QR Code has expired
          </div>
        )}
      </div>
      
      <div className="mt-4 text-xs text-gray-500 text-center">
        Supported payment methods: QRIS, GoPay, OVO, Dana, LinkAja, ShopeePay
      </div>
    </div>
  )
}