import axios from 'axios'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
})

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const data = error.response?.data
    if (data?.message) {
      return Promise.reject(new Error(data.message))
    }
    return Promise.reject(error)
  }
)

export default service

export function request<T = unknown>(config: { url: string; method?: string; data?: unknown }) {
  return service(config) as Promise<T>
}
