const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    // Get token from localStorage if available
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token')
    }
  }

  setToken(token: string) {
    this.token = token
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token)
    }
  }

  removeToken() {
    this.token = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token')
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    skipJsonContentType: boolean = false
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const headers: Record<string, string> = {
      ...((options.headers as Record<string, string>) || {}),
    }
    
    // Only add Content-Type for JSON requests
    if (!skipJsonContentType) {
      headers['Content-Type'] = 'application/json'
    }

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`
    }

    const config: RequestInit = {
      ...options,
      headers,
    }

    const response = await fetch(url, config)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || `HTTP ${response.status}`)
    }

    return response.json()
  }

  // Auth endpoints
  async login(email: string, password: string) {
    const response = await this.request<any>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    
    if (response.data?.token) {
      this.setToken(response.data.token)
    }
    
    return response
  }

  async logout() {
    await this.request('/auth/logout', { method: 'POST' })
    this.removeToken()
  }

  async getCurrentUser() {
    return this.request<any>('/auth/me')
  }

  async register(data: {
    name: string
    email: string
    password: string
    role: string
  }) {
    return this.request<any>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Product endpoints
  async getProducts(params?: {
    category_id?: string
    is_active?: boolean
    search?: string
    limit?: number
    offset?: number
  }) {
    const searchParams = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const endpoint = `/products${queryString ? `?${queryString}` : ''}`
    
    return this.request<any>(endpoint)
  }

  async getProduct(id: string) {
    return this.request<any>(`/products/${id}`)
  }

  async createProduct(data: {
    name: string
    description?: string
    price: number
    stock: number
    category_id: string
    sku?: string
  }) {
    return this.request<any>('/products', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateProduct(id: string, data: {
    name: string
    description?: string
    price: number
    stock: number
    category_id: string
    sku?: string
    is_active?: boolean
  }) {
    return this.request<any>(`/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteProduct(id: string) {
    return this.request<any>(`/products/${id}`, {
      method: 'DELETE',
    })
  }

  async updateProductStock(id: string, quantity: number) {
    return this.request<any>(`/products/${id}/stock`, {
      method: 'PATCH',
      body: JSON.stringify({ quantity }),
    })
  }

  // Category endpoints
  async getCategories(params?: { limit?: number; offset?: number }) {
    const searchParams = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, value.toString())
        }
      })
    }
    
    const queryString = searchParams.toString()
    const endpoint = `/categories${queryString ? `?${queryString}` : ''}`
    
    return this.request<any>(endpoint)
  }

  async createCategory(data: { name: string }) {
    return this.request<any>('/categories', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Generic HTTP methods
  async get(endpoint: string) {
    return this.request<any>(endpoint)
  }

  async post(endpoint: string, data?: any, customHeaders?: Record<string, string>) {
    const isFormData = data instanceof FormData
    
    return this.request<any>(endpoint, {
      method: 'POST',
      body: isFormData ? data : (data ? JSON.stringify(data) : undefined),
      headers: customHeaders,
    }, isFormData)
  }

  async put(endpoint: string, data?: any) {
    return this.request<any>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async patch(endpoint: string, data?: any) {
    return this.request<any>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete(endpoint: string, data?: any) {
    return this.request<any>(endpoint, {
      method: 'DELETE',
      body: data ? JSON.stringify(data) : undefined,
    })
  }
}

export const api = new ApiClient(API_BASE_URL)
export default api