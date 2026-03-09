<template>
  <div class="audio-player">
    <div class="player-controls">
      <el-button
        :icon="isPlaying ? VideoPause : VideoPlay"
        circle
        size="large"
        type="primary"
        @click="togglePlay"
      />
      <el-button :icon="VideoCamera" circle size="large" @click="handleDownload" />
    </div>
    <div class="progress-bar">
      <el-slider
        v-model="currentTime"
        :max="duration"
        :format-tooltip="formatTime"
        @change="handleSeek"
      />
      <div class="time-display">
        <span>{{ formatTime(currentTime) }}</span>
        <span>/</span>
        <span>{{ formatTime(duration) }}</span>
      </div>
    </div>
    <div class="volume-control">
      <el-icon><component :is="'Volume'" /></el-icon>
      <el-slider v-model="volume" :min="0" :max="1" :step="0.1" @change="handleVolumeChange" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { VideoPlay, VideoPause, VideoCamera } from '@element-plus/icons-vue'
import { AudioPlayer as AudioPlayerClass, downloadAudio } from '@/utils/audio'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  audioUrl?: string
}>()

const player = ref<AudioPlayerClass | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(1)

watch(
  () => props.audioUrl,
  (url) => {
    if (url && player.value) {
      player.value.load(url)
      duration.value = player.value.duration
    }
  }
)

onMounted(() => {
  player.value = new AudioPlayerClass()
  if (player.value) {
    player.value.onTimeUpdate((time) => {
      currentTime.value = time
    })
    player.value.onEnded(() => {
      isPlaying.value = false
      currentTime.value = 0
    })
    player.value.onError((error) => {
      ElMessage.error(error.message)
      isPlaying.value = false
    })
  }
})

onUnmounted(() => {
  player.value?.destroy()
})

function togglePlay() {
  if (!player.value || !props.audioUrl) return

  if (isPlaying.value) {
    player.value.pause()
    isPlaying.value = false
  } else {
    player.value.play()
    isPlaying.value = true
  }
}

function handleSeek(time: number) {
  player.value?.seek(time)
}

function handleVolumeChange(value: number) {
  player.value?.setVolume(value)
}

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

async function handleDownload() {
  if (!props.audioUrl) return
  try {
    const filename = `audio_${Date.now()}.mp3`
    await downloadAudio(props.audioUrl, filename)
    ElMessage.success('下载成功')
  } catch {
    ElMessage.error('下载失败')
  }
}
</script>

<style scoped>
.audio-player {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.player-controls {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 16px;
}

.progress-bar {
  margin-bottom: 16px;
}

.time-display {
  display: flex;
  justify-content: center;
  gap: 8px;
  font-size: 12px;
  color: #666;
}

.volume-control {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
