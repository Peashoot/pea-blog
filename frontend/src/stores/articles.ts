import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { articleApi } from '@/api'
import type { Article, SearchParams, ArticleListResponse } from '@/types'

export const useArticleStore = defineStore('article', () => {
  const articles = ref<Article[]>([])
  const currentArticle = ref<Article | null>(null)
  const isLoading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  const hasMore = computed(() => articles.value.length < total.value)

  const fetchArticles = async (params?: SearchParams) => {
    try {
      isLoading.value = true
      const response: ArticleListResponse = await articleApi.getArticles({
        page: currentPage.value,
        pageSize: pageSize.value,
        ...params
      })
      
      if (params?.page === 1 || !params?.page) {
        articles.value = response.articles || []
      } else {
        articles.value.push(...response.articles)
      }
      
      total.value = response.total
      currentPage.value = response.page
      pageSize.value = response.page_size
    } catch (error) {
      console.error('Fetch articles error:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const fetchPublishedArticles = async (params?: SearchParams) => {
    try {
      isLoading.value = true
      const response: ArticleListResponse = await articleApi.getPublishedArticles({
        page: currentPage.value,
        pageSize: pageSize.value,
        ...params
      })
      
      if (params?.page === 1 || !params?.page) {
        articles.value = response.articles || []
      } else {
        articles.value.push(...response.articles)
      }
      
      total.value = response.total
      currentPage.value = response.page
      pageSize.value = response.page_size
    } catch (error) {
      console.error('Fetch published articles error:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const fetchArticleById = async (id: number) => {
    try {
      isLoading.value = true
      const article = await articleApi.getArticleById(id)
      currentArticle.value = article
      return article
    } catch (error) {
      console.error('Fetch article error:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const createArticle = async (articleData: any) => {
    try {
      const newArticle = await articleApi.createArticle(articleData)
      articles.value.unshift(newArticle)
      return newArticle
    } catch (error) {
      console.error('Create article error:', error)
      throw error
    }
  }

  const updateArticle = async (articleData: any) => {
    try {
      const updatedArticle = await articleApi.updateArticle(articleData)
      const index = articles.value.findIndex(a => a.id === updatedArticle.id)
      if (index !== -1) {
        articles.value[index] = updatedArticle
      }
      if (currentArticle.value?.id === updatedArticle.id) {
        currentArticle.value = updatedArticle
      }
      return updatedArticle
    } catch (error) {
      console.error('Update article error:', error)
      throw error
    }
  }

  const publishArticle = async (id: number) => {
    try {
      const updatedArticle = await articleApi.updateArticle({ id, status: 'published' })
      const index = articles.value.findIndex(a => a.id === updatedArticle.id)
      if (index !== -1) {
        articles.value[index] = updatedArticle
      }
      if (currentArticle.value?.id === updatedArticle.id) {
        currentArticle.value = updatedArticle
      }
      return updatedArticle
    } catch (error) {
      console.error('Publish article error:', error)
      throw error
    }
  }

  const deleteArticle = async (id: number) => {
    try {
      await articleApi.deleteArticle(id)
      articles.value = articles.value.filter(a => a.id !== id)
      if (currentArticle.value?.id === id) {
        currentArticle.value = null
      }
    } catch (error) {
      console.error('Delete article error:', error)
      throw error
    }
  }

  const likeArticle = async (id: number) => {
    try {
      await articleApi.likeArticle(id)
      const article = articles.value.find(a => a.id === id)
      if (article) {
        article.like_count++
      }
      if (currentArticle.value?.id === id) {
        currentArticle.value.like_count++
      }
    } catch (error) {
      console.error('Like article error:', error)
      throw error
    }
  }

  const unlikeArticle = async (id: number) => {
    try {
      await articleApi.unlikeArticle(id)
      const article = articles.value.find(a => a.id === id)
      if (article) {
        article.like_count--
      }
      if (currentArticle.value?.id === id) {
        currentArticle.value.like_count--
      }
    } catch (error) {
      console.error('Unlike article error:', error)
      throw error
    }
  }

  const searchArticles = async (keyword: string, params?: Omit<SearchParams, 'keyword'>) => {
    try {
      isLoading.value = true
      const response = await articleApi.searchArticles(keyword, params)
      articles.value = response.articles
      total.value = response.total
      currentPage.value = response.page
      pageSize.value = response.page_size
      return response
    } catch (error) {
      console.error('Search articles error:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const resetState = () => {
    articles.value = []
    currentArticle.value = null
    total.value = 0
    currentPage.value = 1
  }

  return {
    articles,
    currentArticle,
    isLoading,
    total,
    currentPage,
    pageSize,
    hasMore,
    fetchArticles,
    fetchPublishedArticles,
    fetchArticleById,
    createArticle,
    updateArticle,
    publishArticle,
    deleteArticle,
    likeArticle,
    unlikeArticle,
    searchArticles,
    resetState
  }
})