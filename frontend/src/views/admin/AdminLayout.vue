<template>
  <div class="admin-layout">
    <div class="admin-sidebar">
      <div class="admin-header">
        <h2 class="gradient-text">{{ $t('admin_layout.dashboard') }}</h2>
      </div>
      
      <nav class="admin-nav">
        <router-link to="/admin/articles" class="nav-item">
          <el-icon><Document /></el-icon>
          <span>{{ $t('admin_layout.articles') }}</span>
        </router-link>
        <router-link to="/admin/articles/new" class="nav-item">
          <el-icon><EditPen /></el-icon>
          <span>{{ $t('article_editor.create_article') }}</span>
        </router-link>
      </nav>

      <div class="admin-footer">
        <button class="back-btn" @click="$router.push('/')">
          <el-icon><Back /></el-icon>
          <span>{{ $t('admin_layout.back_to_site') }}</span>
        </button>
      </div>
    </div>

    <div class="admin-main">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
// Admin layout component
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
}

.admin-header {
  padding: 2rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
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
  
  .admin-main {
    padding: 1rem;
  }
}
</style>