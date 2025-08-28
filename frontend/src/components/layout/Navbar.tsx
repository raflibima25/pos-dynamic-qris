'use client'

import { useRouter, usePathname } from 'next/navigation'
import { useAuthStore } from '@/store/auth'
import { 
  ArrowLeftIcon, 
  HomeIcon, 
  ShoppingCartIcon, 
  CubeIcon,
  ChartBarIcon,
  UserGroupIcon,
  ArrowRightOnRectangleIcon
} from '@heroicons/react/24/outline'
import Link from 'next/link'

interface NavbarProps {
  title?: string
  showBackButton?: boolean
  backHref?: string
}

export function Navbar({ title, showBackButton = true, backHref }: NavbarProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, logout } = useAuthStore()

  const handleBack = () => {
    if (backHref) {
      router.push(backHref as any)
    } else {
      router.back()
    }
  }

  const handleLogout = async () => {
    await logout()
    router.push('/login')
  }

  // Navigation items based on user role and current page
  const getNavigationItems = () => {
    const items = [
      {
        name: 'Dashboard',
        href: '/dashboard',
        icon: HomeIcon,
        show: true
      },
      {
        name: 'POS',
        href: '/pos',
        icon: ShoppingCartIcon,
        show: true
      },
      {
        name: 'Products',
        href: '/products',
        icon: CubeIcon,
        show: user?.role === 'admin'
      },
      {
        name: 'Analytics',
        href: '/analytics',
        icon: ChartBarIcon,
        show: user?.role === 'admin'
      },
      {
        name: 'Users',
        href: '/users',
        icon: UserGroupIcon,
        show: user?.role === 'admin'
      }
    ]

    return items.filter(item => item.show)
  }

  return (
    <div className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Left side - Back button and title */}
          <div className="flex items-center space-x-4">
            {showBackButton && pathname !== '/dashboard' && (
              <button
                onClick={handleBack}
                className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </button>
            )}
            
            <div>
              <h1 className="text-xl font-semibold text-gray-900">
                {title || getPageTitle(pathname)}
              </h1>
              {user && (
                <p className="text-sm text-gray-500">
                  Welcome, {user.name}
                </p>
              )}
            </div>
          </div>

          {/* Right side - Navigation and user menu */}
          <div className="flex items-center space-x-4">
            {/* Quick Navigation */}
            <div className="hidden md:flex items-center space-x-2">
              {getNavigationItems().map(item => {
                const Icon = item.icon
                const isActive = pathname === item.href
                
                return (
                  <Link
                    key={item.name}
                    href={item.href as any}
                    className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2 ${
                      isActive
                        ? 'bg-blue-100 text-blue-700'
                        : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                    }`}
                  >
                    <Icon className="w-4 h-4" />
                    <span className="hidden lg:inline">{item.name}</span>
                  </Link>
                )
              })}
            </div>

            {/* User info and logout */}
            {user && (
              <div className="flex items-center space-x-3">
                <div className="hidden sm:block text-right">
                  <p className="text-sm font-medium text-gray-900">{user.name}</p>
                  <p className="text-xs text-gray-500 capitalize">{user.role}</p>
                </div>
                <button
                  onClick={handleLogout}
                  className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                  title="Logout"
                >
                  <ArrowRightOnRectangleIcon className="w-5 h-5" />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

// Helper function to get page title from pathname
function getPageTitle(pathname: string): string {
  const titleMap: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/pos': 'Point of Sale',
    '/products': 'Product Management',
    '/analytics': 'Analytics & Reports',
    '/users': 'User Management',
    '/login': 'Login'
  }

  return titleMap[pathname] || 'QRIS POS System'
}