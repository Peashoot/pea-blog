<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card glass-effect">
        <div class="login-header">
          <h1 class="gradient-text">登录</h1>
          <p>欢迎回到 Pea Blog</p>
        </div>

        <el-form 
          ref="loginForm" 
          :model="form" 
          :rules="rules" 
          @submit.prevent="handleLogin"
          size="large"
        >
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="用户名"
              prefix-icon="User"
              :disabled="isLoading"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码"
              prefix-icon="Lock"
              show-password
              :disabled="isLoading"
              @keyup.enter="handleLogin"
            />
          </el-form-item>

          <el-form-item>
            <button 
              type="submit" 
              class="login-btn tech-button"
              :disabled="isLoading"
              style="width: 100%"
            >
              <span v-if="isLoading">登录中...</span>
              <span v-else>登录</span>
            </button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <p class="demo-hint">
            演示账户: admin / password
          </p>
          <router-link to="/" class="back-home">
            ← 返回首页
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import { ElMessage, type FormInstance } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginForm = ref<FormInstance>()
const isLoading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginForm.value) return

  const valid = await loginForm.value.validate().catch(() => false)
  if (!valid) return

  try {
    isLoading.value = true
    await authStore.login({
      username: form.username,
      password: form.password
    })

    ElMessage.success('登录成功')
    
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: var(--dark-bg);
  background-image: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(120, 219, 255, 0.2) 0%, transparent 50%);
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  padding: 2.5rem;
  text-align: center;
}

.login-header {
  margin-bottom: 2rem;
}

.login-header h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
}

.login-header p {
  color: var(--text-secondary);
  margin: 0;
}

.el-form {
  text-align: left;
}

.el-form-item {
  margin-bottom: 1.5rem;
}

.login-btn {
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
}

.login-footer {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.demo-hint {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin: 0 0 1rem 0;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.back-home {
  color: var(--primary-color);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.3s ease;
}

.back-home:hover {
  color: var(--text-primary);
}

@media (max-width: 480px) {
  .login-page {
    padding: 1rem;
  }
  
  .login-card {
    padding: 1.5rem;
  }
  
  .login-header h1 {
    font-size: 1.5rem;
  }
}
</style>