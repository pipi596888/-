<template>
  <div class="generate-page">
    <!-- 左侧音色面板 -->
    <aside class="sidebar sidebar-left">
      <div class="sidebar-card">
        <div class="card-header">
          <div class="header-title">
            <svg viewBox="0 0 24 24" width="20" height="20" class="header-icon">
              <path fill="currentColor" d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
            </svg>
            <h3>音色选择</h3>
          </div>
          <el-button type="primary" size="small" plain @click="handleAddVoice">
            <svg viewBox="0 0 24 24" width="14" height="14" style="margin-right: 4px;">
              <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
            添加
          </el-button>
        </div>
        <div class="voice-list">
          <div
            v-for="voice in voiceList"
            :key="voice.id"
            class="voice-item"
            :class="{ active: selectedVoiceId === voice.id }"
            @click="handleSelectVoice(voice.id)"
          >
            <div class="voice-avatar" :style="{ background: getAvatarColor(voice.name) }">
              {{ voice.name.charAt(0) }}
            </div>
            <div class="voice-info">
              <div class="voice-name">{{ voice.name }}</div>
              <div class="voice-meta">{{ voice.gender }} · {{ voice.tone }}</div>
            </div>
            <el-tag v-if="voice.isDefault" size="small" type="success" class="default-tag">默认</el-tag>
          </div>
          <el-empty v-if="voiceList.length === 0" description="暂无音色" :image-size="80" />
        </div>
      </div>
    </aside>

    <!-- 中间编辑区域 -->
    <main class="main-area">
      <div class="main-card editor-card">
        <div class="card-header">
          <div class="header-title">
            <svg viewBox="0 0 24 24" width="20" height="20" class="header-icon">
              <path fill="currentColor" d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
            </svg>
            <h3>文本编辑</h3>
          </div>
          <el-button type="primary" @click="handleAddSegment">
            <svg viewBox="0 0 24 24" width="14" height="14" style="margin-right: 4px;">
              <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
            添加片段
          </el-button>
        </div>
        <div class="card-body">
          <SegmentEditor />
        </div>
      </div>

      <!-- 音频播放器 -->
      <transition name="slide-up">
        <div v-if="taskStatus?.audioUrl" class="main-card player-card">
          <AudioPlayer :audio-url="taskStatus.audioUrl" />
        </div>
      </transition>
    </main>

    <!-- 右侧设置面板 -->
    <aside class="sidebar sidebar-right">
      <div class="sidebar-card">
        <div class="card-header">
          <div class="header-title">
            <svg viewBox="0 0 24 24" width="20" height="20" class="header-icon">
              <path fill="currentColor" d="M19.14 12.94c.04-.31.06-.63.06-.94 0-.31-.02-.63-.06-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.04.31-.06.63-.06.94s.02.63.06.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>
            </svg>
            <h3>生成设置</h3>
          </div>
        </div>
        <div class="settings-body">
          <div class="setting-group">
            <label class="setting-label">音频格式</label>
            <el-select v-model="format" placeholder="选择格式" class="setting-select">
              <el-option label="MP3" value="mp3">
                <div class="option-item">
                  <span>MP3</span>
                  <span class="option-desc">通用格式</span>
                </div>
              </el-option>
              <el-option label="WAV" value="wav">
                <div class="option-item">
                  <span>WAV</span>
                  <span class="option-desc">无损格式</span>
                </div>
              </el-option>
              <el-option label="FLAC" value="flac">
                <div class="option-item">
                  <span>FLAC</span>
                  <span class="option-desc">高品质</span>
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="setting-group">
            <label class="setting-label">声道</label>
            <el-select v-model="channel" placeholder="选择声道" class="setting-select">
              <el-option label="单声道" value="mono">
                <div class="option-item">
                  <span>单声道</span>
                  <span class="option-desc">适合语音</span>
                </div>
              </el-option>
              <el-option label="双声道" value="stereo">
                <div class="option-item">
                  <span>双声道</span>
                  <span class="option-desc">立体声</span>
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="setting-group">
            <label class="setting-label">文本字数</label>
            <div class="char-count">
              <span class="count-number">{{ totalCharacters }}</span>
              <span class="count-unit">字符</span>
            </div>
          </div>
        </div>

        <div class="actions-area">
          <el-button
            type="primary"
            size="large"
            class="generate-btn"
            :loading="isGenerating"
            @click="handleGenerate"
          >
            <svg v-if="!isGenerating" viewBox="0 0 24 24" width="20" height="20" class="btn-icon">
              <path fill="currentColor" d="M8 5v14l11-7z"/>
            </svg>
            {{ isGenerating ? '生成中...' : '开始生成' }}
          </el-button>
          <el-button size="large" class="clear-btn" @click="handleClear">
            清空内容
          </el-button>
        </div>

        <!-- 生成状态 -->
        <div v-if="taskStatus" class="status-area">
          <div class="progress-wrapper">
            <el-progress
              :percentage="taskStatus.progress"
              :stroke-width="6"
              :status="taskStatus.status === 'success' ? 'success' : taskStatus.status === 'failed' ? 'exception' : undefined"
            />
          </div>
          <p v-if="taskStatus.status === 'success'" class="status-text success">
            <svg viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
            音频生成成功
          </p>
          <p v-if="taskStatus.status === 'failed'" class="status-text error">
            <svg viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
            </svg>
            {{ taskStatus.error || '生成失败，请重试' }}
          </p>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import SegmentEditor from '@/components/SegmentEditor/index.vue'
