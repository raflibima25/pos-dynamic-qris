'use client'

import { useState, useEffect } from 'react'
import { ProductGrid } from '@/components/pos/ProductGrid'
import { ShoppingCart } from '@/components/pos/ShoppingCart'
import { CheckoutSummary } from '@/components/pos/CheckoutSummary'
import { useProductStore } from '@/store/product'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { redirect } from 'next/navigation'

export default function POSPage() {
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const { user } = useAuthStore()
  const { products, categories, listProducts, listCategories, loading } = useProductStore()
  const { items, summary } = useCartStore()

  // Redirect if not authenticated (but wait for auth check to complete)
  useEffect(() => {
    const token = useAuthStore.getState().token
    if (!token && !user) {
      redirect('/login')
    }
  }, [user])

  // Load initial data
  useEffect(() => {
    if (user) {
      listProducts()
      listCategories()
    }
  }, [user, listProducts, listCategories])

  // Filter products by category
  const filteredProducts = selectedCategory === 'all' 
    ? products 
    : products.filter(product => product.categoryId === selectedCategory)

  if (!user) {
    return null // Will redirect
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-4">
              <h1 className="text-xl font-semibold text-gray-900">
                Point of Sale
              </h1>
              <div className="text-sm text-gray-500">
                Welcome, {user.name}
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-600">
                Cart: {summary.itemCount} items • ${summary.total.toFixed(2)}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-12 gap-8">
          {/* Left Panel - Products */}
          <div className="col-span-8">
            {/* Category Filter */}
            <div className="mb-6">
              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => setSelectedCategory('all')}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    selectedCategory === 'all'
                      ? 'bg-blue-500 text-white'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  All Products
                </button>
                {categories.map(category => (
                  <button
                    key={category.id}
                    onClick={() => setSelectedCategory(category.id)}
                    className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                      selectedCategory === category.id
                        ? 'bg-blue-500 text-white'
                        : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {category.name}
                  </button>
                ))}
              </div>
            </div>

            {/* Product Grid */}
            <ProductGrid 
              products={filteredProducts}
              loading={loading}
            />
          </div>

          {/* Right Panel - Cart & Checkout */}
          <div className="col-span-4">
            <div className="sticky top-8 space-y-6">
              {/* Shopping Cart */}
              <div className="bg-white rounded-lg shadow-sm border p-6">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">
                  Shopping Cart
                </h2>
                <ShoppingCart items={items} />
              </div>

              {/* Checkout Summary */}
              {items.length > 0 && (
                <div className="bg-white rounded-lg shadow-sm border p-6">
                  <CheckoutSummary summary={summary} />
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}