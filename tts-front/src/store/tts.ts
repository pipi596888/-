import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Segment, AudioFormat, AudioChannel, TTSTaskResponse } from '@/types/tts'
import { ttsApi } from '@/api/tts'

export const useTTSStore = defineStore('tts', () => {
  const segments = ref<Segment[]>([])
  const selectedVoiceId = ref<number>(0)
  const format = ref<AudioFormat>('mp3')
  const channel = ref<AudioChannel>('mono')
  const currentTaskId = ref<string>('')
  const taskStatus = ref<TTSTaskResponse | null>(null)
  const isGenerating = ref(false)
  const pollingTimer = ref<number | null>(null)

  const totalCharacters = computed(() => {
    return segments.value.reduce((sum, seg) => sum + seg.text.length, 0)
  })

  function addSegment(text: string, voiceId?: number) {
    const id = `seg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    segments.value.push({
      id,
      voiceId: voiceId || selectedVoiceId.value,
      text,
      order: segments.value.length,
    })
  }

  function updateSegment(id: string, updates: Partial<Segment>) {
    const index = segments.value.findIndex((s) => s.id === id)
    if (index !== -1) {
      const existing = segments.value[index]!
      segments.value[index] = {
        id: existing.id,
        voiceId: updates.voiceId ?? existing.voiceId,
        emotion: updates.emotion,
        text: updates.text ?? existing.text,
        order: existing.order,
      }
    }
  }

  function removeSegment(id: string) {
    segments.value = segments.value.filter((s) => s.id !== id)
  }

  function reorderSegments(newOrder: Segment[]) {
    segments.value = newOrder.map((seg, index) => ({
      ...seg,
      order: index,
    }))
  }

  function clearSegments() {
    segments.value = []
  }

  async function generateAudio(): Promise<string> {
    if (segments.value.length === 0) {
      throw new Error('No segments to generate')
    }

    isGenerating.value = true
    try {
      const res = await ttsApi.generate({
        segments: segments.value.map((s) => ({
          voiceId: s.voiceId,
          emotion: s.emotion,
          text: s.text,
        })),
        format: format.value,
        channel: channel.value,
      })

      currentTaskId.value = res.data.taskId
      startPolling(res.data.taskId)
      return res.data.taskId
    } finally {
      isGenerating.value = false
    }
  }

  function startPolling(taskId: string) {
    stopPolling()
    pollingTimer.value = window.setInterval(async () => {
      try {
        const res = await ttsApi.getTask(taskId)
        taskStatus.value = res.data

        if (res.data.status === 'success' || res.data.status === 'failed') {
          stopPolling()
        }
      } catch (error) {
        console.error('Polling error:', error)
      }
    }, 2000)
  }

  function stopPolling() {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
      pollingTimer.value = null
    }
  }

  function resetTask() {
    currentTaskId.value = ''
    taskStatus.value = null
    stopPolling()
  }

  return {
    segments,
    selectedVoiceId,
    format,
    channel,
    currentTaskId,
    taskStatus,
    isGenerating,
    totalCharacters,
    addSegment,
    updateSegment,
    removeSegment,
    reorderSegments,
    clearSegments,
    generateAudio,
    startPolling,
    stopPolling,
    resetTask,
  }
})
