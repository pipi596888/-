<template>
  <div class="works-page">
    <div class="page-card">
      <div class="card-header">
        <h2>我的作品</h2>
      </div>
      <div class="card-body">
        <el-table :data="worksList" v-loading="loading" style="width: 100%">
          <el-table-column prop="taskId" label="任务ID" min-width="280" />
          <el-table-column prop="format" label="格式" width="100" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.status === 'success'" type="success">成功</el-tag>
              <el-tag v-else-if="row.status === 'failed'" type="danger">失败</el-tag>
              <el-tag v-else-if="row.status === 'processing'" type="warning">处理中</el-tag>
              <el-tag v-else type="info">等待中</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="进度" width="180">
            <template #default="{ row }">
              <el-progress :percentage="row.progress" :status="row.status === 'success' ? 'success' : row.status === 'failed' ? 'exception' : undefined" />
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" min-width="180" />
        </el-table>
        <el-empty v-if="!loading && worksList.length === 0" description="暂无作品" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { worksApi, type Work } from '@/api/works'

const worksList = ref<Work[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await worksApi.list()
    worksList.value = res.data.list
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.works-page {
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
