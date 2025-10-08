'use client'

import { useState, useEffect } from 'react'
import { useProductStore } from '@/store/product'
import { useAuthStore } from '@/store/auth'
import { redirect } from 'next/navigation'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'
import { AddProductModal } from '@/components/products/AddProductModal'
import { EditProductModal } from '@/components/products/EditProductModal'
import { formatRupiah } from '@/lib/currency'
import { Product } from '@/types'
import { Navbar } from '@/components/layout/Navbar'
import { MobileNav } from '@/components/layout/MobileNav'

export default function ProductsPage() {
  const { user } = useAuthStore()
  const { products, categories, loading, listProducts, listCategories, deleteProduct } = useProductStore()
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [searchTerm, setSearchTerm] = useState('')
  const [showAddModal, setShowAddModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)

  // Redirect if not admin
  useEffect(() => {
    const token = useAuthStore.getState().token
    if (!token && !user) {
      redirect('/login')
    }
    if (user && user.role !== 'admin') {
      redirect('/dashboard')
    }
  }, [user])

  // Load data
  useEffect(() => {
    if (user && user.role === 'admin') {
      listProducts()
      listCategories()
    }
  }, [user, listProducts, listCategories])

  // Filter products with safety check
  const filteredProducts = products.filter(product => {
    // Safety check: skip undefined or null products
    if (!product || !product.name) return false

    const matchesCategory = selectedCategory === 'all' || product.category_id === selectedCategory
    const matchesSearch = product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         (product.description && product.description.toLowerCase().includes(searchTerm.toLowerCase()))
    return matchesCategory && matchesSearch
  })

  // Handle edit product
  const handleEditProduct = (product: Product) => {
    setSelectedProduct(product)
    setShowEditModal(true)
  }

  // Handle close edit modal
  const handleCloseEditModal = () => {
    setShowEditModal(false)
    setSelectedProduct(null)
  }

  if (!user || user.role !== 'admin') {
    return null
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-16 md:pb-0">
      {/* Navigation */}
      <Navbar backHref="/dashboard" />
      
      {/* Add Product Bar */}
      <div className="bg-white border-b px-4 py-3">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div>
            <p className="text-sm text-gray-500">Manage your inventory and product catalog</p>
          </div>
          <button 
            onClick={() => setShowAddModal(true)}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors"
          >
            <PlusIcon className="w-4 h-4" />
            Add Product
          </button>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Filters */}
        <div className="mb-6 space-y-4 bg-white rounded-lg shadow-sm border p-4">
          <div className="flex flex-col sm:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <label htmlFor="search" className="block text-sm font-medium text-gray-700 mb-1">
                Search Products
              </label>
              <input
                type="text"
                id="search"
                placeholder="Search by name or description..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Category Filter */}
            <div className="sm:w-48">
              <label htmlFor="category" className="block text-sm font-medium text-gray-700 mb-1">
                Category
              </label>
              <select
                id="category"
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Categories</option>
                {categories.map(category => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        {/* Products Grid */}
        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm border p-4 animate-pulse">
                <div className="aspect-square bg-gray-200 rounded-lg mb-3"></div>
                <div className="h-4 bg-gray-200 rounded mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-2/3 mb-2"></div>
                <div className="h-6 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        ) : filteredProducts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredProducts.map(product => (
              <div key={product.id} className="bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow">
                {/* Product Image */}
                <div className="aspect-square relative overflow-hidden rounded-t-lg bg-gray-100">
                  {product.image_url ? (
                    <img
                      src={product.image_url}
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
                  
                  {/* Stock Status */}
                  {product.stock <= 0 && (
                    <div className="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
                      Out of Stock
                    </div>
                  )}
                  {product.stock > 0 && product.stock <= 5 && (
                    <div className="absolute top-2 right-2 bg-orange-500 text-white text-xs px-2 py-1 rounded">
                      Low Stock
                    </div>
                  )}
                </div>

                {/* Product Details */}
                <div className="p-4">
                  <h3 className="font-semibold text-gray-900 mb-1 line-clamp-2">
                    {product.name}
                  </h3>
                  
                  {product.description && (
                    <p className="text-sm text-gray-600 mb-2 line-clamp-2">
                      {product.description}
                    </p>
                  )}
                  
                  <div className="flex items-center justify-between mb-3">
                    <div>
                      <p className="text-lg font-bold text-gray-900">
                        {formatRupiah(product.price)}
                      </p>
                      <p className="text-sm text-gray-500">
                        Stock: {product.stock}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-xs text-gray-500">
                        {categories.find(c => c.id === product.category_id)?.name}
                      </p>
                      <p className="text-xs text-gray-500">
                        {product.sku}
                      </p>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-2">
                    <button 
                      onClick={() => handleEditProduct(product)}
                      className="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded text-sm font-medium transition-colors flex items-center justify-center gap-1"
                    >
                      <PencilIcon className="w-4 h-4" />
                      Edit
                    </button>
                    <button 
                      onClick={() => deleteProduct(product.id)}
                      className="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded text-sm font-medium transition-colors flex items-center justify-center"
                    >
                      <TrashIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <div className="w-24 h-24 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
              <svg className="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No products found</h3>
            <p className="text-gray-500 mb-6">
              {searchTerm || selectedCategory !== 'all' 
                ? 'Try adjusting your search criteria' 
                : 'Get started by adding your first product'
              }
            </p>
            <button 
              onClick={() => setShowAddModal(true)}
              className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 mx-auto transition-colors"
            >
              <PlusIcon className="w-4 h-4" />
              Add Product
            </button>
          </div>
        )}
      </div>

      {/* Add Product Modal */}
      <AddProductModal 
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
      />

      {/* Edit Product Modal */}
      <EditProductModal 
        isOpen={showEditModal}
        onClose={handleCloseEditModal}
        product={selectedProduct}
      />

      {/* Mobile Navigation */}
      <MobileNav />
    </div>
  )
}