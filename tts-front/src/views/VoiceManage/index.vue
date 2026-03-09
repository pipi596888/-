<template>
  <div class="voice-page">
    <div class="page-card">
      <div class="card-header">
        <h2>音色管理</h2>
        <el-button type="primary" @click="handleAddVoice">添加音色</el-button>
      </div>
      <div class="card-body">
        <el-table :data="voiceList" v-loading="loading" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="名称" min-width="120" />
          <el-table-column prop="tone" label="音色" min-width="120" />
          <el-table-column prop="gender" label="性别" width="80" />
          <el-table-column label="默认" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.isDefault" type="success">默认</el-tag>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button v-if="!row.isDefault" size="small" @click="handleSetDefault(row.id)">
                设为默认
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row.id)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useVoiceStore } from '@/store/voice'
import { storeToRefs } from 'pinia'

const voiceStore = useVoiceStore()
const { voiceList, loading } = storeToRefs(voiceStore)

onMounted(() => {
  voiceStore.fetchVoiceList()
})

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
    await ElMessageBox.confirm('确定要删除这个音色吗？', '警告', {
      type: 'warning',
    })
    await voiceStore.deleteVoice(id)
    ElMessage.success('删除成功')
  } catch {
    // cancelled
  }
}

function handleAddVoice() {
  ElMessage.info('添加音色弹窗即将推出')
}
</script>

<style scoped>
.voice-page {
  height: 100%;
}

.page-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.card-header h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-body {
  padding: 20px 24px;
  flex: 1;
  overflow: auto;
}
</style>
