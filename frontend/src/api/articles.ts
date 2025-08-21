import { apiClient } from './client'
import type {
  Article,
  ArticleListResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
  SearchParams
} from '@/types'

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
    return apiClient.post('/articles', data)
  },

  updateArticle: (data: UpdateArticleRequest): Promise<Article> => {
    return apiClient.put(`/articles/${data.id}`, data)
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