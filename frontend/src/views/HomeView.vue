<template>
  <div class="home-page">
    <Navbar />
    
    <main class="main-content">
      <div class="container">
        <div class="hero-section">
          <h1 class="hero-title gradient-text">{{ $t('home.welcome') }}</h1>
          <p class="hero-subtitle">{{ $t('home.subtitle') }}</p>
        </div>

        <SearchBar @search="handleSearch" @clear="handleClear" />

        <div class="content-section">
          <div v-if="articleStore.isLoading" class="loading-container">
            <div class="loading-spinner"></div>
            <p>加载中...</p>
          </div>

          <div v-else-if="articles.length > 0" class="articles-grid">
            <ArticleCard 
              v-for="article in articles" 
              :key="article.id" 
              :article="article"
              @click="goToArticle(article.id)"
            />
          </div>

          <div v-else class="empty-state">
            <el-icon size="48" color="var(--text-secondary)"><Document /></el-icon>
            <p>{{ $t('home.no_articles') }}</p>
          </div>

          <div v-if="hasMoreComments" class="load-more-comments">
              <button class="tech-button" @click="loadComments(true)" :disabled="isLoadingComments">
                {{ isLoadingComments ? '加载中...' : $t('home.load_more') }}
              </button>
            </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useArticleStore, useAuthStore } from '@/stores'
import { useInfiniteScroll } from '@/composables'
import Navbar from '@/components/Navbar.vue'
import SearchBar from '@/components/SearchBar.vue'
import ArticleCard from '@/components/ArticleCard.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const articleStore = useArticleStore()
const authStore = useAuthStore()

const currentSearchParams = ref<any>({})

const articles = computed(() => articleStore.articles)

const handleSearch = async (params: any) => {
  try {
    currentSearchParams.value = params
    if (params.keyword) {
      await articleStore.searchArticles(params.keyword, {
        tags: params.tags.length ? params.tags : undefined,
        sortBy: params.sortBy,
        sortOrder: params.sortOrder,
        page: 1
      })
    } else {
      await articleStore.fetchPublishedArticles({
        tags: params.tags.length ? params.tags.join(',') : undefined,
        sortBy: params.sortBy,
        sortOrder: params.sortOrder,
        page: 1
      })
    }
  } catch (error) {
    ElMessage.error('搜索失败')
  }
}

const handleClear = async () => {
  try {
    currentSearchParams.value = {}
    await articleStore.fetchArticles({ page: 1 })
  } catch (error) {
    ElMessage.error('加载文章失败')
  }
}

const goToArticle = (id: number) => {
  router.push(`/articles/${id}`)
}

const loadMore = async () => {
  try {
    const nextPage = articleStore.currentPage + 1
    if (currentSearchParams.value.keyword) {
      await articleStore.searchArticles(currentSearchParams.value.keyword, {
        ...currentSearchParams.value,
        page: nextPage
      })
    } else {
      await articleStore.fetchPublishedArticles({
        ...currentSearchParams.value,
        page: nextPage
      })
    }
  } catch (error) {
    ElMessage.error('加载更多文章失败')
  }
}

useInfiniteScroll(() => {
  if (articleStore.hasMore && !articleStore.isLoading) {
    loadMore()
  }
})

onMounted(async () => {
  try {
    await authStore.initAuth()
    await articleStore.fetchPublishedArticles({ page: 1 })
  } catch (error) {
    ElMessage.error('初始化失败')
  }
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
}

.main-content {
  padding-top: 80px;
  padding-bottom: 2rem;
}

.hero-section {
  text-align: center;
  padding: 3rem 0;
  margin-bottom: 2rem;
}

.hero-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
  line-height: 1.2;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  margin: 0;
}

.content-section {
  margin-top: 2rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  gap: 1rem;
  color: var(--text-secondary);
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  color: var(--text-secondary);
  gap: 1rem;
}

.load-more {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }
  
  .hero-subtitle {
    font-size: 1rem;
  }
  
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .main-content {
    padding-top: 70px;
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 2rem 0;
  }
  
  .hero-title {
    font-size: 1.75rem;
  }
}
</style>
