import { request } from '@/utils/request'
import type { User } from '@/types/api'

export interface LoginResponse {
  token: string
  user: User
}

export const userApi = {
  login(username: string, password: string) {
    return request<LoginResponse>({
      url: '/user/login',
      method: 'POST',
      data: { username, password },
    })
  },

  register(username: string, password: string, email: string) {
    return request<LoginResponse>({
      url: '/user/register',
      method: 'POST',
      data: { username, password, email },
    })
  },

  getInfo() {
    return request<User>({
      url: '/user/info',
      method: 'GET',
    })
  },
}
