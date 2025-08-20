<template>
  <header class="navbar glass-effect">
    <div class="container">
      <div class="navbar-content">
        <router-link to="/" class="logo">
          <h1 class="gradient-text">Pea Blog</h1>
        </router-link>
        
        <nav class="nav-menu" :class="{ active: isMobileMenuOpen }">
          <router-link to="/" class="nav-link" @click="closeMobileMenu">{{ $t('navbar.home') }}</router-link>
          <div class="language-switcher">
            <el-dropdown @command="handleLanguageCommand">
              <span class="el-dropdown-link">
                <el-icon><Switch /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="zh-CN">中文</el-dropdown-item>
                  <el-dropdown-item command="en">English</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div v-if="authStore.isLoggedIn" class="nav-user">
            <el-dropdown @command="handleUserCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="authStore.user?.avatar">
                  {{ authStore.user?.username?.charAt(0).toUpperCase() }}
                </el-avatar>
                <span class="username">{{ authStore.user?.username }}</span>
                <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="authStore.isAdmin" command="admin">
                    <el-icon><Setting /></el-icon>
                    {{ $t('navbar.admin') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    {{ $t('navbar.logout') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <router-link 
            v-else 
            to="/login" 
            class="nav-link login-btn tech-button"
            @click="closeMobileMenu"
          >
            {{ $t('navbar.login') }}
          </router-link>
        </nav>

        <button 
          class="mobile-menu-btn"
          @click="toggleMobileMenu"
          v-if="isMobile"
        >
          <el-icon><Menu /></el-icon>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores'
import { useResponsive } from '@/composables'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const authStore = useAuthStore()
const { isMobile } = useResponsive()
const { locale } = useI18n()

const isMobileMenuOpen = ref(false)

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

const handleLanguageCommand = (command: string) => {
  locale.value = command
  localStorage.setItem('locale', command)
}

const handleUserCommand = async (command: string) => {
  switch (command) {
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      await authStore.logout()
      ElMessage.success(t('navbar.logout_success'))
      router.push('/')
      break
  }
}

const handleClickOutside = (event: Event) => {
  const navbar = document.querySelector('.navbar')
  if (navbar && !navbar.contains(event.target as Node)) {
    closeMobileMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-color);
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-link {
  color: var(--text-primary);
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-1px);
}

.nav-link.router-link-active {
  background: var(--tech-gradient);
  color: white;
}

.login-btn {
  margin-left: 1rem;
}

.nav-user .user-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.nav-user .user-info:hover {
  background: rgba(255, 255, 255, 0.1);
}

.username {
  color: var(--text-primary);
  font-weight: 500;
}

.dropdown-icon {
  font-size: 14px;
  color: var(--text-secondary);
}

.mobile-menu-btn {
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.mobile-menu-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

@media (max-width: 768px) {
  .nav-menu {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--card-bg);
    backdrop-filter: blur(20px);
    border: 1px solid var(--border-color);
    border-top: none;
    flex-direction: column;
    gap: 0;
    padding: 1rem 0;
    transform: translateY(-100%);
    opacity: 0;
    pointer-events: none;
    transition: all 0.3s ease;
  }

  .nav-menu.active {
    transform: translateY(0);
    opacity: 1;
    pointer-events: all;
  }

  .nav-link {
    padding: 1rem 2rem;
    width: 100%;
    text-align: left;
    border-radius: 0;
  }

  .login-btn {
    margin: 0.5rem 2rem;
    width: calc(100% - 4rem);
    justify-content: center;
  }
}
</style>