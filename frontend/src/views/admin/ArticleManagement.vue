<template>
  <div class="article-management">
    <div class="page-header">
      <h1>{{ $t('article_management.title') }}</h1>
      <router-link to="/admin/articles/new" class="tech-button">
        <el-icon><Plus /></el-icon>
        {{ $t('article_management.new_article') }}
      </router-link>
    </div>

    <div class="articles-list glass-effect">
      <div v-if="isLoading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>{{ $t('article_management_page.loading') }}</p>
      </div>

      <div v-else-if="articles.length > 0">
        <div class="article-item" v-for="article in articles" :key="article.id">
          <div class="article-info">
            <div class="article-title">{{ article.title }}</div>
            <div class="article-meta">
              <span class="status" :class="article.status">{{ article.status === 'published' ? $t('article_management_page.published') : $t('article_management_page.draft') }}</span>
              <span class="date">{{ formatDate(article.created_at) }}</span>
              <span class="stats">
                <el-icon><View /></el-icon>
                {{ article.view_count }}
                <el-icon><Star /></el-icon>
                {{ article.like_count }}
              </span>
            </div>
          </div>
          
          <div class="article-actions">
            <button v-if="article.status === 'draft'" class="action-btn publish-btn" @click="handlePublish(article.id)">
              <el-icon><Promotion /></el-icon>
              {{ $t('article_management.publish') }}
            </button>
            <router-link :to="`/articles/${article.id}`" class="action-btn view-btn">
              <el-icon><View /></el-icon>
              {{ $t('article_management.preview') }}
            </router-link>
            <router-link :to="`/admin/articles/${article.id}/edit`" class="action-btn edit-btn">
              <el-icon><Edit /></el-icon>
              {{ $t('article_management.edit') }}
            </router-link>
            <button class="action-btn delete-btn" @click="handleDelete(article.id)">
              <el-icon><Delete /></el-icon>
              {{ $t('article_management.delete') }}
            </button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <el-icon size="48"><Document /></el-icon>
        <p>{{ $t('article_management.no_articles') }}</p>
        <router-link to="/admin/articles/new" class="tech-button">{{ $t('article_management.create_first_article') }}</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useArticleStore } from '@/stores'
import { formatDate } from '@/utils'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const articleStore = useArticleStore()
const isLoading = ref(false)
const articles = computed(() => articleStore.articles)

const loadArticles = async () => {
  try {
    isLoading.value = true
    await articleStore.fetchArticles({ page: 1, pageSize: 50 })
  } catch (error) {
    ElMessage.error(t('article_management_page.load_fail'))
  } finally {
    isLoading.value = false
  }
}

const handlePublish = async (id: number) => {
  try {
    await ElMessageBox.confirm(t('article_management_page.publish_confirm_text'), t('article_management_page.publish_confirm_title'), {
      confirmButtonText: t('article_management_page.confirm'),
      cancelButtonText: t('article_management_page.cancel'),
      type: 'info'
    })
    
    await articleStore.publishArticle(id)
    ElMessage.success(t('article_management_page.publish_success'))
    await loadArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('article_management_page.publish_fail'))
    }
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm(t('article_management_page.delete_confirm_text'), t('article_management_page.delete_confirm_title'), {
      confirmButtonText: t('article_management_page.confirm'),
      cancelButtonText: t('article_management_page.cancel'),
      type: 'warning'
    })
    
    await articleStore.deleteArticle(id)
    ElMessage.success(t('article_management_page.delete_success'))
    await loadArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('article_management_page.delete_fail'))
    }
  }
}

const handleUnpublish = async (id: number) => {
  try {
    await ElMessageBox.confirm(t('article_management_page.unpublish_confirm_text'), t('article_management_page.unpublish_confirm_title'), {
      confirmButtonText: t('article_management_page.confirm'),
      cancelButtonText: t('article_management_page.cancel'),
      type: 'warning'
    })
    
    await articleStore.unpublishArticle(id)
    ElMessage.success(t('article_management_page.unpublish_success'))
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('article_management_page.unpublish_fail'))
    }
  }
}

const fileInput = ref<HTMLInputElement | null>(null)

const triggerImport = () => {
  fileInput.value?.click()
}

const handleExport = async () => {
  try {
    await articleStore.exportArticles()
    ElMessage.success(t('article_management_page.export_success'))
  } catch (error) {
    ElMessage.error(t('article_management_page.export_fail'))
  }
}

const handleImport = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    try {
      await articleStore.importArticles(file)
      ElMessage.success(t('article_management_page.import_success'))
      await loadArticles()
    } catch (error) {
      ElMessage.error(t('article_management_page.import_fail'))
    }
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
.article-management {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  color: var(--text-primary);
  margin: 0;
}

.articles-list {
  padding: 1.5rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 0;
  gap: 1rem;
  color: var(--text-secondary);
}

.article-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 0;
  border-bottom: 1px solid var(--border-color);
}

.article-item:last-child {
  border-bottom: none;
}

.article-info {
  flex: 1;
}

.article-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
}

.status.published {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.status.draft {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
}

.stats {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.article-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.875rem;
  text-decoration: none;
  transition: all 0.3s ease;
  cursor: pointer;
  background: transparent;
}

.view-btn {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.view-btn:hover {
  background: var(--primary-color);
  color: white;
}

.edit-btn {
  color: #fbbf24;  
  border-color: #fbbf24;
}

.edit-btn:hover {
  background: #fbbf24;
  color: white;
}

.publish-btn {
  color: #22c55e;
  border-color: #22c55e;
}

.publish-btn:hover {
  background: #22c55e;
  color: white;
}

.delete-btn {
  color: #ef4444;
  border-color: #ef4444;
}

.delete-btn:hover {
  background: #ef4444;
  color: white;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 0;
  gap: 1rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .article-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .article-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>