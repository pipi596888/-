<template>
  <div class="login-page">
    <!-- 全屏背景 -->
    <div class="bg-layer"></div>
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <!-- 居中卡片 -->
    <div class="login-card">
      <div class="card-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" width="32" height="32">
              <path fill="currentColor" d="M12 3a9 9 0 0 0-9 9v7a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2v-3a3 3 0 0 1 3-3h2a3 3 0 0 1 3 3v3a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2v-7a9 9 0 0 0-9-9z"/>
            </svg>
          </div>
          <div class="brand-text">
            <h1>TTS 语音合成</h1>
            <p>智能文本转语音解决方案</p>
          </div>
        </div>
      </div>

      <div class="card-body">
        <!-- 登录表单 -->
        <transition name="fade" mode="out-in">
          <div v-if="!isRegister" key="login">
            <h2 class="welcome-text">欢迎回来</h2>
            <p class="sub-text">请登录您的账号以继续</p>

            <el-form @submit.prevent="handleLogin" class="login-form">
              <el-form-item>
                <el-input
                  v-model="username"
                  placeholder="请输入用户名"
                  size="large"
                  class="custom-input"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item>
                <el-input
                  v-model="password"
                  type="password"
                  placeholder="请输入密码"
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="handleLogin"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <div class="form-extra">
                <el-checkbox v-model="rememberMe">记住我</el-checkbox>
                <a href="javascript:void(0)" class="forgot-link">忘记密码？</a>
              </div>

              <el-button
                type="primary"
                :loading="loading"
                class="login-btn"
                size="large"
                @click="handleLogin"
              >
                {{ loading ? '登录中...' : '登 录' }}
              </el-button>
            </el-form>

            <div class="card-footer">
              <span>还没有账号？</span>
              <a href="javascript:void(0)" class="register-link" @click="switchToRegister">立即注册</a>
            </div>
          </div>

          <!-- 注册表单 -->
          <div v-else key="register">
            <h2 class="welcome-text">创建账号</h2>
            <p class="sub-text">填写以下信息完成注册</p>

            <el-form
              ref="registerFormRef"
              :model="registerForm"
              :rules="registerRules"
              class="login-form"
            >
              <el-form-item prop="username">
                <el-input
                  v-model="registerForm.username"
                  placeholder="请输入用户名"
                  size="large"
                  class="custom-input"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item prop="email">
                <el-input
                  v-model="registerForm.email"
                  placeholder="请输入邮箱"
                  size="large"
                  class="custom-input"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item prop="password">
                <el-input
                  v-model="registerForm.password"
                  type="password"
                  placeholder="请输入密码"
                  size="large"
                  show-password
                  class="custom-input"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item prop="confirmPassword">
                <el-input
                  v-model="registerForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入密码"
                  size="large"
                  show-password
                  class="custom-input"
                  @keyup.enter="handleRegisterSubmit"
                >
                  <template #prefix>
                    <div class="input-prefix">
                      <svg viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
                      </svg>
                    </div>
                  </template>
                </el-input>
              </el-form-item>

              <el-button
                type="primary"
                :loading="registerLoading"
                class="login-btn"
                size="large"
                @click="handleRegisterSubmit"
              >
                {{ registerLoading ? '注册中...' : '立即注册' }}
              </el-button>
            </el-form>

            <div class="card-footer">
              <span>已有账号？</span>
              <a href="javascript:void(0)" class="register-link" @click="switchToLogin">立即登录</a>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { userApi } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const rememberMe = ref(false)

// 注册相关
const isRegister = ref(false)
const registerLoading = ref(false)
const registerFormRef = ref<FormInstance>()
const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

function switchToRegister() {
  isRegister.value = true
}

function switchToLogin() {
  isRegister.value = false
}

// 密码验证规则
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

async function handleLogin() {
  if (!username.value || !password.value) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const res = await userApi.login(username.value, password.value)
    const token = res.token
    const user = res.user
    localStorage.setItem('token', token)
    userStore.setToken(token)
    userStore.setUserInfo(user)

    if (rememberMe.value) {
      localStorage.setItem('username', username.value)
    } else {
      localStorage.removeItem('username')
    }

    ElMessage.success('登录成功')
    router.push('/generate')
  } catch (e: any) {
    ElMessage.error(e.message || '登录失败')
  } finally {
    loading.value = false
  }
}

async function handleRegisterSubmit() {
  if (!registerFormRef.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      registerLoading.value = true
      try {
        await userApi.register(
          registerForm.username,
          registerForm.password,
          registerForm.email
        )

        ElMessage.success('注册成功，请登录！')
        // 切换回登录表单
        isRegister.value = false
        // 自动填充用户名
        username.value = registerForm.username
        password.value = ''
      } catch (e: any) {
        ElMessage.error(e.message || '注册失败')
      } finally {
        registerLoading.value = false
      }
    }
  })
}

const savedUsername = localStorage.getItem('username')
if (savedUsername) {
  username.value = savedUsername
  rememberMe.value = true
}
</script>

<style scoped>
/* 页面容器 */
.login-page {
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

/* 全屏背景 */
.bg-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
}

.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(255,255,255,0.08) 0%, rgba(255,255,255,0.02) 100%);
}

.circle-1 {
  width: 700px;
  height: 700px;
  top: -250px;
  left: -150px;
  animation: float 15s ease-in-out infinite;
}

.circle-2 {
  width: 500px;
  height: 500px;
  bottom: -150px;
  right: -100px;
  animation: float 18s ease-in-out infinite reverse;
}

.circle-3 {
  width: 300px;
  height: 300px;
  top: 40%;
  right: 15%;
  animation: pulse 10s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-50px) rotate(5deg); }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.2); opacity: 0.25; }
}

/* 登录卡片 */
.login-card {
  position: relative;
  width: 100%;
  max-width: 440px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
  overflow: hidden;
  animation: cardIn 0.6s ease-out;
}

@keyframes cardIn {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.card-header {
  padding: 36px 40px 32px;
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  text-align: center;
}

.logo-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  backdrop-filter: blur(10px);
}

.brand-text h1 {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 4px;
}

.brand-text p {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  letter-spacing: 1px;
}

.card-body {
  padding: 36px 40px 40px;
}

.welcome-text {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 6px;
  text-align: center;
}

.sub-text {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0 0 32px;
  text-align: center;
}

/* 表单样式 */
.login-form :deep(.el-form-item) {
  margin-bottom: 20px;
}

.custom-input :deep(.el-input__wrapper) {
  padding: 6px 14px;
  border-radius: 10px;
  box-shadow: 0 0 0 1px #e5e5e5 inset;
  transition: all 0.3s ease;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #b3b3b3 inset;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #1976d2 inset;
}

.custom-input :deep(.el-input__inner) {
  height: 42px;
  font-size: 15px;
}

.input-prefix {
  display: flex;
  align-items: center;
  color: #8c8c8c;
  margin-right: 8px;
}

.form-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.form-extra :deep(.el-checkbox__label) {
  color: #595959;
  font-size: 14px;
}

.forgot-link {
  font-size: 14px;
  color: #1976d2;
  text-decoration: none;
  font-weight: 500;
}

.forgot-link:hover {
  color: #1565c0;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  border: none;
  box-shadow: 0 4px 16px rgba(25, 118, 210, 0.3);
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(25, 118, 210, 0.4);
}

/* 底部 */
.card-footer {
  text-align: center;
  margin-top: 28px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.card-footer span {
  color: #8c8c8c;
  font-size: 14px;
}

.register-link {
  color: #1976d2;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  margin-left: 4px;
}

.register-link:hover {
  color: #1565c0;
}

/* 表单切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
