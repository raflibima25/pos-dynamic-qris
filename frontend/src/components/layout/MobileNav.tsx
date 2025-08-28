'use client'

import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/store/auth'
import { 
  HomeIcon, 
  ShoppingCartIcon, 
  CubeIcon,
  ChartBarIcon,
  UserGroupIcon
} from '@heroicons/react/24/outline'
import Link from 'next/link'

export function MobileNav() {
  const pathname = usePathname()
  const { user } = useAuthStore()

  // Don't show on login page
  if (pathname === '/login' || !user) {
    return null
  }

  const navigationItems = [
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

  const visibleItems = navigationItems.filter(item => item.show)

  return (
    <div className="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50">
      <div className="grid grid-cols-5 gap-1 px-2 py-2">
        {visibleItems.map(item => {
          const Icon = item.icon
          const isActive = pathname === item.href
          
          return (
            <Link
              key={item.name}
              href={item.href as any}
              className={`flex flex-col items-center justify-center py-2 px-1 rounded-lg text-xs transition-colors ${
                isActive
                  ? 'bg-blue-100 text-blue-700'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
              }`}
            >
              <Icon className="w-5 h-5 mb-1" />
              <span className="truncate">{item.name}</span>
            </Link>
          )
        })}
      </div>
    </div>
  )
}