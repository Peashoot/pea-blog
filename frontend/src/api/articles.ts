import { apiClient } from './client'
import type {
  Article,
  ArticleListResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
  SearchParams
} from '@/types'
import { parseDateTimeFromLocal, formatDateTimeForAPI } from '@/utils'

export const articleApi = {
  getArticles: (params?: SearchParams): Promise<ArticleListResponse> => {
    return apiClient.get('/articles', { params })
  },

  getPublishedArticles: (params?: SearchParams): Promise<ArticleListResponse> => {
    return apiClient.get('/articles/published', { params })
  },

  getArticleById: (id: number): Promise<Article> => {
    return apiClient.get(`/articles/${id}`)
  },

  getArticleByTitle: (title: string): Promise<Article> => {
    return apiClient.get(`/articles/title/${encodeURIComponent(title)}`)
  },

  createArticle: (data: CreateArticleRequest): Promise<Article> => {
    // 转换数据格式以匹配后端期望的格式
    const requestData = { ...data }
    
    // 如果状态为 scheduled 且有 published_at，确保是有效的日期格式并包含时区
    if (requestData.status === 'scheduled' && requestData.published_at) {
      // 解析本地时间字符串并转换为带时区的ISO格式
      const localDate = parseDateTimeFromLocal(requestData.published_at)
      if (!isNaN(localDate.getTime())) {
        requestData.published_at = formatDateTimeForAPI(localDate)
      }
    }
    
    return apiClient.post('/articles', requestData)
  },

  updateArticle: (data: UpdateArticleRequest): Promise<Article> => {
    // 转换数据格式以匹配后端期望的格式
    const requestData = { ...data }
    
    // 如果状态为 scheduled 且有 published_at，确保是有效的日期格式并包含时区
    if (requestData.status === 'scheduled' && requestData.published_at) {
      // 解析本地时间字符串并转换为带时区的ISO格式
      const localDate = parseDateTimeFromLocal(requestData.published_at)
      if (!isNaN(localDate.getTime())) {
        requestData.published_at = formatDateTimeForAPI(localDate)
      }
    }
    
    return apiClient.put(`/articles/${data.id}`, requestData)
  },

  deleteArticle: (id: number): Promise<void> => {
    return apiClient.delete(`/articles/${id}`)
  },

  likeArticle: (id: number): Promise<void> => {
    return apiClient.post(`/articles/${id}/like`)
  },

  unlikeArticle: (id: number): Promise<void> => {
    return apiClient.delete(`/articles/${id}/like`)
  },

  searchArticles: (keyword: string, params?: Omit<SearchParams, 'keyword'>): Promise<ArticleListResponse> => {
    return apiClient.get('/articles/search', {
      params: { keyword, ...params }
    })
  },

  exportArticles: (): Promise<Blob> => {
    return apiClient.get('/articles/export', { responseType: 'blob' })
  },

  importArticles: (file: File): Promise<void> => {
    const formData = new FormData()
    formData.append('file', file)
    return apiClient.post('/articles/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  unpublishArticle: (id: number): Promise<void> => {
    return apiClient.post(`/articles/${id}/unpublish`)
  },

  uploadImage: (file: File): Promise<{ url: string }> => {
    const formData = new FormData()
    formData.append('file', file)
    return apiClient.post('/images/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}