import { request } from '@/utils/request'
import type { Voice, VoiceCreateParams, VoiceListResponse } from '@/types/voice'

export const voiceApi = {
  list() {
    return request<VoiceListResponse>({
      url: '/voice/list',
      method: 'GET',
    })
  },

  create(params: VoiceCreateParams) {
    return request<Voice>({
      url: '/voice/create',
      method: 'POST',
      data: params,
    })
  },

  delete(id: number) {
    return request<void>({
      url: `/voice/${id}`,
      method: 'DELETE',
    })
  },

  setDefault(id: number) {
    return request<void>({
      url: `/voice/default/${id}`,
      method: 'PUT',
    })
  },
}
