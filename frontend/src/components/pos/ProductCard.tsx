'use client'

import { Product } from '@/types'
import { useCartStore } from '@/store/cart'
import { useState } from 'react'
import { ShoppingCartIcon, PlusIcon, MinusIcon } from '@heroicons/react/24/outline'

interface ProductCardProps {
  product: Product
}

export function ProductCard({ product }: ProductCardProps) {
  const { addItem, getItem, updateQuantity } = useCartStore()
  const [isAdding, setIsAdding] = useState(false)
  
  const cartItem = getItem(product.id)
  const quantity = cartItem?.quantity || 0

  const handleAddToCart = async () => {
    if (product.stock <= 0) return
    
    setIsAdding(true)
    addItem(product)
    setTimeout(() => setIsAdding(false), 200) // Brief animation
  }

  const handleUpdateQuantity = (newQuantity: number) => {
    if (newQuantity < 0) return
    if (newQuantity > product.stock) return
    updateQuantity(product.id, newQuantity)
  }

  const isOutOfStock = product.stock <= 0
  const isLowStock = product.stock <= 5 && product.stock > 0

  return (
    <div className="bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow">
      {/* Product Image */}
      <div className="aspect-square relative overflow-hidden rounded-t-lg bg-gray-100">
        {product.image ? (
          <img
            src={product.image}
            alt={product.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <svg className="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
        )}
        
        {/* Stock Badge */}
        {isOutOfStock && (
          <div className="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
            Out of Stock
          </div>
        )}
        {isLowStock && (
          <div className="absolute top-2 right-2 bg-orange-500 text-white text-xs px-2 py-1 rounded">
            Low Stock
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="p-4">
        <h3 className="font-medium text-gray-900 mb-1 line-clamp-2">
          {product.name}
        </h3>
        
        <div className="flex items-center justify-between mb-3">
          <div>
            <p className="text-lg font-semibold text-gray-900">
              ${product.price.toFixed(2)}
            </p>
            <p className="text-sm text-gray-500">
              Stock: {product.stock}
            </p>
          </div>
        </div>

        {/* Add to Cart Controls */}
        <div className="space-y-2">
          {quantity === 0 ? (
            <button
              onClick={handleAddToCart}
              disabled={isOutOfStock || isAdding}
              className={`w-full flex items-center justify-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                isOutOfStock
                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                  : isAdding
                  ? 'bg-green-500 text-white'
                  : 'bg-blue-500 hover:bg-blue-600 text-white'
              }`}
            >
              {isAdding ? (
                <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              ) : (
                <ShoppingCartIcon className="w-4 h-4" />
              )}
              {isAdding ? 'Added!' : 'Add to Cart'}
            </button>
          ) : (
            <div className="flex items-center justify-between bg-gray-50 rounded-lg p-2">
              <button
                onClick={() => handleUpdateQuantity(quantity - 1)}
                className="w-8 h-8 flex items-center justify-center rounded-full bg-white border hover:bg-gray-50 transition-colors"
              >
                <MinusIcon className="w-4 h-4" />
              </button>
              
              <span className="font-medium text-gray-900 mx-3">
                {quantity}
              </span>
              
              <button
                onClick={() => handleUpdateQuantity(quantity + 1)}
                disabled={quantity >= product.stock}
                className={`w-8 h-8 flex items-center justify-center rounded-full border transition-colors ${
                  quantity >= product.stock
                    ? 'bg-gray-100 border-gray-200 text-gray-400 cursor-not-allowed'
                    : 'bg-white border-gray-300 hover:bg-gray-50'
                }`}
              >
                <PlusIcon className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}