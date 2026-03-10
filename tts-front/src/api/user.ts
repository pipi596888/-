import { request } from '@/utils/request'
import type { User } from '@/types/api'
import { encryptPassword } from '@/utils/passwordCipher'

export interface LoginResponse {
  token: string
  user: User
}

export const userApi = {
  async login(username: string, password: string) {
    const passwordEncrypted = await encryptPassword(password)
    return request<LoginResponse>({
      url: '/user/login',
      method: 'POST',
      data: { username, passwordEncrypted },
    })
  },

  async register(username: string, password: string, email: string) {
    const passwordEncrypted = await encryptPassword(password)
    return request<LoginResponse>({
      url: '/user/register',
      method: 'POST',
      data: { username, passwordEncrypted, email },
    })
  },

  getInfo() {
    return request<User>({
      url: '/user/info',
      method: 'GET',
    })
  },
}

