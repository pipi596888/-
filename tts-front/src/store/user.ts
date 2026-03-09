import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/api'
import { userApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<User | null>(null)
  const token = ref<string>('')

  async function fetchUserInfo() {
    try {
      const res = await userApi.getInfo()
      userInfo.value = res.data
    } catch (error) {
      console.error('Failed to fetch user info:', error)
    }
  }

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUserInfo(info: User) {
    userInfo.value = info
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    userInfo,
    token,
    fetchUserInfo,
    setToken,
    setUserInfo,
    logout,
  }
})