import AudioPlayer from '@/components/AudioPlayer/index.vue'
import { useTTSStore } from '@/store/tts'
import { useVoiceStore } from '@/store/voice'
import { storeToRefs } from 'pinia'

const ttsStore = useTTSStore()
const voiceStore = useVoiceStore()
const { voiceList } = storeToRefs(voiceStore)
const { selectedVoiceId, format, channel, totalCharacters, isGenerating, taskStatus } = storeToRefs(ttsStore)

onMounted(() => {
  voiceStore.fetchVoiceList()
})

function handleSelectVoice(voiceId: number) {
  ttsStore.selectedVoiceId = voiceId
}

function handleAddVoice() {
  ElMessage.info('添加音色功能开发中')
}

function handleAddSegment() {
  const defaultVoice = voiceStore.getDefaultVoice()
  ttsStore.addSegment('', defaultVoice?.id)
}

async function handleGenerate() {
  if (ttsStore.segments.length === 0) {
    ElMessage.warning('请先添加文本片段')
    return
  }

  try {
    await ttsStore.generateAudio()
    ElMessage.success('开始生成音频')
  } catch (error: any) {
    ElMessage.error(error.message || '生成失败')
  }
}

function handleClear() {
  ttsStore.clearSegments()
  ttsStore.resetTask()
  ElMessage.success('已清空')
}

// 生成随机颜色
function getAvatarColor(name: string) {
  const colors = [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
    'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)',
  ]
  const index = name.charCodeAt(0) % colors.length
  return colors[index]
}
</script>

<style scoped>
.generate-page {
  height: 100%;
  display: flex;
  gap: 16px;
  padding: 16px;
  box-sizing: border-box;
}

/* 侧边栏 */
.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.sidebar-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  color: #1976d2;
}

.card-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1a1a1a;
}

/* 音色列表 */
.voice-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.voice-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 8px;
}

.voice-item:hover {
  background: #f8f9fc;
}

.voice-item.active {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  box-shadow: 0 2px 8px rgba(25, 118, 210, 0.15);
}

.voice-avatar {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.voice-info {
  flex: 1;
  min-width: 0;
}

.voice-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 2px;
}

.voice-meta {
  font-size: 12px;
  color: #8c8c8c;
}

.default-tag {
  flex-shrink: 0;
}

/* 中间区域 */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
  min-height: 0;
}

.main-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.editor-card {
  flex: 1;
  min-height: 200px;
  display: flex;
  flex-direction: column;
}

.player-card {
  flex-shrink: 0;
  padding: 20px;
}

.card-body {
  padding: 16px;
  flex: 1;
  overflow: auto;
  min-height: 0;
}

/* 右侧设置 */
.sidebar-right .card-header {
  border-bottom: none;
  padding-bottom: 0;
}

.settings-body {
  padding: 20px;
}

.setting-group {
  margin-bottom: 24px;
}

.setting-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #595959;
  margin-bottom: 10px;
}

.setting-select {
  width: 100%;
}

.option-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.option-desc {
  font-size: 12px;
  color: #8c8c8c;
}

.char-count {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.count-number {
  font-size: 32px;
  font-weight: 700;
  color: #1976d2;
}

.count-unit {
  font-size: 14px;
  color: #8c8c8c;
}

/* 操作按钮 */
.actions-area {
  padding: 0 20px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.generate-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  border: none;
  box-shadow: 0 4px 16px rgba(25, 118, 210, 0.3);
  transition: all 0.3s ease;
}

.generate-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(25, 118, 210, 0.4);
}

.btn-icon {
  margin-right: 8px;
}

.clear-btn {
  width: 100%;
  height: 40px;
  border-radius: 10px;
}

/* 状态区域 */
.status-area {
  padding: 16px 20px;
  background: #fafafa;
  border-top: 1px solid #f0f0f0;
}

.progress-wrapper {
  margin-bottom: 12px;
}

.status-text {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  margin: 0;
}

.status-text.success {
  color: #52c41a;
}

.status-text.error {
  color: #ff4d4f;
}

/* 动画 */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
