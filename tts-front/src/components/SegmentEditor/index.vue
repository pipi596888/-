<template>
  <div class="segment-editor">
    <el-table :data="segments" row-key="id" :drag="true" @row-drop="handleDrop">
      <el-table-column label="序号" width="60">
        <template #default="{ $index }">
          <span>{{ $index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="音色" width="140">
        <template #default="{ row }">
          <VoiceSelector v-model="row.voiceId" />
        </template>
      </el-table-column>
      <el-table-column label="情绪" width="120">
        <template #default="{ row }">
          <el-select v-model="row.emotion" placeholder="情绪" clearable>
            <el-option label="中性" value="neutral" />
            <el-option label="开心" value="happy" />
            <el-option label="悲伤" value="sad" />
            <el-option label="愤怒" value="angry" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="文本">
        <template #default="{ row }">
          <el-input
            v-model="row.text"
            type="textarea"
            :rows="2"
            placeholder="请输入文本"
            @input="handleUpdate(row.id, { text: row.text })"
          />
        </template>
      </el-table-column>
      <el-table-column label="字符" width="60">
        <template #default="{ row }">
          {{ row.text.length }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="80">
        <template #default="{ row }">
          <el-button type="danger" :icon="Delete" circle size="small" @click="handleDelete(row.id)" />
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import VoiceSelector from '@/components/VoiceSelector/index.vue'
import { useTTSStore } from '@/store/tts'
import { storeToRefs } from 'pinia'

const ttsStore = useTTSStore()
const { segments } = storeToRefs(ttsStore)

function handleUpdate(id: string, updates: any) {
  ttsStore.updateSegment(id, updates)
}

function handleDelete(id: string) {
  ttsStore.removeSegment(id)
}

function handleDrop(newSegments: any[]) {
  ttsStore.reorderSegments(newSegments)
}
</script>

<style scoped>
.segment-editor {
  width: 100%;
  height: 100%;
}

.segment-editor :deep(.el-table) {
  height: 100%;
}

.segment-editor :deep(.el-table__inner-wrapper) {
  height: 100%;
}

.segment-editor :deep(.el-table__body-wrapper) {
  overflow: auto;
}
</style>
