<template>
  <div class="admin-layout">
    <div class="admin-sidebar" :class="{ 'sidebar-collapsed': isSidebarCollapsed }">
      <div class="admin-header">
        <div class="header-content">
          <el-icon class="dashboard-icon"><Odometer /></el-icon>
          <h2 class="gradient-text" v-show="isSidebarTextVisible">{{ $t('admin_layout.dashboard') }}</h2>
        </div>
      </div>
      
      <nav class="admin-nav">
        <router-link to="/admin/articles" class="nav-item">
          <el-icon><Document /></el-icon>
          <span v-show="isSidebarTextVisible">{{ $t('admin_layout.articles') }}</span>
        </router-link>
        <router-link to="/admin/articles/new" class="nav-item">
          <el-icon><EditPen /></el-icon>
          <span v-show="isSidebarTextVisible">{{ $t('article_editor.create_article') }}</span>
        </router-link>
      </nav>

      <div class="admin-footer">
        <button class="back-btn" @click="$router.push('/')">
          <el-icon><Back /></el-icon>
          <span v-show="isSidebarTextVisible">{{ $t('admin_layout.back_to_site') }}</span>
        </button>
      </div>
    </div>

    <!-- Sidebar Toggle Button -->
    <div class="sidebar-toggle" @click="toggleSidebar">
      <el-icon :class="{ 'rotate-180': isSidebarCollapsed }">
        <ArrowLeft />
      </el-icon>
    </div>

    <div class="admin-main" :class="{ 'main-expanded': isSidebarCollapsed }">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ArrowLeft, Odometer } from '@element-plus/icons-vue'

const isSidebarCollapsed = ref(false)
const isSidebarTextVisible = ref(true)

const toggleSidebar = () => {
  if (!isSidebarCollapsed.value) {
    isSidebarCollapsed.value = true
    isSidebarTextVisible.value = false
  } else {
    isSidebarCollapsed.value = false
    setTimeout(() => {
      isSidebarTextVisible.value = true
    }, 200);
  }
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--dark-bg);
}

.admin-sidebar {
  width: 250px;
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
}

.admin-sidebar.sidebar-collapsed {
  width: 60px;
}

.sidebar-toggle {
  position: absolute;
  left: 250px;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 60px;
  background: var(--primary-color);
  border: 1px solid var(--border-color);
  border-radius: 0 8px 8px 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: left 0.3s ease, all 0.3s ease;
  z-index: 10;
}

.admin-sidebar.sidebar-collapsed + .sidebar-toggle {
  left: 60px;
}

.sidebar-toggle:hover {
  background: var(--primary-hover);
}

.sidebar-toggle .el-icon {
  color: white;
  font-size: 14px;
  transition: transform 0.3s ease;
}

.sidebar-toggle .rotate-180 {
  transform: rotate(180deg);
}

.admin-main {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  transition: margin-left 0.3s ease;
  margin-left: 0;
}

.admin-main.main-expanded {
  margin-left: 0;
}

.admin-header {
  padding: 2rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  height: 89px;
  box-sizing: border-box;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.dashboard-icon {
  font-size: 1.5rem;
  color: var(--primary-color);
  flex-shrink: 0;
}

.admin-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

.admin-nav {
  flex: 1;
  padding: 1rem 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  color: var(--text-primary);
  text-decoration: none;
  transition: all 0.3s ease;
}

/* Text hiding animations - when collapsing */
.admin-sidebar.sidebar-collapsed .header-content h2,
.admin-sidebar.sidebar-collapsed .nav-item span,
.admin-sidebar.sidebar-collapsed .back-btn span {
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.2s ease, visibility 0.2s ease;
}

/* Text showing animations - when expanding */
.admin-sidebar:not(.sidebar-collapsed) .header-content h2,
.admin-sidebar:not(.sidebar-collapsed) .nav-item span,
.admin-sidebar:not(.sidebar-collapsed) .back-btn span {
  opacity: 1;
  visibility: visible;
  transition: opacity 0.2s ease 0.3s, visibility 0.2s ease 0.3s;
}

/* Initially hidden when collapsed */
.admin-sidebar.sidebar-collapsed .header-content h2,
.admin-sidebar.sidebar-collapsed .nav-item span,
.admin-sidebar.sidebar-collapsed .back-btn span {
  opacity: 0;
  visibility: hidden;
}

/* Layout changes when collapsing */
.admin-sidebar.sidebar-collapsed .nav-item {
  padding: 1rem;
  justify-content: center;
  transition: padding 0.3s ease 0.2s, justify-content 0.3s ease 0.2s;
}

.admin-sidebar.sidebar-collapsed .admin-header {
  padding: 1.5rem 0.5rem;
  text-align: center;
  height: 89px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  transition: padding 0.3s ease 0.2s;
}

.admin-sidebar.sidebar-collapsed .header-content {
  justify-content: center;
  transition: justify-content 0.3s ease 0.2s;
}

.admin-sidebar.sidebar-collapsed .admin-footer {
  padding: 1.5rem 0.5rem;
  transition: padding 0.3s ease 0.2s;
}

.admin-sidebar.sidebar-collapsed .back-btn {
  padding: 0.75rem 0.5rem;
  justify-content: center;
  transition: padding 0.3s ease 0.2s, justify-content 0.3s ease 0.2s;
}

/* Layout changes when expanding */
.admin-sidebar:not(.sidebar-collapsed) .nav-item {
  transition: padding 0.3s ease, justify-content 0.3s ease;
}

.admin-sidebar:not(.sidebar-collapsed) .admin-header {
  transition: padding 0.3s ease;
}

.admin-sidebar:not(.sidebar-collapsed) .header-content {
  transition: justify-content 0.3s ease;
}

.admin-sidebar:not(.sidebar-collapsed) .admin-footer {
  transition: padding 0.3s ease;
}

.admin-sidebar:not(.sidebar-collapsed) .back-btn {
  transition: padding 0.3s ease, justify-content 0.3s ease;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.nav-item.router-link-active {
  background: var(--tech-gradient);
  color: white;
}

.admin-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.admin-main {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

@media (max-width: 768px) {
  .admin-sidebar {
    width: 200px;
  }
  
  .admin-sidebar.sidebar-collapsed {
    width: 50px;
  }
  
  .sidebar-toggle {
    left: 200px;
  }
  
  .admin-sidebar.sidebar-collapsed + .sidebar-toggle {
    left: 50px;
  }
  
  .admin-main {
    padding: 1rem;
  }
}
</style>