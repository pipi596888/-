import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Voice } from '@/types/voice'
import { voiceApi } from '@/api/voice'

export const useVoiceStore = defineStore('voice', () => {
  const voiceList = ref<Voice[]>([])
  const loading = ref(false)

  async function fetchVoiceList() {
    loading.value = true
    try {
      const res = await voiceApi.list()
      voiceList.value = res.data.list
    } catch (error) {
      console.error('Failed to fetch voice list:', error)
    } finally {
      loading.value = false
    }
  }

  async function createVoice(params: { name: string; tone: string; gender: string }) {
    const res = await voiceApi.create(params)
    voiceList.value.push(res.data)
    return res.data
  }

  async function deleteVoice(id: number) {
    await voiceApi.delete(id)
    voiceList.value = voiceList.value.filter((v) => v.id !== id)
  }

  async function setDefaultVoice(id: number) {
    await voiceApi.setDefault(id)
    voiceList.value = voiceList.value.map((v) => ({
      ...v,
      isDefault: v.id === id,
    }))
  }

  function getDefaultVoice(): Voice | undefined {
    return voiceList.value.find((v) => v.isDefault)
  }

  function getVoiceById(id: number): Voice | undefined {
    return voiceList.value.find((v) => v.id === id)
  }

  return {
    voiceList,
    loading,
    fetchVoiceList,
    createVoice,
    deleteVoice,
    setDefaultVoice,
    getDefaultVoice,
    getVoiceById,
  }
})
