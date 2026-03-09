<template>
  <div class="generate-panel">
    <div class="panel-header">
      <h3>生成设置</h3>
    </div>
    <div class="settings">
      <div class="setting-item">
        <label>音频格式</label>
        <el-select v-model="format">
          <el-option label="MP3" value="mp3" />
          <el-option label="WAV" value="wav" />
          <el-option label="FLAC" value="flac" />
        </el-select>
      </div>
      <div class="setting-item">
        <label>声道</label>
        <el-select v-model="channel">
          <el-option label="单声道" value="mono" />
          <el-option label="双声道" value="stereo" />
        </el-select>
      </div>
      <div class="setting-item">
        <label>总字符数</label>
        <span class="character-count">{{ totalCharacters }}</span>
      </div>
    </div>
    <div class="actions">
      <el-button type="primary" size="large" :loading="isGenerating" @click="handleGenerate">
        {{ isGenerating ? '生成中...' : '生成音频' }}
      </el-button>
      <el-button size="large" @click="handleClear">清空</el-button>
    </div>
    <div v-if="taskStatus" class="task-status">
      <el-progress
        :percentage="taskStatus.progress"
        :status="taskStatus.status === 'success' ? 'success' : taskStatus.status === 'failed' ? 'exception' : undefined"
      />
      <p v-if="taskStatus.status === 'success'" class="success-text">
        生成成功！
      </p>
      <p v-if="taskStatus.status === 'failed'" class="error-text">
        {{ taskStatus.error || '生成失败' }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { useTTSStore } from '@/store/tts'
import { storeToRefs } from 'pinia'

const emit = defineEmits<{
  (e: 'generate'): void
}>()

const ttsStore = useTTSStore()
const { format, channel, totalCharacters, isGenerating, taskStatus } = storeToRefs(ttsStore)

async function handleGenerate() {
  if (ttsStore.segments.length === 0) {
    ElMessage.warning('请先添加文本片段')
    return
  }

  try {
    await ttsStore.generateAudio()
    emit('generate')
  } catch (error: any) {
    ElMessage.error(error.message || '生成失败')
  }
}

function handleClear() {
  ttsStore.clearSegments()
  ttsStore.resetTask()
}
</script>

<style scoped>
.generate-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.settings {
  padding: 16px;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #666;
}

.character-count {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
}

.actions {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: auto;
}

.task-status {
  padding: 16px;
  border-top: 1px solid #eee;
}

.success-text {
  color: #67c23a;
  margin-top: 8px;
}

.error-text {
  color: #f56c6c;
  margin-top: 8px;
}
</style>
