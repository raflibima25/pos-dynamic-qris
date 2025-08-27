'use client'

import { CartItem } from '@/types'
import { useCartStore } from '@/store/cart'
import { TrashIcon, PlusIcon, MinusIcon } from '@heroicons/react/24/outline'

interface ShoppingCartProps {
  items: CartItem[]
}

export function ShoppingCart({ items }: ShoppingCartProps) {
  const { updateQuantity, removeItem } = useCartStore()

  if (items.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
          <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4m0 0L7 13m0 0l-1.5 6M7 13l-1.5 6m0 0h9m0 0l2.5-5M6 19a2 2 0 100 4 2 2 0 000-4zm10 0a2 2 0 100 4 2 2 0 000-4z" />
          </svg>
        </div>
        <p className="text-gray-500 text-sm">Your cart is empty</p>
        <p className="text-gray-400 text-xs mt-1">Add some products to get started</p>
      </div>
    )
  }

  return (
    <div className="space-y-3 max-h-96 overflow-y-auto">
      {items.map(item => (
        <CartItemRow
          key={item.id || `${item.product.id}-${Date.now()}`}
          item={item}
          onUpdateQuantity={updateQuantity}
          onRemove={removeItem}
        />
      ))}
    </div>
  )
}

interface CartItemRowProps {
  item: CartItem
  onUpdateQuantity: (productId: string, quantity: number) => void
  onRemove: (productId: string) => void
}

function CartItemRow({ item, onUpdateQuantity, onRemove }: CartItemRowProps) {
  const handleQuantityChange = (newQuantity: number) => {
    if (newQuantity <= 0) {
      onRemove(item.product.id)
    } else {
      onUpdateQuantity(item.product.id, newQuantity)
    }
  }

  return (
    <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
      {/* Product Image */}
      <div className="w-12 h-12 rounded bg-gray-200 flex-shrink-0 overflow-hidden">
        {item.product.image ? (
          <img
            src={item.product.image}
            alt={item.product.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-gray-900 truncate">
          {item.product.name}
        </p>
        <p className="text-sm text-gray-500">
          ${item.product.price.toFixed(2)} each
        </p>
      </div>

      {/* Quantity Controls */}
      <div className="flex items-center gap-2">
        <button
          onClick={() => handleQuantityChange(item.quantity - 1)}
          className="w-7 h-7 flex items-center justify-center rounded-full bg-white border border-gray-300 hover:bg-gray-50 transition-colors"
        >
          <MinusIcon className="w-3 h-3" />
        </button>
        
        <span className="w-8 text-center text-sm font-medium text-gray-900">
          {item.quantity}
        </span>
        
        <button
          onClick={() => handleQuantityChange(item.quantity + 1)}
          className="w-7 h-7 flex items-center justify-center rounded-full bg-white border border-gray-300 hover:bg-gray-50 transition-colors"
        >
          <PlusIcon className="w-3 h-3" />
        </button>
      </div>

      {/* Subtotal */}
      <div className="text-right">
        <p className="text-sm font-medium text-gray-900">
          ${(item.product.price * item.quantity).toFixed(2)}
        </p>
      </div>

      {/* Remove Button */}
      <button
        onClick={() => onRemove(item.product.id)}
        className="w-7 h-7 flex items-center justify-center rounded-full text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors"
      >
        <TrashIcon className="w-4 h-4" />
      </button>
    </div>
  )
}