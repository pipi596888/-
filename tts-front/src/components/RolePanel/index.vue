<template>
  <div class="role-panel">
    <div class="panel-header">
      <h3>音色列表</h3>
      <el-button size="small" @click="handleAddRole">添加</el-button>
    </div>
    <div class="role-list">
      <div
        v-for="voice in voiceList"
        :key="voice.id"
        class="role-item"
        :class="{ active: selectedVoiceId === voice.id }"
        @click="handleSelect(voice.id)"
      >
        <div class="role-info">
          <div class="role-avatar">
            <span>{{ voice.name.charAt(0) }}</span>
          </div>
          <div class="role-text">
            <span class="role-name">{{ voice.name }}</span>
            <span class="role-desc">{{ voice.gender }} · {{ voice.tone }}</span>
          </div>
          <el-tag v-if="voice.isDefault" size="small" type="success">默认</el-tag>
        </div>
        <div class="role-actions">
          <el-button
            v-if="!voice.isDefault"
            size="small"
            type="primary"
            link
            @click.stop="handleSetDefault(voice.id)"
          >
            设为默认
          </el-button>
          <el-button
            size="small"
            type="danger"
            link
            @click.stop="handleDelete(voice.id)"
          >
            删除
          </el-button>
        </div>
      </div>
      <el-empty v-if="voiceList.length === 0" description="暂无音色" :image-size="60" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useVoiceStore } from '@/store/voice'
import { useTTSStore } from '@/store/tts'
import { storeToRefs } from 'pinia'

const voiceStore = useVoiceStore()
const ttsStore = useTTSStore()
const { voiceList } = storeToRefs(voiceStore)
const { selectedVoiceId } = storeToRefs(ttsStore)

onMounted(() => {
  voiceStore.fetchVoiceList()
})

function handleSelect(voiceId: number) {
  ttsStore.selectedVoiceId = voiceId
}

async function handleSetDefault(id: number) {
  try {
    await voiceStore.setDefaultVoice(id)
    ElMessage.success('默认音色已更新')
  } catch {
    ElMessage.error('设置默认音色失败')
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个音色吗？', '提示', {
      type: 'warning',
    })
    await voiceStore.deleteVoice(id)
    ElMessage.success('删除成功')
  } catch {
    // 用户取消
  }
}

function handleAddRole() {
  ElMessage.info('添加音色弹窗即将推出')
}
</script>

<style scoped>
.role-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.role-list {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
}

.role-item {
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.role-item:hover {
  background: #f5f7fa;
}

.role-item.active {
  background: #ecf5ff;
  border-color: #409eff;
}

.role-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.role-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #409eff, #1677ff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
}

.role-text {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.role-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.role-desc {
  font-size: 12px;
  color: #909399;
}

.role-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.role-item:hover .role-actions {
  opacity: 1;
}
</style>
