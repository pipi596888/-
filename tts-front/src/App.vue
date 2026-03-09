<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const showHeader = computed(() => route.path !== '/login')

function handleLogout() {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    userStore.logout()
    router.push('/login')
  }).catch(() => {})
}

onMounted(() => {
  const token = localStorage.getItem('token')
  if (token && !userStore.token) {
    userStore.setToken(token)
  }
})
</script>

<template>
  <div id="app">
    <el-container class="app-container">
      <el-header v-if="showHeader" class="main-header">
        <div class="header-left">
          <div class="logo">
            <svg viewBox="0 0 24 24" width="20" height="20">
              <path fill="currentColor" d="M12 3a9 9 0 0 0-9 9v7a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2v-3a3 3 0 0 1 3-3h2a3 3 0 0 1 3 3v3a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2v-7a9 9 0 0 0-9-9z"/>
            </svg>
          </div>
          <h1>TTS 语音合成系统</h1>
        </div>
        <el-menu 
          mode="horizontal" 
          :router="true" 
          :default-active="route.path"
          class="main-menu"
        >
          <el-menu-item index="/generate">生成音频</el-menu-item>
          <el-menu-item index="/works">我的作品</el-menu-item>
          <el-menu-item index="/voice">音色管理</el-menu-item>
        </el-menu>
        <div class="header-right">
          <span class="username">{{ userStore.userInfo?.username || '用户' }}</span>
          <el-button size="small" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <RouterView />
      </el-main>
    </el-container>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
}

#app {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f7fa;
}

.app-container {
  height: 100vh;
}

.main-header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 32px;
  height: 64px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-right: 48px;
}

.logo {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #1677ff, #409eff);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.header-left h1 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.main-menu {
  flex: 1;
  border-bottom: none !important;
  background: transparent;
}

.main-menu .el-menu-item {
  height: 64px;
  line-height: 64px;
  font-size: 15px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.username {
  color: #606266;
  font-size: 14px;
}

.main-content {
  padding: 20px;
  height: calc(100vh - 64px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.main-content > * {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
</style>
