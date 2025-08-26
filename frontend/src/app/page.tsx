'use client'

import Link from 'next/link'
import { ShoppingCart, BarChart3, Package, Users } from 'lucide-react'

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            QRIS POS System
          </h1>
          <p className="text-xl text-gray-600">
            Dynamic QRIS Point of Sale System for Modern Cashiers
          </p>
        </div>

        {/* Feature Cards */}
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
          <div className="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <ShoppingCart className="h-12 w-12 mx-auto text-blue-600 mb-4" />
            <h3 className="text-lg font-semibold mb-2">Point of Sale</h3>
            <p className="text-gray-600 text-sm">
              Fast and intuitive POS interface for cashiers
            </p>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <Package className="h-12 w-12 mx-auto text-green-600 mb-4" />
            <h3 className="text-lg font-semibold mb-2">Product Management</h3>
            <p className="text-gray-600 text-sm">
              Manage inventory and product catalog
            </p>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <BarChart3 className="h-12 w-12 mx-auto text-purple-600 mb-4" />
            <h3 className="text-lg font-semibold mb-2">Analytics</h3>
            <p className="text-gray-600 text-sm">
              Real-time sales reports and insights
            </p>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
            <Users className="h-12 w-12 mx-auto text-orange-600 mb-4" />
            <h3 className="text-lg font-semibold mb-2">User Management</h3>
            <p className="text-gray-600 text-sm">
              Manage staff and access controls
            </p>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            href="/pos"
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-8 rounded-lg text-center transition-colors"
          >
            Start POS Session
          </Link>
          
          <Link
            href="/dashboard"
            className="bg-white hover:bg-gray-50 text-gray-900 font-semibold py-3 px-8 rounded-lg border border-gray-300 text-center transition-colors"
          >
            Go to Dashboard
          </Link>
        </div>

        {/* Status Cards */}
        <div className="grid md:grid-cols-3 gap-6 mt-12">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h4 className="text-lg font-semibold mb-3 text-gray-900">
              🚀 Key Features
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>• Dynamic QRIS generation</li>
              <li>• Real-time payment monitoring</li>
              <li>• Automatic receipt generation</li>
              <li>• Stock management</li>
            </ul>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h4 className="text-lg font-semibold mb-3 text-gray-900">
              💳 Payment Methods
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>• QRIS (All Indonesian Banks)</li>
              <li>• GoPay, OVO, Dana</li>
              <li>• LinkAja, ShopeePay</li>
              <li>• Bank Transfer</li>
            </ul>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h4 className="text-lg font-semibold mb-3 text-gray-900">
              📊 Reports Available
            </h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>• Daily sales summary</li>
              <li>• Product performance</li>
              <li>• Payment method analytics</li>
              <li>• Monthly/Yearly reports</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}