import { saveAs } from 'file-saver'

export class AudioPlayer {
  private audio: HTMLAudioElement
  private _onTimeUpdate: ((time: number) => void) | null = null
  private _onEnded: (() => void) | null = null
  private _onError: ((error: Error) => void) | null = null

  constructor() {
    this.audio = new Audio()
    this.setupListeners()
  }

  private setupListeners() {
    this.audio.addEventListener('timeupdate', () => {
      this._onTimeUpdate?.(this.audio.currentTime)
    })

    this.audio.addEventListener('ended', () => {
      this._onEnded?.()
    })

    this.audio.addEventListener('error', (e) => {
      this._onError?.(new Error(`Audio error: ${(e.target as HTMLAudioElement)?.error?.message || 'Unknown error'}`))
    })
  }

  load(url: string) {
    this.audio.src = url
    this.audio.load()
  }

  play() {
    return this.audio.play()
  }

  pause() {
    this.audio.pause()
  }

  stop() {
    this.audio.pause()
    this.audio.currentTime = 0
  }

  seek(time: number) {
    this.audio.currentTime = time
  }

  setVolume(volume: number) {
    this.audio.volume = Math.max(0, Math.min(1, volume))
  }

  get duration(): number {
    return this.audio.duration || 0
  }

  get currentTime(): number {
    return this.audio.currentTime
  }

  get paused(): boolean {
    return this.audio.paused
  }

  onTimeUpdate(callback: (time: number) => void) {
    this._onTimeUpdate = callback
  }

  onEnded(callback: () => void) {
    this._onEnded = callback
  }

  onError(callback: (error: Error) => void) {
    this._onError = callback
  }

  destroy() {
    this.audio.pause()
    this.audio.src = ''
    this._onTimeUpdate = null
    this._onEnded = null
    this._onError = null
  }
}

export async function downloadAudio(url: string, filename: string) {
  try {
    const response = await fetch(url)
    const blob = await response.blob()
    saveAs(blob, filename)
  } catch (error) {
    throw new Error('Download failed')
  }
}
