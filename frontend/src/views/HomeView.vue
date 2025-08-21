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
            <p>{{ $t('common.loading') }}</p>
          </div>

          <div v-else-if="articles.length > 0" class="articles-grid">
            <ArticleCard 
              v-for="article in articles" 
              :key="article.id" 
              :article="article"
              @click="goToArticle(article.title)"
            />
          </div>

          <div v-else class="empty-state">
            <el-icon size="48" color="var(--text-secondary)"><Document /></el-icon>
            <p>{{ $t('home.no_articles') }}</p>
          </div>

          <div v-if="articleStore.hasMore" class="load-more-comments">
              <button class="tech-button" @click="loadMore" :disabled="articleStore.isLoading">
                {{ articleStore.isLoading ? $t('common.loading') : $t('home.load_more') }}
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
import { useI18n } from 'vue-i18n'

const router = useRouter()
const articleStore = useArticleStore()
const authStore = useAuthStore()
const { t } = useI18n()

const currentSearchParams = ref<any>({})

const articles = computed(() => articleStore.articles)

const handleSearch = async (params: any) => {
  try {
    currentSearchParams.value = params
    if (params.keyword) {
      await articleStore.searchArticles(params.keyword, {
        tags: params.tags.length ? params.tags : undefined,
        sort_by: params.sortBy,
        sort_order: params.sortOrder,
        page: 1
      })
    } else {
      await articleStore.fetchPublishedArticles({
        tags: params.tags.length ? params.tags.join(',') : undefined,
        sort_by: params.sortBy,
        sort_order: params.sortOrder,
        page: 1
      })
    }
  } catch (error) {
    ElMessage.error(t('common.search_failed'))
  }
}

const handleClear = async () => {
  try {
    currentSearchParams.value = {}
    await articleStore.fetchArticles({ page: 1 })
  } catch (error) {
    ElMessage.error(t('common.load_articles_failed'))
  }
}

const goToArticle = (title: string) => {
  router.push(`/articles/${encodeURIComponent(title)}`)
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
    ElMessage.error(t('common.load_more_articles_failed'))
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
    ElMessage.error(t('common.init_failed'))
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
