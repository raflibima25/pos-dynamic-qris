'use client'

import { useState, useRef } from 'react'
import { PhotoIcon, XMarkIcon, ArrowUpTrayIcon } from '@heroicons/react/24/outline'
import { api } from '@/lib/api'

interface ImageUploadProps {
  value?: string
  onChange: (imageUrl: string) => void
  onRemove: () => void
  disabled?: boolean
  maxSizeDisplay?: string
}

export function ImageUpload({ 
  value, 
  onChange, 
  onRemove, 
  disabled = false,
  maxSizeDisplay = "2MB max"
}: ImageUploadProps) {
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [dragActive, setDragActive] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = async (file: File) => {
    if (!file) return

    // Validate file type
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif']
    if (!allowedTypes.includes(file.type)) {
      setError('Please upload a valid image file (JPEG, PNG, WebP, or GIF)')
      return
    }

    // Validate file size (2MB = 2 * 1024 * 1024 bytes)
    const maxSize = 2 * 1024 * 1024
    if (file.size > maxSize) {
      setError('File size must be less than 2MB')
      return
    }

    setError(null)
    setUploading(true)

    try {
      // Create FormData
      const formData = new FormData()
      formData.append('file', file)

      // Upload to backend
      const response = await api.post('/images/upload', formData, {
        'Content-Type': 'multipart/form-data'
      })

      if (response.data && response.data.image_url) {
        onChange(response.data.image_url)
      } else {
        throw new Error('Invalid response from server')
      }
    } catch (error: any) {
      console.error('Upload error:', error)
      setError(error.message || 'Failed to upload image')
    } finally {
      setUploading(false)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!disabled && !uploading) {
      setDragActive(true)
    }
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)
    
    if (disabled || uploading) return

    const files = e.dataTransfer.files
    if (files && files[0]) {
      handleFileSelect(files[0])
    }
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (files && files[0]) {
      handleFileSelect(files[0])
    }
    // Clear input value so same file can be selected again
    e.target.value = ''
  }

  const handleRemoveImage = async () => {
    if (!value) return

    try {
      await api.delete('/images/delete', { image_url: value })
      onRemove()
      setError(null)
    } catch (error: any) {
      console.error('Delete error:', error)
      // Still remove from UI even if backend delete fails
      onRemove()
    }
  }

  if (value) {
    return (
      <div className="relative">
        <div className="relative aspect-square w-full max-w-sm mx-auto bg-gray-100 rounded-lg overflow-hidden">
          <img
            src={value}
            alt="Product image"
            className="w-full h-full object-cover"
          />
          {!disabled && (
            <button
              type="button"
              onClick={handleRemoveImage}
              className="absolute top-2 right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors"
            >
              <XMarkIcon className="w-4 h-4" />
            </button>
          )}
        </div>
        {error && (
          <p className="mt-2 text-sm text-red-600 text-center">{error}</p>
        )}
      </div>
    )
  }

  return (
    <div>
      <div
        className={`relative border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
          dragActive
            ? 'border-blue-500 bg-blue-50'
            : error
            ? 'border-red-300 bg-red-50'
            : 'border-gray-300 hover:border-gray-400'
        } ${disabled || uploading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => {
          if (!disabled && !uploading && fileInputRef.current) {
            fileInputRef.current.click()
          }
        }}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileInputChange}
          className="hidden"
          disabled={disabled || uploading}
        />
        
        <div className="space-y-3">
          {uploading ? (
            <>
              <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                <svg className="w-6 h-6 text-blue-500 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <p className="text-sm text-gray-600">Uploading...</p>
            </>
          ) : (
            <>
              <div className={`mx-auto w-12 h-12 rounded-full flex items-center justify-center ${
                error ? 'bg-red-100' : 'bg-gray-100'
              }`}>
                {error ? (
                  <XMarkIcon className="w-6 h-6 text-red-500" />
                ) : dragActive ? (
                  <ArrowUpTrayIcon className="w-6 h-6 text-blue-500" />
                ) : (
                  <PhotoIcon className="w-6 h-6 text-gray-400" />
                )}
              </div>
              <div>
                <p className={`text-sm font-medium ${error ? 'text-red-600' : 'text-gray-900'}`}>
                  {error || (dragActive ? 'Drop image here' : 'Click to upload or drag and drop')}
                </p>
                {!error && (
                  <p className="text-xs text-gray-500 mt-1">
                    PNG, JPG, WebP, GIF up to {maxSizeDisplay}
                  </p>
                )}
              </div>
            </>
          )}
        </div>
      </div>
      
      {error && (
        <button
          type="button"
          onClick={() => setError(null)}
          className="mt-2 text-sm text-blue-600 hover:text-blue-500 underline"
        >
          Try again
        </button>
      )}
    </div>
  )
}