import { request } from '@/utils/request'
import type { TTSGenerateParams, TTSGenerateResponse, TTSTaskResponse } from '@/types/tts'

export const ttsApi = {
  generate(params: TTSGenerateParams) {
    return request<TTSGenerateResponse>({
      url: '/tts/generate',
      method: 'POST',
      data: params,
    })
  },

  getTask(taskId: string) {
    return request<TTSTaskResponse>({
      url: `/tts/task/${taskId}`,
      method: 'GET',
    })
  },
}
