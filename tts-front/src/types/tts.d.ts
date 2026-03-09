export interface Segment {
  id: string
  voiceId: number
  emotion?: string
  text: string
  order: number
}

export type AudioFormat = 'mp3' | 'wav' | 'flac'
export type AudioChannel = 'mono' | 'stereo'
export type TaskStatus = 'pending' | 'processing' | 'success' | 'failed'

export interface TTSGenerateParams {
  segments: {
    voiceId: number
    emotion?: string
    text: string
  }[]
  format: AudioFormat
  channel: AudioChannel
}

export interface TTSGenerateResponse {
  taskId: string
}

export interface TTSTaskResponse {
  taskId: string
  status: TaskStatus
  progress: number
  audioUrl?: string
  error?: string
}

export interface TTSState {
  segments: Segment[]
  selectedVoiceId: number
  format: AudioFormat
  channel: AudioChannel
}
