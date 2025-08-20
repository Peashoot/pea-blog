import { apiClient } from './client'
import type { Comment, CreateCommentRequest, CommentListResponse } from '@/types'

export const commentApi = {
  getCommentsByArticleId: (articleId: number, page: number, pageSize: number): Promise<CommentListResponse> => {
    return apiClient.get(`/articles/${articleId}/comments`, { params: { page, page_size: pageSize } })
  },

  getRepliesByCommentId: (commentId: number, page: number, pageSize: number): Promise<CommentListResponse> => {
    return apiClient.get(`/comments/${commentId}/replies`, { params: { page, page_size: pageSize } })
  },

  createComment: (data: CreateCommentRequest): Promise<Comment> => {
    return apiClient.post('/comments', data)
  },

  deleteComment: (id: number, fingerprint?: string): Promise<void> => {
    const config = fingerprint ? { data: { fingerprint } } : {}
    return apiClient.delete(`/comments/${id}`, config)
  }
}