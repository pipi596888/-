import type { Voice, VoiceCreateParams, VoiceListResponse } from './voice'
import type { TTSGenerateParams, TTSGenerateResponse, TTSTaskResponse } from './tts'

export interface User {
  id: number
  username: string
  balance: number
  characterCount: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export type { Voice, VoiceCreateParams, VoiceListResponse }
export type { TTSGenerateParams, TTSGenerateResponse, TTSTaskResponse }
