import axios from 'axios'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
})

type ApiEnvelope<T = unknown> = {
  code: number
  message?: string
  msg?: string
  data: T
}

function isEnvelope(value: unknown): value is ApiEnvelope {
  if (!value || typeof value !== 'object') return false
  return Object.prototype.hasOwnProperty.call(value, 'code') && Object.prototype.hasOwnProperty.call(value, 'data')
}

const SENSITIVE_KEYS = new Set([
  'password',
  'confirmPassword',
  'oldPassword',
  'newPassword',
  'token',
  'accessToken',
  'refreshToken',
])

function redactObject(value: any): any {
  if (!value || typeof value !== 'object') return value
  if (Array.isArray(value)) return value.map(redactObject)
  const out: Record<string, any> = {}
  for (const [k, v] of Object.entries(value)) {
    if (SENSITIVE_KEYS.has(k)) {
      out[k] = '[REDACTED]'
    } else {
      out[k] = redactObject(v)
    }
  }
  return out
}

function redactAxiosConfigData(error: any) {
  const cfg = error?.config
  if (!cfg || cfg.data == null) return
  try {
    if (typeof cfg.data === 'string') {
      const trimmed = cfg.data.trim()
      if (trimmed.startsWith('{') || trimmed.startsWith('[')) {
        const parsed = JSON.parse(cfg.data)
        cfg.data = JSON.stringify(redactObject(parsed))
      }
      return
    }
    if (typeof cfg.data === 'object') {
      cfg.data = redactObject(cfg.data)
    }
  } catch {
    // best-effort redaction; never break error handling
  }
}

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
    const payload = response.data

    // Supports both:
    // 1) { code, message, data } style envelopes
    // 2) bare JSON payloads (go-zero httpx.OkJsonCtx style)
    if (isEnvelope(payload)) {
      const ok = payload.code === 0 || payload.code === 200
      if (!ok) {
        const msg = payload.message || payload.msg || 'Request failed'
        throw new Error(msg)
      }
      return payload.data
    }

    return payload
  },
  (error) => {
    const data = error.response?.data
    const status = error.response?.status

    // Avoid leaking credentials into logs if downstream code prints axios errors.
    redactAxiosConfigData(error)

    if (status === 401) {
      // Global auth fail: clear token and go back to login.
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      if (window.location.pathname !== '/login') {
        window.location.assign('/login')
      }
      return Promise.reject(new Error('Unauthorized'))
    }

    if (data?.message) {
      return Promise.reject(new Error(data.message))
    }
    return Promise.reject(error)
  }
)

export default service

export function request<T = unknown>(config: { url: string; method?: string; data?: unknown; params?: Record<string, any> }) {
  return service(config) as Promise<T>
}

